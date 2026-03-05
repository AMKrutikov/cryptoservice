package cases_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"githab.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/stretchr/testify/mock"
)

// Запуск и покрытие тестов: go test ./... -cover

// Визуализируем вывод тестов в рамках пакета entities_test
func TestMain(m *testing.M) {
	fmt.Println("___START test Package cases_test___")
	res := m.Run()
	fmt.Println("___THE END test Package cases_test___")

	os.Exit(res)
}

// CryptoProvider Mock
type MockCryptoProvider struct {
	mock.Mock
}

func (c *MockCryptoProvider) GetActualRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	argumets := c.Called(ctx, titles)
	return argumets.Get(0).([]entities.Coin), argumets.Error(1)
}

// Storage Mock

func TestGetActualRates(t *testing.T) {
	fmt.Println("_START TEST TestGetActualRates")

	mockCryptoProvider := new(MockCryptoProvider)
	t.Cleanup(func() {
		mockCryptoProvider.AssertExpectations(t)
		fmt.Println("THE END TEST TestRunService")
	})

	dataCryptoProvider := []entities.Coin{
		{Title: "Bitcoin"},
		{Title: "Ethereum"},
	}
	mockCryptoProvider.On("GetActualRates", mock.Anything, mock.Anything).Return(dataCryptoProvider, nil)

	//сравнение с тестами

	ctx := context.Background()
	titles := []string{
		"Bitcoin", "Ethereum",
	}

	t.Run("Test One", func(t *testing.T) {
		mockCryptoProvider.GetActualRates(ctx, titles)
	})

	fmt.Println("_STOP TEST TestGetActualRates_")
}
