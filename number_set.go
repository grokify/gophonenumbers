package gophonenumbers

import (
	"errors"
	"io/fs"

	"github.com/grokify/simplego/encoding/jsonutil"
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

func (set *NumberSet) WriteFileJSON(filename, prefix, indent string, perm fs.FileMode) error {
	return jsonutil.WriteFile(filename, set, prefix, indent, perm)
}

type MapStringString map[string]string
