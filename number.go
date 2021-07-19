package gophonenumbers

import (
	"fmt"
	"strings"
)

type Number struct {
	E164Number        string            `json:"e164Number"`
	CountryCode       string            `json:"countryCode"`
	CarrierNumberInfo CarrierNumberInfo `json:"carrier"`
	Lookups           LookupSet         `json:"lookups"`
}

func NewNumber() Number {
	return Number{
		Lookups: NewLookupSet()}
}

func (num *Number) SetLatest(source string) error {
	latest, err := num.Lookups.Latest(source)
	if err != nil {
		return err
	}
	latest.CarrierNumberInfo.E164Number = strings.TrimSpace(latest.CarrierNumberInfo.E164Number)
	num.E164Number = strings.TrimSpace(num.E164Number)
	if num.E164Number != latest.CarrierNumberInfo.E164Number {
		return fmt.Errorf("lookup number mismatch: number [%s] lookup [%s]",
			num.E164Number,
			latest.CarrierNumberInfo.E164Number)
	}
	num.CarrierNumberInfo = latest.CarrierNumberInfo
	return nil
}

func (num *Number) RemoveLookups() {
	num.Lookups = NewLookupSet()
}
