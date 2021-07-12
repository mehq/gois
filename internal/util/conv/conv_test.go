package conv

import (
	"testing"

	"github.com/mzbaulhaque/gois/internal/util/testutil"
)

var atoiTestCases = []testutil.TestCase{
	{
		In:  "123",
		Out: 123,
	},
	{
		In:  "0",
		Out: 0,
	},
	{
		In:  "-1",
		Out: -1,
	},
}

func TestAtoi(t *testing.T) {
	for _, test := range atoiTestCases {
		got := Atoi(test.In.(string))
		if test.Out != got {
			t.Errorf("conv#Atoi invalid output %d for input %s, expected %d", got, test.In, test.Out)
		}
	}

	// edge case
	got := Atoi("ABC")
	if got != 0 {
		t.Errorf("conv#Atoi invalid output %d for input %s, expected %d", got, "ABC", 0)
	}
}

func TestItoa(t *testing.T) {
	for _, test := range atoiTestCases {
		got := Itoa(test.Out.(int))
		if test.In != got {
			t.Errorf("conv#Itoa invalid output %s for input %d, expected %s", got, test.Out, test.In)
		}
	}
}
