package public

import (
	"net/http"

	"github.com/AMKrutikov/cryptoservice/internal/port"
	"github.com/gorilla/mux"
)

type Server struct {
	service port.Service
}

type HTTPServer struct {
	server *Server
	addres string
}

func NewHTTPServer(server *Server, addres string) *HTTPServer {
	return &HTTPServer{
		server: server,
		addres: addres,
	}
}

func NewServer(service port.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/coins/rates").Methods(http.MethodPost).HandlerFunc(s.server.GetLastRates)
	router.Path("/coins/rates/aggregate").Methods(http.MethodPost).HandlerFunc(s.server.GetAggregateRates)

	return http.ListenAndServe(s.addres, router) //
}
