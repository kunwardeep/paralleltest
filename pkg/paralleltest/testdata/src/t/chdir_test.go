package t

import (
	"fmt"
	"testing"
)

func TestFunctionWithChdir(t *testing.T) {
	// unable to call t.Parallel with t.Chdir
	t.Chdir("foo")
}

func TestFunctionWithChdirLookalike(t *testing.T) { // want "Function TestFunctionWithChdirLookalike missing the call to method parallel"
	var other notATest
	other.Chdir("foo")
}

func TestFunctionWithChdirChild(t *testing.T) {
	// ancestor of Chdir cant call t.Parallel
	t.Run("1", func(t *testing.T) {
		// unable to call t.Parallel with t.Chdir
		t.Chdir("foo")
		fmt.Println("1")
	})
}

func TestFunctionWithChdirBuilder(t *testing.T) {
	t.Run("1", builderWithChdir())
}

func builderWithChdir() func(t *testing.T) {
	return func(t *testing.T) {
		// unable to call t.Parallel with t.Chdir
		t.Chdir("foo")
	}
}

func TestFunctionWithChdirFunc(t *testing.T) {
	t.Run("1", scenarioWithChdir)
}

func scenarioWithChdir(t *testing.T) {
	// unable to call t.Parallel with t.Chdir
	t.Chdir("foo")
}

func TestFunctionChdirChildrenCanBeParallel(t *testing.T) {
	// unable to call t.Parallel with t.Chdir
	t.Chdir("foo")
	t.Run("1", func(t *testing.T) { // want "Function literal missing the call to method parallel in the t.Run\n"
		fmt.Println("1")
	})
	t.Run("2", func(t *testing.T) { // want "Function literal missing the call to method parallel in the t.Run\n"
		fmt.Println("2")
	})
}

func TestFunctionRunWithChdirSibling(t *testing.T) {
	// ancestor of Chdir cant call t.Parallel
	t.Run("1", func(t *testing.T) {
		// unable to call t.Parallel with t.Chdir
		t.Chdir("foo")
		fmt.Println("1")
	})
	t.Run("2", func(t *testing.T) { // want "Function literal missing the call to method parallel in the t.Run\n"
		fmt.Println("2")
	})
}

func TestFunctionWithChdirRange(t *testing.T) {
	// ancestor of Chdir cant call t.Parallel
	testCases := []struct {
		name string
	}{{name: "foo"}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// unable to call t.Parallel with t.Chdir
			t.Chdir("foo")
		})
	}
}

func ChdirHelper(t *testing.T) {
	t.Chdir("foo")
}

func TestFunctionWithHelperChdir(t *testing.T) {
	ChdirHelper(t)
}

func TestFunctionWithNestedHelperChdir(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		ChdirHelper(t)
	})
}

func TestFunctionWithNestedSiblingHelperChdir(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		ChdirHelper(t)
	})
	t.Run("2", func(t *testing.T) {
		t.Parallel()
	})
}

// ChdirHelper prevents the caller but not the returning literal from being parallel
func builderWithChdirHelper(t *testing.T) func(t *testing.T) {
	ChdirHelper(t)
	return func(t *testing.T) {
		t.Parallel()
	}
}

func TestFunctionWithBuilderHelperChdir(t *testing.T) {
	t.Run("1", builderWithChdirHelper(t))
}
