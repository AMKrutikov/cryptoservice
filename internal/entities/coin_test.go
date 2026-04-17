package entities_test

import (
	"testing"
	"time"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/stretchr/testify/require"
)

func TestNewCoin(t *testing.T) {

	testCases := []struct {
		name    string
		title   string
		price   float64
		wantErr bool
		resErr  error
	}{
		{
			name:  "correct data",
			title: "bitcoin",
			price: 100.00,
		},
		{
			name:    "title empty",
			title:   "",
			price:   100.00,
			wantErr: true,
			resErr:  entities.ErrInvalidParam,
		},
		{
			name:    "price zero",
			title:   "bitcoin",
			price:   0.00,
			wantErr: true,
			resErr:  entities.ErrInvalidParam,
		},
		{
			name:    "price negative",
			title:   "bitcoin",
			price:   -1.00,
			wantErr: true,
			resErr:  entities.ErrInvalidParam,
		},
		{
			name:    "title empty AND price zero or negative",
			title:   " ",
			price:   -1.00,
			wantErr: true,
			resErr:  entities.ErrInvalidParam,
		},
	}
	for _, elem := range testCases {
		t.Run(elem.name, func(t *testing.T) {
			t.Parallel()

			now := time.Now()
			coin, err := entities.NewCoin(elem.title, elem.price, now)
			if elem.wantErr {
				require.Nil(t, coin)
				require.ErrorIs(t, err, entities.ErrInvalidParam)
				return
			}
			require.NotNil(t, coin)
			require.NoError(t, err)

			require.Equal(t, elem.title, coin.Title())
			require.Equal(t, elem.price, coin.Price())
			require.Equal(t, now, coin.ActualAt())
		})
	}
}
