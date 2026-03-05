package cases

import (
	"context"

	"githab.com/AMKrutikov/cryptoservice/internal/entities"
)

type Storage interface {
	Store(ctx context.Context, coins []entities.Coin) error                          // Xранилище
	GetCoinsList(ctx context.Context) ([]string, error)                              // Список монет
	GetActualCoins(ctx context.Context, titles []string) ([]entities.Coin, error)    // ???
	GetAggregateCoins(ctx context.Context, titles []string) ([]entities.Coin, error) // Агрегированный запрос над монетами
}
