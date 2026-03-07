package entities_test

import (
	"testing"
	"time"

	"githab.com/AMKrutikov/cryptoservice/internal/entities"
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
			resErr:  entities.ERRInvalidParam,
		},
		{
			name:    "price zero",
			title:   "bitcoin",
			price:   0.00,
			wantErr: true,
			resErr:  entities.ERRInvalidParam,
		},
		{
			name:    "price negative",
			title:   "bitcoin",
			price:   -1.00,
			wantErr: true,
			resErr:  entities.ERRInvalidParam,
		},
		{
			name:    "title empty AND price zero or negative",
			title:   " ",
			price:   -1.00,
			wantErr: true,
			resErr:  entities.ERRInvalidParam,
		},
	}
	for _, elem := range testCases {
		t.Run(elem.name, func(t *testing.T) {
			coin, err := entities.NewCoin(elem.title, elem.price, time.Now())
			if elem.wantErr {
				require.Nil(t, coin)
				require.ErrorIs(t, err, entities.ERRInvalidParam)
				return
			}
			require.NotNil(t, coin)
			require.NoError(t, err)
		})
	}
}
