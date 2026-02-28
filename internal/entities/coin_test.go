package entities_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"githab.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/stretchr/testify/assert"
)

// Запуск и покрытие тестов: go test ./... -cover

// Визуализируем вывод тестов в рамках пакета entities_test
func TestMain(m *testing.M) {
	fmt.Println("TEST Package entities_test")
	res := m.Run()
	fmt.Println("The end test Package entities_test")

	os.Exit(res)
}

//

func TestNewCoin(t *testing.T) {
	t.Run("title empty", func(t *testing.T) {
		title := ""
		price := 100.00
		actualAT := time.Now()

		Coin, err := entities.NewCoin(title, price, actualAT)
		assert.Equal(t, Coin != nil, Coin != nil,
			fmt.Sprintf("coin not validation: err = %v", err))
	})

	t.Run("price zero or negative", func(t *testing.T) {
		title := "bitcoin"
		price := 0.00
		actualAT := time.Now()

		Coin, err := entities.NewCoin(title, price, actualAT)
		assert.Equal(t, Coin != nil, Coin != nil,
			fmt.Sprintf("coin not validation: err = %v", err))

		price = -1.00
		Coin, err = entities.NewCoin(title, price, actualAT)
		assert.Equal(t, Coin != nil, Coin != nil,
			fmt.Sprintf("coin not validation: err = %v", err))
	})

	t.Run("title empty AND price zero or negative", func(t *testing.T) {
		title := " "
		price := 0.00
		actualAT := time.Now()

		Coin, err := entities.NewCoin(title, price, actualAT)
		assert.Equal(t, Coin != nil, Coin != nil,
			fmt.Errorf("coin not validation: err = %v", err))
	})
	t.Run("current data", func(t *testing.T) {
		title := "Bitcoin"
		price := 100.00
		actualAT := time.Now()

		Coin, err := entities.NewCoin(title, price, actualAT)
		assert.Equal(t, Coin == nil, Coin == nil,
			fmt.Sprintf("coin not validation: err = %v", err))
	})
}
