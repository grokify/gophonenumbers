package gophonenumbers

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/type/stringsutil"
)

type NumberSet struct {
	Numbers map[string]Number `json:"numbers"`
}

func NewNumberSet() NumberSet {
	return NumberSet{Numbers: map[string]Number{}}
}

func (set *NumberSet) Add(num Number) error {
	if len(num.E164Number) == 0 {
		return errors.New("no E.164 number")
	}
	set.Numbers[num.E164Number] = num
	return nil
}

func (set *NumberSet) Validate() error {
	for e164, num := range set.Numbers {
		nums := stringsutil.SliceCondenseSpace([]string{
			e164,
			num.E164Number,
			num.CarrierNumberInfo.E164Number}, true, true)
		if len(nums) == 0 {
			return errors.New("no phone number")
		} else if len(nums) > 1 {
			return fmt.Errorf("mismmatched numbers [%s]", strings.Join(nums, ","))
		}
	}
	return nil
}

func (set *NumberSet) RemoveLookups() {
	for e164, num := range set.Numbers {
		num.RemoveLookups()
		set.Numbers[e164] = num
	}
}

func (set *NumberSet) MapNumberCarrierName() map[string]string {
	mss := map[string]string{}
	for _, num := range set.Numbers {
		mss[num.E164Number] = num.CarrierNumberInfo.Name
	}
	return mss
}

func (set *NumberSet) MapNumberLineType() map[string]string {
	mss := map[string]string{}
	for _, num := range set.Numbers {
		mss[num.E164Number] = num.CarrierNumberInfo.LineType
	}
	return mss
}

func (set *NumberSet) HistogramCarrierName() map[string]int {
	stats := map[string]int{}
	for _, num := range set.Numbers {
		stats[num.CarrierNumberInfo.Name] += 1
	}
	return stats
}

func (set *NumberSet) WriteFileJSON(filename, prefix, indent string, perm fs.FileMode) error {
	return jsonutil.MarshalFile(filename, set, prefix, indent, perm)
}

type MapStringString map[string]string
