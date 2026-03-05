package cases

import "context"

type Coin struct {
	Ttl string
} // временная заглушка

type Service struct {
}

func (s Service) GetLastRates(ctx context.Context, titles []string) ([]Coin, error) { // Получить последние цены
	return []Coin{}, nil
}

func (s Service) GetMaxRates(ctx context.Context, titles []string) ([]Coin, error) { // Получить максимальные цены
	return []Coin{}, nil
}

func (s Service) GetMinRates(ctx context.Context, titles []string) ([]Coin, error) { // Получить минимальные цены
	return []Coin{}, nil
}

func (s Service) GetAvgRates(ctx context.Context, titles []string) ([]Coin, error) { // Получить средние цены
	return []Coin{}, nil
}

// func (s Service) ActualizeRates(ctx context.Context, opts ...option) error {      // Получить актуальные цены
// 	return nil                                                                       // Демон подгружает данные
// }                                                                                 // с некой периодичностью
