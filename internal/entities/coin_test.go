package entities_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"githab.com/AMKrutikov/cryptoservice/internal/entities"
)

// Запуск и покрытие тестов: go test ./... -cover

// Визуализируем вывод тестов в рамках пакета Coin
func TestMain(m *testing.M) {
	fmt.Println("TEST Package Coin")
	res := m.Run()
	fmt.Println("The end test Package Coin")

	os.Exit(res)
}

//

func TestNewCoin(t *testing.T) {
	t.Run("current data", func(t *testing.T) {
		title := "Bitcoin"
		price := 100.00
		actualAT := time.Now()

		Coin, err := entities.NewCoin(title, price, actualAT)
		if Coin == nil {
			t.Errorf("coin not validation: %v", err)
		}

	})
	t.Run("title empty", func(t *testing.T) {
		title := ""
		price := 100.00
		actualAT := time.Now()

		Coin, err := entities.NewCoin(title, price, actualAT)
		if Coin != nil {
			t.Errorf("coin not validation: %v", err)
		}

	})

	t.Run("price zero or negative", func(t *testing.T) {
		title := "bitcoin"
		price := -1.00
		actualAT := time.Now()

		Coin, err := entities.NewCoin(title, price, actualAT)
		if Coin != nil {
			t.Errorf("coin not validation: %v", err)
		}

	})

	t.Run("title empty AND price zero or negative", func(t *testing.T) {
		title := " "
		price := 0.00
		actualAT := time.Now()

		Coin, err := entities.NewCoin(title, price, actualAT)
		if Coin != nil {
			t.Errorf("coin not validation: %v", err)
		}
	})

}
