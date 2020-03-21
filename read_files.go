package gophonenumbers

import (
	"strings"

	"github.com/grokify/gophonenumbers/numverify"
	"github.com/grokify/gophonenumbers/twilio"
)

func ReadLookupsNumverify(dir, rxPattern string, numsSet *NumbersSet) (*NumbersSet, error) {
	if numsSet == nil {
		newNumsSet := NewNumbersSet()
		numsSet = &newNumsSet
	}
	nmr, err := numverify.NewMultiResultsFiles(dir, rxPattern)
	if err != nil {
		return numsSet, err
	}
	for _, nRes := range nmr.Responses {
		if nRes == nil || nRes.Success == nil || nRes.StatusCode >= 300 {
			continue
		}
		e164 := strings.TrimSpace(nRes.Success.InternationalFormat)
		if len(e164) == 0 {
			continue
		}
		lookup, err := NewNumberLookupNumverify(nRes)
		if err != nil {
			return numsSet, err
		}
		err = numsSet.AddLookup(lookup)
		if err != nil {
			return numsSet, err
		}
	}
	return numsSet, nil
}

func ReadLookupsTwilio(dir, rxPattern string, numsSet *NumbersSet) (*NumbersSet, error) {
	if numsSet == nil {
		newNumsSet := NewNumbersSet()
		numsSet = &newNumsSet
	}
	tmr, err := twilio.NewMultiResultsFiles(dir, rxPattern)
	if err != nil {
		return numsSet, err
	}
	for _, tnum := range tmr.Responses {
		if tnum == nil {
			continue
		}
		if tnum.ApiResponseInfo.StatusCode >= 300 {
			continue
		}
		e164 := strings.TrimSpace(tnum.PhoneNumber)
		if len(e164) == 0 {
			continue
		}
		lookup, err := NewNumberLookupTwilio(tnum)
		if err != nil {
			return numsSet, err
		}
		err = numsSet.AddLookup(lookup)
		if err != nil {
			return numsSet, err
		}
	}
	return numsSet, nil
}
