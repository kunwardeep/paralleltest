package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:     "paralleltest",
	Doc:      "Checks that tests have t.Parallel enabled",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		funcDecl := node.(*ast.FuncDecl)
		funcHasParallelMethod := false
		rangeStatementExists := false
		rangeHasParallelMethod := false

		// TODO : Only run for test files
		if strings.HasPrefix(funcDecl.Name.Name, "Test_") {
			for _, l := range funcDecl.Body.List {
				switch v := l.(type) {

				case *ast.ExprStmt:
					ast.Inspect(v, func(n ast.Node) bool {
						if funcHasParallelMethod == false {
							funcHasParallelMethod = callExprCallsMethodParallel(n)
						}
						return true
					})

				case *ast.RangeStmt:
					// TODO: Check range statements is over testcases

					rangeStatementExists= true
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
				pass.Reportf(node.Pos(), "Range statement %s missing the call to method parallel \n", funcDecl.Name.Name)
			}
		}
	})

	return nil, nil
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
