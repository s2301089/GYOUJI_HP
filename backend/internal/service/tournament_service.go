package service

import (
	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"
)

// テスト用インターフェース
type TournamentServiceInterface interface {
	GetTournamentsBySport(sport string) (interface{}, error)
}

type TournamentService struct {
	Repo *repository.TournamentRepository
}

func NewTournamentService(r *repository.TournamentRepository) *TournamentService {
	return &TournamentService{Repo: r}
}

// GetTournamentsBySport godoc
// @Summary      指定競技のトーナメント情報を取得
// @Description  競技名をパスパラメータとして受け取り、関連するトーナメント（晴天時・雨天時など）の情報を返します。
// @Tags         Tournaments
// @Accept       json
// @Produce      json
// @Param        sport   path      string  true  "競技名 (volleyball, table_tennis, soccer)"
// @Success      200  {object}  interface{}
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/tournaments/{sport} [get]
func (s *TournamentService) GetTournamentsBySport(sport string) (interface{}, error) {
	return s.Repo.GetTournamentsBySport(sport)
}
