package gophonenumbers

import (
	"strconv"

	"github.com/grokify/mogo/errors/errorsutil"
)

type Direction int

const (
	To Direction = iota
	From
)

type ToFromNumbersSets struct {
	All          ToFromNumbersSet
	AreaCodesMap map[string]ToFromNumbersSet
}

type ToFromNumbersSetsStats struct {
	ToAreaCodeCount   uint
	ToNumberCount     uint
	FromAreaCodeCount uint
	FromNumberCount   uint
}

func NewToFromNumbersSets() ToFromNumbersSets {
	return ToFromNumbersSets{
		AreaCodesMap: map[string]ToFromNumbersSet{}}
}

type ToFromNumbersSet struct {
	To   NumbersSet
	From NumbersSet
}

func NewToFromNumbersSet() ToFromNumbersSet {
	return ToFromNumbersSet{
		To:   NewNumbersSet(),
		From: NewNumbersSet()}
}

func (tfns *ToFromNumbersSet) AddNumber(ni NumberInfo, direction Direction) error {
	if direction == To {
		return tfns.To.AddNumber(ni)
	}
	return tfns.From.AddNumber(ni)
}

func (tfSets *ToFromNumbersSets) AddNumber(num NumberInfo, direction Direction, addAreaCode bool) error {
	if direction == To {
		err := tfSets.All.To.AddNumber(num)
		if err != nil {
			return errorsutil.Wrap(err, "AreaCodeNumbersSets.AddNumber")
		}
		if addAreaCode {
			num.InflateComponents()
			areaCode := strconv.Itoa(int(num.Components.NANPAreaCode))
			tfSetAC, ok := tfSets.AreaCodesMap[areaCode]
			if !ok {
				tfSetAC = NewToFromNumbersSet()
			}
			tfSetAC.AddNumber(num, direction)
			tfSets.AreaCodesMap[areaCode] = tfSetAC
		}
		return nil
	}
	err := tfSets.All.From.AddNumber(num)
	if err != nil {
		return errorsutil.Wrap(err, "AreaCodeNumbersSets.AddNumber")
	}
	if addAreaCode {
		num.InflateComponents()
		areaCode := strconv.Itoa(int(num.Components.NANPAreaCode))
		tfSetAC, ok := tfSets.AreaCodesMap[areaCode]
		if !ok {
			tfSetAC = NewToFromNumbersSet()
		}
		tfSetAC.AddNumber(num, direction)
		tfSets.AreaCodesMap[areaCode] = tfSetAC
	}
	return nil
}

func (tfSets *ToFromNumbersSets) Stats() ToFromNumbersSetsStats {
	toAreaCodes := tfSets.All.To.AreaCodes()
	fromAreaCodes := tfSets.All.From.AreaCodes()
	return ToFromNumbersSetsStats{
		ToAreaCodeCount:   uint(len(toAreaCodes.Bins)),
		FromAreaCodeCount: uint(len(fromAreaCodes.Bins)),
		ToNumberCount:     uint(len(tfSets.All.To.NumbersMap)),
		FromNumberCount:   uint(len(tfSets.All.From.NumbersMap))}
}
