package models

var (
	StatusUnassigned = "UNASSIGNED"
	StatusTaken      = "TAKEN"
)

type Order struct {
	Id           int64     `json:"id"`
	Origins      []float64 `json:"-"`
	Destinations []float64 `json:"-"`
	Distance     int       `json:"distance"`
	Status       string    `json:"status"`
}
