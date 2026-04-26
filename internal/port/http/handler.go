package http

import (
	"encoding/json"
	"net/http"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/AMKrutikov/cryptoservice/internal/port"
	"github.com/pkg/errors"
)

type Server struct {
	service port.Service
}

func NewServer(service port.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) GetLastRates(w http.ResponseWriter, r *http.Request) {
	var coinDTO cryptoDTO

	if err := json.NewDecoder(r.Body).Decode(&coinDTO); err != nil {
		err := errors.Wrapf(entities.ErrInvalidParam, "invalid json format: %v", err)
		responseError(w, err, http.StatusBadRequest)
		return
	}

	titles := coinDTO.Titles

	if len(titles) == 0 {
		err := errors.Wrap(entities.ErrInvalidParam, "empty request")
		responseError(w, err, http.StatusBadRequest)
		return
	}

	coins, err := s.service.GetLastRates(r.Context(), titles)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidParam) {
			responseError(w, err, http.StatusBadRequest)
			return
		}
		responseError(w, err, http.StatusInternalServerError)
		return
	}

	responseJSON(w, coins)

}

func (s *Server) GetAggregateRates(w http.ResponseWriter, r *http.Request) {
	var aggDTO aggregateDTO

	if err := json.NewDecoder(r.Body).Decode(&aggDTO); err != nil {
		err := errors.Wrapf(entities.ErrInvalidParam, "invalid json format: %v", err)
		responseError(w, err, http.StatusBadRequest)
		return
	}

	titles := aggDTO.Titles
	aggType := aggDTO.AggType

	if len(titles) == 0 {
		err := errors.Wrap(entities.ErrInvalidParam, "empty request")
		responseError(w, err, http.StatusBadRequest)
		return
	}

	coins, err := s.service.GetAggregateRates(r.Context(), titles, aggType)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidParam) {
			responseError(w, err, http.StatusBadRequest)
			return
		}
		responseError(w, err, http.StatusInternalServerError)
		return
	}

	responseJSON(w, coins)
}
