package gophonenumbers

import (
	"errors"
	"strings"

	"github.com/grokify/gocharts/data/frequency"
)

type NumbersSet struct {
	NumbersMap map[string]NumberInfo
}

func NewNumbersSet() NumbersSet {
	return NumbersSet{NumbersMap: map[string]NumberInfo{}}
}

func (set *NumbersSet) AddLookup(lookup NumberInfoLookup) error {
	lookup.NumberE164 = strings.TrimSpace(lookup.NumberE164)
	if len(lookup.NumberE164) == 0 {
		return errors.New("E_NO_NUMBER")
	}
	if numberInfo, ok := set.NumbersMap[lookup.NumberE164]; ok {
		numberInfo.Lookups = append(numberInfo.Lookups, lookup)
		numberInfo.Inflate()
		set.NumbersMap[lookup.NumberE164] = numberInfo
	} else {
		numberInfo := NewNumberInfo()
		numberInfo.Lookups = append(numberInfo.Lookups, lookup)
		numberInfo.Inflate()
		set.NumbersMap[lookup.NumberE164] = numberInfo
	}
	return nil
}

func (set *NumbersSet) Inflate() {
	for e164, ni := range set.NumbersMap {
		e164 = strings.TrimSpace(e164)
		ni.NumberE164 = strings.TrimSpace(ni.NumberE164)
		if len(ni.NumberE164) == 0 {
			ni.NumberE164 = e164
		}
		if ni.Components.CountryCode == 0 {
			ni.Components = ParseE164(ni.NumberE164)
		}
	}
}

func NumbersSetAreaCodesNANP(numSet *NumbersSet) frequency.FrequencyStats {
	numSet.Inflate()
	fs := frequency.NewFrequencyStats("AreaCodes")
	for _, ni := range numSet.NumbersMap {
		if ni.Components.NANPAreaCode > 0 {
			fs.AddInt(ni.Components.NANPAreaCode)
		}
	}
	return fs
}

func NumbersSetNumbersE164(numSet *NumbersSet) frequency.FrequencyStats {
	numSet.Inflate()
	fs := frequency.NewFrequencyStats("NumbersE164")
	for e164, ni := range numSet.NumbersMap {
		if len(ni.NumberE164) > 0 && e164 == ni.NumberE164 {
			fs.AddString(e164)
		}
	}
	return fs
}
