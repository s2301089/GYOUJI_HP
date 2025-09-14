package model

type RelayResult struct {
	ID        int
	RelayType string // 'A' or 'B'
	RelayRank int    // 1～6
	ClassID   int
	CreatedAt string
}

type RelayRegisterRequest struct {
	RelayType string `json:"relay_type"` // 'A' or 'B'
	ClassIDs  []int  `json:"class_ids"`  // 順位順（1位～6位）
}
