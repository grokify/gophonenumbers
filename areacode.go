package gophonenumbers

import (
	"github.com/nyaruka/phonenumbers"
)

// ParseAreaCode will attempt to extract an areacode from a phone
// number string.
func ParseAreaCodeString(numberToParse, defaultRegion string) (string, error) {
	num, err := phonenumbers.Parse(numberToParse, defaultRegion)
	if err != nil {
		return "", err
	}
	pn := PhoneNumber{Number: num}
	return pn.GetAreaCodeString(), nil
}
