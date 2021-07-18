package gophonenumbers

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

type Lookup struct {
	CarrierNumberInfo CarrierNumberInfo      `json:"carrierNumberInfo"`
	LookupSource      string                 `json:"lookupSource"`
	LookupTime        time.Time              `json:"lookupTime"`
	LookupResponse    map[string]interface{} `json:"lookupResponse"`
}

type LookupSet struct {
	LookupMap map[string]Lookup `json:"lookupMap"`
}

func NewLookupSet() LookupSet {
	return LookupSet{
		LookupMap: map[string]Lookup{}}
}

func (set LookupSet) Add(lookup Lookup) {
	if set.LookupMap == nil {
		set.LookupMap = map[string]Lookup{}
	}
	set.LookupMap[lookup.LookupTime.Format(time.RFC3339)] = lookup
}

func (set LookupSet) Latest(source string) (Lookup, error) {
	if set.LookupMap == nil || len(set.LookupMap) == 0 {
		return Lookup{}, errors.New("no latest lookup. lookupSet is empty.")
	}
	times := []string{}
	for dtKey, lookup := range set.LookupMap {
		if len(source) > 0 && lookup.LookupSource != source {
			continue
		}
		times = append(times, dtKey)
	}
	if len(times) == 0 {
		return Lookup{}, fmt.Errorf("no lookup for source [%s]", string(source))
	}
	if len(times) == 1 {
		return set.LookupMap[times[0]], nil
	}
	sort.Strings(times)
	return set.LookupMap[times[len(times)-1]], nil
}

func (set LookupSet) Validate() {
	if set.LookupMap == nil || len(set.LookupMap) == 0 {
		return
	}
	for dtKey, lookup := range set.LookupMap {
		_, err := time.Parse(time.RFC3339, dtKey)
		if err != nil {
			delete(set.LookupMap, dtKey)
			set.LookupMap[lookup.LookupTime.Format(time.RFC3339)] = lookup
		}
	}
}
