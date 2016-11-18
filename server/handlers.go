package server

import (
	"net/http"

	"github.com/gamegos/jsend"
	"github.com/gopherskatowice/vatcheck-svc"
	"github.com/labstack/echo"
	"github.com/mattes/vat"
)

// getVatID checks in the VIES database if the given number is valid.
func (srv *Instance) getVatID(c echo.Context) error {
	vatid := c.Param("vatid")
	w := c.Response()

	valid, err := vatcheck.IsValid(vatid)
	if err != nil {
		jsend.Wrap(w).
			Status(getCodeForError(err)).
			Message(err.Error()).
			Send()
		return err
	}

	jsend.Wrap(w).
		Status(http.StatusOK).
		Data(valid).
		Send()
	return nil
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
