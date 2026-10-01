//go:build go1.24

package t

import "testing"

func (notATest) Chdir(string) {}

func TestFunctionWithChdir(t *testing.T) {
	t.Chdir(".")
}

func TestFunctionWithChdirOtherTestingVar(q *testing.T) {
	q.Chdir(".")
}

func TestFunctionWithChdirLookalike(t *testing.T) { // want "Function TestFunctionWithChdirLookalike missing the call to method parallel"
	var other notATest
	other.Chdir(".")
}

func TestFunctionWithChdirChild(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		t.Chdir(".")
	})
}

func TestFunctionChdirChildrenCanBeParallel(t *testing.T) {
	t.Chdir(".")
	t.Run("1", func(t *testing.T) { // want "Function TestFunctionChdirChildrenCanBeParallel missing the call to method parallel in the test run"
		t.Log("1")
	})
	t.Run("2", func(t *testing.T) { // want "Function TestFunctionChdirChildrenCanBeParallel missing the call to method parallel in the test run"
		t.Log("2")
	})
}

func TestFunctionRunWithChdirSibling(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		t.Chdir(".")
	})
	t.Run("2", func(t *testing.T) { // want "Function TestFunctionRunWithChdirSibling missing the call to method parallel in the test run"
		t.Log("2")
	})
}

func TestFunctionWithChdirRange(t *testing.T) {
	for _, name := range []string{"1", "2"} {
		t.Run(name, func(q *testing.T) {
			q.Chdir(".")
		})
	}
}
