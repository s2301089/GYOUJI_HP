package repository

import (
	"database/sql"

	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

type TournamentRepository struct {
	DB *sql.DB
}

func NewTournamentRepository(db *sql.DB) *TournamentRepository {
	return &TournamentRepository{DB: db}
}

// 指定競技のトーナメント情報を取得
func (r *TournamentRepository) GetTournamentsBySport(sport string) (interface{}, error) {
	// DBから指定sportのトーナメント一覧を取得
	tournamentsRows, err := r.DB.Query("SELECT id, name, sport, weather_condition FROM tournaments WHERE sport = ?", sport)
	if err != nil {
		return nil, err
	}
	defer tournamentsRows.Close()

	var brackets []model.TournamentBracket

	for tournamentsRows.Next() {
		var t model.TournamentBracket
		if err := tournamentsRows.Scan(&t.Tournament.ID, &t.Tournament.Name, &t.Tournament.Sport, &t.Tournament.WeatherCondition); err != nil {
			return nil, err
		}

		// チーム取得
		teamsRows, err := r.DB.Query("SELECT id, name, tournament_id FROM teams WHERE tournament_id = ?", t.Tournament.ID)
		if err != nil {
			return nil, err
		}
		for teamsRows.Next() {
			var team model.Team
			if err := teamsRows.Scan(&team.ID, &team.Name, &team.TournamentID); err != nil {
				teamsRows.Close()
				return nil, err
			}
			t.Teams = append(t.Teams, team)
		}
		teamsRows.Close()

		// 試合取得
		matchesRows, err := r.DB.Query("SELECT id, tournament_id, round, match_number_in_round, team1_id, team2_id, team1_score, team2_score, winner_team_id, next_match_id FROM matches WHERE tournament_id = ?", t.Tournament.ID)
		if err != nil {
			return nil, err
		}
		for matchesRows.Next() {
			var match model.Match
			if err := matchesRows.Scan(&match.ID, &match.TournamentID, &match.Round, &match.MatchNumber, &match.Team1ID, &match.Team2ID, &match.Team1Score, &match.Team2Score, &match.WinnerTeamID, &match.NextMatchID); err != nil {
				matchesRows.Close()
				return nil, err
			}
			t.Matches = append(t.Matches, match)
		}
		matchesRows.Close()

		brackets = append(brackets, t)
	}

	return brackets, nil
}
