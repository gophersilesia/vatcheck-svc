package server

import (
	"net/http"

	"github.com/gamegos/jsend"
	"github.com/gopherskatowice/vatcheck-svc"
	"github.com/labstack/echo"
	"github.com/mattes/vat"
)

// vatHandler checks in the VIES database if the given number is valid.
func (srv *Instance) vatHandler(c echo.Context) (err error) {
	vatid := c.Param("vatid")
	w := c.Response()

	checker := vatcheck.New(vat.CheckVAT)
	valid, err := checker.IsValid(vatid)
	if err != nil {
		_, err := jsend.Wrap(w).
			Status(getCodeForError(err)).
			Message(err.Error()).
			Send()
		return err
	}

	_, err = jsend.Wrap(w).
		Status(http.StatusOK).
		Data(valid).
		Send()
	return err
}

// healthHandler checks in the VIES database if the given number is valid.
func (srv *Instance) healthHandler(c echo.Context) (err error) {
	w := c.Response()
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"alive": true}`))
	return err
}

// getCodeForError returns the matching HTTP code for a given error.
func getCodeForError(err error) (code int) {
	// 500 unless specified otherwise
	code = http.StatusInternalServerError
	switch err {
	case vatcheck.ErrInvalidFormat:
		code = http.StatusBadRequest
	case vatcheck.ErrCircuitTripped:
		fallthrough
	case vat.ErrVATnumberNotValid:
		fallthrough
	case vat.ErrVATserviceUnreachable:
		code = http.StatusServiceUnavailable
	}
	return
}
