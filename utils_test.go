package main

import "testing"

type Test struct {
	in  string
	out int
}

var AtoiTests = []Test{
	{
		in:  "123",
		out: 123,
	},
	{
		in:  "0",
		out: 0,
	},
	{
		in:  "-1",
		out: -1,
	},
}

func TestAtoi(t *testing.T) {
	for _, test := range AtoiTests {
		got := Atoi(test.in)
		if test.out != got {
			t.Errorf("Atoi invalid output %d for input %s, expected %d", got, test.in, test.out)
		}
	}
}

func TestItoa(t *testing.T) {
	for _, test := range AtoiTests {
		got := Itoa(test.out)
		if test.in != got {
			t.Errorf("Atoi invalid output %s for input %d, expected %s", got, test.out, test.in)
		}
	}
}
