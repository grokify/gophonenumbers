package gophonenumbers

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	CarrierATT     = "att.com"
	CarrierSprint  = "sprint.com"
	CarrierTMobile = "t-mobile.com"
	CarrierVerizon = "verizon.com"
)

type LineType int

const (
	LineLocal LineType = iota
	LineMobile
	LineTollFree
)

type NumberInfo struct {
	NumberE164       string
	Components       Components
	Carrier          Carrier
	Lookups          []NumberLookup
	CarrierNamesEach []string
	LineTypesEach    []string
}

func NewNumberInfo() NumberInfo {
	return NumberInfo{
		Lookups:          []NumberLookup{},
		CarrierNamesEach: []string{},
		LineTypesEach:    []string{}}
}

func (ni *NumberInfo) Inflate() error {
	ni.NumberE164 = strings.TrimSpace(ni.NumberE164)

	if len(ni.Lookups) > 0 {
		lookups := ni.Lookups
		if 1 == 0 {
			sort.Slice(
				lookups,
				func(i, j int) bool {
					// Ascending
					// return lookups[i].LookupTime.Before(lookups[j].LookupTime)
					// Descending
					return lookups[i].LookupTime.After(lookups[j].LookupTime)
				},
			)
		}
		sort.SliceStable(lookups, func(i, j int) bool {
			return sortLookupCompareString(lookups[i]) >
				sortLookupCompareString(lookups[j])
		})
		ni.Lookups = lookups
	}
	carrierNames := []string{}
	carrierNamesSources := map[string]int{}
	lineTypes := []string{}
	lineTypesSources := map[string]int{}
	for _, lookup := range ni.Lookups {
		lookup.NumberE164 = strings.TrimSpace(lookup.NumberE164)
		lookup.Carrier.Name = strings.TrimSpace(lookup.Carrier.Name)
		lookup.Carrier.LineType = strings.TrimSpace(lookup.Carrier.LineType)
		if len(lookup.NumberE164) == 0 {
			return errors.New("E_LOOKUP_NO_E164")
		}
		// Set Top Level E164 number.s
		if len(ni.NumberE164) == 0 {
			ni.NumberE164 = lookup.NumberE164
		} else if ni.NumberE164 != lookup.NumberE164 {
			// E.164 should be the same for all lookups.
			return fmt.Errorf("E_LOOKUP_E164_MISMATCH TOP[%v] LOOKUP[%v]",
				ni.NumberE164, lookup.NumberE164)
		}
		// Set Top Level Carrier info.
		ni.Carrier.Name = strings.TrimSpace(ni.Carrier.Name)
		ni.Carrier.LineType = strings.TrimSpace(ni.Carrier.LineType)
		if len(ni.Carrier.Name) == 0 {
			ni.Carrier.Name = lookup.Carrier.Name
		}
		if len(ni.Carrier.LineType) == 0 {
			ni.Carrier.LineType = lookup.Carrier.LineType
		}
		source := strings.TrimSpace(string(lookup.LookupSource))
		if _, ok := carrierNamesSources[source]; !ok {
			carrierNames = append(carrierNames,
				fmt.Sprintf("%s(%s)", ni.Carrier.Name, source))
			carrierNamesSources[source] = 1
		}
		if _, ok := lineTypesSources[source]; !ok {
			lineTypes = append(lineTypes,
				fmt.Sprintf("%s(%s)", ni.Carrier.LineType, source))
			lineTypesSources[source] = 1
		}
	}
	ni.CarrierNamesEach = carrierNames
	ni.LineTypesEach = lineTypes
	ni.Components = ParseE164(ni.NumberE164)
	return nil
}

func (ni *NumberInfo) InflateComponents() {
	ni.Components = ParseE164(ni.NumberE164)
}
