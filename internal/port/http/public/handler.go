package public

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/AMKrutikov/cryptoservice/pkg/dto"
	"github.com/pkg/errors"
)

// @Summary      Getting rates cryptocurrencies
// @Description  Accepts coin names and returns the latest current prices
// @Tags         Crypto
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CryptoDTO  true  "list of coin titles"
// @Success      200      {array}   dto.CoinResponseDTO  "list crypto coins"
// @Failure      400      {object}  dto.ErrorDTO         "invalid JSON, empty list or unknown coins"
// @Failure      500      {object}  dto.ErrorDTO         "internal server error"
// @Router       /coins/rates [post]
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

// @Summary      Getting aggregated rates cryptocurrencies
// @Description  Accepts coin names and returns min, max or avg prices
// @Tags         Crypto
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AggregateDTO  true  "list titles coins and aggtype(min/max/avg)"
// @Success      200      {array}   dto.CoinResponseDTO  "list crypto coins and aggregated type"
// @Failure      400      {object}  dto.ErrorDTO         "invalid JSON, empty list, invalid aggType or unknown coins"
// @Failure      500      {object}  dto.ErrorDTO         "internal server error"
// @Router       /coins/rates/aggregate [post]
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
