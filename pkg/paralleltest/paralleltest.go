package paralleltest

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

// TODO add ignoring ability flag
func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "paralleltest",
		Doc:      "Checks that tests have t.Parallel enabled",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Run only for test files
	inspector := inspector.New(getTestFiles(pass))

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		funcDecl := node.(*ast.FuncDecl)
		funcHasParallelMethod := false
		rangeStatementExists := false
		rangeHasParallelMethod := false

		// Check runs for test functions only
		if !isItATestFunction(funcDecl) {
			return
		}

		for _, l := range funcDecl.Body.List {
			switch v := l.(type) {

			// Check if the test method is calling t.parallel
			case *ast.ExprStmt:
				ast.Inspect(v, func(n ast.Node) bool {
					if funcHasParallelMethod == false {
						funcHasParallelMethod = callExprCallsMethodParallel(n)
					}
					return true
				})

			case *ast.RangeStmt:
				// Check if the range over testcases is calling t.parallel
				// TODO: Check range statements is over testcases and not any other ranges

				rangeStatementExists = true
				// TODO: Also check for the assignment tc:tc
				ast.Inspect(v, func(n ast.Node) bool {
					if rangeHasParallelMethod == false {
						rangeHasParallelMethod = callExprCallsMethodParallel(n)
					}
					return true
				})
			}
		}

		if !funcHasParallelMethod {
			pass.Reportf(node.Pos(), "Function %s missing the call to method parallel \n", funcDecl.Name.Name)
		}
		if rangeStatementExists && !rangeHasParallelMethod {
			pass.Reportf(node.Pos(), "Range statement for test %s missing the call to method parallel \n", funcDecl.Name.Name)
		}

	})

	return nil, nil
}

func getTestFiles(pass *analysis.Pass) []*ast.File {
	testFileSuffix := "_test.go"

	var testFiles []*ast.File
	for _, f := range pass.Files {
		fileName := pass.Fset.Position(f.Pos()).Filename
		if strings.HasSuffix(fileName,testFileSuffix ) {
			testFiles = append(testFiles, f)
		}
	}
	return testFiles
}

func callExprCallsMethodParallel(node ast.Node) bool {
	methodName := "Parallel"

	switch n := node.(type) {
	default:
	case *ast.CallExpr:
		if fun, ok := n.Fun.(*ast.SelectorExpr); ok {
			return fun.Sel.Name == methodName
		}
	}
	return false
}

func isItATestFunction(funcDecl *ast.FuncDecl) bool {
	testMethodParamName := "t"
	testMethodPackageType := "testing"
	testMethodStruct := "T"

	if funcDecl.Type.Params != nil {
		if len(funcDecl.Type.Params.List) > 1 {
			return false
		}

		for _, param := range funcDecl.Type.Params.List {
			if param.Names[0].String() == testMethodParamName {
				if starExp, ok := param.Type.(*ast.StarExpr); ok {
					if selectExpr, ok := starExp.X.(*ast.SelectorExpr); ok {
						return fmt.Sprint(selectExpr.X) == testMethodPackageType &&
							fmt.Sprint(selectExpr.Sel) == testMethodStruct
					}
				}
			}
		}
	}
	return false
}
