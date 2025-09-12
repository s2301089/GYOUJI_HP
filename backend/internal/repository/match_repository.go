package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

type MatchRepository struct {
	DB *sql.DB
}

func NewMatchRepository(db *sql.DB) *MatchRepository {
	return &MatchRepository{DB: db}
}

// UpdateMatchScore は試合のスコアを更新し、勝者を次の試合へ進出させます。
func (r *MatchRepository) UpdateMatchScore(matchID int, team1Score int, team2Score int, user interface{}) (interface{}, error) {
	// ユーザー情報の検証
	userMap, ok := user.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	role, _ := userMap["role"].(string)
	assignedSport, _ := userMap["assigned_sport"].(string)

	// データベース操作をトランザクション内で実行
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// 何か問題が発生した場合はロールバックする
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err)
			tx.Rollback()
		}
	}()

	// 試合情報と次の試合IDを取得
	// next_match_idも一緒に取得
	var tournamentID int
	var team1ID, team2ID sql.NullInt64 // NULLを許容する型に変更
	var nextMatchID sql.NullInt64      // NULLを許容する型
	var matchSport string
	var currentRound int
	query := `
        SELECT m.tournament_id, m.round, m.team1_id, m.team2_id, m.next_match_id, t.sport
        FROM matches m JOIN tournaments t ON m.tournament_id = t.id
        WHERE m.id = ?`
	err = tx.QueryRow(query, matchID).Scan(&tournamentID, &currentRound, &team1ID, &team2ID, &nextMatchID, &matchSport)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to query match: %w", err)
	}

	// 権限判定
	if !(role == "superroot" || (role == "admin" && assignedSport == matchSport)) {
		return nil, fmt.Errorf("forbidden")
	}

	// 勝者判定
	var winnerTeamID, loserTeamID sql.NullInt64
	if team1Score > team2Score {
		winnerTeamID = team1ID
		loserTeamID = team2ID
	} else if team2Score > team1Score {
		winnerTeamID = team2ID
		loserTeamID = team1ID
	}
	// 引き分けの場合は winnerTeamID は NULL (invalid) のまま

	// 現在の試合のスコアと勝者を更新
	_, err = tx.Exec(`UPDATE matches SET team1_score = ?, team2_score = ?, winner_team_id = ?, status = 'finished' WHERE id = ? `, team1Score, team2Score, winnerTeamID, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to update current match: %w", err)
	}

	// 勝者を次の試合に進めるロジックを追加
	if winnerTeamID.Valid && nextMatchID.Valid {
		// 次の試合の情報を取得して、どちらのスロットが空いているか確認
		var nextTeam1ID, nextTeam2ID sql.NullInt64
		err = tx.QueryRow(`SELECT team1_id, team2_id FROM matches WHERE id = ?`, nextMatchID.Int64).Scan(&nextTeam1ID, &nextTeam2ID)
		if err != nil {
			// ここでエラーが発生しても、スコア更新は完了しているためロールバック
			return nil, fmt.Errorf("failed to query next match: %w", err)
		}

		// 空いているスロットに勝者IDを設定
		var updateNextMatchQuery string
		if !nextTeam1ID.Valid {
			// team1_idが空いている場合
			updateNextMatchQuery = `UPDATE matches SET team1_id = ? WHERE id = ?`
		} else if !nextTeam2ID.Valid {
			// team2_idが空いている場合
			updateNextMatchQuery = `UPDATE matches SET team2_id = ? WHERE id = ?`
		}

		if updateNextMatchQuery != "" {
			_, err = tx.Exec(updateNextMatchQuery, winnerTeamID.Int64, nextMatchID.Int64)
			if err != nil {
				return nil, fmt.Errorf("failed to update next match: %w", err)
			}
		}
	}

	// 準決勝のみ敗者を三位決定戦にすすめる
	if loserTeamID.Valid {
		var maxRound int
		err = tx.QueryRow(`SELECT MAX(round) FROM matches WHERE tournament_id = ?`, tournamentID).Scan(&maxRound)
		if err != nil {
			return nil, fmt.Errorf("failed to get max round for tournament: %w", err)
		}

		// 現在の試合が準決勝の場合のみ処理
		if currentRound == maxRound-1 {
			var thirdPlaceMatchID sql.NullInt64
			// 3位決定戦を特定 (最終ラウンドで試合番号が一番若い試合と想定)
			query3rd := `
                SELECT id FROM matches 
                WHERE tournament_id = ? AND round = ? 
                ORDER BY match_number_in_round ASC LIMIT 1`
			err = tx.QueryRow(query3rd, tournamentID, maxRound).Scan(&thirdPlaceMatchID)

			if err == nil && thirdPlaceMatchID.Valid {
				var p1, p2 sql.NullInt64
				err = tx.QueryRow(`SELECT team1_id, team2_id FROM matches WHERE id = ?`, thirdPlaceMatchID.Int64).Scan(&p1, &p2)
				if err != nil {
					return nil, fmt.Errorf("failed to query third place match teams: %w", err)
				}

				var updateQuery string
				if !p1.Valid {
					updateQuery = `UPDATE matches SET team1_id = ? WHERE id = ?`
				} else if !p2.Valid {
					updateQuery = `UPDATE matches SET team2_id = ? WHERE id = ?`
				}

				if updateQuery != "" {
					_, err = tx.Exec(updateQuery, loserTeamID.Int64, thirdPlaceMatchID.Int64)
					if err != nil {
						return nil, fmt.Errorf("failed to update third place match for loser: %w", err)
					}
				}
			} else if err != nil && err != sql.ErrNoRows {
				return nil, fmt.Errorf("failed to find third place match: %w", err)
			}
		}
	}

	// トランザクションをコミット
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 更新後の試合情報を取得して返す（トランザクションの外で実行）
	return r.getMatchByID(matchID)
}

// getMatchByID は指定されたIDの試合情報を取得するヘルパー関数です。
func (r *MatchRepository) getMatchByID(matchID int) (interface{}, error) {
	var match model.Match
	query := `
        SELECT id, tournament_id, round, match_number_in_round,
               team1_id, team2_id, team1_score, team2_score, winner_team_id, next_match_id
        FROM matches
        WHERE id = ?`
	err := r.DB.QueryRow(query, matchID).Scan(
		&match.ID, &match.TournamentID, &match.Round, &match.MatchNumber,
		&match.Team1ID, &match.Team2ID, &match.Team1Score, &match.Team2Score,
		&match.WinnerTeamID, &match.NextMatchID,
	)
	if err != nil {
		return nil, err
	}
	return match, nil
}

// GetMatchesBySport は指定された競技の試合一覧を取得します。
func (r *MatchRepository) GetMatchesBySport(sport string) ([]model.MatchResponse, error) {
	query := `
		SELECT m.id, m.match_number_in_round, m.round,
			   m.team1_id, t1.name, m.team2_id, t2.name,
			   m.team1_score, m.team2_score, m.winner_team_id, m.status, m.next_match_id, tr.name as tournament_name
		FROM matches m
		JOIN tournaments tr ON m.tournament_id = tr.id
		LEFT JOIN teams t1 ON m.team1_id = t1.id
		LEFT JOIN teams t2 ON m.team2_id = t2.id
		WHERE tr.sport = ?
		ORDER BY m.round, m.match_number_in_round`
	rows, err := r.DB.Query(query, sport)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.MatchResponse
	for rows.Next() {
		var res model.MatchResponse
		var team1Name, team2Name, status, tournamentName sql.NullString
		var team1Score, team2Score sql.NullInt64

		err := rows.Scan(
			&res.ID, &res.MatchNumberInRound, &res.Round,
			&res.Team1ID, &team1Name, &res.Team2ID, &team2Name,
			&team1Score, &team2Score, &res.WinnerTeamID, &status, &res.NextMatchID, &tournamentName,
		)
		if err != nil {
			return nil, err
		}

		// sql.Null* 型からポインタ型へ変換
		if team1Name.Valid {
			res.Team1Name = &team1Name.String
		}
		if team2Name.Valid {
			res.Team2Name = &team2Name.String
		}
		if status.Valid {
			res.Status = &status.String
		}
		if tournamentName.Valid {
			res.TournamentName = &tournamentName.String
		}
		if team1Score.Valid {
			score := int(team1Score.Int64)
			res.Team1Score = &score
		}
		if team2Score.Valid {
			score := int(team2Score.Int64)
			res.Team2Score = &score
		}

		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []model.MatchResponse{}, nil
	}

	return results, nil
}
