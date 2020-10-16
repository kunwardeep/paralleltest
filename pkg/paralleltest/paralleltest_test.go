package paralleltest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(wd),"paralleltest", "testdata")
	fmt.Println("Working dir ---",wd)
	fmt.Println("filepath.Dir(wd) ---",filepath.Dir(wd))

	analysistest.Run(t, testdata, NewAnalyzer(), "t")
}
