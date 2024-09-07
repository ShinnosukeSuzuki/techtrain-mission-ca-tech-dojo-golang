package models

// Characterの構造体
type Character struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Probability float64 `json:"probability"`
}

type GachaResult struct {
	CharacterID string
	Name        string
}
