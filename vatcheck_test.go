package vatcheck

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	. "github.com/gopherskatowice/vatcheck-svc/config"
	"github.com/mattes/vat"
)

var (
	errMock       = errors.New("mock error to trip the circuit breaker")
	validNumber   = "PL6492291353"
	invalidNumber = "XX1234567890"
)

func init() {
	// Disable logs
	log.SetOutput(ioutil.Discard)

	// Set a TTL for the cache or it'll default to 0
	Config.CacheDuration = time.Duration(1 * time.Minute)
}

func getChecker(resp *vat.VATresponse, err error) *Checker {
	method := func(id string) (*vat.VATresponse, error) {
		return resp, err
	}
	return New(method)
}

func TestFormatInvalid(t *testing.T) {
	t.Parallel()

	c := getChecker(nil, nil)
	if _, err := c.IsValid("fooo"); err != ErrInvalidFormat {
		t.Error("Format should be invalid")
	}
}

func TestIsValid(t *testing.T) {
	t.Parallel()

	c := getChecker(&vat.VATresponse{Valid: true}, nil)
	if isValid, _ := c.IsValid(validNumber); !isValid {
		t.Errorf("The number %s should be valid", validNumber)
	}
}

func TestInvalid(t *testing.T) {
	t.Parallel()

	c := getChecker(nil, errMock)
	if _, err := c.IsValid(invalidNumber); err == nil {
		t.Error("An error should be returned")
	}
}

func TestCircuitTripped(t *testing.T) {
	t.Parallel()

	// How many failures before the circuit trips
	countToTrip := 2

	c := getChecker(nil, errMock)

	// Trip the circuit
	for i := 0; i < countToTrip; i++ {
		if _, err := c.IsValid(validNumber); err != errMock {
			t.Error("The wrong error has been returned")
		}
	}

	if _, err := c.IsValid(validNumber); err != ErrCircuitTripped {
		t.Error("The wrong error has been returned")
	}
}

func TestIsCached(t *testing.T) {
	t.Parallel()

	c := getChecker(&vat.VATresponse{Valid: true}, nil)
	if _, exists := c.cache.Get(validNumber); exists {
		t.Errorf("The number %s shouldn't be in cache", validNumber)
	}
	if valid, _ := c.IsValid(validNumber); !valid {
		t.Errorf("The number %s should be valid", validNumber)
	}
	if _, exists := c.cache.Get(validNumber); !exists {
		t.Errorf("The number %s should be in cache", validNumber)
	}
}

func BenchmarkIsValid(b *testing.B) {

	c := getChecker(&vat.VATresponse{Valid: true}, nil)
	for n := 0; n < b.N; n++ {
		_, _ = c.IsValid(validNumber)
	}
}
