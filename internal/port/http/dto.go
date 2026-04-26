package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/pkg/errors"
)

type cryptoDTO struct {
	Titles []string `json:"titles"`
}

type aggregateDTO struct {
	Titles  []string `json:"titles"`
	AggType string   `json:"agg_type"`
}

type errorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func newErrorDTO(message string) *errorDTO {
	return &errorDTO{
		Message: message,
		Time:    time.Now(),
	}
}

func (e *errorDTO) toString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(b)
}

func responseJSON(w http.ResponseWriter, coins []*entities.Coin) {

	coinsJson, err := json.MarshalIndent(coins, "", "    ")
	if err != nil {
		err := errors.Wrapf(entities.ErrInternal, "failed to json response: %v", err)
		responseError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(coinsJson); err != nil {
		errors.Wrapf(entities.ErrInternal, "failed to write http response: %v", err)
		return
	}
}

func responseError(w http.ResponseWriter, message error, statusCode int) {
	errDTO := newErrorDTO(message.Error())
	http.Error(w, errDTO.toString(), statusCode)
}
