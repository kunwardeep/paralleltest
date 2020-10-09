package main

import (
	"github.com/kunwardeep/paralleltest/pkg/paralleltest"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(paralleltest.NewAnalyzer())
}
