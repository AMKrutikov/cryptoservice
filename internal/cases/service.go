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
	if err := s.processCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "failed to process coins")
	}

	actualCoins, err := s.storage.GetActualCoins(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get actual rates from storage")
	}
	return actualCoins, nil
}

func (s *Service) GetAgregetedRates(ctx context.Context, titles []string, aggType string) ([]*entities.Coin, error) {
	if err := s.processCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "failed to process coins")
	}

	aggTypeUpper := strings.ToUpper(aggType)
	if aggTypeUpper != "MIN" && aggTypeUpper != "MAX" && aggTypeUpper != "AVG" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "incorrect value for agregeted rates")
	}

	aggCoin, err := s.storage.GetAggregateCoins(ctx, titles, aggTypeUpper)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get aggregated coins from storage")
	}

	return aggCoin, nil
}

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

func (s *Service) processCoins(ctx context.Context, titles []string) error {
	if len(titles) == 0 {
		return errors.Wrap(entities.ErrInvalidParam, "empty titles")
	}

	coinsStorage, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get coins list from storage")
	}

	missingCoins := make([]string, 0, len(titles))

	for _, elem := range titles {
		if !slices.Contains(coinsStorage, elem) {
			missingCoins = append(missingCoins, elem)
		}
	}

	if len(missingCoins) == 0 {
		return nil
	}

	providerCoins, err := s.provider.GetActualRates(ctx, missingCoins)
	if err != nil {
		return errors.Wrap(err, "failed to get actual rates from provider")
	}

	if len(providerCoins) != 0 {
		err = s.storage.Store(ctx, providerCoins)
		if err != nil {
			return errors.Wrap(err, "failed to save coins to storage")
		}
	}

	return nil
}
