package httpServer

import (
	"net/http"

	"github.com/wDRxxx/test-task/internal/api"
	"github.com/wDRxxx/test-task/internal/service"
)

type server struct {
	mux        http.Handler
	apiService service.ApiService
}

func NewHTTPServer(apiService service.ApiService) api.HTTPServer {
	s := &server{
		apiService: apiService,
	}
	s.setRoutes()

	return s
}

func (s *server) Handler() http.Handler {
	return s.mux
}
