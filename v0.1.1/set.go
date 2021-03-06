package gophonenumbers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/data/histogram"
	"github.com/grokify/simplego/encoding/jsonutil"
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
	for e164, num := range set.NumbersMap {
		e164 = strings.TrimSpace(e164)
		if e164 != num.NumberE164 {
			panic("E_MISMATCH")
		}
		num.Inflate()
		if num.Components.CountryCode == 0 {
			num.Components = ParseE164(num.NumberE164)
		}
	}
}

func (set *NumbersSet) AreaCodes() *histogram.Histogram {
	return NumbersSetAreaCodesNANP(set)
}

func NumbersSetAreaCodesNANP(numSet *NumbersSet) *histogram.Histogram {
	numSet.Inflate()
	fs := histogram.NewHistogram("AreaCodes")
	for _, num := range numSet.NumbersMap {
		if num.Components.NANPAreaCode == 0 {
			num.Inflate()
		}
		if num.Components.NANPAreaCode > 0 {
			fs.Add(strconv.Itoa(int(num.Components.NANPAreaCode)), 1)
		}
	}
	fs.Inflate()
	return fs
}

func NumbersSetNumbersE164(numSet *NumbersSet) *histogram.Histogram {
	numSet.Inflate()
	fs := histogram.NewHistogram("NumbersE164")
	for e164, ni := range numSet.NumbersMap {
		if len(ni.NumberE164) > 0 && e164 == ni.NumberE164 {
			fs.Add(e164, 1)
		}
	}
	return fs
}
