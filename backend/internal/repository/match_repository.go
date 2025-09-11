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
	query := `
        SELECT m.tournament_id, m.team1_id, m.team2_id, m.next_match_id, t.sport
        FROM matches m JOIN tournaments t ON m.tournament_id = t.id
        WHERE m.id = ?`
	err = tx.QueryRow(query, matchID).Scan(&tournamentID, &team1ID, &team2ID, &nextMatchID, &matchSport)
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
	var winnerTeamID sql.NullInt64
	if team1Score > team2Score {
		winnerTeamID = team1ID
	} else if team2Score > team1Score {
		winnerTeamID = team2ID
	}
	// 引き分けの場合は winnerTeamID は NULL (invalid) のまま

	// 現在の試合のスコアと勝者を更新
	_, err = tx.Exec(`UPDATE matches SET team1_score = ?, team2_score = ?, winner_team_id = ? WHERE id = ?`, team1Score, team2Score, winnerTeamID, matchID)
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
