package gophonenumbers

import (
	"strings"

	"github.com/nyaruka/phonenumbers"
)

func E164Format(numberToParse, defaultRegion string, numberFormat phonenumbers.PhoneNumberFormat) (string, error) {
	defaultRegion = strings.ToUpper(strings.TrimSpace(defaultRegion))
	// empty string region is used when number is already in E.164 format.
	phone, err := phonenumbers.Parse(numberToParse, defaultRegion)
	if err != nil {
		return "", err
	}
	return phonenumbers.Format(phone, numberFormat), nil
}

func MustE164Format(numberToParse, defaultRegion string, numberFormat phonenumbers.PhoneNumberFormat) string {
	pn, err := E164Format(numberToParse, defaultRegion, numberFormat)
	if err != nil {
		panic(err)
	}
	return pn
}

func FormatsParse(number, region string) (Formats, error) {
	pn, err := phonenumbers.Parse(number, region)
	if err != nil {
		return Formats{}, err
	}
	return FormatsFromPhoneNumber(pn), nil
}

type Formats struct {
	E164          string
	International string
	National      string
	RFC3966       string
}

func FormatsFromPhoneNumber(pn *phonenumbers.PhoneNumber) Formats {
	return Formats{
		E164:          phonenumbers.Format(pn, phonenumbers.E164),
		International: phonenumbers.Format(pn, phonenumbers.INTERNATIONAL),
		National:      phonenumbers.Format(pn, phonenumbers.NATIONAL),
		RFC3966:       phonenumbers.Format(pn, phonenumbers.RFC3966),
	}
}
