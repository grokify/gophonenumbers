package gophonenumbers

import "testing"

var numberTests = []struct {
	strval string
	intval int
	isErr  bool
}{
	{"  16505550100  ", 16505550100, false},
	{"+16505550100", 16505550100, false},
	{"+ 16505550100", 16505550100, true},
}

func TestNumber(t *testing.T) {
	for _, tt := range numberTests {
		got, err := E164Atoi(tt.strval)
		if err != nil {
			if tt.isErr {
				continue
			}
			t.Errorf("E164Atoi error [%s]", err.Error())
		}
		if got != tt.intval {
			t.Errorf("E164Atoi failed want [%d] got [%d]", tt.intval, got)
		}
	}
}
