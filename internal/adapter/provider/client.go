package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/pkg/errors"
)

// type Client struct {
// 	*http.Client
// }

type coinGeckoClient struct {
	apiAuthorization string //apiAuthorization := "x-cg-demo-api-key"
	apiKey           string //apiKey := "CG-SdGMn7C5Rv2F4hTMwLdJ1Pk6"
	vs_currencies    string // usd
	client           *http.Client
	//client           *http.Client
}

func NewProviderClient(apiAuthorization string, apiKey string, vs_currencies string) (*coinGeckoClient, error) {
	if strings.TrimSpace(apiAuthorization) == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "apiAuthorization cannot be empty")
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "apiKey cannot be empty")
	}
	if strings.TrimSpace(vs_currencies) == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "vs_currencies cannot be empty")
	}
	return &coinGeckoClient{
		apiAuthorization: apiAuthorization,
		apiKey:           apiKey,
		vs_currencies:    strings.ToLower(vs_currencies),
		//client:           client,
	}, nil
}

func (с *coinGeckoClient) GetActualRates(ctx context.Context, titles []string) ([]*entities.Coin, error) {
	///
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "empty titles")
	}

	qureyTitles := make([]string, len(titles))
	for idx, elem := range titles {
		qureyTitles[idx] = strings.ToLower(elem)
	}

	url := buildUrl(qureyTitles, с.vs_currencies)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "request error: %v", err)
	}
	request.Header.Add(с.apiAuthorization, с.apiKey)
	///

	response, err := http.DefaultClient.Do(request) // создать в конструкторе , добавить паттерн функциональные опции для конфигурирования таймаута клиента
	//response, err := c.

	if err != nil {
		return nil, errors.Wrap(entities.ErrInternal, "response error")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.Wrap(entities.ErrInternal, "provider api error")
	}
	///
	parsedResponse := make(map[string]map[string]float64)

	if err := json.NewDecoder(response.Body).Decode(&parsedResponse); err != nil {
		return nil, errors.Wrap(entities.ErrInternal, "invalid JSON response")
	}

	resultCoins := make([]*entities.Coin, 0, len(titles))
	tmNow := time.Now()
	for title, elem := range parsedResponse {
		if price, ok := elem[с.vs_currencies]; ok {
			if coin, err := entities.NewCoin(title, price, tmNow); err == nil {
				resultCoins = append(resultCoins, coin)
			}
		}
	}

	return resultCoins, nil
}

func buildUrl(titles []string, vs_currencies string) string {

	coins := strings.Join(titles, "%2C")
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?vs_currencies=%s&ids=%s", vs_currencies, coins)
	return url
}
