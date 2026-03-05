package cases

import "context"

type CryptoProvider interface {
	GetActualRates(ctx context.Context, titles []string) ([]Provider, error)
}

type Provider struct {
	Titl string
}

func (p *Provider) GetActualRates(ctx context.Context, titles []string) ([]Provider, error) {
	result := []Provider{}
	for _, elem := range titles {
		result = append(result, Provider{Titl: elem})
	}
	return result, nil
}
