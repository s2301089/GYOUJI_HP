package repository

import (
	"database/sql"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

type RelayRepository struct {
	DB *sql.DB
}

func NewRelayRepository(db *sql.DB) *RelayRepository {
	return &RelayRepository{DB: db}
}

// 得点登録（順位ごとにclass_idを登録し、team_pointsに加算）
func (r *RelayRepository) RegisterRelayResults(relayType string, classIDs []int) error {
	// 順位ごとの得点
	scores := []int{30, 25, 20, 15, 10, 5}
	for i, classID := range classIDs {
		if i >= len(scores) {
			break
		}
		// relay_resultsテーブルに登録
		_, err := r.DB.Exec("INSERT INTO relay_results (relay_type, relay_rank, class_id) VALUES (?, ?, ?)", relayType, i+1, classID)
		if err != nil {
			return err
		}
		// team_pointsテーブルに加算
		var scoreCol string
		if relayType == "A" {
			scoreCol = "relay_A_score"
		} else {
			scoreCol = "relay_B_score"
		}
		_, err = r.DB.Exec("UPDATE team_points SET "+scoreCol+" = ? WHERE class_id = ?", scores[i], classID)
		if err != nil {
			return err
		}
	}
	return nil
}

// 得点・順位取得
func (r *RelayRepository) GetRelayResults(relayType string) ([]model.RelayResult, error) {
	rows, err := r.DB.Query("SELECT id, relay_type, relay_rank, class_id, created_at FROM relay_results WHERE relay_type = ? ORDER BY relay_rank ASC", relayType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []model.RelayResult
	for rows.Next() {
		var rr model.RelayResult
		if err := rows.Scan(&rr.ID, &rr.RelayType, &rr.RelayRank, &rr.ClassID, &rr.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, rr)
	}
	return results, nil
}
