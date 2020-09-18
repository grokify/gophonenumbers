package gophonenumbers

import (
	"strings"

	"github.com/ttacon/libphonenumber"
)

// DefaultRegion is the region to be used when none is provided.
const DefaultRegion = "US"

// ParseAreaCode will attempt to extract an areacode from a phone
// number string.
func ParseAreaCode(numberToParse, defaultRegion string) (string, error) {
	defaultRegion = strings.TrimSpace(defaultRegion)
	if len(defaultRegion) == 0 {
		defaultRegion = DefaultRegion
	}
	num, err := libphonenumber.Parse(numberToParse, defaultRegion)
	if err != nil {
		return "", err
	}
	return GetAreaCode(num), nil
}

// GetAreaCode will return an area code given a `libphonenumber.PhoneNumber`.
func GetAreaCode(num *libphonenumber.PhoneNumber) string {
	// Get the cleaned number and the length of the area code.
	natSigNumber := libphonenumber.GetNationalSignificantNumber(num)
	geoCodeLength := libphonenumber.GetLengthOfGeographicalAreaCode(num)
	// Extract the area code.
	areaCode := ""
	if geoCodeLength > 0 {
		areaCode = natSigNumber[0:geoCodeLength]
	}
	return areaCode
}
