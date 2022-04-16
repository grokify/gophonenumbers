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
	nmr, err := numverify.ReadFilesMultiResults(dir, rxPattern)
	if err != nil {
		return numsSet, err
	}
	for _, nRes := range nmr.Responses {
		if nRes == nil ||
			nRes.Success == nil ||
			nRes.StatusCode >= 300 ||
			len(strings.TrimSpace(nRes.Success.InternationalFormat)) == 0 {
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
	for _, tRes := range tmr.Responses {
		if tRes == nil ||
			tRes.APIResponseInfo.StatusCode >= 300 ||
			len(strings.TrimSpace(tRes.PhoneNumber)) == 0 {
			continue
		}
		lookup, err := NewNumberLookupTwilio(tRes)
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
