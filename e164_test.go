package gophonenumbers

import (
	"testing"

	"github.com/nyaruka/phonenumbers"
)

var e164FormatTests = []struct {
	v            string
	wantNational string
}{
	{"+16505550100", "(650) 555-0100"},
}

func TestE164Format(t *testing.T) {
	for _, tt := range e164FormatTests {
		nat, err := E164Format(tt.v, "", phonenumbers.NATIONAL)
		if err != nil {
			t.Errorf("phonenumber.E164Format(\"%s\") Error: (%s)", tt.v, err.Error())
		}
		if nat != tt.wantNational {
			t.Errorf("phonenumber.E164Format(\"%v\") Mismatch: want (%s), got (%vs",
				tt.v, tt.wantNational, nat)
		}
		_, err = phonenumbers.Parse(tt.v, "")
		if err != nil {
			t.Errorf("phonenumbers.Parse(\"%s\") Error: (%s)", tt.v, err.Error())
		}
	}
}
