package entities

import "time"

type Coin struct {
	Title    string
	Price    float32
	ActualAT time.Time
}

// далее пишем контсруктор
