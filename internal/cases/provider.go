package cases

import (
	"context"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
)
//go:generate mockery
type CryptoProvider interface {
	GetActualRates(ctx context.Context, titles []string) ([]*entities.Coin, error)
}
