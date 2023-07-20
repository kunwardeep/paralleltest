package paralleltest

import (
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
	analysistest.Run(t, testdata, NewAnalyzer(), "t")
}

func TestIgnoreMissingOption(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	a := newParallelAnalyzer()
	a.ignoreMissing = true

	testdata := filepath.Join(filepath.Dir(wd), "paralleltest", "testdata")
	analysistest.Run(t, testdata, a.analyzer, "i")
}

func TestIgnoreMissingSubtestsOption(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	a := newParallelAnalyzer()
	a.ignoreMissingSubtests = true

	testdata := filepath.Join(filepath.Dir(wd), "paralleltest", "testdata")
	analysistest.Run(t, testdata, a.analyzer, "ignoremissingsubtests")
}
