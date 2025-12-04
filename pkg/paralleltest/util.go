package paralleltest

import (
	"bytes"
	//nolint:gosec // G401: Weak cryptographic primitive
	"crypto/md5"
	"encoding/hex"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// goFmt returns the position plus the string representation of an AST subtree.
func goFmt(node ast.Node) string {
	buf := bytes.Buffer{}
	fs := token.NewFileSet()
	buf.WriteString(strconv.Itoa(int(node.Pos())))
	gofmtConfig := &printer.Config{Tabwidth: 8}
	_ = gofmtConfig.Fprint(&buf, fs, node)
	return buf.String()
}

// nodeHash yields the MD5 hash of the given AST node.
func nodeHash(node ast.Node) string {
	hasher := func(in string) string {
		//nolint:gosec // G401: Weak cryptographic primitive
		binHash := md5.Sum([]byte(in))
		return hex.EncodeToString(binHash[:])
	}
	str := goFmt(node)
	return hasher(str)
}

// isTestFunction checks if a function declaration is a test function
// A test function must:
// 1. Start with "Test"
// 2. Have exactly one parameter
// 3. Have that parameter be of type *testing.T
// Returns true if it is a test function, otherwise false.
func isTestFunction(funcDecl *ast.FuncDecl) bool {
	testPrefix := "Test"

	if !strings.HasPrefix(funcDecl.Name.Name, testPrefix) {
		return false
	}

	if !hasExactlyOneParameter(funcDecl) {
		return false
	}

	return findTestParamName(funcDecl.Type.Params) != ""
}

func hasExactlyOneParameter(funcDecl *ast.FuncDecl) bool {
	return funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) == 1
}

// getCallName returns the name of the function called in a call expression.
func getCallName(callExpr *ast.CallExpr) string {
	switch fun := callExpr.Fun.(type) {
	case *ast.Ident:
		return fun.Name
	case *ast.SelectorExpr:
		return fun.Sel.Name
	default:
		return ""
	}
}

// findFunction looks for the function declaration across all input files by name.
// This is slightly incomplete, as we don't handle methods.
func findFunction(pass *analysis.Pass, name string) *ast.FuncDecl {
	if name == "" {
		return nil
	}
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == name {
				return funcDecl
			}
		}
	}
	return nil
}

// findTestParamName returns the first parameter of type *testing.T.
// This is for analyzing test helper functions.
func findTestParamName(params *ast.FieldList) string {
	for i := range params.List {
		param := params.List[i]
		if starExp, ok := param.Type.(*ast.StarExpr); ok {
			if selectExpr, ok := starExp.X.(*ast.SelectorExpr); ok {
				if selectExpr.Sel.Name == testMethodStruct {
					if s, ok := selectExpr.X.(*ast.Ident); ok {
						if s.Name == testMethodPackageType && len(param.Names) > 0 {
							return param.Names[0].Name
						}
					}
				}
			}
		}
	}
	return ""
}

func isParallelCall(node *ast.CallExpr, testVar string) bool {
	return exprCallHasMethod(node, testVar, "Parallel")
}

func isTestRunCall(node *ast.CallExpr, testVar string) bool {
	return exprCallHasMethod(node, testVar, "Run")
}

func isSetenvCall(node *ast.CallExpr, testVar string) bool {
	return exprCallHasMethod(node, testVar, "Setenv")
}

func isChdirCall(node *ast.CallExpr, testVar string) bool {
	return exprCallHasMethod(node, testVar, "Chdir")
}

func exprCallHasMethod(callExpr *ast.CallExpr, receiverName, methodName string) bool {
	// nolint: gocritic
	if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		if receiver, ok := fun.X.(*ast.Ident); ok {
			return receiver.Name == receiverName && fun.Sel.Name == methodName
		}
	}
	return false
}
