package public

import (
	"fmt"
	"net/http"

	"github.com/AMKrutikov/cryptoservice/internal/port"
	"github.com/gorilla/mux"
)

const (
	basePath      = "/cryptoservice/v1"
	ratesPath     = "/coins/rates"
	aggregatePath = "/aggregate"
)

type Server struct {
	service port.Service
	router  *mux.Router
}

func NewServer(service port.Service) *Server {
	return &Server{
		service: service,
		router:  mux.NewRouter(),
	}
}

func (s *Server) StartServer(address string) error {
	s.router.Path(fmt.Sprintf("%s%s", basePath, ratesPath)).Methods(http.MethodPost).HandlerFunc(s.GetLastRates)
	s.router.Path(fmt.Sprintf("%s%s%s", basePath, ratesPath, aggregatePath)).Methods(http.MethodPost).HandlerFunc(s.GetAggregateRates)

	return http.ListenAndServe(address, s.router)
}
