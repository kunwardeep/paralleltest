package paralleltest

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMissing(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(wd), "paralleltest", "testdata")
	analysistest.Run(t, testdata, Analyzer, "t")
}

func TestIgnoreMissing(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	options := flag.NewFlagSet("", flag.ExitOnError)
	options.Bool("i", true, "")

	analyzer := *Analyzer
	analyzer.Flags = *options

	testdata := filepath.Join(filepath.Dir(wd), "paralleltest", "testdata")
	analysistest.Run(t, testdata, &analyzer, "i")
}
