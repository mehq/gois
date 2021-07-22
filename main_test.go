package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/mzbaulhaque/gois/internal/util"
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
	cmd := exec.Command("go", "build", "-o", goisBinaryAbsPath)
	output, err := cmd.Output()

	if err != nil {
		log.Fatal(err, string(output))
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
			In:  []string{"google", "dogs"},
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
			In:  []string{"google", "-c", "dogs"},
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
