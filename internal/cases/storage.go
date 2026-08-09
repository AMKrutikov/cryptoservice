package cases

import (
	"context"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
)

//go:generate mockery --name CryptoStorage
type CryptoStorage interface {
	Store(ctx context.Context, coins []*entities.Coin) error                                          // положить - метод сохранения в хранилище
	GetCoinsList(ctx context.Context) ([]string, error)                                               // список имен монет, которые представлены
	GetActualCoins(ctx context.Context, titles []string) ([]*entities.Coin, error)                    // получение последних монет по title-ам
	GetAggregateCoins(ctx context.Context, titles []string, aggType string) ([]*entities.Coin, error) // Агрегированный запрос над монетами

	Close()
}
