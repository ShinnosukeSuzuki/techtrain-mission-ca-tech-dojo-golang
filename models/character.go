package models

// Characterの構造体
type Character struct {
	ID          string
	Name        string
	Probability float64
}

type GachaResult struct {
	CharacterID string
	Name        string
}
