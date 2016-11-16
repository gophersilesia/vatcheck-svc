package vatcheck

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	. "github.com/jgautheron/workshop/vat/config"
	"github.com/mattes/vat"
)

var (
	errMock     = errors.New("mock error to trip the circuit breaker")
	validNumber = "PL6492291353"
)

func init() {
	// Disable logs
	log.SetOutput(ioutil.Discard)
}

func mockCheckVAT(id string) (*vat.VATresponse, error) {
	return &vat.VATresponse{Valid: true}, nil
}

func TestFormatInvalid(t *testing.T) {
	t.Parallel()

	if _, err := IsValid("fooo"); err != ErrInvalidFormat {
		t.Error("Format should be invalid")
	}
}

func TestIsValid(t *testing.T) {
	t.Parallel()

	c.checkMethod = func(id string) (*vat.VATresponse, error) {
		return &vat.VATresponse{Valid: true}, nil
	}

	if isValid, _ := IsValid(validNumber); !isValid {
		t.Error("The number %s should be valid", validNumber)
	}
}

func TestInvalid(t *testing.T) {
	t.Parallel()

	c.checkMethod = func(id string) (*vat.VATresponse, error) {
		return nil, errMock
	}

	if _, err := IsValid(validNumber); err == nil {
		t.Error("An error should be returned")
	}
}

func TestCircuitTripped(t *testing.T) {
	// Re-initialise the circuit breaker
	initCircuitBreaker()

	// How many failures before the circuit trips
	countToTrip := 2

	c.checkMethod = func(id string) (*vat.VATresponse, error) {
		return nil, errMock
	}

	// Trip the circuit
	for i := 0; i < countToTrip; i++ {
		if _, err := IsValid(validNumber); err != errMock {
			t.Error("The wrong error has been returned")
		}
	}

	if _, err := IsValid(validNumber); err != ErrCircuitTripped {
		t.Error("The wrong error has been returned")
	}
}

func TestIsCached(t *testing.T) {
	// Re-initialise the circuit breaker
	initCircuitBreaker()

	// Set a TTL for the cache
	ch = nil
	Config.CacheDuration = time.Duration(2) * time.Minute

	// Use another VAT number for testing since the other one is cached
	validNumber := "FR33440953859"

	if _, exists := cache().Get(validNumber); exists {
		t.Error("The number %s shouldn't be in cache", validNumber)
	}

	c.checkMethod = func(id string) (*vat.VATresponse, error) {
		return &vat.VATresponse{Valid: true}, nil
	}

	if isValid, _ := IsValid(validNumber); !isValid {
		t.Error("The number %s should be valid", validNumber)
	}
	if _, exists := cache().Get(validNumber); !exists {
		t.Error("The number %s should be in cache", validNumber)
	}
}

func BenchmarkIsValid(b *testing.B) {
	ch = nil
	Config.CacheDuration = time.Duration(2) * time.Minute

	c.checkMethod = func(id string) (*vat.VATresponse, error) {
		return &vat.VATresponse{Valid: true}, nil
	}
	for n := 0; n < b.N; n++ {
		_, _ = IsValid(validNumber)
	}
}
