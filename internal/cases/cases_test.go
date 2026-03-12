package cases_test

// gomock или  Mockery

import (
	"context"
	"testing"
	"time"

	"github.com/AMKrutikov/cryptoservice/internal/cases"
	"github.com/AMKrutikov/cryptoservice/internal/cases/mocks"
	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	// Пакеты из вашего go.mod
	// Пути к вашим внутренним пакетам
)

func TestCryptoProvider_Mock(t *testing.T) {
	// 1. Создаем экземпляр мока
	// Конструктор NewCryptoProvider генерируется mockery автоматически
	mockProvider := mocks.NewCryptoProvider(t)

	// 2. Подготавливаем тестовые данные
	ctx := context.Background()
	titles := []string{"BTC", "ETH"}

	coin1, _ := entities.NewCoin("BTC", 60000, time.Now())
	coin2, _ := entities.NewCoin("ETH", 3000, time.Now())

	// Создаем ожидаемый результат, который должен вернуть мок
	expectedCoins := []*entities.Coin{coin1, coin2}

	// 3. Настраиваем "Ожидание" (Expectation)
	// Мы говорим моку: "Когда у тебя вызовут GetActualRates с этими аргументами,
	// верни нам список монет и пустую ошибку (nil)"
	mockProvider.On("GetActualRates", ctx, titles).
		Return(expectedCoins, nil)

	// 4. Вызываем метод непосредственно у мока
	// В реальном коде этот вызов сделает ваш сервис, но здесь мы вызываем сами
	result, err := mockProvider.GetActualRates(ctx, titles)

	// 5. Проверяем, что мок вернул именно то, что мы в него заложили
	assert.NoError(t, err)                    // Ошибки быть не должно
	assert.Equal(t, expectedCoins, result)    // Результаты должны совпадать
	assert.Equal(t, "BTC", result[0].Title()) // Сравнение первого элемента в слайсе
}

func TestAAA(t *testing.T) {
	t.Run("test number one", func(t *testing.T) {
		mockCryptoProvider := mocks.NewCryptoProvider(t)
		tmr := time.Now()
		ctx := context.Background()
		expectedTitles := []string{"A1", "A2", "A3"}
		c1, _ := entities.NewCoin("A1", 0, tmr)
		c2, _ := entities.NewCoin("A2", 0, tmr)
		c3, _ := entities.NewCoin("A3", 0, tmr)
		expectedCoin := []*entities.Coin{c1, c2, c3}

		mockCryptoProvider.On("GetActualRates", mock.Anything, expectedTitles).
			Return(expectedCoin, nil)

		result, err := cases.AAA(ctx, expectedTitles, mockCryptoProvider)
		require.NoError(t, err)
		require.Equal(t, expectedCoin, result)
	})
}
