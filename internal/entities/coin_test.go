package entities_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"githab.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/stretchr/testify/require"
)

// Запуск и покрытие тестов: go test ./... -cover

// Визуализируем вывод тестов в рамках пакета entities_test
func TestMain(m *testing.M) {
	fmt.Println("___START test Package entities_test___")
	res := m.Run()
	fmt.Println("___THE END test Package entities_test___")

	os.Exit(res)
}

//

func TestNewCoin(t *testing.T) {

	testCases := []struct {
		name        string
		title       string
		price       float64
		create_coin bool
	}{
		{
			name:        "current data",
			title:       "bitcoin",
			price:       100.00,
			create_coin: true,
		},
		{
			name:        "title empty",
			title:       "",
			price:       100.00,
			create_coin: false,
		},
		{
			name:        "price zero",
			title:       "bitcoin",
			price:       0.00,
			create_coin: false,
		},
		{
			name:        "price negative",
			title:       "bitcoin",
			price:       -1.00,
			create_coin: false,
		},
		{
			name:        "title empty AND price zero or negative",
			title:       " ",
			price:       -1.00,
			create_coin: false,
		},
	}

	for _, elem := range testCases {
		t.Run(elem.name, func(t *testing.T) {
			if elem.create_coin == true {
				Coin, err := entities.NewCoin(elem.title, elem.price, time.Now())
				require.Equal(t, Coin == nil, false,
					fmt.Sprintf("coin not validation: err = %v", err))
				fmt.Println("Coin created")
			} else {
				Coin, err := entities.NewCoin(elem.title, elem.price, time.Now())
				require.Equal(t, Coin == nil, true,
					fmt.Sprintf("coin not validation: err = %v", err))
				fmt.Println("Coin not created")
			}
		})
	}
}
