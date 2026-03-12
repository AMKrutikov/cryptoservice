package cases

import (
	"context"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
)

//go:generate mockery
type CryptoProvider interface {
	GetActualRates(ctx context.Context, titles []string) ([]*entities.Coin, error)
}

func AAA(ctx context.Context, titles []string, c CryptoProvider) ([]*entities.Coin, error) {
	return c.GetActualRates(ctx, titles)
}
