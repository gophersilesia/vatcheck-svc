package server

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gopherskatowice/vatcheck-svc"
	"github.com/labstack/echo"
	"github.com/mattes/vat"
)

// Response defines the standard JSON output.
type Response struct {
	Success bool `json:"success"`
}

// vatHandler checks in the VIES database if the given number is valid.
func (srv *Instance) vatHandler(c echo.Context) (err error) {
	vatid := c.Param("vatid")
	status := http.StatusOK

	valid, err := srv.Checker.IsValid(vatid)
	if err != nil {
		log.WithField("vatid", vatid).Error(err)
		status = getCodeForError(err)
	}

	return c.JSON(status, Response{valid})
}

// healthHandler allows to check externally if the service is reachable.
func (srv *Instance) healthHandler(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, Response{true})
}

// getCodeForError returns the matching HTTP code for a given error.
func getCodeForError(err error) (code int) {
	// 500 unless specified otherwise
	code = http.StatusInternalServerError
	switch err {
	case vatcheck.ErrInvalidFormat:
		code = http.StatusBadRequest
	case vatcheck.ErrCircuitTripped, vat.ErrVATnumberNotValid, vat.ErrVATserviceUnreachable:
		code = http.StatusServiceUnavailable
	}
	return
}
