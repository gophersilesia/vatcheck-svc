// Package vatcheck uses the official VIES VAT number validation SOAP
// web service to validate european VAT numbers.
package vatcheck

import (
	"errors"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	. "github.com/gopherskatowice/vatcheck-svc/config"
	"github.com/mattes/vat"
	"github.com/sony/gobreaker"
	"github.com/wunderlist/ttlcache"
)

var (
	// ErrInvalidFormat holds error for invalid VAT number
	ErrInvalidFormat = errors.New("Invalid format")
	// ErrCircuitTripped holds error for circuit
	ErrCircuitTripped = errors.New("The circuit is tripped")
)

type Checker struct {
	method  func(id string) (*vat.VATresponse, error)
	cache   *ttlcache.Cache
	circuit *gobreaker.CircuitBreaker
}

// New creates a new checker instance with the given check method.
// Eventually the cache and circuit could also be mocked by replacing
// the types with interfaces.
func New(method func(id string) (*vat.VATresponse, error)) *Checker {
	ch := ttlcache.NewCache(Config.CacheDuration)
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "VIES",
		Timeout: time.Duration(20) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	})
	return &Checker{method, ch, cb}
}

// IsValid checks if the given number is a valid EU VAT number.
// The format is first checked and then the actual validity is
// checked against the VIES EU database.
func (c *Checker) IsValid(id string) (bool, error) {
	// Strip whitespace
	id = strings.Replace(id, " ", "", -1)

	// Prepare the logger
	lgr := log.WithField("vatid", id)

	if vfmt := IsValidFormat(id); !vfmt {
		lgr.Debug("Invalid format")
		return false, ErrInvalidFormat
	}

	if value, exists := c.cache.Get(id); exists {
		lgr.Debug("Cache HIT")
		return strconv.ParseBool(value)
	}

	valid, err := c.circuit.Execute(func() (interface{}, error) {
		data, err := c.method(id)
		if err != nil {
			return false, err
		}
		// Only return the number validity.
		// We could send back the data for prefilling/matching in a later iteration.
		return data.Valid, nil
	})
	if err != nil {
		// Check if the circuit is tripped
		if c.circuit.State() == gobreaker.StateOpen {
			lgr.Warn("Circuit tripped")
			return false, ErrCircuitTripped
		}
		return false, err
	}

	// Keep in cache the result for a while
	c.cache.Set(id, strconv.FormatBool(valid.(bool)))

	return valid.(bool), nil
}
