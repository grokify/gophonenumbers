package gophonenumbers

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/grokify/gophonenumbers/numverify"
	"github.com/grokify/gophonenumbers/twilio"
)

type Source string

const (
	Ekata     Source = "ekata"
	Numverify Source = "numverify"
	Twilio    Source = "twilio"
)

type NumberLookup struct {
	NumberE164      string
	Components      Components
	Carrier         Carrier
	LookupSource    Source
	LookupTime      time.Time
	SourceNumverify *numverify.ResponseSuccess
	SourceTwilio    *twilio.NumberInfo
}

type Carrier struct {
	MobileCountryCode string `json:"mobileCountryCode,omitempty"`
	MobileNetworkCode string `json:"mobileNetworkCode,omitempty"`
	Name              string `json:"name,omitempty"`
	LineType          string `json:"lineType,omitempty"`
	ErrorCode         string `json:"errorCode,omitempty"`
}

func NewNumberLookupNumverify(src *numverify.Response) (NumberLookup, error) {
	if src.Success == nil {
		return NumberLookup{}, errors.New("numverify api response is failure")
	} else if src.StatusCode >= 300 {
		return NumberLookup{}, fmt.Errorf("numverify api response has bad status code [%d]", src.StatusCode)
	}

	lookup := NumberLookup{
		NumberE164: src.Success.InternationalFormat,
		Carrier: Carrier{
			Name:     src.Success.Carrier,
			LineType: src.Success.LineType},
		LookupSource:    Numverify,
		LookupTime:      src.Time,
		SourceNumverify: src.Success}
	return lookup, nil
}

func NewNumberLookupTwilio(src *twilio.NumberInfo) (NumberLookup, error) {
	if src == nil {
		return NumberLookup{}, errors.New("no twilio source")
	}
	lookup := NumberLookup{
		NumberE164:   src.PhoneNumber,
		Carrier:      TwilioCarrierToCommon(src.Carrier),
		LookupSource: Twilio,
		LookupTime:   src.APIResponseInfo.Time,
		SourceTwilio: src}
	return lookup, nil
}

func TwilioCarrierToCommon(c twilio.Carrier) Carrier {
	return Carrier{
		MobileCountryCode: c.MobileCountryCode,
		MobileNetworkCode: c.MobileNetworkCode,
		Name:              c.Name,
		LineType:          c.Type}
}

func SortLookupsDesc(lookups []NumberLookup) []NumberLookup {
	sort.SliceStable(lookups, func(i, j int) bool {
		return sortLookupCompareString(lookups[i]) >
			sortLookupCompareString(lookups[j])
	})
	return lookups
}

func sortLookupCompareString(lookup NumberLookup) string {
	ranks := map[Source]int{
		Numverify: 1,
		Twilio:    2,
		Ekata:     3}
	rank := 0
	if try, ok := ranks[lookup.LookupSource]; ok {
		rank = try
	}
	return fmt.Sprintf("%d.%s", rank, lookup.LookupTime.Format(time.RFC3339))
}
