package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	vatcheck "github.com/gopherskatowice/vatcheck-svc"
	. "github.com/gopherskatowice/vatcheck-svc/config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Instance type
type Instance struct {
	Checker *vatcheck.Checker
	bind    string
	port    int
}

// New creates a new server instance.
func New(checker *vatcheck.Checker, bind string, port int) *Instance {
	return &Instance{checker, bind, port}
}

// ListenAndServe listens on the TCP network address and then
// calls handlers to handle requests.
func (srv *Instance) ListenAndServe() error {
	// Log information about server
	log.WithFields(log.Fields{
		"bind": srv.bind,
		"port": srv.port,
	}).Debug("HTTP server is listening")

	// Start the HTTP server
	return http.ListenAndServe(
		fmt.Sprintf("%s:%d", srv.bind, srv.port),
		srv.handlers(),
	)
}

// handlers returns httprouter handlers.
func (srv *Instance) handlers() http.Handler {
	e := echo.New()

	// Add middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: Config.CorsOrigins,
		AllowHeaders: Config.CorsHeaders,
		AllowMethods: Config.CorsMethods,
	}))

	// Enable debug mode
	e.Debug = log.GetLevel() == log.DebugLevel

	// Declare routes
	e.GET("/:vatid", srv.vatHandler)
	e.GET("/_health", srv.healthHandler)

	return e
}
