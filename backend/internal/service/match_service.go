package service

import "github.com/saku0512/GYOUJI_HP/backend/internal/repository"

// インターフェース（テスト用）
type MatchServiceInterface interface {
	UpdateMatchScore(matchID int, team1Score int, team2Score int, user interface{}) (interface{}, error)
}

type MatchService struct {
	Repo *repository.MatchRepository
}

func NewMatchService(r *repository.MatchRepository) *MatchService {
	return &MatchService{Repo: r}
}

func (s *MatchService) UpdateMatchScore(matchID int, team1Score int, team2Score int, user interface{}) (interface{}, error) {
	return s.Repo.UpdateMatchScore(matchID, team1Score, team2Score, user)
}
