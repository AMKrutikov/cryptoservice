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
	t.Run("TestSuccessGetCoins", func(t *testing.T) {
		mockCryptoProvider := mocks.NewCryptoProvider(t)

		tNow := time.Now()
		ctx := context.Background()
		titles := []string{"BTC", "ETH", "XPR"}

		coin1, err := entities.NewCoin("BTC", 1000.00, tNow)
		require.NoError(t, err)
		coin2, err := entities.NewCoin("ETH", 500.00, tNow)
		require.NoError(t, err)
		coin3, err := entities.NewCoin("XPR", 100.00, tNow)
		require.NoError(t, err)

		expectedCoins := []*entities.Coin{coin1, coin2, coin3}

		mockCryptoProvider.On("GetActualRates", ctx, titles).
			Return(expectedCoins, nil)
		// mockCryptoProvider.EXPECT().GetActualRates(ctx, expectedTitles).
		// 	Return(expectedCoin, nil)

		actualCoins, err := cases.GetCoins(ctx, titles, mockCryptoProvider)
		require.NoError(t, err)
		require.Equal(t, expectedCoins, actualCoins)
	})
}
