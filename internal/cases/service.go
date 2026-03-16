package cases

import (
	"context"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
)

type Service struct {
	provider CryptoProvider
	storage  CryptoStorage
}

func NewService(provider CryptoProvider, storage CryptoStorage) (*Service, error) {
	if provider == nil || storage == nil {
		return nil, entities.ErrService
	}
	return &Service{
		provider: provider,
		storage:  storage,
	}, nil
}
func (s Service) GetLastRates(ctx context.Context, titles []string) ([]*entities.Coin, error) { // Получить последние цены
	return []*entities.Coin{}, nil
	// 1 - получить последние цены
	//   есть ли в базе(storage) такие монеты (GetCoinsList),
	// - проверяем что все они есть, выводим слайс монет!
	// - если нет запрашиваем в провайдере и сохраняем в базу!
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
