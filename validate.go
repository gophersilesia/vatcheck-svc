package vatcheck

import "regexp"

var euCountries = map[string]string{
	"BE": "^0[0-9]{9}$",
	"BG": "^[0-9]{9,10}$",
	"CZ": "^[0-9]{8,10}$",
	"DK": "^[0-9]{8}$",
	"DE": "^[0-9]{9}$",
	"EE": "^[0-9]{9}$",
	"EL": "^[0-9]{9}$",
	"ES": "^[A-Z0-9][0-9]{7}[A-Z0-9]$",
	"FR": "^[0-9A-Z]{2}[0-9]{9}$",
	"HR": "^[0-9]{11}$",
	"IE": "^(([0-9][0-9A-Z][0-9]{5}[A-Z])|([0-9]{7}WI))$",
	"IT": "^[0-9]{11}$",
	"CY": "^[0-9]{8}[A-Z]$",
	"LV": "^[0-9]{11}$",
	"LT": "^([0-9]{9}|[0-9]{12})$",
	"LU": "^[0-9]{8}$",
	"HU": "^[0-9]{8}$",
	"MT": "^[0-9]{8}$",
	"NL": "^[0-9]{9}B[0-9]{2}$",
	"AT": "^U[0-9]{8}$",
	"PL": "^[0-9]{10}$",
	"PT": "^[0-9]{9}$",
	"RO": "^[0-9]{2,10}$",
	"SI": "^[0-9]{8}$",
	"SK": "^[0-9]{10}$",
	"FI": "^[0-9]{8}$",
	"SE": "^[0-9]{12}$",
	"GB": "^([0-9]{9}|[0-9]{12}|[0-9A-Z]{5})$",
}

// IsValidFormat checks if the given EU VAT number format is valid.
func IsValidFormat(id string) bool {
	// Prevent any slice out of bounds issue
	if len(id) < 4 {
		return false
	}

	reg, found := euCountries[id[:2]]
	if !found {
		return false
	}

	match, _ := regexp.MatchString(reg, id[2:])
	return match
}
