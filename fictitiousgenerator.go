package gophonenumbers

import (
	"github.com/grokify/mogo/crypto/randutil"
)

const (
	fakeLineNumberMin = uint16(100)
	fakeLineNumberMax = uint16(199)
)

type FakeNumberGenerator struct {
	AreaCodes []uint16
}

func NewFakeNumberGenerator(areacodes []uint16) FakeNumberGenerator {
	fng := FakeNumberGenerator{
		AreaCodes: areacodes,
	}
	return fng
}

// RandomAreaCode generates a random area code.
func (fng *FakeNumberGenerator) RandomAreaCode() uint16 {
	num := randutil.Int64n(uint(len(fng.AreaCodes)))
	return fng.AreaCodes[num]
}

// RandomLineNumber generates a random line number
func (fng *FakeNumberGenerator) RandomLineNumber() uint16 {
	return fng.RandomLineNumberMinMax(fakeLineNumberMin, fakeLineNumberMax)
}

// RandomLineNumber generates a random line number
func (fng *FakeNumberGenerator) RandomLineNumberMinMax(min, max uint16) uint16 {
	num := randutil.Int64n(uint(max) - uint(min))
	return uint16(num) + min
}

// RandomLocalNumberUS returns a US E.164 number
// AreaCode + Prefix + Line Number
func (fng *FakeNumberGenerator) RandomLocalNumberUS() (uint64, error) {
	num := fng.RandomLineNumber()
	num2 := fng.RandomAreaCode()
	return fng.LocalNumberUS(num2, num), nil
}

// RandomLocalNumberUS returns a US E.164 number
// AreaCode + Prefix + Line Number
func (fng *FakeNumberGenerator) RandomLocalNumberUSAreaCodes(acs []uint16) uint64 {
	idx := randutil.Int64n(uint(len(acs)))
	num := fng.RandomLineNumber()
	return fng.LocalNumberUS(acs[idx], num)
}

// LocalNumberUS returns a US E.164 number given an areacode and line number
func (fng *FakeNumberGenerator) LocalNumberUS(ac uint16, ln uint16) uint64 {
	return 10000000000 + (uint64(ac) * 10000000) + (5550000) + uint64(ln)
}

// RandomLocalNumberUSUnique returns a US E.164 number
// AreaCode + Prefix + Line Number
func (fng *FakeNumberGenerator) RandomLocalNumberUSUnique(set map[uint64]int8) (uint64, map[uint64]int8, error) {
	try, err := fng.RandomLocalNumberUS()
	if err != nil {
		return 0, set, err
	}
	_, ok := set[try]
	for ok {
		try, err := fng.RandomLocalNumberUS()
		if err != nil {
			return 0, set, err
		}
		_, ok = set[try]
	}
	set[try] = 1
	return try, set, nil
}

// RandomLocalNumberUSUnique returns a US E.164 number
// AreaCode + Prefix + Line Number
func (fng *FakeNumberGenerator) RandomLocalNumberUSUniqueAreaCodeSet(set map[uint64]int8, acs []uint16) (uint64, map[uint64]int8, error) {
	try := fng.RandomLocalNumberUSAreaCodes(acs)
	_, ok := set[try]
	for ok {
		try = fng.RandomLocalNumberUSAreaCodes(acs)
	}
	set[try] = 1
	return try, set, nil
}

// LocalNumberUS returns a US E.164 number given an areacode and line number
func LocalNumberUS(ac uint16, ln uint16) uint64 {
	return 10000000000 + (uint64(ac) * 10000000) + (5550000) + uint64(ln)
}

type AreaCodeIncrementor struct {
	Counter map[uint16]uint16
	Base    uint16
}

func NewAreaCodeIncrementor(base uint16) AreaCodeIncrementor {
	return AreaCodeIncrementor{Counter: map[uint16]uint16{}, Base: base}
}

func (aci *AreaCodeIncrementor) GetNext(ac uint16) uint64 {
	if count, ok := aci.Counter[ac]; ok {
		count++
		aci.Counter[ac] = count
		return LocalNumberUS(ac, count)
	}
	count := aci.Base + 1
	aci.Counter[ac] = count
	return LocalNumberUS(ac, count)
}
