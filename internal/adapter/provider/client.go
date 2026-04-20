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
	defaultCostIn = "usd"

	baseURL        = "https://api.coingecko.com/api/v3"
	priceMultiPath = "/simple/price"

	queryVs_currencies = "vs_currencies"
	queryIds           = "ids"
)

type client struct {
	apiAuthorization string // apiAuthorization := "x-cg-demo-api-key"
	apiKey           string // apiKey := "CG-SdGMn7C5Rv2F4hTMwLdJ1Pk6"
	vs_currencies    string // usd
	httpClient       *http.Client
}

type Option func(c *client)

func WithTimeout(timeout time.Duration) Option {
	return func(c *client) {
		if timeout > 0 {
			c.httpClient.Timeout = timeout
		}
	}
}

func WithCurrencies(vs_currencies string) Option {
	return func(c *client) {
		if strings.TrimSpace(vs_currencies) != "" {
			c.vs_currencies = strings.ToLower(vs_currencies)
		}
	}
}

func (c *client) setOptions(opts ...Option) {
	for _, opt := range opts {
		if opt != nil {
			opt(c)
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
		vs_currencies:    strings.ToLower(defaultCostIn),
		httpClient:       &http.Client{Timeout: time.Second * 10},
	}

	c.setOptions(opts...)

	return c, nil

}

func (c *client) GetActualRates(ctx context.Context, titles []string) ([]*entities.Coin, error) {

	request, err := c.prepareRequest(ctx, titles)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInternal, "response error")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.Wrap(entities.ErrInternal, "provider api error")
	}

	resultCoins, err := c.decodeResponse(response)
	if err != nil {
		return nil, err
	}

	return resultCoins, nil
}

func buildUrl(titles []string, vs_currencies string) (string, error) {

	urlRaw, err := url.Parse(fmt.Sprintf("%s%s", baseURL, priceMultiPath))
	if err != nil {
		return "", errors.Wrap(entities.ErrInternal, "url parse failed")
	}
	query := urlRaw.Query() // map[string][]string
	query.Set(queryVs_currencies, vs_currencies)
	query.Set(queryIds, strings.Join(titles, ","))
	urlRaw.RawQuery = query.Encode()

	return urlRaw.String(), nil
}

func (c *client) prepareRequest(ctx context.Context, titles []string) (*http.Request, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "empty titles")
	}

	qureyTitles := make([]string, len(titles))
	for idx, elem := range titles {
		qureyTitles[idx] = strings.ToLower(elem)
	}

	url, err := buildUrl(qureyTitles, c.vs_currencies)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "request error: %v", err)
	}
	request.Header.Set(c.apiAuthorization, c.apiKey)

	return request, nil
}

func (c *client) decodeResponse(response *http.Response) ([]*entities.Coin, error) {
	parsedResponse := make(map[string]map[string]float64)

	if err := json.NewDecoder(response.Body).Decode(&parsedResponse); err != nil {
		return nil, errors.Wrap(entities.ErrInternal, "invalid JSON response")
	}

	resultCoins := make([]*entities.Coin, 0, len(parsedResponse))
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
