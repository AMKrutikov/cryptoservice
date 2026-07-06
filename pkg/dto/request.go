package dto

type CryptoDTO struct {
	Titles []string `json:"titles" example:"bitcoin,ethereum"`
}

type AggregateDTO struct {
	Titles  []string `json:"titles" example:"bitcoin,ethereum"`
	AggType string   `json:"agg_type" example:"avg"`
}

// omitempty
