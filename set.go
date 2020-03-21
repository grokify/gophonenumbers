package gophonenumbers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grokify/gocharts/data/frequency"
	"github.com/grokify/gotilla/encoding/jsonutil"
)

type NumbersSet struct {
	NumbersMap map[string]NumberInfo
}

func NewNumbersSet() NumbersSet {
	return NumbersSet{NumbersMap: map[string]NumberInfo{}}
}

func (set *NumbersSet) AddNumber(num NumberInfo) error {
	num.NumberE164 = strings.TrimSpace(num.NumberE164)
	if len(num.NumberE164) == 0 {
		return fmt.Errorf("E_NO_PHONE_NUMBER NumbersSet.AddNumber [%v]",
			jsonutil.MustMarshal(num, true))
	}
	if set.NumbersMap == nil {
		set.NumbersMap = map[string]NumberInfo{}
	}
	set.NumbersMap[num.NumberE164] = num
	return nil
}

func (set *NumbersSet) AddLookup(lookup NumberLookup) error {
	lookup.NumberE164 = strings.TrimSpace(lookup.NumberE164)
	if len(lookup.NumberE164) == 0 {
		return errors.New("E_NO_NUMBER")
	}
	numberInfo, ok := set.NumbersMap[lookup.NumberE164]
	if !ok {
		numberInfo = NewNumberInfo()
	}
	numberInfo.Lookups = append(numberInfo.Lookups, lookup)
	numberInfo.Inflate()
	set.NumbersMap[lookup.NumberE164] = numberInfo
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

func (set *NumbersSet) AreaCodes() frequency.FrequencyStats {
	return NumbersSetAreaCodesNANP(set)
}

func NumbersSetAreaCodesNANP(numSet *NumbersSet) frequency.FrequencyStats {
	numSet.Inflate()
	fs := frequency.NewFrequencyStats("AreaCodes")
	for _, num := range numSet.NumbersMap {
		if num.Components.NANPAreaCode == 0 {
			num.Inflate()
		}
		if num.Components.NANPAreaCode > 0 {
			fs.AddInt(int(num.Components.NANPAreaCode))
		}
	}
	fs.Inflate()
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
