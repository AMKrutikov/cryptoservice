package postgres

import "time"

type CryptoModel struct {
	Id        int
	Title     string
	Price     float64
	Actual_at time.Time
}
