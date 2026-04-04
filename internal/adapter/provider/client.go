package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/pkg/errors"
)

const (
	baseURL        = "https://pro-api.coingecko.com/api/v3"
	priceMultiPath = "/simple/price"

	queryVs_currencies = "vs_currencies"
	queryIds           = "ids"
)

type client struct {
	apiAuthorization string //apiAuthorization := "x-cg-demo-api-key"
	apiKey           string //apiKey := "CG-SdGMn7C5Rv2F4hTMwLdJ1Pk6"
	vs_currencies    string // usd
	httpClient       *http.Client
}

type Option func(c *client)

func WithTimeout(timeout time.Duration) Option {
	return func(c *client) {
		c.httpClient.Timeout = timeout
	}
}

func WithCurrencies(vs_currencies string) Option {
	return func(c *client) {
		if strings.TrimSpace(vs_currencies) != "" {
			c.vs_currencies = vs_currencies
		}
	}
}

func NewProviderClient(apiAuthorization string, apiKey string, opts ...Option) (*client, error) {
	if strings.TrimSpace(apiAuthorization) == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "apiAuthorization cannot be empty")
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "apiKey cannot be empty")
	}
	c := &client{
		apiAuthorization: apiAuthorization,
		apiKey:           apiKey,
		vs_currencies:    "usd",
		httpClient:       &http.Client{Timeout: time.Second * 10},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil

}

func (c *client) GetActualRates(ctx context.Context, titles []string) ([]*entities.Coin, error) {
	///
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "empty titles")
	}

	qureyTitles := make([]string, len(titles))
	for idx, elem := range titles {
		qureyTitles[idx] = strings.ToLower(elem)
	}

	url := buildUrl(qureyTitles, c.vs_currencies)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "request error: %v", err)
	}
	request.Header.Add(c.apiAuthorization, c.apiKey)
	///

	response, err := c.httpClient.Do(request)
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
		if price, ok := elem[c.vs_currencies]; ok {
			if coin, err := entities.NewCoin(title, price, tmNow); err == nil {
				resultCoins = append(resultCoins, coin)
			}
		}
	}

	return resultCoins, nil
}

func buildUrl(titles []string, vs_currencies string) string {
	// url := "https://api.coingecko.com/api/v3/simple/price?vs_currencies=usd&ids=bitcoin%2Cethereum%2Ctron"

	url, err := url.Parse(fmt.Sprintf("%s%s", baseURL, priceMultiPath))
	if err != nil {
		err = errors.Wrap(entities.ErrInternal, "url parse")
	}
	query := url.Query()
	query.Set(queryVs_currencies, vs_currencies)
	query.Set(queryIds, strings.Join(titles, ","))
	url.RawQuery = query.Encode()

	return url.String()
}

// func (c *client) prepareRequest(ctx context.Context, titles []string) (*http.Request, error) {
// 	baseURL := "sacs"
// 	priceMultiPath := "jkj"
// 	// slog.Info("prepareRequest")
// 	// ctx, span, cancel := tracer.Start(ctx, "prepareRequest adapter provider")
// 	// defer cancel()
// 	url, _ := url.Parse(fmt.Sprintf("%s%s", baseURL, priceMultiPath))
// 	query := url.Query()
// 	query.Set()

// 	// url, err := url.Parse(fmt.Sprintf("%s%s", baseURL, priceMultiPath))
// 	// if err != nil {
// 	// 	err := errors.Wrap(entities.ErrInternal, "url parse")
// 	// 	span.SetError(err)
// 	// 	slog.Error("Parse", "err", err)
// 	// 	return nil, err
// 	// }

// 	// query := url.Query()
// 	// query.Set(queryFsyms, strings.Join(titles, ","))
// 	// query.Set(queryTsyms, c.costIn)
// 	// url.RawQuery = query.Encode()

// 	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
// 	if err != nil {
// 		err := errors.Wrap(entities.ErrInternal, "new request with context")
// 		span.SetError(err)
// 		slog.Error("NewRequestWithContext", "err", err)
// 		return nil, err
// 	}

// 	return request, nil
// }
