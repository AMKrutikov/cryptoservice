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

func (s *Service) GetLastRates(ctx context.Context, titles []string) ([]*entities.Coin, error) { // Получить последние цены
	listCoinMap := make(map[string]struct{})
	containsAll := true
	// запрашиваем список монет с базы (GetCoinsList) и сохраняем в переменную (storeTitles)
	storeTitles, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		return nil, entities.ErrStorage
	}
	// заполняем listCoinMap пока просто списком монет, которые есть в базе
	for _, elem := range storeTitles {
		listCoinMap[elem] = struct{}{}
	}
	// проверяем есть ли в базе(storage) запрашиваемые монеты(titles)
	for _, elem := range titles {
		if _, exists := listCoinMap[elem]; !exists {
			containsAll = false
		}
	}
	// если есть, выводим слайс монет
	if containsAll {
		resultCoins, err := s.storage.GetActualCoins(ctx, titles)
		if err != nil {
			return nil, entities.ErrStorage
		}
		return resultCoins, nil
	}

	// если нет запрашиваем у провайдера
	providerCoins, err := s.provider.GetActualRates(ctx, titles)
	if err != nil {
		return nil, entities.ErrProvider
	}
	// и сохраняем в базу
	err = s.storage.Store(ctx, providerCoins)
	if err != nil {
		return nil, entities.ErrStorage
	}
	// выводим слайс монет, по запросу через провайдера! без обращения к базе.
	return providerCoins, nil
}

func (s *Service) GetMaxRates(ctx context.Context, titles []string) ([]entities.Coin, error) { // Получить максимальные цены
	return []entities.Coin{}, nil
}

func (s *Service) GetMinRates(ctx context.Context, titles []string) ([]entities.Coin, error) { // Получить минимальные цены
	return []entities.Coin{}, nil
}

func (s *Service) GetAvgRates(ctx context.Context, titles []string) ([]entities.Coin, error) { // Получить средние цены
	return []entities.Coin{}, nil
}

// func (s Service) ActualizeRates(ctx context.Context, opts ...option) error {      // Получить актуальные цены
// 	return nil                                                                       // Демон подгружает данные
// }                                                                                 // с некой периодичностью
