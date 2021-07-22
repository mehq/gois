package util

import (
	"regexp"
	"testing"
)

// TestCase represents a generic test case.
type TestCase struct {
	In  interface{}
	Out interface{}
}

// CheckErr can be used to report generic test errors.
func CheckErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

// CheckCmdOutput can be used to match output of a command with a target Regexp.
func CheckCmdOutput(t *testing.T, output []byte, matchWith *regexp.Regexp) {
	if !matchWith.Match(output) {
		t.Error("command output not matching")
	}
}
