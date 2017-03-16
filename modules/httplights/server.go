package httplights

import (
	"net/http"

	"github.com/collinux/gohue"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPLightsModule is an interface that can be implemented to allow a type to
// be registered as a http controlling lights module in the httplights.Server.
type HTTPLightsModule interface {
	// SetHueBridge configures the hue bridge used by the module.
	SetHueBridge(*hue.Bridge)

	// RegisterInRouter asks the module to register itself in a mux router.
	RegisterInRouter(*mux.Router)

	// ServeHTTP response to the HTTP request triggering the module
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// Server provides a means to register http controlled light modules that may
// be triggered over http.
type Server struct {
	// HueBridge configures the bridge to control lights through.
	HueBridge *hue.Bridge

	// ServerAddr configures the address that the HTTP server should listen on.
	ServerAddr string

	modules []HTTPLightsModule

	logger logrus.FieldLogger
}

// RegisterModule registers a http lights module.
func (s *Server) RegisterModule(module HTTPLightsModule) {
	s.modules = append(s.modules, module)
}

// Start starts the http light server
func (s *Server) Start() {
	s.logger = logrus.WithFields(logrus.Fields{
		"module":  "httplights",
		"address": s.ServerAddr,
	})

	router := mux.NewRouter()

	for _, module := range s.modules {
		module.RegisterInRouter(router)
		module.SetHueBridge(s.HueBridge)
	}

	go http.ListenAndServe(s.ServerAddr, router)

	s.logger.Info("Started http lights server.")
}
