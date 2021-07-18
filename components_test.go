package gophonenumbers

import (
	"testing"
)

var parseE164Tests = []struct {
	e164Number   string
	e164Uint     uint
	countryCode  uint
	areaCode     uint
	exchangeCode uint
	lineNumber   uint
}{
	{"+16505550100", uint(16505550100), uint(1), uint(650), uint(555), uint(100)},
	{"+14155550199", uint(14155550199), uint(1), uint(415), uint(555), uint(199)},
}

func TestParseE164(t *testing.T) {
	for _, tt := range parseE164Tests {
		num := Number{E164Number: tt.e164Number}
		got, err := num.NANPComponents()
		if err != nil {
			t.Errorf("Number.NANPComponents(\"%s\") Want: [%v] Got: [%v]",
				tt.e164Number, tt.e164Uint, got.E164NumberUint)
		}
		if got.E164NumberUint != tt.e164Uint {
			t.Errorf("Number.NANPComponents() Want: [%v] Got: [%v]",
				tt.e164Uint, got.E164NumberUint)
		}
		if got.CountryCode != tt.countryCode {
			t.Errorf("Number.NANPComponents() Want: [%v] Got: [%v]",
				tt.countryCode, got.CountryCode)
		}
		if got.NANPAreaCode != tt.areaCode {
			t.Errorf("Number.NANPComponents(\"%s\") Want: [%v] Got: [%v]",
				tt.e164Number, tt.areaCode, got.NANPAreaCode)
		}
		if got.NANPExchangeCode != tt.exchangeCode {
			t.Errorf("Number.NANPComponents(\"%s\") Want: [%v] Got: [%v]",
				tt.e164Number, tt.exchangeCode, got.NANPExchangeCode)
		}
		if got.NANPLineNumber != tt.lineNumber {
			t.Errorf("Number.NANPComponents(\"%s\") Want: [%v] Got: [%v]",
				tt.e164Number, tt.lineNumber, got.NANPLineNumber)
		}
	}
}
