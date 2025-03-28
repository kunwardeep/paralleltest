package paralleltest

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMissing(t *testing.T) {
	t.Parallel()

	analysistest.Run(t, analysistest.TestData(), NewAnalyzer(), "t")
}

func TestIgnoreMissingOption(t *testing.T) {
	t.Parallel()

	a := newParallelAnalyzer()
	a.ignoreMissing = true

	analysistest.Run(t, analysistest.TestData(), a.analyzer, "i")
}

func TestIgnoreMissingSubtestsOption(t *testing.T) {
	t.Parallel()

	a := newParallelAnalyzer()
	a.ignoreMissingSubtests = true

	analysistest.Run(t, analysistest.TestData(), a.analyzer, "ignoremissingsubtests")
}
