package cases

import (
	"context"

	"githab.com/AMKrutikov/cryptoservice/internal/entities"
)

type CryptoProvider interface {
	GetActualRates(ctx context.Context, titles []string) ([]entities.Coin, error)
}

type myCoin entities.Coin

func (c *myCoin) GetActualRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	result := []entities.Coin{}
	for _, elem := range titles {
		result = append(result, entities.Coin{Title: elem})
	}
	return result, nil
}
