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

type ProviderClient struct {
	apiAuthorization string
	apiKey           string
	//apiAuthorization := "x-cg-demo-api-key"
	//apiKey := "CG-SdGMn7C5Rv2F4hTMwLdJ1Pk6"
}

func NewProviderClient(apiAuthorization string, apiKey string) (*ProviderClient, error) {
	if strings.TrimSpace(apiAuthorization) == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "apiAuthorization cannot be empty")
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "apiKey cannot be empty")
	}
	return &ProviderClient{
		apiAuthorization: apiAuthorization,
		apiKey:           apiKey,
	}, nil
}

func (p *ProviderClient) GetActualRates(ctx context.Context, titles []string) ([]*entities.Coin, error) {

	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "empty titles")
	}

	qureyTitles := make([]string, len(titles))
	for idx, elem := range titles {
		qureyTitles[idx] = strings.ToLower(elem)
	}

	url := buildUrl(qureyTitles, "usd")

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParam, "request error")
	}
	request.Header.Add(p.apiAuthorization, p.apiKey)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParam, "response error")
	}
	defer response.Body.Close()

	parsedResponse := make(map[string]struct {
		Usd float64 `json:"usd"`
	})

	if err := json.NewDecoder(response.Body).Decode(&parsedResponse); err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParam, "invalid JSON response")
	}
	//
	// body, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	return nil, errors.Wrap(entities.ErrInvalidParam, "error reading response body")
	// }

	// if err := json.Unmarshal(body, &parsedResponse); err != nil {
	// 	return nil, errors.Wrap(entities.ErrInvalidParam, "invalid JSON response")
	// }
	//

	resultCoins := make([]*entities.Coin, 0, len(titles))
	for title, elem := range parsedResponse {
		coin, err := entities.NewCoin(title, elem.Usd, time.Now())
		if err != nil {
			errors.Wrapf(entities.ErrInvalidParam, "coin creation error: %s\n", title) //???
			continue
		}
		resultCoins = append(resultCoins, coin)
	}

	return resultCoins, nil
}

func buildUrl(titles []string, vs_currencies string) string {

	coins := strings.Join(titles, "%2C")
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?vs_currencies=%s&ids=%s", vs_currencies, coins)
	return url
}
