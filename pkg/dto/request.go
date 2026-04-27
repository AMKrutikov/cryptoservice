package dto

type CryptoDTO struct {
	Titles []string `json:"titles"`
}

type AggregateDTO struct {
	Titles  []string `json:"titles"`
	AggType string   `json:"agg_type"`
}