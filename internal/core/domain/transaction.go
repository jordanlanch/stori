package domain

type Transaction struct {
	ID     int     `json:"id" gorm:"primaryKey"`
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}
