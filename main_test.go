package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/mzbaulhaque/gois/internal/util/testutil"
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

	var testCases = []testutil.TestCase{
		{
			In:  []string{"bing", "cats"},
			Out: nil,
		},
		{
			In:  []string{"google", "cats"},
			Out: nil,
		},
		{
			In:  []string{"bing", "-x", "cats"},
			Out: nil,
		},
		{
			In:  []string{"google", "-x", "cats"},
			Out: nil,
		},
		//{
		//	In:  []string{"bing", "-g", "cats"},
		//	Out: nil,
		//},
		//{
		//	In:  []string{"google", "-g", "cats"},
		//	Out: nil,
		//},
		//{
		//	In:  []string{"bing", "-B", "cats"},
		//	Out: nil,
		//},
		//{
		//	In:  []string{"google", "-B", "cats"},
		//	Out: nil,
		//},
		{
			In:  []string{"bing", "-H", "1080", "-w", "1920", "cats"},
			Out: nil,
		},
		{
			In:  []string{"google", "-H", "1080", "-w", "1920", "cats"},
			Out: nil,
		},
	}

	for _, test := range testCases {
		cmd := exec.Command(goisBinaryAbsPath, test.In.([]string)...)
		output, err := cmd.CombinedOutput()
		testutil.CheckErr(t, err)
		testutil.CheckCmdOutput(t, output, re)
	}
}

func TestGois_Compact(t *testing.T) {
	re := regexp.MustCompile("^http.+\n")

	var testCases = []testutil.TestCase{
		{
			In:  []string{"bing", "-c", "cats"},
			Out: nil,
		},
		{
			In:  []string{"google", "-c", "cats"},
			Out: nil,
		},
		{
			In:  []string{"bing", "-c", "-x", "cats"},
			Out: nil,
		},
		{
			In:  []string{"google", "-c", "-x", "cats"},
			Out: nil,
		},
		//{
		//	In:  []string{"google", "-c", "-g", "cats"},
		//	Out: nil,
		//},
		//{
		//	In:  []string{"google", "-c", "-B", "cats"},
		//	Out: nil,
		//},
		{
			In:  []string{"bing", "-c", "-H", "1080", "-w", "1920", "cats"},
			Out: nil,
		},
		{
			In:  []string{"google", "-c", "-H", "1080", "-w", "1920", "cats"},
			Out: nil,
		},
	}

	for _, test := range testCases {
		cmd := exec.Command(goisBinaryAbsPath, test.In.([]string)...)
		output, err := cmd.CombinedOutput()
		testutil.CheckErr(t, err)
		testutil.CheckCmdOutput(t, output, re)
	}
}
