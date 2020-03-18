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
		got := ParseE164(tt.e164Number)
		if got.E164Uint != tt.e164Uint {
			t.Errorf("gophonenumbers.ParseE164(\"%s\") Want: [%v] Got: [%v]",
				tt.e164Number, tt.e164Uint, got.E164Uint)
		}
		if got.CountryCode != tt.countryCode {
			t.Errorf("gophonenumbers.ParseE164(\"%s\") Want: [%v] Got: [%v]",
				tt.e164Number, tt.countryCode, got.CountryCode)
		}
		if got.NANPAreaCode != tt.areaCode {
			t.Errorf("gophonenumbers.ParseE164(\"%s\") Want: [%v] Got: [%v]",
				tt.e164Number, tt.areaCode, got.NANPAreaCode)
		}
		if got.NANPExchangeCode != tt.exchangeCode {
			t.Errorf("gophonenumbers.ParseE164(\"%s\") Want: [%v] Got: [%v]",
				tt.e164Number, tt.exchangeCode, got.NANPExchangeCode)
		}
		if got.NANPLineNumber != tt.lineNumber {
			t.Errorf("gophonenumbers.ParseE164(\"%s\") Want: [%v] Got: [%v]",
				tt.e164Number, tt.lineNumber, got.NANPLineNumber)
		}
	}
}
