package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	userMap, ok := user.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	role, _ := userMap["role"].(string)
	assignedSport, _ := userMap["assigned_sport"].(string)

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Printf("Rolling back transaction due to error: %v", err)
			tx.Rollback()
		}
	}()

	// loser_next_match_idとトーナメント名も取得する
	var tournamentID int
	var team1ID, team2ID, nextMatchID, loserNextMatchID sql.NullInt64
	var matchSport, tournamentName string
	var currentRound int
	query := `
        SELECT m.tournament_id, m.round, t.name, m.team1_id, m.team2_id, 
               m.next_match_id, m.loser_next_match_id, t.sport
        FROM matches m JOIN tournaments t ON m.tournament_id = t.id
        WHERE m.id = ?`
	err = tx.QueryRow(query, matchID).Scan(&tournamentID, &currentRound, &tournamentName, &team1ID, &team2ID, &nextMatchID, &loserNextMatchID, &matchSport)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not found")
		}
		return nil, fmt.Errorf("failed to query initial match: %w", err)
	}

	if !(role == "superroot" || (role == "admin" && assignedSport == matchSport)) {
		return nil, fmt.Errorf("forbidden")
	}

	var winnerTeamID, loserTeamID sql.NullInt64
	if team1Score > team2Score {
		winnerTeamID = team1ID
		loserTeamID = team2ID
	} else if team2Score > team1Score {
		winnerTeamID = team2ID
		loserTeamID = team1ID
	}

	// DRY原則に基づき、更新ロジックを関数化
	err = r.updateMatchAndProgress(tx, matchID, team1Score, team2Score, winnerTeamID, loserTeamID, nextMatchID, loserNextMatchID)
	if err != nil {
		return nil, err
	}

	// --- ポイント付与ロジック ---
	if winnerTeamID.Valid {
		// 1) 勝利ポイント +10（決勝は除く）
		isFinalRound := false
		isBronzeMatch := false
		{
			var cnt int
			_ = tx.QueryRow(`SELECT COUNT(*) FROM matches WHERE tournament_id = ? AND round = ?`, tournamentID, currentRound).Scan(&cnt)
			if cnt == 1 {
				isFinalRound = true
			} else if cnt >= 2 {
				var maxNum, thisNum int
				_ = tx.QueryRow(`SELECT MAX(match_number_in_round) FROM matches WHERE tournament_id = ? AND round = ?`, tournamentID, currentRound).Scan(&maxNum)
				_ = tx.QueryRow(`SELECT match_number_in_round FROM matches WHERE id = ?`, matchID).Scan(&thisNum)
				if thisNum != maxNum {
					isBronzeMatch = true
				} else {
					isFinalRound = true
				}
			}
		}
		if !isFinalRound {
			if err := r.awardPoints(tx, int(winnerTeamID.Int64), matchSport, tournamentName, "win", 10, matchID); err != nil {
				return nil, err
			}
		}

		// 2) ベスト4ボーナス
		if isFinalRound {
			if err := r.awardPoints(tx, int(winnerTeamID.Int64), matchSport, tournamentName, "final_bonus_winner", 80, matchID); err != nil {
				return nil, err
			}
			if loserTeamID.Valid {
				if err := r.awardPoints(tx, int(loserTeamID.Int64), matchSport, tournamentName, "final_bonus_runnerup", 60, matchID); err != nil {
					return nil, err
				}
			}
		} else if isBronzeMatch {
			if err := r.awardPoints(tx, int(winnerTeamID.Int64), matchSport, tournamentName, "bronze_bonus_winner", 50, matchID); err != nil {
				return nil, err
			}
			if loserTeamID.Valid {
				if err := r.awardPoints(tx, int(loserTeamID.Int64), matchSport, tournamentName, "bronze_bonus_runnerup", 40, matchID); err != nil {
					return nil, err
				}
			}
		}

		// 3) 雨天時 敗者復活戦 優勝ボーナス +10
		if matchSport == "table_tennis" && (strings.Contains(tournamentName, "敗者戦左側") || strings.Contains(tournamentName, "敗者戦右側")) {
			if currentRound >= 2 {
				if err := r.awardPoints(tx, int(winnerTeamID.Int64), matchSport, tournamentName, "rainy_loser_champion", 10, matchID); err != nil {
					return nil, err
				}
			}
		}
	}

	// 卓球（晴天時）の場合、雨天時も同期更新するロジック
	if matchSport == "table_tennis" && tournamentName == "卓球（晴天時）" {
		var rainyMatchID int
		var rainyNextMatchID, rainyLoserNextMatchID sql.NullInt64
		var rainyTeam1ID, rainyTeam2ID sql.NullInt64

		rainyQuery := `
			SELECT m_rainy.id, m_rainy.next_match_id, m_rainy.loser_next_match_id, m_rainy.team1_id, m_rainy.team2_id, m_rainy.match_number_in_round
			FROM matches m_sunny
			JOIN tournaments t_rainy ON t_rainy.name = '卓球（雨天時）'
			JOIN matches m_rainy ON m_rainy.tournament_id = t_rainy.id 
				AND m_rainy.round = m_sunny.round 
				AND m_rainy.match_number_in_round = m_sunny.match_number_in_round
			WHERE m_sunny.id = ?`

		var rainyMatchNumberInRound sql.NullInt64
		err = tx.QueryRow(rainyQuery, matchID).Scan(&rainyMatchID, &rainyNextMatchID, &rainyLoserNextMatchID, &rainyTeam1ID, &rainyTeam2ID, &rainyMatchNumberInRound)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to find corresponding rainy match: %w", err)
		}

		if rainyMatchID != 0 {
			// 雨天側の勝者/敗者は雨天トーナメントのチームIDで決定する
			var rainyWinnerTeamID, rainyLoserTeamID sql.NullInt64
			if team1Score > team2Score {
				rainyWinnerTeamID = rainyTeam1ID
				rainyLoserTeamID = rainyTeam2ID
			} else if team2Score > team1Score {
				rainyWinnerTeamID = rainyTeam2ID
				rainyLoserTeamID = rainyTeam1ID
			}

			err = r.updateMatchAndProgress(tx, rainyMatchID, team1Score, team2Score, rainyWinnerTeamID, rainyLoserTeamID, rainyNextMatchID, rainyLoserNextMatchID)
			if err != nil {
				return nil, fmt.Errorf("failed to sync rainy match: %w", err)
			}

			// 1回戦の敗者を雨天の敗者復活戦トーナメントへ登録（重複防止チェックつき）
			if currentRound == 1 && rainyLoserTeamID.Valid && rainyMatchNumberInRound.Valid {
				var loserTournamentName string
				var loserEntryMatchNumber int
				switch rainyMatchNumberInRound.Int64 {
				case 1, 2:
					loserTournamentName = "卓球（雨天時・敗者戦左側）"
					loserEntryMatchNumber = 13
				case 3, 4:
					loserTournamentName = "卓球（雨天時・敗者戦左側）"
					loserEntryMatchNumber = 14
				case 5, 6:
					loserTournamentName = "卓球（雨天時・敗者戦右側）"
					loserEntryMatchNumber = 13
				case 7, 8:
					loserTournamentName = "卓球（雨天時・敗者戦右側）"
					loserEntryMatchNumber = 14
				}
				if loserTournamentName != "" {
					var loserEntryMatchID sql.NullInt64
					q := `
						SELECT m.id
						FROM tournaments tr
						JOIN matches m ON m.tournament_id = tr.id
						WHERE tr.name = ? AND m.round = 1 AND m.match_number_in_round = ?`
					err = tx.QueryRow(q, loserTournamentName, loserEntryMatchNumber).Scan(&loserEntryMatchID)
					if err == nil && loserEntryMatchID.Valid {
						// 既に同じチームが登録済みか確認し、未登録なら割り当て
						var existingTeam1, existingTeam2 sql.NullInt64
						if err := tx.QueryRow(`SELECT team1_id, team2_id FROM matches WHERE id = ?`, loserEntryMatchID.Int64).Scan(&existingTeam1, &existingTeam2); err == nil {
							if (existingTeam1.Valid && existingTeam1.Int64 == rainyLoserTeamID.Int64) || (existingTeam2.Valid && existingTeam2.Int64 == rainyLoserTeamID.Int64) {
								// すでに登録済み
							} else {
								if err := r.progressTeamToNextMatch(tx, rainyLoserTeamID, loserEntryMatchID); err != nil {
									return nil, fmt.Errorf("failed to progress loser to rainy consolation: %w", err)
								}
							}
						}
					}
				}
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.getMatchByID(matchID)
}

func (r *MatchRepository) updateMatchAndProgress(tx *sql.Tx, matchID int, team1Score int, team2Score int, winnerTeamID, loserTeamID, nextMatchID, loserNextMatchID sql.NullInt64) error {
	_, err := tx.Exec(`UPDATE matches SET team1_score = ?, team2_score = ?, winner_team_id = ?, status = 'finished' WHERE id = ?`,
		team1Score, team2Score, winnerTeamID, matchID)
	if err != nil {
		return fmt.Errorf("failed to update match id %d: %w", matchID, err)
	}
	if winnerTeamID.Valid && nextMatchID.Valid {
		err := r.progressTeamToNextMatch(tx, winnerTeamID, nextMatchID)
		if err != nil {
			return fmt.Errorf("failed to progress winner from match id %d: %w", matchID, err)
		}
	}
	if loserTeamID.Valid && loserNextMatchID.Valid {
		err := r.progressTeamToNextMatch(tx, loserTeamID, loserNextMatchID)
		if err != nil {
			return fmt.Errorf("failed to progress loser from match id %d: %w", matchID, err)
		}
	}
	return nil
}

// awardPoints は重複を避けつつ team_points にポイント行を追加します
func (r *MatchRepository) awardPoints(tx *sql.Tx, teamID int, sport, tournamentName, pointType string, points int, sourceMatchID int) error {
	// 重複チェックはDBのユニークキーでも担保されるが、先に存在確認で冪等に
	var exists int
	err := tx.QueryRow(`SELECT COUNT(1) FROM team_points WHERE team_id = ? AND source_match_id = ? AND point_type = ?`, teamID, sourceMatchID, pointType).Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}
	_, err = tx.Exec(`INSERT INTO team_points (team_id, sport, tournament_name, point_type, points, source_match_id) VALUES (?,?,?,?,?,?)`, teamID, sport, tournamentName, pointType, points, sourceMatchID)
	return err
}

func (r *MatchRepository) progressTeamToNextMatch(tx *sql.Tx, teamToProgress, nextMatchID sql.NullInt64) error {
	var team1, team2 sql.NullInt64
	err := tx.QueryRow(`SELECT team1_id, team2_id FROM matches WHERE id = ?`, nextMatchID.Int64).Scan(&team1, &team2)
	if err != nil {
		return err
	}

	var updateQuery string
	if !team1.Valid {
		updateQuery = `UPDATE matches SET team1_id = ? WHERE id = ?`
	} else if !team2.Valid {
		updateQuery = `UPDATE matches SET team2_id = ? WHERE id = ?`
	}

	if updateQuery != "" {
		_, err = tx.Exec(updateQuery, teamToProgress.Int64, nextMatchID.Int64)
		if err != nil {
			return err
		}
	}
	return nil
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
