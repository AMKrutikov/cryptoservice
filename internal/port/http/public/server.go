package public

import (
	"context"
	"net/http"

	_ "github.com/AMKrutikov/cryptoservice/docs"
	"github.com/AMKrutikov/cryptoservice/internal/port"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	basePath      = "/cryptoservice/v1"
	ratesPath     = basePath + "/coins/rates"
	aggregatePath = ratesPath + "/aggregate"
)

type Server struct {
	service    port.Service
	router     *mux.Router
	httpServer *http.Server
}

func NewServer(service port.Service) *Server {
	return &Server{
		service: service,
		router:  mux.NewRouter(),
	}
}

func (s *Server) StartServer(address string) error {
	s.router.HandleFunc(ratesPath, s.GetLastRates).Methods(http.MethodPost)
	s.router.HandleFunc(aggregatePath, s.GetAggregateRates).Methods(http.MethodPost)

	s.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	s.httpServer = &http.Server{
		Addr:    address,
		Handler: s.router,
	}

	return http.ListenAndServe(address, s.router)
}

func (s *Server) Stop(ctx context.Context) error {

	if s.httpServer == nil {
		return nil
	}

	return s.httpServer.Shutdown(ctx)
}
