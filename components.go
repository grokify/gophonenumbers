package gophonenumbers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var rxNANPFormat = regexp.MustCompile(`^\+(1([0-9]{3})([0-9]{3})([0-9]{4}))$`)

type Components struct {
	E164Number       string
	E164NumberUint   uint
	RegionCode       string
	CountryCode      uint
	NANPAreaCode     uint // NPA - Numbering plan area code
	NANPExchangeCode uint // NXX - Central office (exchange) code
	NANPLineNumber   uint // xxxx - Line number or subscriber number
}

func (num *Number) NANPComponents() (Components, error) {
	num.E164Number = strings.TrimSpace(num.E164Number)
	comp := Components{E164Number: num.E164Number}
	m := rxNANPFormat.FindAllStringSubmatch(num.E164Number, -1)
	if len(m) == 0 {
		return comp, fmt.Errorf("number is not E.164 [%s]", num.E164Number)
	}
	comp.CountryCode = 1
	e164int, err := strconv.Atoi(m[0][1])
	if err != nil {
		panic(fmt.Sprintf("ParseE164 [%v]", err.Error()))
	}
	comp.E164NumberUint = uint(e164int)
	areaCode, err := strconv.Atoi(m[0][2])
	if err != nil {
		panic(fmt.Sprintf("ParseE164 [%v]", err.Error()))
	}
	comp.NANPAreaCode = uint(areaCode)
	exchangeCode, err := strconv.Atoi(m[0][3])
	if err != nil {
		panic(fmt.Sprintf("ParseE164 [%v]", err.Error()))
	}
	comp.NANPExchangeCode = uint(exchangeCode)
	lineNumber, err := strconv.Atoi(m[0][4])
	if err != nil {
		panic(fmt.Sprintf("ParseE164 [%v]", err.Error()))
	}
	comp.NANPLineNumber = uint(lineNumber)
	return comp, nil
}
