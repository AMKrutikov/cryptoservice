package cases

import (
	"context"
	"slices"
	"strings"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/pkg/errors"
)

type Service struct {
	provider CryptoProvider
	storage  CryptoStorage
}

var (
	min = "min"
	max = "max"
	avg = "avg"
)

func NewService(provider CryptoProvider, storage CryptoStorage) (*Service, error) {
	if provider == nil {
		return nil, errors.Wrap(entities.ErrInvalidParam, "provider cannot be nil")
	}
	if storage == nil {
		return nil, errors.Wrap(entities.ErrInvalidParam, "storage cannot be nil")
	}
	return &Service{
		provider: provider,
		storage:  storage,
	}, nil
}

func (s *Service) GetLastRates(ctx context.Context, titles []string) ([]*entities.Coin, error) { // Получить последние цены
	missingCoins := make([]string, 0, len(titles))

	coinsStorage, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get coins list from storage")
	}

	for _, elem := range titles {
		if !slices.Contains(coinsStorage, elem) {
			//providerCoins, err := s.provider.GetActualRates(ctx, titles) - обновить список всех монет
			missingCoins = append(missingCoins, elem)
		}
	}
	if len(missingCoins) > 0 {
		providerCoins, err := s.provider.GetActualRates(ctx, missingCoins)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get actual rates from provider")
		}

		err = s.storage.Store(ctx, providerCoins)
		if err != nil { //providerCoins, err := s.provider.GetActualRates(ctx, titles)
			return nil, errors.Wrap(err, "failed to save coins to storage")
		}
	}

	actualCoins, err := s.storage.GetActualCoins(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get actual rates from storage")
	}
	return actualCoins, nil
}

func (s *Service) GetAgregetedRates(ctx context.Context, titles []string, aggType string) ([]*entities.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "incorrect value for list of coins")
	}

	aggTypeLower := strings.ToLower(aggType)
	if aggTypeLower != min && aggTypeLower != max && aggTypeLower != avg {
		return nil, errors.Wrap(entities.ErrInvalidParam, "incorrect value for agregeted rates")
	}

	if _, err := s.GetLastRates(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "failed to ensure coins presence")
	}

	aggCoin, err := s.storage.GetAggregateCoins(ctx, titles, aggTypeLower)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get aggregated coins from storage")
	}

	return aggCoin, nil
}

// Получить актуальные цены Демон подгружает данные с некой периодичностью
func (s *Service) ActualizeRates(ctx context.Context) error {
	coinsList, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get coins list from storage")
	}

	if len(coinsList) == 0 {
		return nil
	}

	updatedCoins, err := s.provider.GetActualRates(ctx, coinsList)
	if err != nil {
		return errors.Wrap(err, "failed to get actual rates from provider")
	}

	if err := s.storage.Store(ctx, updatedCoins); err != nil {
		return errors.Wrap(err, "failed to save updated coins to storage")
	}

	return nil
}
