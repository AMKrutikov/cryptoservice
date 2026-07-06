package public

import (
	"fmt"
	"net/http"

	_ "github.com/AMKrutikov/cryptoservice/docs"
	"github.com/AMKrutikov/cryptoservice/internal/port"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	basePath      = "/cryptoservice/v1"
	ratesPath     = "/cryptoservice/v1/coins/rates/aggregate"
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

	s.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return http.ListenAndServe(address, s.router)
}
