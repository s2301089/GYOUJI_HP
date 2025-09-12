package repository

import (
	"database/sql"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

type ScoreRepository struct {
	DB *sql.DB
}

func NewScoreRepository(db *sql.DB) *ScoreRepository { return &ScoreRepository{DB: db} }

func (r *ScoreRepository) GetScores() ([]model.ScoreBreakdown, error) {
	// 各クラス（teams）ごとにポイント内訳を集計
	query := `
        SELECT t.name AS class_name,
               COALESCE(t.init_score,0) AS init_score,
               COALESCE(t.attendance_score,0) AS attendance_score,
               COALESCE(SUM(CASE WHEN tp.point_type='win' THEN tp.points END),0) AS win_points,
               COALESCE(SUM(CASE WHEN tp.point_type='final_bonus_winner' THEN tp.points END),0) AS final_winner_bonus,
               COALESCE(SUM(CASE WHEN tp.point_type='final_bonus_runnerup' THEN tp.points END),0) AS final_runnerup_bonus,
               COALESCE(SUM(CASE WHEN tp.point_type='bronze_bonus_winner' THEN tp.points END),0) AS bronze_winner_bonus,
               COALESCE(SUM(CASE WHEN tp.point_type='bronze_bonus_runnerup' THEN tp.points END),0) AS bronze_runnerup_bonus,
               COALESCE(SUM(CASE WHEN tp.point_type='rainy_loser_champion' THEN tp.points END),0) AS rainy_loser_champion_bonus
        FROM teams t
        LEFT JOIN team_points tp ON tp.team_id = t.id
        GROUP BY t.id, t.name, t.init_score, t.attendance_score
    `
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.ScoreBreakdown
	for rows.Next() {
		var s model.ScoreBreakdown
		if err := rows.Scan(&s.ClassName, &s.InitScore, &s.AttendanceScore, &s.WinPoints, &s.FinalWinnerBonus, &s.FinalRunnerupBonus, &s.BronzeWinnerBonus, &s.BronzeRunnerupBonus, &s.RainyLoserChampion); err != nil {
			return nil, err
		}
		s.TotalExcludingInit = s.AttendanceScore + s.WinPoints + s.FinalWinnerBonus + s.FinalRunnerupBonus + s.BronzeWinnerBonus + s.BronzeRunnerupBonus + s.RainyLoserChampion
		s.TotalIncludingInit = s.InitScore + s.TotalExcludingInit
		results = append(results, s)
	}
	return results, nil
}
