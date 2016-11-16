package vatcheck

import (
	"errors"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mattes/vat"
	"github.com/sony/gobreaker"
)

var (
	c *check

	ErrInvalidFormat  = errors.New("Invalid format")
	ErrCircuitTripped = errors.New("The circuit is tripped")
)

type check struct {
	checkMethod func(id string) (*vat.VATresponse, error)
}

func init() {
	c = &check{
		checkMethod: vat.CheckVAT,
	}
}

// IsValid checks if the given number is a valid EU VAT number.
// The format is first checked and then the actual validity is
// checked against the VIES EU database.
func IsValid(id string) (bool, error) {
	// Strip whitespace
	id = strings.Replace(id, " ", "", -1)

	// Prepare the logger
	lgr := log.WithField("vatid", id)

	if vfmt := IsValidFormat(id); !vfmt {
		lgr.Debug("Invalid format")
		return false, ErrInvalidFormat
	}

	if value, exists := cache().Get(id); exists {
		lgr.Debug("Cache HIT")
		return strconv.ParseBool(value)
	}

	valid, err := cb.Execute(func() (interface{}, error) {
		data, err := c.checkMethod(id)
		if err != nil {
			return false, err
		}

		// Only return the number validity.
		// We could send back the data for prefilling/matching in a later iteration.
		return data.Valid, nil
	})
	if err != nil {
		// Check if the circuit is tripped
		if cb.State() == gobreaker.StateOpen {
			lgr.Warn("Circuit tripped")
			return false, ErrCircuitTripped
		}
		return false, err
	}

	// Keep in cache the result for a while
	cache().Set(id, strconv.FormatBool(valid.(bool)))

	return valid.(bool), nil
}
