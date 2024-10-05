package dto

// Characterの構造体
type Character struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	Probability float64 `db:"probability"`
}
