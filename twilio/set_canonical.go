package twilio

import (
	"github.com/grokify/gophonenumbers"
)

func (mr *MultiResults) Canonical() (gophonenumbers.NumberSet, error) {
	set := gophonenumbers.NewNumberSet()
	for _, numInfo := range mr.Responses {
		numInfoCan, err := numInfo.Canonical()
		if err != nil {
			return set, err
		}
		set.Add(numInfoCan)
	}
	return set, nil
}
