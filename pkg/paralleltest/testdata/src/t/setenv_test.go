package t

import (
	"fmt"
	"testing"
)

func TestFunctionWithSetenv(t *testing.T) {
	// unable to call t.Parallel with t.Setenv
	t.Setenv("foo", "bar")
}

func TestFunctionWithSetenvLookalike(t *testing.T) { // want "Function TestFunctionWithSetenvLookalike missing the call to method parallel"
	var other notATest
	other.Setenv("foo", "bar")
}

func TestFunctionWithSetenvChild(t *testing.T) {
	// ancestor of setenv cant call t.Parallel
	t.Run("1", func(t *testing.T) {
		// unable to call t.Parallel with t.Setenv
		t.Setenv("foo", "bar")
		fmt.Println("1")
	})
}

func TestFunctionWithSetenvBuilder(t *testing.T) {
	t.Run("1", builderWithSetenv())
}

func builderWithSetenv() func(t *testing.T) {
	return func(t *testing.T) {
		// unable to call t.Parallel with t.Setenv
		t.Setenv("foo", "bar")
	}
}

func TestFunctionWithSetenvFunc(t *testing.T) {
	t.Run("1", scenarioWithSetenv)
}

func scenarioWithSetenv(t *testing.T) {
	// unable to call t.Parallel with t.Setenv
	t.Setenv("foo", "bar")
}

func TestFunctionSetenvChildrenCanBeParallel(t *testing.T) {
	// unable to call t.Parallel with t.Setenv
	t.Setenv("foo", "bar")
	t.Run("1", func(t *testing.T) { // want "Function literal missing the call to method parallel in the t.Run\n"
		fmt.Println("1")
	})
	t.Run("2", func(t *testing.T) { // want "Function literal missing the call to method parallel in the t.Run\n"
		fmt.Println("2")
	})
}

func TestFunctionRunWithSetenvSibling(t *testing.T) {
	// ancestor of setenv cant call t.Parallel
	t.Run("1", func(t *testing.T) {
		// unable to call t.Parallel with t.Setenv
		t.Setenv("foo", "bar")
		fmt.Println("1")
	})
	t.Run("2", func(t *testing.T) { // want "Function literal missing the call to method parallel in the t.Run\n"
		fmt.Println("2")
	})
}

func TestFunctionWithSetenvRange(t *testing.T) {
	// ancestor of setenv cant call t.Parallel
	testCases := []struct {
		name string
	}{{name: "foo"}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// unable to call t.Parallel with t.Setenv
			t.Setenv("foo", "bar")
		})
	}
}

func setEnvHelper(t *testing.T) {
	t.Setenv("foo", "bar")
}

func TestFunctionWithHelperSetenv(t *testing.T) {
	setEnvHelper(t)
}

func TestFunctionWithNestedHelperSetenv(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		setEnvHelper(t)
	})
}

func TestFunctionWithNestedSiblingHelperSetenv(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		setEnvHelper(t)
	})
	t.Run("2", func(t *testing.T) {
		t.Parallel()
	})
}

// setEnvHelper prevents the caller but not the returning literal from being parallel
func builderWithSetenvHelper(t *testing.T) func(t *testing.T) {
	setEnvHelper(t)
	return func(t *testing.T) {
		t.Parallel()
	}
}

func TestFunctionWithBuilderHelperSetenv(t *testing.T) {
	t.Run("1", builderWithSetenvHelper(t))
}
