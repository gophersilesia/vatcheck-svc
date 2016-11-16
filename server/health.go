package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	. "github.com/jgautheron/workshop/vat/config"
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/health"
	"github.com/nbio/httpcontext"
)

const (
	// healthJobKey is a generic name per request
	healthJobKey = "health_job"

	// Jobs per endpoint
	jobCheckVat = "check_vat"
)

func initHealthStream() (stream *health.Stream) {
	stream = health.NewStream()
	stream.AddSink(&health.WriterSink{os.Stdout})
	return
}

// initHealthServer configures and launches the health check server.
// As a result the /health endpoint will be exposed and accessible.
// See https://github.com/gocraft/health#start-healthd
func (srv *Instance) initHealthServer() {
	sink := health.NewJsonPollingSink(time.Minute, time.Minute*5)
	log.WithFields(log.Fields{
		"bind": Config.HealthBind,
		"port": Config.HealthPort,
	}).Debug("Health server is listening")
	srv.stream.AddSink(sink)
	sink.StartServer(fmt.Sprintf("%s:%d", Config.HealthBind, Config.HealthPort))
}

// initHealthJob initialises a health job for the current request.
func (srv *Instance) initHealthJob(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// There is one endpoint so there's a single job
		httpcontext.Set(req, healthJobKey, srv.stream.NewJob(jobCheckVat))
		next.ServeHTTP(w, req)
	})
}

// healthJob extracts the job from the request.
func healthJob(req *http.Request) *health.Job {
	return httpcontext.Get(req, healthJobKey).(*health.Job)
}
