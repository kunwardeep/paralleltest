package analyzer

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goprintffuncname",
	Doc:      "Checks that printf-like functions are named with `f` at the end.",
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

		// TODO : Only run for test files
		if strings.HasPrefix(funcDecl.Name.Name, "Test_") {
			fmt.Println(funcDecl.Name.Name)
			//funcName := extractFuncCallInFunc(funcDecl.Body.List)
			for _, l := range funcDecl.Body.List {

				switch v := l.(type) {
				case *ast.ExprStmt:

					expressionStatementHasParallel := doesItHaveParallelMethod(v)
					fmt.Println("expressionStatementHasParallel", expressionStatementHasParallel)
				case *ast.RangeStmt:
					// TODO: Also check for the assignment tc:tc

					rangeStatementHasParallel := doesItHaveParallelMethod(v)
					fmt.Println("rangeStatementHasParallel", rangeStatementHasParallel)

					//funcName := extractFuncCallFromRangeStmt(v)
					//fmt.Println("RangeStmt",funcName)
				}
			}
			return
		}
		//pass.Reportf(node.Pos(), "printf-like formatting function '%s' should be named '%sf'",
		//	funcDecl.Name.Name, funcDecl.Name.Name)
	})

	return nil, nil
}

func extractFuncCallFromExprStmt(stmt ast.Stmt) string {
	if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
		if call, ok := exprStmt.X.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				return fun.Sel.Name
			}
		}
	}
	return ""
}

func checkParallelMethod(stmt ast.Stmt)  {
	ast.Inspect(stmt, isParallelCalled)
}

func isParallelCalled(node ast.Node) bool {
	methodName := "Parallel"

	switch n := node.(type) {
	case *ast.CallExpr:
		if fun, ok := n.Fun.(*ast.SelectorExpr); ok {
			return fun.Sel.Name == methodName // prints every func call expression
		}
	}
	return false
}

//func extractFuncCallFromRangeStmt(stmt ast.Stmt) string {
//	if blockStmt, ok := stmt.(*ast.BlockStmt); ok {
//		for _,blkStmt := range blockStmt.List{
//			if exprStmt, ok := blkStmt.(*ast.ExprStmt); ok {
//				if call, ok := exprStmt.X.(*ast.CallExpr); ok {
//					for _, arg := range call.Args{
//						if funcLit, ok := arg.(*ast.FuncLit); ok {
//							if bb, ok := funcLit.Body.(*ast.BlockStmt); ok {
//							}
//
//						}
//					}
//				}
//				//if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
//				//	return fun.Sel.Name
//				//}
//			}
//		}
//
//	}
//	return ""
//}
