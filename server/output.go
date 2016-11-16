package server

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gamegos/jsend"
	"github.com/gocraft/health"
)

// writeSuccess outputs a JSend compliant message.
func (srv *Instance) writeSuccess(w http.ResponseWriter, req *http.Request, data interface{}, s int) {
	healthJob(req).Complete(health.Success)
	if _, err := jsend.Wrap(w).Status(s).Data(data).Send(); err != nil {
		srv.writeError(
			w, req, err, s,
		)
	}
}

// writeError outputs a JSend compliant error.
func (srv *Instance) writeError(w http.ResponseWriter, req *http.Request, err error, s int) {
	healthJob(req).Complete(health.Error)
	log.Error(err)
	if _, err = jsend.Wrap(w).Status(s).Message(err.Error()).Send(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
