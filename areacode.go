package gophonenumbers

import "github.com/ttacon/libphonenumber"

func ParseAreaCode(numberToParse, defaultRegion string) (string, error) {
	num, err := libphonenumber.Parse(numberToParse, defaultRegion)
	if err != nil {
		return "", err
	}
	return GetAreaCode(num), nil
}

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
