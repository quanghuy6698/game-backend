package tetGame

type Player struct {
	ID        string    `json:"id"`
	PublicID  string    `json:"public_id"`
	Name      string    `json:"name"`
}

type Score struct {
	PlayerID string `json:"player_id"`
	Score    int64  `json:"score"`
	Skin     string `json:"skin"`
}

type LeaderboardItem struct {
	PublicID string `json:"public_id"`
	Name     string `json:"name"`
	Score    int    `json:"score"`
	Skin     string `json:"skin"`
}
