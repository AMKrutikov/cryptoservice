package dto

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/pkg/errors"
)

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type CoinResponseDTO struct {
	Title    string    `json:"title"`
	Price    float64   `json:"price"`
	ActualAt time.Time `json:"actual_at"`
}

func ResponseJSON(w http.ResponseWriter, coins []*entities.Coin) {
	responseDTO := make([]CoinResponseDTO, 0, len(coins))

	for _, elem := range coins {
		responseDTO = append(responseDTO, CoinResponseDTO{
			Title:    elem.Title(),
			Price:    elem.Price(),
			ActualAt: elem.ActualAt(),
		})
	}

	coinsJson, err := json.MarshalIndent(responseDTO, "", "    ")
	if err != nil {
		err := errors.Wrapf(entities.ErrInternal, "failed to json response: %v", err)
		ResponseError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(coinsJson); err != nil {
		errors.Wrapf(entities.ErrInternal, "failed to write http response: %v", err)
		return
	}
}

func ResponseError(w http.ResponseWriter, message error, statusCode int) {
	ErrorDTO := ErrorDTO{Message: message.Error(), Time: time.Now()}
	b, err := json.MarshalIndent(ErrorDTO, "", "    ")
	if err != nil {
		panic(err)
	}
	http.Error(w, string(b), statusCode)
}
