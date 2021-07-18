package twilio

import (
	"encoding/json"

	"github.com/grokify/gophonenumbers"
	"github.com/grokify/simplego/net/httputilmore"
)

type NumberInfo struct {
	CallerName      map[string]string         `json:"caller_name,omitempty"`
	CountryCode     string                    `json:"country_code,omitempty"`
	PhoneNumber     string                    `json:"phone_number,omitempty"`
	NationalFormat  string                    `json:"national_format,omitempty"`
	URL             string                    `json:"url,omitempty"`
	Carrier         Carrier                   `json:"carrier,omitempty"`
	ApiResponseInfo httputilmore.ResponseInfo `json:"api_response_info,omitempty"`
}

func (num *NumberInfo) Canonical() (gophonenumbers.Number, error) {
	canNum := gophonenumbers.NewNumber()
	canNum.E164Number = num.PhoneNumber
	canNum.CountryCode = num.CountryCode
	canNum.CarrierNumberInfo = num.Carrier.Canonical()
	canNum.CarrierNumberInfo.E164Number = canNum.E164Number
	lookup := gophonenumbers.Lookup{
		CarrierNumberInfo: canNum.CarrierNumberInfo,
		LookupSource:      gophonenumbers.SourceTwilio,
		LookupTime:        num.ApiResponseInfo.Time.UTC()}
	if len(num.ApiResponseInfo.Body) > 0 {
		msi := map[string]interface{}{}
		err := json.Unmarshal([]byte(num.ApiResponseInfo.Body), &msi)
		if err != nil {
			return canNum, err
		}
		lookup.LookupResponse = msi
	}
	canNum.Lookups.Add(lookup)
	return canNum, nil
}

type Carrier struct {
	MobileCountryCode string `json:"mobile_country_code,omitempty"`
	MobileNetworkCode string `json:"mobile_network_code,omitempty"`
	Name              string `json:"name,omitempty"`
	Type              string `json:"type,omitempty"`
	ErrorCode         string `json:"error_code,omitempty"`
}

func (car *Carrier) Canonical() gophonenumbers.CarrierNumberInfo {
	return gophonenumbers.CarrierNumberInfo{
		E164Number:        "",
		MobileCountryCode: car.MobileCountryCode,
		MobileNetworkCode: car.MobileNetworkCode,
		Name:              car.Name,
		LineType:          car.Type,
		ErrorCode:         car.ErrorCode}
}
