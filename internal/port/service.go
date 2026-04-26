package port

import (
	"context"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
)

type Service interface {
	GetLastRates(ctx context.Context, titles []string) ([]*entities.Coin, error)
	GetAggregateRates(ctx context.Context, titles []string, aggType string) ([]*entities.Coin, error)
}
