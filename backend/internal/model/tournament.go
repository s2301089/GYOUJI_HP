package model

// このファイルは、フロントエンドのライブラリ 'bracketry' が要求するデータ構造を定義します。

type BracketryData struct {
	Name        string                         `json:"name"`
	Rounds      []BracketryRound               `json:"rounds"`
	Matches     []BracketryMatch               `json:"matches"`
	Contestants map[string]BracketryContestant `json:"contestants"`
}

// BracketryRound defines the name of a round.
type BracketryRound struct {
	Name string `json:"name"`
}

// BracketryMatch represents a single match.
type BracketryMatch struct {
	RoundIndex  int             `json:"roundIndex"`
	Order       int             `json:"order"`
	Sides       []BracketrySide `json:"sides"`
	MatchStatus *string         `json:"matchStatus,omitempty"`
}

// BracketrySide represents one side of a match (a player or a BYE).
type BracketrySide struct {
	ContestantID *string `json:"contestantId,omitempty"`
	Scores       []Score `json:"scores,omitempty"`
	IsWinner     *bool   `json:"isWinner,omitempty"`
	Title        *string `json:"title,omitempty"` // Used for BYEs
}

// Score represents the score within a set or game. The target format requires string.
type Score struct {
	MainScore string `json:"mainScore"`
	Subscore  *int   `json:"subscore,omitempty"`
	IsWinner  *bool  `json:"isWinner,omitempty"`
}

// BracketryContestant represents a participant (a team or player).
type BracketryContestant struct {
	EntryStatus *string  `json:"entryStatus,omitempty"`
	Players     []Player `json:"players"`
}

// Player represents a single player within a contestant entry.
type Player struct {
	Title string `json:"title"`
}

// --- データベースから読み込むための構造体 ---

// Team represents the data fetched from the 'teams' table.
type Team struct {
	ID           int64
	Name         string
	TournamentID int64
	EntryStatus  *string // Assumed new column in 'teams' table
}

// Match represents the data fetched from the 'matches' table.
type Match struct {
	ID           int64
	TournamentID int64
	Round        int
	MatchNumber  int // Corresponds to match_number_in_round
	Team1ID      *int64
	Team2ID      *int64
	Team1Score   *int
	Team2Score   *int
	WinnerTeamID *int64
	NextMatchID  *int64
	Status       *string // Assumed new column in 'matches' table
}
