package paralleltest

import (
	"testing"

	_ "github.com/stretchr/testify"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMissing(t *testing.T) {
	t.Parallel()

	analyzer := NewAnalyzer(Config{CheckCleanup: true})

	analysistest.Run(t, analysistest.TestData(), analyzer, "t")
}

func TestIgnoreMissingOption(t *testing.T) {
	t.Parallel()

	analyzer := NewAnalyzer(Config{IgnoreMissing: true})

	analysistest.Run(t, analysistest.TestData(), analyzer, "i")
}

func TestIgnoreMissingSubtestsOption(t *testing.T) {
	t.Parallel()

	analyzer := NewAnalyzer(Config{IgnoreMissingSubtests: true})

	analysistest.Run(t, analysistest.TestData(), analyzer, "ignoremissingsubtests")
}

func TestCheckCleanupOption(t *testing.T) {
	t.Parallel()

	analyzer := NewAnalyzer(Config{CheckCleanup: true})

	analysistest.Run(t, analysistest.TestData(), analyzer, "cleanup")
}

func TestExtraSigsOptions(t *testing.T) {
	t.Parallel()

	analyzer := NewAnalyzer(Config{ExtraSigs: []string{"ExtraSigs"}})

	analysistest.Run(t, analysistest.TestData(), analyzer, "skip")
}
