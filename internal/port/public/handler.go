package public

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/AMKrutikov/cryptoservice/pkg/dto"
	"github.com/pkg/errors"
)

func (s *Server) GetLastRates(w http.ResponseWriter, r *http.Request) {
	var coinDTO dto.CryptoDTO

	if err := json.NewDecoder(r.Body).Decode(&coinDTO); err != nil {
		err := errors.Wrapf(entities.ErrInvalidParam, "invalid json format: %v", err)
		dto.ResponseError(w, err, http.StatusBadRequest)
		return
	}

	titles := coinDTO.Titles

	if len(titles) == 0 {
		err := errors.Wrap(entities.ErrInvalidParam, "empty request")
		dto.ResponseError(w, err, http.StatusBadRequest)
		return
	}

	for i, title := range coinDTO.Titles {
		coinDTO.Titles[i] = strings.ToLower(title)
	}

	coins, err := s.service.GetLastRates(r.Context(), titles)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidParam) {
			dto.ResponseError(w, err, http.StatusBadRequest)
			return
		}
		dto.ResponseError(w, err, http.StatusInternalServerError)
		return
	}

	dto.ResponseJSON(w, coins)

}

func (s *Server) GetAggregateRates(w http.ResponseWriter, r *http.Request) {
	var aggDTO dto.AggregateDTO

	if err := json.NewDecoder(r.Body).Decode(&aggDTO); err != nil {
		err := errors.Wrapf(entities.ErrInvalidParam, "invalid json format: %v", err)
		dto.ResponseError(w, err, http.StatusBadRequest)
		return
	}

	titles := aggDTO.Titles
	aggType := aggDTO.AggType

	if len(titles) == 0 {
		err := errors.Wrap(entities.ErrInvalidParam, "empty request")
		dto.ResponseError(w, err, http.StatusBadRequest)
		return
	}

	for i, title := range titles {
		titles[i] = strings.ToLower(title)
	}

	coins, err := s.service.GetAggregateRates(r.Context(), titles, aggType)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidParam) {
			dto.ResponseError(w, err, http.StatusBadRequest)
			return
		}
		dto.ResponseError(w, err, http.StatusInternalServerError)
		return
	}

	dto.ResponseJSON(w, coins)
}
