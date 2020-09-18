package gophonenumbers

import (
	"strconv"
	"strings"

	"github.com/ttacon/libphonenumber"
)

type PhoneNumber struct {
	Number *libphonenumber.PhoneNumber
}

func (num *PhoneNumber) GetRegionCode() string {
	return libphonenumber.GetRegionCodeForNumber(num.Number)
}

func (num *PhoneNumber) GetCountryCode() int32 {
	return num.Number.GetCountryCode()
}

func (num *PhoneNumber) GetNumberE164() string {
	return libphonenumber.Format(num.Number, libphonenumber.E164)
}

func (num *PhoneNumber) GetNumberE164Uint() uint {
	e164 := num.GetNumberE164()
	if len(e164) == 0 {
		return 0
	}
	e164 = strings.TrimLeft(e164, "+")
	e164int, err := strconv.Atoi(e164)
	if err != nil {
		panic(err)
	}
	return uint(e164int)
}

func (num *PhoneNumber) GetAreaCode() int {
	acint, err := strconv.Atoi(num.GetAreaCodeString())
	if err != nil {
		return -1
	}
	return acint
}

func (num *PhoneNumber) GetAreaCodeString() string {
	// Get the cleaned number and the length of the area code.
	natSigNumber := libphonenumber.GetNationalSignificantNumber(num.Number)
	geoCodeLength := libphonenumber.GetLengthOfGeographicalAreaCode(num.Number)
	// Extract the area code.
	areaCode := ""
	if geoCodeLength > 0 {
		areaCode = natSigNumber[0:geoCodeLength]
	}
	return areaCode
}

func (num *PhoneNumber) Meta() Components {
	comp := Components{
		E164:         num.GetNumberE164(),
		E164Uint:     num.GetNumberE164Uint(),
		RegionCode:   num.GetRegionCode(),
		CountryCode:  uint(num.GetCountryCode()),
		NANPAreaCode: uint(num.GetAreaCode())}
	return comp
}
