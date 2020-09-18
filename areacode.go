package gophonenumbers

import (
	"github.com/ttacon/libphonenumber"
)

// ParseAreaCode will attempt to extract an areacode from a phone
// number string.
func ParseAreaCodeString(numberToParse, defaultRegion string) (string, error) {
	num, err := libphonenumber.Parse(numberToParse, defaultRegion)
	if err != nil {
		return "", err
	}
	pn := PhoneNumber{Number: num}
	return pn.GetAreaCodeString(), nil
}
