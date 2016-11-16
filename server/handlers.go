package server

import (
	"net/http"

	"github.com/jgautheron/workshop/vat"
	"github.com/mattes/vat"
	"github.com/nbio/hitch"
)

// getVatID checks in the VIES database if the given number is valid.
func (srv *Instance) getVatID(w http.ResponseWriter, r *http.Request) {
	p := hitch.Params(r)

	valid, err := vatcheck.IsValid(p.ByName("vatid"))
	if err != nil {
		srv.writeError(w, r, err, getCodeForError(err))
		return
	}

	srv.writeSuccess(w, r, valid, http.StatusOK)
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
