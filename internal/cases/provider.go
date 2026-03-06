package cases

import (
	"context"

	"githab.com/AMKrutikov/cryptoservice/internal/entities"
)

type CryptoProvider interface {
	GetActualRates(ctx context.Context, titles []string) ([]*entities.Coin, error)
}

// gomock либо mokkery
