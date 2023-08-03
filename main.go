package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/kunwardeep/paralleltest/pkg/paralleltest"
)

func main() {
	singlechecker.Main(paralleltest.NewAnalyzer())
}
