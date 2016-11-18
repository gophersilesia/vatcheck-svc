package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/health"
	. "github.com/gopherskatowice/vatcheck-svc/config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Instance type
type Instance struct {
	bind        string
	port        int
	disableCors bool

	stream *health.Stream
}

// New creates a new server instance.
func New(bind string, port int, disableCors bool) *Instance {
	return &Instance{bind, port, false, initHealthStream()}
}

// ListenAndServe listens on the TCP network address and then
// calls handlers to handle requests.
func (srv *Instance) ListenAndServe() {
	// Log information about server
	log.WithFields(log.Fields{
		"bind": srv.bind,
		"port": srv.port,
	}).Debug("Server is listening")

	// Start the health server
	if !Config.DisableHealth {
		srv.initHealthServer()
	}

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%d", srv.bind, srv.port),
		srv.handlers(),
	))
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
	e.GET("/:vatid", srv.getVatID)

	return e
}
