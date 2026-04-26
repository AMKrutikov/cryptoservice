package porthttp

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	server *Server
}

func NewHTTPServer(server *Server) *HTTPServer {
	return &HTTPServer{
		server: server,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/coins/rates").Methods("POST").HandlerFunc(s.server.GetLastRates)
	router.Path("/coins/rates/aggregate").Methods("POST").HandlerFunc(s.server.GetAggregateRates)

	return http.ListenAndServe(":9091", router)
}
