package server

import (
	"fmt"
	"net/http"

	. "github.com/jgautheron/workshop/vat/config"
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/health"
	"github.com/nbio/hitch"
	"github.com/rs/cors"
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
	h := hitch.New()

	if !srv.disableCors {
		// The below is necessary when interacting with the API using a regular
		// browser and not in s2s scenarios.
		h.Use(cors.New(cors.Options{
			AllowedOrigins: Config.CorsOrigins,
			AllowedHeaders: Config.CorsHeaders,
			AllowedMethods: Config.CorsMethods,
			Debug:          log.GetLevel() == log.DebugLevel,
		}).Handler)
	}

	// Middlewares
	h.Use(srv.initContentType)
	h.Use(srv.initHealthJob)

	// Declare routes
	h.Get("/:vatid", http.HandlerFunc(srv.getVatID))

	return h.Handler()
}

// initContentType middleware sets the HTTP content-type.
// We consider that the output will always be JSON.
func (srv *Instance) initContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, req)
	})
}
