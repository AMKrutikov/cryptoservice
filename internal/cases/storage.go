package cases

import (
	"context"

	"githab.com/AMKrutikov/cryptoservice/internal/entities"
)

type Storage interface {
	Store(ctx context.Context, coins []*entities.Coin) error                                          // положить - метод
	GetCoinsList(ctx context.Context) ([]string, error)                                               // Список title монет, которые предст только имена
	GetActualCoins(ctx context.Context, titles []string) ([]*entities.Coin, error)                    // получение последних монет по title-ам
	GetAggregateCoins(ctx context.Context, titles []string, aggType string) ([]*entities.Coin, error) // Агрегированный запрос над монетами
}
