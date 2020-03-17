package numberinfo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var rxNANPFormat = regexp.MustCompile(`^+1(([0-9]{3})([0-9]{3})([0-9]{4}))$`)

type Components struct {
	E164             string
	E164Int          int
	CountryCode      int
	NANPAreaCode     int // NPA - Numbering plan area code
	NANPExchangeCode int // NXX - Central office (exchange) code
	NANPLineNumber   int // xxxx - Line number or subscriber number
}

func ParseE164(e164 string) Components {
	e164 = strings.TrimSpace(e164)
	comp := Components{E164: e164}
	m := rxNANPFormat.FindAllStringSubmatch(e164, -1)
	if m != nil && len(m) > 0 {
		comp.CountryCode = 1
		e164int, err := strconv.Atoi(m[0][1])
		if err != nil {
			panic(fmt.Sprintf("ParseE164 [%v]", err.Error()))
		}
		comp.E164Int = e164int
		areaCode, err := strconv.Atoi(m[0][2])
		if err != nil {
			panic(fmt.Sprintf("ParseE164 [%v]", err.Error()))
		}
		comp.NANPAreaCode = areaCode
		exchangeCode, err := strconv.Atoi(m[0][3])
		if err != nil {
			panic(fmt.Sprintf("ParseE164 [%v]", err.Error()))
		}
		comp.NANPExchangeCode = exchangeCode
		lineNumber, err := strconv.Atoi(m[0][4])
		if err != nil {
			panic(fmt.Sprintf("ParseE164 [%v]", err.Error()))
		}
		comp.NANPExchangeCode = lineNumber
	}
	return comp
}
