package entities

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Coin struct {
	title    string
	price    float64
	actualAT time.Time
}

func NewCoin(title string, price float64, actualAT time.Time) (*Coin, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.Wrap(ErrInvalidParam, "Title cannot be empty")
	}
	if price <= 0 {
		return nil, errors.Wrap(ErrInvalidParam, "Price cannot be zero or negative")
	}
	return &Coin{
		title:    title,
		price:    price,
		actualAT: actualAT,
	}, nil
}

func (c *Coin) Title() string {
	return c.title
}

func (c *Coin) Price() float64 {
	return c.price
}

func (c *Coin) ActualAt() time.Time {
	return c.actualAT
}
