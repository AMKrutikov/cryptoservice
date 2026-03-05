package cases

import (
	"context"

	"githab.com/AMKrutikov/cryptoservice/internal/entities"
)

type Service struct {
}

func (s Service) GetLastRates(ctx context.Context, titles []string) ([]entities.Coin, error) { // Получить последние цены
	return []entities.Coin{}, nil
}

func (s Service) GetMaxRates(ctx context.Context, titles []string) ([]entities.Coin, error) { // Получить максимальные цены
	return []entities.Coin{}, nil
}

func (s Service) GetMinRates(ctx context.Context, titles []string) ([]entities.Coin, error) { // Получить минимальные цены
	return []entities.Coin{}, nil
}

func (s Service) GetAvgRates(ctx context.Context, titles []string) ([]entities.Coin, error) { // Получить средние цены
	return []entities.Coin{}, nil
}

// func (s Service) ActualizeRates(ctx context.Context, opts ...option) error {      // Получить актуальные цены
// 	return nil                                                                       // Демон подгружает данные
// }                                                                                 // с некой периодичностью
