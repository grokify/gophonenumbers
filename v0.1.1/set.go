package gophonenumbers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/mogo/encoding/jsonutil"
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
		return fmt.Errorf("no phone number provide forNumbersSet.AddNumber [%v]",
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
		return errors.New("number must have NumberE164")
	}
	numberInfo, ok := set.NumbersMap[lookup.NumberE164]
	if !ok {
		numberInfo = NewNumberInfo()
	}
	numberInfo.Lookups = append(numberInfo.Lookups, lookup)
	err := numberInfo.Inflate()
	if err != nil {
		return err
	}
	set.NumbersMap[lookup.NumberE164] = numberInfo
	return nil
}

func (set *NumbersSet) Inflate() error {
	for e164, num := range set.NumbersMap {
		e164 = strings.TrimSpace(e164)
		if e164 != num.NumberE164 {
			panic("number mismatch")
		}
		err := num.Inflate()
		if err != nil {
			return err
		}
		if num.Components.CountryCode == 0 {
			num.Components = ParseE164(num.NumberE164)
		}
	}
	return nil
}

func (set *NumbersSet) AreaCodes() (*histogram.Histogram, error) {
	return NumbersSetAreaCodesNANP(set)
}

func NumbersSetAreaCodesNANP(numSet *NumbersSet) (*histogram.Histogram, error) {
	err := numSet.Inflate()
	if err != nil {
		return nil, err
	}
	fs := histogram.NewHistogram("AreaCodes")
	for _, num := range numSet.NumbersMap {
		if num.Components.NANPAreaCode == 0 {
			err := num.Inflate()
			if err != nil {
				return nil, err
			}
		}
		if num.Components.NANPAreaCode > 0 {
			fs.Add(strconv.Itoa(int(num.Components.NANPAreaCode)), 1)
		}
	}
	fs.Inflate()
	return fs, nil
}

func NumbersSetNumbersE164(numSet *NumbersSet) (*histogram.Histogram, error) {
	err := numSet.Inflate()
	if err != nil {
		return nil, err
	}
	fs := histogram.NewHistogram("NumbersE164")
	for e164, ni := range numSet.NumbersMap {
		if len(ni.NumberE164) > 0 && e164 == ni.NumberE164 {
			fs.Add(e164, 1)
		}
	}
	return fs, nil
}
