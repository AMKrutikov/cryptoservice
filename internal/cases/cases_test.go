package cases_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"githab.com/AMKrutikov/cryptoservice/internal/cases"
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

func (c *MockCryptoProvider) GetActualRates(ctx context.Context, titles []string) ([]cases.Provider, error) {
	argumets := c.Called(ctx, titles)
	return argumets.Get(0).([]cases.Provider), argumets.Error(1)
}

// Storage Mock

// 	dataCryptoProvider := []cases.Coin{
// 		{Titl: "Bitcoin"},
// 		{Titl: "Ethereum"},
// 	}
// 	mockCryptoProvider.On("GetActualRates").Return(dataCryptoProvider, nil)

func TestGetActualRates(t *testing.T) {
	fmt.Println("_START TEST TestRunService_")

	mockCryptoProvider := new(MockCryptoProvider)
	t.Cleanup(func() {
		mockCryptoProvider.AssertExpectations(t)
		fmt.Println("THE END TEST TestRunService")
	})

	dataCryptoProvider := []cases.Provider{
		{Titl: "Bitcoin"},
		{Titl: "Ethereum"},
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

	fmt.Println("_STOP TEST TestRunService_")
}
