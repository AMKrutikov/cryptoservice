package cases_test

// gomock или  Mockery

import (
	"context"
	"testing"
	"time"

	"github.com/AMKrutikov/cryptoservice/internal/cases"
	"github.com/AMKrutikov/cryptoservice/internal/cases/mocks"
	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/stretchr/testify/require"
)

func TestGetCoinsMock(t *testing.T) {
	// CryptoProvider interface
	expectedErr := entities.ErrProvider
	tNow := time.Now()
	ctx := context.Background()
	titles := []string{"BTC", "ETH"}

	coin1, err := entities.NewCoin("BTC", 1000.00, tNow)
	require.NoError(t, err)
	coin2, err := entities.NewCoin("ETH", 500.00, tNow)
	require.NoError(t, err)
	expectedCoins := []*entities.Coin{coin1, coin2}
	t.Run("TestSuccessGetCoins", func(t *testing.T) {
		mockCryptoProvider := mocks.NewCryptoProvider(t)

		mockCryptoProvider.On("GetActualRates", ctx, titles).
			Return(expectedCoins, nil)
		// mockCryptoProvider.EXPECT().GetActualRates(ctx, expectedTitles).
		// 	Return(expectedCoin, nil)

		actualCoins, err := cases.GetCoins(ctx, titles, mockCryptoProvider)
		require.NoError(t, err)
		require.Equal(t, expectedCoins, actualCoins)
	})
	t.Run("TestFailGetCoins", func(t *testing.T) {
		mockCryptoProvider := mocks.NewCryptoProvider(t)

		mockCryptoProvider.On("GetActualRates", ctx, titles).
			Return(nil, expectedErr)

		_, err := cases.GetCoins(ctx, titles, mockCryptoProvider)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
}

func TestSavingStorageMock(t *testing.T) {
	// CryptoStorage interface
	ctx := context.Background()
	tNow := time.Now()
	expectedErr := entities.ErrStorage

	coin1, err := entities.NewCoin("BTC", 1000.00, tNow)
	require.NoError(t, err)
	coin2, err := entities.NewCoin("ETH", 500.00, tNow)
	require.NoError(t, err)
	coins := []*entities.Coin{coin1, coin2}
	t.Run("TestSuccessSavingToStorage", func(t *testing.T) {
		mockCryptoStorage := mocks.NewCryptoStorage(t)

		mockCryptoStorage.On("Store", ctx, coins).Return(nil)
		require.NoError(t, cases.SavingToStorage(ctx, coins, mockCryptoStorage))

	})

	t.Run("TestFailSavingToStorage", func(t *testing.T) {
		mockCryptoStorage := mocks.NewCryptoStorage(t)

		mockCryptoStorage.On("Store", ctx, coins).Return(expectedErr)

		err = cases.SavingToStorage(ctx, coins, mockCryptoStorage)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
}

func TestGetCoinsNamesMock(t *testing.T) {
	// CryptoStorage interface
	ctx := context.Background()
	titles := []string{"BTC", "ETH"}
	expectedErr := entities.ErrStorage
	t.Run("TestSuccessCoinsNames", func(t *testing.T) {
		mockGetCoinsNames := mocks.NewCryptoStorage(t)

		mockGetCoinsNames.On("GetCoinsList", ctx).Return(titles, nil)

		actualTitles, err := cases.GetCoinsNames(ctx, mockGetCoinsNames)
		require.NoError(t, err)
		require.Equal(t, titles, actualTitles)

	})
	t.Run("TestFailCoinsNames", func(t *testing.T) {
		mockGetCoinsNames := mocks.NewCryptoStorage(t)

		mockGetCoinsNames.On("GetCoinsList", ctx).Return(nil, expectedErr)

		_, err := cases.GetCoinsNames(ctx, mockGetCoinsNames)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
}

func TestGetLastCoinsMock(t *testing.T) {
	// CryptoStorage interface
	ctx := context.Background()
	tNow := time.Now()
	expectedErr := entities.ErrStorage

	coin1, err := entities.NewCoin("BTC", 1000.00, tNow)
	require.NoError(t, err)
	coin2, err := entities.NewCoin("ETH", 500.00, tNow)
	require.NoError(t, err)
	coins := []*entities.Coin{coin1, coin2}

	titles := []string{"BTC", "ETH"}
	t.Run("TestSuccessGetLastCoins", func(t *testing.T) {
		mockGetLastCoins := mocks.NewCryptoStorage(t)

		mockGetLastCoins.On("GetActualCoins", ctx, titles).Return(coins, nil)

		actualCoins, err := cases.GetLastCoins(ctx, titles, mockGetLastCoins)
		require.NoError(t, err)
		require.Equal(t, coins, actualCoins)

	})
	t.Run("TestFailGetLastCoins", func(t *testing.T) {
		mockGetLastCoins := mocks.NewCryptoStorage(t)

		mockGetLastCoins.On("GetActualCoins", ctx, titles).Return(nil, expectedErr)

		_, err := cases.GetLastCoins(ctx, titles, mockGetLastCoins)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
}

func TestGetAggregatedRequestMock(t *testing.T) {
	// CryptoStorage interface
	ctx := context.Background()
	tNow := time.Now()
	expectedErr := entities.ErrStorage

	coin1, err := entities.NewCoin("BTC", 1000.00, tNow)
	require.NoError(t, err)
	coin2, err := entities.NewCoin("ETH", 500.00, tNow)
	require.NoError(t, err)
	coins := []*entities.Coin{coin1, coin2}
	aggType := ""

	titles := []string{"BTC", "ETH"}

	t.Run("TestSuccessGetAggRequest", func(t *testing.T) {
		mockGetAggRequest := mocks.NewCryptoStorage(t)

		mockGetAggRequest.On("GetAggregateCoins", ctx, titles, aggType).Return(coins, nil)

		actualCoins, err := cases.GetAggregatedRequest(ctx, titles, aggType, mockGetAggRequest)
		require.NoError(t, err)
		require.Equal(t, coins, actualCoins)

	})
	t.Run("TestFailGetAggRequest", func(t *testing.T) {
		mockGetAggRequest := mocks.NewCryptoStorage(t)

		mockGetAggRequest.On("GetAggregateCoins", ctx, titles, aggType).Return(nil, expectedErr)

		_, err := cases.GetAggregatedRequest(ctx, titles, aggType, mockGetAggRequest)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)

	})
}
