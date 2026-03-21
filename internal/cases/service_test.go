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

func TestNewService(t *testing.T) {
	t.Run("CryptoProviderNilFail", func(t *testing.T) {
		mockStorage := mocks.NewCryptoStorage(t)
		expectedError := entities.ErrInvalidParam

		_, err := cases.NewService(nil, mockStorage)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedError)
	})
	t.Run("CryptoStorageNilFail", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		expectedError := entities.ErrInvalidParam

		_, err := cases.NewService(mockProvider, nil)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedError)
	})
	t.Run("NewServiceSuccess", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

	})
}

func TestGetLastRates(t *testing.T) {
	expectedErr := entities.ErrInvalidParam
	ctx := context.Background()
	titles := []string{"BTC", "ETH"}
	tNow := time.Now()
	coin1, err := entities.NewCoin("BTC", 1000.00, tNow)
	require.NoError(t, err)
	coin2, err := entities.NewCoin("ETH", 500.00, tNow)
	require.NoError(t, err)
	expectedCoins := []*entities.Coin{coin1, coin2}
	t.Run("FailLenTitles", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		_, err = service.GetLastRates(ctx, []string{})
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)

	})
	t.Run("AllCoinsStorageSuccess", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return(titles, nil)
		mockStorage.EXPECT().GetActualCoins(ctx, titles).Return(expectedCoins, nil)

		actualCoins, err := service.GetLastRates(ctx, titles)
		require.NoError(t, err)
		require.Equal(t, expectedCoins, actualCoins)
	})

	t.Run("AllCoinsProviderSuccess", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return([]string{}, nil)
		mockProvider.EXPECT().GetActualRates(ctx, titles).Return(expectedCoins, nil)
		mockStorage.EXPECT().Store(ctx, expectedCoins).Return(nil)
		mockStorage.EXPECT().GetActualCoins(ctx, titles).Return(expectedCoins, nil)

		actualCoins, err := service.GetLastRates(ctx, titles)
		require.NoError(t, err)
		require.Equal(t, expectedCoins, actualCoins)
	})
	t.Run("FailGetCoinsList", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return([]string{}, expectedErr)

		_, err = service.GetLastRates(ctx, titles)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("FailGetActualCoins", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return(titles, nil)
		mockStorage.EXPECT().GetActualCoins(ctx, titles).Return(nil, expectedErr)

		_, err = service.GetLastRates(ctx, titles)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("FailStorage", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return([]string{}, nil)
		mockProvider.EXPECT().GetActualRates(ctx, titles).Return(expectedCoins, nil)
		mockStorage.EXPECT().Store(ctx, expectedCoins).Return(expectedErr)

		_, err = service.GetLastRates(ctx, titles)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("FailProvider", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return([]string{}, nil)
		mockProvider.EXPECT().GetActualRates(ctx, titles).Return(nil, expectedErr)

		_, err = service.GetLastRates(ctx, titles)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
}

func TestGetAgregetedRates(t *testing.T) {
	expectedErr := entities.ErrInvalidParam
	ctx := context.Background()
	titles := []string{"BTC", "ETH"}
	tNow := time.Now()
	coin1, err := entities.NewCoin("BTC", 1000.00, tNow)
	require.NoError(t, err)
	coin2, err := entities.NewCoin("ETH", 500.00, tNow)
	require.NoError(t, err)
	expectedCoins := []*entities.Coin{coin1, coin2}
	avg := "avg"
	invalidAgg := "invalidAgg"
	t.Run("FailLenTitles", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		_, err = service.GetAgregetedRates(ctx, []string{}, avg)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)

	})
	t.Run("FailInvalidAggTypeLower", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		_, err = service.GetAgregetedRates(ctx, titles, invalidAgg)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)

	})
	t.Run("FailGetLastRates", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return([]string{}, expectedErr)

		_, err = service.GetAgregetedRates(ctx, titles, avg)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)

	})
	t.Run("FailGetAggregateCoins", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return(titles, nil)
		mockStorage.EXPECT().GetActualCoins(ctx, titles).Return(expectedCoins, nil)
		mockStorage.EXPECT().GetAggregateCoins(ctx, titles, avg).Return(nil, expectedErr)
		_, err = service.GetAgregetedRates(ctx, titles, avg)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("FailGetAggregateCoins", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return(titles, nil)
		mockStorage.EXPECT().GetActualCoins(ctx, titles).Return(expectedCoins, nil)
		mockStorage.EXPECT().GetAggregateCoins(ctx, titles, avg).Return(nil, expectedErr)

		_, err = service.GetAgregetedRates(ctx, titles, avg)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("SuccessGetAggregateCoins", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return(titles, nil)
		mockStorage.EXPECT().GetActualCoins(ctx, titles).Return(expectedCoins, nil)
		mockStorage.EXPECT().GetAggregateCoins(ctx, titles, avg).Return(expectedCoins, nil)

		actualCoins, err := service.GetAgregetedRates(ctx, titles, avg)
		require.NoError(t, err)
		require.Equal(t, expectedCoins, actualCoins)
	})
}

func TestActualizeRates(t *testing.T) {
	expectedErr := entities.ErrInvalidParam
	ctx := context.Background()
	titles := []string{"BTC", "ETH"}
	tNow := time.Now()
	coin1, err := entities.NewCoin("BTC", 1000.00, tNow)
	require.NoError(t, err)
	coin2, err := entities.NewCoin("ETH", 500.00, tNow)
	require.NoError(t, err)
	expectedCoins := []*entities.Coin{coin1, coin2}
	t.Run("FailGetCoinsList", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return([]string{}, expectedErr)

		err = service.ActualizeRates(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("LenCoinsListEmpty", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return([]string{}, nil)

		err = service.ActualizeRates(ctx)
		require.NoError(t, err)
		require.Nil(t, err)
	})
	t.Run("FailGetActualRates", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return(titles, nil)
		mockProvider.EXPECT().GetActualRates(ctx, titles).Return(nil, expectedErr)

		err = service.ActualizeRates(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("FailStore", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return(titles, nil)
		mockProvider.EXPECT().GetActualRates(ctx, titles).Return(expectedCoins, nil)
		mockStorage.EXPECT().Store(ctx, expectedCoins).Return(expectedErr)

		err = service.ActualizeRates(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
	t.Run("SuccessActualizeRates", func(t *testing.T) {
		mockProvider := mocks.NewCryptoProvider(t)
		mockStorage := mocks.NewCryptoStorage(t)

		service, err := cases.NewService(mockProvider, mockStorage)
		require.NoError(t, err)
		require.NotNil(t, service)

		mockStorage.EXPECT().GetCoinsList(ctx).Return(titles, nil)
		mockProvider.EXPECT().GetActualRates(ctx, titles).Return(expectedCoins, nil)
		mockStorage.EXPECT().Store(ctx, expectedCoins).Return(nil)

		err = service.ActualizeRates(ctx)
		require.NoError(t, err)
		require.Nil(t, err)
	})
}

func TestCryptoStorage(t *testing.T) {
	ctx := context.Background()
	tNow := time.Now()
	expectedErr := entities.ErrInvalidParam
	coin1, err := entities.NewCoin("BTC", 1000.00, tNow)
	require.NoError(t, err)
	coin2, err := entities.NewCoin("ETH", 500.00, tNow)
	require.NoError(t, err)
	coins := []*entities.Coin{coin1, coin2}
	t.Run("TestSuccessSavingToStorage", func(t *testing.T) {
		mockCryptoStorage := mocks.NewCryptoStorage(t)

		mockCryptoStorage.EXPECT().Store(ctx, coins).Return(nil)
		require.NoError(t, mockCryptoStorage.Store(ctx, coins))

	})

	t.Run("TestFailSavingToStorage", func(t *testing.T) {
		mockCryptoStorage := mocks.NewCryptoStorage(t)

		mockCryptoStorage.EXPECT().Store(ctx, coins).Return(expectedErr)

		err = mockCryptoStorage.Store(ctx, coins)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
}

// func TestGetCoinsNamesMock(t *testing.T) {
// 	// CryptoStorage interface
// 	ctx := context.Background()
// 	titles := []string{"BTC", "ETH"}
// 	expectedErr := entities.ErrInvalidParam
// 	t.Run("TestSuccessCoinsNames", func(t *testing.T) {
// 		mockGetCoinsNames := mocks.NewCryptoStorage(t)

// 		mockGetCoinsNames.On("GetCoinsList", ctx).Return(titles, nil)

// 		actualTitles, err := cases.GetCoinsNames(ctx, mockGetCoinsNames)
// 		require.NoError(t, err)
// 		require.Equal(t, titles, actualTitles)

// 	})
// 	t.Run("TestFailCoinsNames", func(t *testing.T) {
// 		mockGetCoinsNames := mocks.NewCryptoStorage(t)

// 		mockGetCoinsNames.On("GetCoinsList", ctx).Return(nil, expectedErr)

// 		_, err := cases.GetCoinsNames(ctx, mockGetCoinsNames)
// 		require.Error(t, err)
// 		require.ErrorIs(t, err, expectedErr)
// 	})
// }

// func TestGetLastCoinsMock(t *testing.T) {
// 	// CryptoStorage interface
// 	ctx := context.Background()
// 	tNow := time.Now()
// 	expectedErr := entities.ErrStorage

// 	coin1, err := entities.NewCoin("BTC", 1000.00, tNow)
// 	require.NoError(t, err)
// 	coin2, err := entities.NewCoin("ETH", 500.00, tNow)
// 	require.NoError(t, err)
// 	coins := []*entities.Coin{coin1, coin2}

// 	titles := []string{"BTC", "ETH"}
// 	t.Run("TestSuccessGetLastCoins", func(t *testing.T) {
// 		mockGetLastCoins := mocks.NewCryptoStorage(t)

// 		mockGetLastCoins.On("GetActualCoins", ctx, titles).Return(coins, nil)

// 		actualCoins, err := cases.GetLastCoins(ctx, titles, mockGetLastCoins)
// 		require.NoError(t, err)
// 		require.Equal(t, coins, actualCoins)

// 	})
// 	t.Run("TestFailGetLastCoins", func(t *testing.T) {
// 		mockGetLastCoins := mocks.NewCryptoStorage(t)

// 		mockGetLastCoins.On("GetActualCoins", ctx, titles).Return(nil, expectedErr)

// 		_, err := cases.GetLastCoins(ctx, titles, mockGetLastCoins)
// 		require.Error(t, err)
// 		require.ErrorIs(t, err, expectedErr)
// 	})
// }

// func TestGetAggregatedRequestMock(t *testing.T) {
// 	// CryptoStorage interface
// 	ctx := context.Background()
// 	tNow := time.Now()
// 	expectedErr := entities.ErrStorage

// 	coin1, err := entities.NewCoin("BTC", 1000.00, tNow)
// 	require.NoError(t, err)
// 	coin2, err := entities.NewCoin("ETH", 500.00, tNow)
// 	require.NoError(t, err)
// 	coins := []*entities.Coin{coin1, coin2}
// 	aggType := ""

// 	titles := []string{"BTC", "ETH"}

// 	t.Run("TestSuccessGetAggRequest", func(t *testing.T) {
// 		mockGetAggRequest := mocks.NewCryptoStorage(t)

// 		mockGetAggRequest.On("GetAggregateCoins", ctx, titles, aggType).Return(coins, nil)

// 		actualCoins, err := cases.GetAggregatedRequest(ctx, titles, aggType, mockGetAggRequest)
// 		require.NoError(t, err)
// 		require.Equal(t, coins, actualCoins)

// 	})
// 	t.Run("TestFailGetAggRequest", func(t *testing.T) {
// 		mockGetAggRequest := mocks.NewCryptoStorage(t)

// 		mockGetAggRequest.On("GetAggregateCoins", ctx, titles, aggType).Return(nil, expectedErr)

// 		_, err := cases.GetAggregatedRequest(ctx, titles, aggType, mockGetAggRequest)
// 		require.Error(t, err)
// 		require.ErrorIs(t, err, expectedErr)

// 	})
// }
