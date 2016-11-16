package vatcheck_test

import (
	"strings"
	"testing"

	"github.com/jgautheron/workshop/vat"
	"github.com/jmcvetta/randutil"
)

var numbers = []struct {
	vatNumber   string
	countryCode string
	isValid     bool
}{
	{"0" + randNumeric(9), "BE", true},
	{"1" + randNumeric(9), "BE", false},
	{"0" + randNumeric(10), "BE", false},
	{"0" + randNumeric(5), "BE", false},
	{"A" + randNumeric(9), "BE", false},
	{"0" + randNumeric(8) + "A", "BE", false},
	{randNumeric(9) + "A", "BE", false},

	{randNumeric(9), "BG", true},
	{randNumeric(10), "BG", true},
	{randNumeric(11), "BG", false},
	{randNumeric(5), "BG", false},
	{"A" + randNumeric(9), "BG", false},
	{randNumeric(8) + "A", "BG", false},

	{randNumeric(8), "CZ", true},
	{randNumeric(9), "CZ", true},
	{randNumeric(10), "CZ", true},
	{randNumeric(11), "CZ", false},
	{randNumeric(5), "CZ", false},
	{"A" + randNumeric(9), "CZ", false},
	{randNumeric(8) + "A", "CZ", false},

	{randNumeric(8), "DK", true},
	{randNumeric(11), "DK", false},
	{randNumeric(5), "DK", false},
	{"A" + randNumeric(9), "DK", false},
	{randNumeric(8) + "A", "DK", false},

	{randNumeric(9), "DE", true},
	{randNumeric(11), "DE", false},
	{randNumeric(5), "DE", false},
	{"A" + randNumeric(9), "DE", false},
	{randNumeric(8) + "A", "DE", false},

	{randNumeric(9), "EE", true},
	{randNumeric(11), "EE", false},
	{randNumeric(5), "EE", false},
	{"A" + randNumeric(9), "EE", false},
	{randNumeric(8) + "A", "EE", false},

	{randNumeric(9), "EL", true},
	{randNumeric(11), "EL", false},
	{randNumeric(5), "EL", false},
	{"A" + randNumeric(9), "EL", false},
	{randNumeric(8) + "A", "EL", false},

	{randNumeric(11), "HR", true},
	{randNumeric(15), "HR", false},
	{randNumeric(5), "HR", false},
	{"A" + randNumeric(9), "HR", false},
	{randNumeric(8) + "A", "HR", false},

	{randNumeric(11), "IT", true},
	{randNumeric(15), "IT", false},
	{randNumeric(5), "IT", false},
	{"A" + randNumeric(9), "IT", false},
	{randNumeric(8) + "A", "IT", false},

	{randNumeric(11), "LV", true},
	{randNumeric(15), "LV", false},
	{randNumeric(5), "LV", false},
	{"A" + randNumeric(9), "LV", false},
	{randNumeric(8) + "A", "LV", false},

	{randNumeric(8), "LU", true},
	{randNumeric(11), "LU", false},
	{randNumeric(5), "LU", false},
	{"A" + randNumeric(9), "LU", false},
	{randNumeric(8) + "A", "LU", false},

	{randNumeric(8), "HU", true},
	{randNumeric(11), "HU", false},
	{randNumeric(5), "HU", false},
	{"A" + randNumeric(9), "HU", false},
	{randNumeric(8) + "A", "HU", false},

	{randNumeric(8), "MT", true},
	{randNumeric(11), "MT", false},
	{randNumeric(5), "MT", false},
	{"A" + randNumeric(9), "MT", false},
	{randNumeric(8) + "A", "MT", false},

	{randNumeric(8), "FI", true},
	{randNumeric(11), "FI", false},
	{randNumeric(5), "FI", false},
	{"A" + randNumeric(9), "FI", false},
	{randNumeric(8) + "A", "FI", false},

	{randNumeric(8), "SI", true},
	{randNumeric(11), "SI", false},
	{randNumeric(5), "SI", false},
	{"A" + randNumeric(9), "SI", false},
	{randNumeric(8) + "A", "SI", false},

	{randNumeric(12), "SE", true},
	{randNumeric(15), "SE", false},
	{randNumeric(5), "SE", false},
	{"A" + randNumeric(9), "SE", false},
	{randNumeric(8) + "A", "SE", false},

	{randNumeric(10), "SK", true},
	{randNumeric(15), "SK", false},
	{randNumeric(5), "SK", false},
	{"A" + randNumeric(9), "SK", false},
	{randNumeric(8) + "A", "SK", false},

	{randNumeric(10), "PL", true},
	{randNumeric(15), "PL", false},
	{randNumeric(5), "PL", false},
	{"A" + randNumeric(9), "PL", false},
	{randNumeric(8) + "A", "PL", false},

	{randNumeric(9), "PT", true},
	{randNumeric(15), "PT", false},
	{randNumeric(5), "PT", false},
	{"A" + randNumeric(9), "PT", false},
	{randNumeric(8) + "A", "PT", false},

	{randNumeric(2), "RO", true},
	{randNumeric(3), "RO", true},
	{randNumeric(4), "RO", true},
	{randNumeric(5), "RO", true},
	{randNumeric(6), "RO", true},
	{randNumeric(7), "RO", true},
	{randNumeric(8), "RO", true},
	{randNumeric(9), "RO", true},
	{randNumeric(10), "RO", true},
	{randNumeric(15), "RO", false},
	{randNumeric(1), "RO", false},
	{"A" + randNumeric(9), "RO", false},
	{randNumeric(8) + "A", "RO", false},

	{randString(1) + randNumeric(7) + randString(1), "ES", true},
	{randNumeric(8) + randString(1), "ES", true},
	{randString(1) + randNumeric(8), "ES", true},
	{randString(2) + randNumeric(6) + randString(1), "ES", false},
	{randNumeric(9), "ES", true},
	{randNumeric(15), "ES", false},
	{randNumeric(5), "ES", false},

	{randString(1) + randNumeric(8) + randString(1), "FR", false},
	{randString(2) + randNumeric(9), "FR", true},
	{randString(1) + randNumeric(10), "FR", true},
	{randString(2) + randNumeric(6) + randString(1), "FR", false},
	{randNumeric(11), "FR", true},
	{randNumeric(15), "FR", false},
	{randNumeric(5), "FR", false},

	{randString(1) + randNumeric(8) + randString(1), "CY", false},
	{randString(2) + randNumeric(9), "CY", false},
	{randNumeric(8) + randString(1), "CY", true},
	{randString(2) + randNumeric(6) + randString(1), "CY", false},
	{randNumeric(9), "CY", false},

	{randNumeric(9), "LT", true},
	{randNumeric(12), "LT", true},
	{randNumeric(10), "LT", false},
	{randNumeric(15), "LT", false},
	{randNumeric(5), "LT", false},
	{"A" + randNumeric(9), "LT", false},
	{randNumeric(8) + "A", "LT", false},

	{randNumeric(9) + "B" + randNumeric(2), "NL", true},
	{randNumeric(9) + "C" + randNumeric(2), "NL", false},
	{randNumeric(12), "NL", false},
	{randNumeric(15), "NL", false},
	{randNumeric(5), "NL", false},
	{"A" + randNumeric(9), "NL", false},
	{randNumeric(8) + "A", "NL", false},

	{randNumeric(1) + randString(1) + randNumeric(5) + randString(1), "IE", true},
	{randNumeric(2) + randString(1) + randNumeric(5) + randString(1), "IE", false},
	{randNumeric(1) + randString(2), "IE", false},
	{randNumeric(7) + "C" + randNumeric(2), "IE", false},
	{randNumeric(7) + randString(1), "IE", true},
	{randNumeric(7) + "WI", "IE", true},
	{randNumeric(12), "IE", false},
	{randNumeric(15), "IE", false},
	{randNumeric(5), "IE", false},
	{"A" + randNumeric(9), "IE", false},
	{randNumeric(8) + "A", "IE", false},
}

func TestValidFormat(t *testing.T) {
	var isValid bool
	for _, tt := range numbers {
		isValid = vatcheck.IsValidFormat(tt.countryCode + tt.vatNumber)
		if tt.isValid != isValid {
			t.Errorf(`Test failed for "%s%s"`, tt.countryCode, tt.vatNumber)
		}
	}
}

func TestInvalidCountry(t *testing.T) {
	isValid := vatcheck.IsValidFormat("US123")
	if isValid {
		t.Error("A non-EU country should not validate")
	}
}

func randNumeric(n int) (str string) {
	str, _ = randutil.String(n, randutil.Numerals)
	return
}

func randString(n int) (str string) {
	str, _ = randutil.String(n, randutil.Alphabet)
	return strings.ToUpper(str)
}
