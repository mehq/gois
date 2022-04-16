// E2E testing.

package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/mehq/gois/internal/util"
)

const (
	goisCmdBinary = "gois"
	tmpDir        = ".tmp"
)

var (
	goisBinaryRelPath    = path.Join(tmpDir, goisCmdBinary)
	goisBinaryAbsPath, _ = filepath.Abs(goisBinaryRelPath)
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

// setup any prerequisites for the tests.
func setup() {
	// Build test binary with mocked data.
	cmd := exec.Command("go", "build", "-tags", "test", "-o", goisBinaryAbsPath)
	output, err := cmd.Output()

	if err != nil {
		log.Fatalf("error building test binary: %v (%s)", err, string(output))
	}
}

func TestGois(t *testing.T) {
	re := regexp.MustCompile("^Title: .+\nWebpage: http.*\nResolution: [0-9]+x[0-9]+\nURL: http.+\nThumbnail: http.+\n")

	var testCases = []util.TestCase{
		{
			In:  []string{"bing", "cats"},
			Out: nil,
		},
		{
			In:  []string{"flickr", "cats"},
			Out: nil,
		},
		{
			In:  []string{"google", "cats"},
			Out: nil,
		},
		{
			In:  []string{"yahoo", "cats"},
			Out: nil,
		},
		{
			In:  []string{"yandex", "cats"},
			Out: nil,
		},
	}

	for _, test := range testCases {
		cmd := exec.Command(goisBinaryAbsPath, test.In.([]string)...)
		output, err := cmd.CombinedOutput()
		util.CheckErr(t, err)
		util.CheckCmdOutput(t, output, re)
	}
}

func TestGois_Compact(t *testing.T) {
	re := regexp.MustCompile("^http.+\n")

	var testCases = []util.TestCase{
		{
			In:  []string{"bing", "-c", "cats"},
			Out: nil,
		},
		{
			In:  []string{"flickr", "-c", "cats"},
			Out: nil,
		},
		{
			In:  []string{"google", "-c", "cats"},
			Out: nil,
		},
		{
			In:  []string{"yahoo", "-c", "cats"},
			Out: nil,
		},
		{
			In:  []string{"yandex", "-c", "cats"},
			Out: nil,
		},
	}

	for _, test := range testCases {
		cmd := exec.Command(goisBinaryAbsPath, test.In.([]string)...)
		output, err := cmd.CombinedOutput()
		util.CheckErr(t, err)
		util.CheckCmdOutput(t, output, re)
	}
}
