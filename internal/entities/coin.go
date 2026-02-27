package entities

import (
	"errors"
	"strings"
	"time"
)

type Coin struct {
	Title    string
	Price    float64
	ActualAT time.Time
}

func NewCoin(title string, price float64, actualAT time.Time) (*Coin, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("Title cannot be empty")
	}
	if price <= 0 {
		return nil, errors.New("Price cannot be zero or negative")
	}
	return &Coin{
		Title:    title,
		Price:    price,
		ActualAT: actualAT,
	}, nil
}
