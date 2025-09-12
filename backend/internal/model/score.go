package model

type ScoreBreakdown struct {
	ClassName           string `json:"class_name"`
	InitScore           int    `json:"init_score"`
	AttendanceScore     int    `json:"attendance_score"`
	WinPoints           int    `json:"win_points"`
	FinalWinnerBonus    int    `json:"final_winner_bonus"`
	FinalRunnerupBonus  int    `json:"final_runnerup_bonus"`
	BronzeWinnerBonus   int    `json:"bronze_winner_bonus"`
	BronzeRunnerupBonus int    `json:"bronze_runnerup_bonus"`
	RainyLoserChampion  int    `json:"rainy_loser_champion_bonus"`
	TotalExcludingInit  int    `json:"total_excluding_init"`
	TotalIncludingInit  int    `json:"total_including_init"`
}
