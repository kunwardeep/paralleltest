package main

import (
	"github.com/jirfag/go-printf-func-name/pkg/paralleltest"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(paralleltest.NewAnalyzer())
}
