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
		Doc:      "Checks that tests have t.Parallel enabled and that range loop variable is reinitialised",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Run only for test files
	inspector := inspector.New(testFiles(pass))

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		funcDecl := node.(*ast.FuncDecl)
		var funcHasParallelMethod,
			rangeStatementOverTestCasesExists,
			rangeStatementHasParallelMethod,

			testLoopVariableReinitialised bool
		var testRunLoopIdentifier string

		var rangeNode ast.Node

		// Check runs for test functions only
		if !testFunction(funcDecl) {
			return
		}

		for _, l := range funcDecl.Body.List {
			switch v := l.(type) {

			// Check if the test method is calling t.parallel
			case *ast.ExprStmt:
				ast.Inspect(v, func(n ast.Node) bool {
					if !funcHasParallelMethod {
						funcHasParallelMethod = methodParallelIsCalledInTestFunction(n)
					}
					return true
				})

			// Check if the range over testcases is calling t.parallel
			case *ast.RangeStmt:
				rangeNode = v

				ast.Inspect(v, func(n ast.Node) bool {
					switch r := n.(type) {
					case *ast.ExprStmt:
						if methodRunIsCalledInRangeStatement(r.X) {
							rangeStatementOverTestCasesExists = true
							testRunLoopIdentifier = methodRunFirstArgumentObjectName(r.X)

							if !rangeStatementHasParallelMethod {
								rangeStatementHasParallelMethod = methodParallelIsCalledInMethodRun(r.X)
							}
						}
					}
					return true
				})

				// Check for the range loop value identifier re assignment
				// More info here https://gist.github.com/kunwardeep/80c2e9f3d3256c894898bae82d9f75d0
				if rangeStatementOverTestCasesExists {
					var rangeValueIdentifier string
					if i, ok := v.Value.(*ast.Ident); ok {
						rangeValueIdentifier = i.Name
					}

					testLoopVariableReinitialised = testCaseLoopVariableReinitialised(v.Body.List, rangeValueIdentifier, testRunLoopIdentifier)

				}
			}
		}

		if !funcHasParallelMethod {
			pass.Reportf(node.Pos(), "Function %s missing the call to method parallel \n", funcDecl.Name.Name)
		}

		if rangeStatementOverTestCasesExists && rangeNode != nil {
			if !rangeStatementHasParallelMethod {
				pass.Reportf(rangeNode.Pos(), "Range statement for test %s missing the call to method parallel \n", funcDecl.Name.Name)
			}
			if !testLoopVariableReinitialised {
				pass.Reportf(rangeNode.Pos(), "Range statement for test %s does not reinitialise the variable %s  \n", funcDecl.Name.Name, testRunLoopIdentifier)
			}
		}
	})

	return nil, nil
}

func testCaseLoopVariableReinitialised(statements []ast.Stmt, rangeValueIdentifier string, testRunLoopIdentifier string) bool {
	if len(statements) > 1 {
		for _, s := range statements {
			leftIdentifier, rightIdentifier := getLeftAndRightIdentifier(s)
			if leftIdentifier == testRunLoopIdentifier && rightIdentifier == rangeValueIdentifier {
				return true
			}
		}
	}
	return false
}

// Return the left hand side and the right hand side identifiers name
func getLeftAndRightIdentifier(s ast.Stmt) (string, string) {
	var leftIdentifier, rightIdentifier string
	switch v := s.(type) {
	case *ast.AssignStmt:
		if len(v.Rhs) == 1 {
			if i, ok := v.Rhs[0].(*ast.Ident); ok {
				rightIdentifier = i.Name
			}
		}
		if len(v.Lhs) == 1 {
			if i, ok := v.Lhs[0].(*ast.Ident); ok {
				leftIdentifier = i.Name
			}
		}
	}
	return leftIdentifier, rightIdentifier
}

func testFiles(pass *analysis.Pass) []*ast.File {
	testFileSuffix := "_test.go"

	var testFiles []*ast.File
	for _, f := range pass.Files {
		fileName := pass.Fset.Position(f.Pos()).Filename
		if strings.HasSuffix(fileName, testFileSuffix) {
			testFiles = append(testFiles, f)
		}
	}
	return testFiles
}

func methodParallelIsCalledInMethodRun(node ast.Node) bool {
	var methodParallelCalled bool
	switch callExp := node.(type) {
	case *ast.CallExpr:
		for _, arg := range callExp.Args {
			if !methodParallelCalled {
				ast.Inspect(arg, func(n ast.Node) bool {
					if !methodParallelCalled {
						methodParallelCalled = methodParallelIsCalledInRunMethod(n)
						return true
					}
					return false
				})
			}
		}
	}
	return methodParallelCalled
}

func methodParallelIsCalledInRunMethod(node ast.Node) bool {
	return exprCallHasMethod(node, "Parallel")
}

func methodParallelIsCalledInTestFunction(node ast.Node) bool {
	return exprCallHasMethod(node, "Parallel")
}

func methodRunIsCalledInRangeStatement(node ast.Node) bool {
	return exprCallHasMethod(node, "Run")
}

func exprCallHasMethod(node ast.Node, methodName string) bool {
	switch n := node.(type) {
	case *ast.CallExpr:
		if fun, ok := n.Fun.(*ast.SelectorExpr); ok {
			return fun.Sel.Name == methodName
		}
	}
	return false
}

func methodRunFirstArgumentObjectName(node ast.Node) string {
	switch n := node.(type) {
	case *ast.CallExpr:
		for _, arg := range n.Args {
			if s, ok := arg.(*ast.SelectorExpr); ok {
				if i, ok := s.X.(*ast.Ident); ok {
					return i.Name
				}
			}
		}
	}
	return ""
}

func testFunction(funcDecl *ast.FuncDecl) bool {
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
