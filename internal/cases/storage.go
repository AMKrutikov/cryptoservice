package cases

import "context"

type Storage interface {
	Store(ctx context.Context, coins []Coin) error                          // Xранилище
	GetCoinsList(ctx context.Context) ([]string, error)                     // Список монет
	GetActualCoins(ctx context.Context, titles []string) ([]Coin, error)    // ???
	GetAggregateCoins(ctx context.Context, titles []string) ([]Coin, error) // Агрегированный запрос над монетами
}
