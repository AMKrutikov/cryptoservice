package porthttp

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/AMKrutikov/cryptoservice/internal/port"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	service port.Service
	router  mux.Router ////
}

func NewServer(service port.Service) *Server {
	//return &Server{
	s := &Server{ ///
		service: service,
		router:  *mux.NewRouter(), ////
	}
	s.routes() ///
	return s   ///
}

// //
func (s *Server) routes() { ////
	s.router.Path("/coins/rates").Methods("POST").HandlerFunc(s.GetLastRates)                ////
	s.router.Path("/coins/rates/aggregate").Methods("POST").HandlerFunc(s.GetAggregateRates) ////
}

func (s *Server) Start() error { ////
	return http.ListenAndServe(":9091", &s.router) ////
}

// //

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

	for i, title := range coinDTO.Titles {
		coinDTO.Titles[i] = strings.ToLower(title)
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

	for i, title := range titles {
		titles[i] = strings.ToLower(title)
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
