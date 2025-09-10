package model

type Tournament struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Sport            string `json:"sport"`
	WeatherCondition string `json:"weather_condition"`
}

type Team struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	TournamentID int64  `json:"tournament_id"`
}

type Match struct {
	ID           int64  `json:"id"`
	TournamentID int64  `json:"tournament_id"`
	Round        int    `json:"round"`
	MatchNumber  int    `json:"match_number_in_round"`
	Team1ID      int64  `json:"team1_id"`
	Team2ID      int64  `json:"team2_id"`
	Team1Score   *int   `json:"team1_score"`
	Team2Score   *int   `json:"team2_score"`
	WinnerTeamID *int64 `json:"winner_team_id"`
	NextMatchID  *int64 `json:"next_match_id"`
}

// Bracketry用レスポンス
// 複数トーナメントを返す場合は[]TournamentBracket

type TournamentBracket struct {
	Tournament Tournament `json:"tournament"`
	Teams      []Team     `json:"teams"`
	Matches    []Match    `json:"matches"`
}
