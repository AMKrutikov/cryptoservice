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
}

func SavingToStorage(ctx context.Context, coins []*entities.Coin, c CryptoStorage) error {
	return c.Store(ctx, coins)
}

func GetCoinsNames(ctx context.Context, c CryptoStorage) ([]string, error) {
	return c.GetCoinsList(ctx)
}

func GetLastCoins(ctx context.Context, titles []string, c CryptoStorage) ([]*entities.Coin, error) {
	return c.GetActualCoins(ctx, titles)
}

func GetAggregatedRequest(ctx context.Context, titles []string, aggType string, c CryptoStorage) ([]*entities.Coin, error) {
	return c.GetAggregateCoins(ctx, titles, aggType)
}
