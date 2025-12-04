package paralleltest

import (
	"flag"
	"go/ast"
	"strings"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

const testMethodPackageType = "testing"
const testMethodStruct = "T"
const Doc = `check that tests use t.Parallel() and use it properly`

type Config struct {
	// IgnoreMissing check that missing calls to t.Parallel are not reported
	IgnoreMissing bool `json:"ignoreMissing"`
	// IgnoreMissingSubtests check that missing calls to t.Parallel in subtests are not reported
	IgnoreMissingSubtests bool `json:"ignoreMissingSubtests"`
	// CheckCleanup check that defer is not used with t.Parallel (use t.Cleanup instead)
	CheckCleanup bool `json:"checkCleanup"`
	// ExtraSigs is a list of extra functions that cannot be used with t.Parallel
	ExtraSigs []string `json:"extraSigs"`
}

func NewAnalyzer(config Config) *analysis.Analyzer {
	a := &parallelAnalyzer{
		config:  config,
		mu:      &sync.RWMutex{},
		visited: make(map[string]*testAnalysis),
	}

	var ignoreLoopVar bool
	var flags flag.FlagSet
	flags.BoolVar(&config.IgnoreMissing, "i", false, "ignore missing calls to t.Parallel")
	flags.BoolVar(&config.IgnoreMissingSubtests, "ignoremissingsubtests", false, "ignore missing calls to t.Parallel in subtests")
	flags.BoolVar(&ignoreLoopVar, "ignoreloopVar", false, "ignore loop variable detection <deprecated with go 1.22>")
	flags.BoolVar(&config.CheckCleanup, "checkcleanup", false, "check that defer is not used with t.Parallel (use t.Cleanup instead)")

	return &analysis.Analyzer{
		Name:  "paralleltest",
		Doc:   Doc,
		Run:   a.run,
		Flags: flags,
	}
}

// parallelAnalyzer is an internal analyzer that makes options available to a
// run pass. It wraps an `analysis.Analyzer` that should be returned for
// linters.
type parallelAnalyzer struct {
	config  Config
	mu      *sync.RWMutex
	visited map[string]*testAnalysis
}

type testAnalysis struct {
	hasParallel,
	cantParallel,
	funcHasDeferStatement bool
	numberOfTestRun int
	deferStatements []ast.Node
}

func (a *testAnalysis) merge(other *testAnalysis) {
	a.hasParallel = a.hasParallel || other.hasParallel
	a.cantParallel = a.cantParallel || other.cantParallel
	a.numberOfTestRun += other.numberOfTestRun
}

// getAnalysis returns the cached analysis for the given node, or nil if it has not been visited yet.
func (a *parallelAnalyzer) getAnalysis(node ast.Node) (string, *testAnalysis) {
	hash := nodeHash(node)
	a.mu.RLock()
	defer a.mu.RUnlock()
	if v, ok := a.visited[hash]; ok {
		return hash, v
	}
	return hash, nil
}

// cacheAnalysis caches the given analysis for the given node.
func (a *parallelAnalyzer) cacheAnalysis(hash string, analysis *testAnalysis) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.visited[hash] = analysis
}

func (a *parallelAnalyzer) run(pass *analysis.Pass) (any, error) {
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspector.New(pass.Files).Preorder(nodeFilter, func(node ast.Node) {
		funcDecl := node.(*ast.FuncDecl)
		// Only process _test.go files
		if !strings.HasSuffix(pass.Fset.File(funcDecl.Pos()).Name(), "_test.go") {
			return
		}

		// Check runs for test functions only
		if isTestFunction(funcDecl) {
			a.analyzeTestFunction(pass, funcDecl)
		}
	})

	return nil, nil
}

func (a *parallelAnalyzer) analyzeTestFunction(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	analysis := a.analyzeFunction(pass, funcDecl)

	if !a.config.IgnoreMissing && !analysis.hasParallel && !analysis.cantParallel {
		pass.Reportf(funcDecl.Pos(), "Function %s missing the call to method parallel\n", funcDecl.Name.Name)
	}

	a.reportDefer(pass, analysis, funcDecl.Name.Name)
}

func (a *parallelAnalyzer) reportDefer(pass *analysis.Pass, analysis *testAnalysis, name string) {
	if a.config.CheckCleanup && analysis.hasParallel && analysis.funcHasDeferStatement && analysis.numberOfTestRun > 0 {
		for _, deferStmt := range analysis.deferStatements {
			pass.Reportf(deferStmt.Pos(), "Function %s uses defer with t.Parallel, use t.Cleanup instead to ensure cleanup runs after parallel subtests complete\n", name)
		}
	}
}

func (a *parallelAnalyzer) reportParallelSubtest(pass *analysis.Pass, analysis *testAnalysis, node ast.Node, name string) {
	if !a.config.IgnoreMissing && !a.config.IgnoreMissingSubtests && !analysis.hasParallel && !analysis.cantParallel {
		pass.Reportf(node.Pos(), "Function %s missing the call to method parallel in the t.Run\n", name)
	}
}

// analyzeTestRun analyzes the three types of t.Run calls:
// 1. Inline function: t.Run("name", new func(t *testing.T) {...})
// 2. Direct function identifier: t.Run("name", myFunc)
// 3. Builder function: t.Run("name", builder(t))
func (a *parallelAnalyzer) analyzeTestRun(pass *analysis.Pass, callExpr *ast.CallExpr, testVar string) *testAnalysis {
	if isTestRunCall(callExpr, testVar) && len(callExpr.Args) > 1 {
		if funcLit, ok := callExpr.Args[1].(*ast.FuncLit); ok {
			// Case 1: Inline function: t.Run("name", new func(t *testing.T) {...})
			analysis := a.analyzeFuncLit(pass, funcLit)

			a.reportDefer(pass, analysis, "literal")
			a.reportParallelSubtest(pass, analysis, funcLit, "literal")
			analysis.numberOfTestRun++

			return analysis
		} else if ident, ok := callExpr.Args[1].(*ast.Ident); ok {
			// Case 2: Direct function identifier: t.Run("name", myFunc)
			funcDecl := findFunction(pass, ident.Name)
			if funcDecl != nil && hasExactlyOneParameter(funcDecl) {
				analysis := a.analyzeFunction(pass, funcDecl)

				a.reportDefer(pass, analysis, ident.Name)
				a.reportParallelSubtest(pass, analysis, callExpr, ident.Name)
				analysis.numberOfTestRun++

				return analysis
			}
		} else if builderCall, ok := callExpr.Args[1].(*ast.CallExpr); ok {
			// Case 3: Builder function: t.Run("name", builder(t))
			funcName := getCallName(builderCall)
			funcDecl := findFunction(pass, funcName)

			if funcDecl != nil {
				parentAnalysis, builderAnalysis := a.analyzeBuilderCall(pass, funcDecl)

				a.reportDefer(pass, builderAnalysis, funcName)
				a.reportParallelSubtest(pass, builderAnalysis, callExpr, funcName)
				parentAnalysis.merge(builderAnalysis)
				parentAnalysis.numberOfTestRun++

				return parentAnalysis
			}
		}
	}

	return &testAnalysis{}
}

func (a *parallelAnalyzer) analyzeFunctionCall(pass *analysis.Pass, callExpr *ast.CallExpr) *testAnalysis {
	funcName := getCallName(callExpr)
	funcDecl := findFunction(pass, funcName)
	if funcDecl == nil {
		return &testAnalysis{}
	}

	return a.analyzeFunction(pass, funcDecl)
}

func (a *parallelAnalyzer) visitExprStmt(pass *analysis.Pass, analysis *testAnalysis, testVar string) func(n ast.Node) bool {
	return func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			a.analyzeCallExpr(pass, analysis, testVar, callExpr)
			return false
		}
		return true
	}
}

func (a *parallelAnalyzer) analyzeCallExpr(pass *analysis.Pass, analysis *testAnalysis, testVar string, callExpr *ast.CallExpr) {
	// Edge case, check each parameter of the call to analyze.
	for _, arg := range callExpr.Args {
		if nestedCallExpr, ok := arg.(*ast.CallExpr); ok {
			a.analyzeCallExpr(pass, analysis, testVar, nestedCallExpr)
		}
	}

	analysis.hasParallel = analysis.hasParallel || isParallelCall(callExpr, testVar)
	analysis.cantParallel = analysis.cantParallel || isSetenvCall(callExpr, testVar)
	analysis.cantParallel = analysis.cantParallel || isChdirCall(callExpr, testVar)
	if fnIdent, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		signature := pass.TypesInfo.ObjectOf(fnIdent.Sel).String()
		analysis.cantParallel = analysis.cantParallel || contains(a.config.ExtraSigs, signature)
	}
	if fnIdent, ok := callExpr.Fun.(*ast.Ident); ok {
		signature := pass.TypesInfo.ObjectOf(fnIdent).String()
		analysis.cantParallel = analysis.cantParallel || contains(a.config.ExtraSigs, signature)
	}
	analysis.merge(a.analyzeTestRun(pass, callExpr, testVar))
	analysis.merge(a.analyzeFunctionCall(pass, callExpr))
}

func (a *parallelAnalyzer) analyzeFuncLit(pass *analysis.Pass, funcLit *ast.FuncLit) *testAnalysis {
	return a.analyzeFunctionF(pass, funcLit.Type, funcLit.Body)
}

func (a *parallelAnalyzer) analyzeFunction(pass *analysis.Pass, funcDecl *ast.FuncDecl) *testAnalysis {
	return a.analyzeFunctionF(pass, funcDecl.Type, funcDecl.Body)
}

func (a *parallelAnalyzer) analyzeFunctionF(pass *analysis.Pass, funcType *ast.FuncType, body *ast.BlockStmt) *testAnalysis {
	analysis := &testAnalysis{}

	testVar := findTestParamName(funcType.Params)
	if testVar == "" {
		return analysis
	}
	hash, v := a.getAnalysis(body)
	if v != nil {
		return v
	}

	for _, l := range body.List {
		switch v := l.(type) {
		case *ast.DeferStmt:
			if a.config.CheckCleanup {
				analysis.funcHasDeferStatement = true
				analysis.deferStatements = append(analysis.deferStatements, v)
			}
		default:
			ast.Inspect(v, a.visitExprStmt(pass, analysis, testVar))
		}
	}

	a.cacheAnalysis(hash, analysis)

	return analysis
}

// analyzeBuilderCall analyzes a function call that returns a test function
// to see if the returned function contains t.Parallel()
func (a *parallelAnalyzer) analyzeBuilderCall(pass *analysis.Pass, funcDecl *ast.FuncDecl) (*testAnalysis, *testAnalysis) {
	parentAnalysis := &testAnalysis{}
	builderAnalysis := &testAnalysis{}
	testVar := findTestParamName(funcDecl.Type.Params)

	// Found the builder function, analyze it and return immediately
	ast.Inspect(funcDecl, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.ExprStmt, *ast.RangeStmt:
			// We only need to analyze if we have a test variable to actually check for t.Parallel()
			if testVar != "" {
				ast.Inspect(v, a.visitExprStmt(pass, parentAnalysis, testVar))
			}
		case *ast.ReturnStmt:
			// Check if the return value is a function literal
			for _, result := range v.Results {
				if funcLit, ok := result.(*ast.FuncLit); ok {
					innerAnalysis := a.analyzeFuncLit(pass, funcLit)
					builderAnalysis.merge(innerAnalysis)
				}
			}
		}
		return true
	})
	return parentAnalysis, builderAnalysis
}

func contains(slice []string, el string) bool {
	for _, s := range slice {
		if strings.Contains(el, s) {
			return true
		}
	}

	return false
}
