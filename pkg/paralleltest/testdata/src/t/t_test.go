package t

import (
	"fmt"
	"testing"
)

func notTestFunction()      {}
func notTestFunctionWithParam(i int) { }
func TestFunctionSuccessful(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
	}{{name: "foo"}}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fmt.Println(tc.name)
		})
	}
}

func TestFunctionMissingCallToParallel(t *testing.T) {} // want "Function TestFunctionMissingCallToParallel missing the call to method parallel"
func TestFunctionRangeMissingCallToParallel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
	}{{name: "foo"}}

	// this range loop should be okay as it does not have t.Run
	for _, tc := range testCases {
		fmt.Println(tc.name)
	}

	for _, tc := range testCases { // want "Range statement for test TestFunctionRangeMissingCallToParallel missing the call to method parallel in t.Run"
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
		})
	}
}

func TestFunctionMissingCallToParallelAndRangeNotUsingRangeValueInTDotRun(t *testing.T) {// want "Function TestFunctionMissingCallToParallelAndRangeNotUsingRangeValueInTDotRun missing the call to method parallel"
	testCases := []struct {
		name string
	}{{name: "foo"}}

	for _, tc := range testCases { // want "Range statement for test TestFunctionMissingCallToParallelAndRangeNotUsingRangeValueInTDotRun missing the call to method parallel in t.Run"
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
		})
	}
}

func TestFunctionRangeNotUsingRangeValueInTDotRun(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
	}{{name: "foo"}}
	for _, tc := range testCases { // want "Range statement for test TestFunctionRangeNotUsingRangeValueInTDotRun does not use range value in t.Run"
		t.Run("tc.name", func(t *testing.T) {
			t.Parallel()
			fmt.Println(tc.name)
		})
	}
}

func TestFunctionRangeNotReInitialisingVariable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
	}{{name: "foo"}}
	for _, tc := range testCases { // want "Range statement for test TestFunctionRangeNotReInitialisingVariable does not reinitialise the variable tc"
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fmt.Println(tc.name)
		})
	}
}

// TODO this test should fail when the missing functionality is implemented as t.Run does not call t.Parallel
func TestFunctionSuccessful_maybe(t *testing.T) {
	t.Parallel()

	t.Run("1", func(t *testing.T) {
		fmt.Println("1")
	})
	t.Run("2", func(t *testing.T) {
		fmt.Println("2")
	})
}

