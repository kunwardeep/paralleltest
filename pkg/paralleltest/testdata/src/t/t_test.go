package t

import (
	"fmt"
	"testing"
)

// Non-test functions to verify linter doesn't flag them
func NoATestFunction()                                  {}
func TestingFunctionLooksLikeATestButIsNotWithParam()   {}
func TestingFunctionLooksLikeATestButIsWithParam(i int) {}
func AbcFunctionSuccessful(t *testing.T)                {}

// Basic test cases - successful parallel tests
func TestFunctionSuccessfulRangeTest(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
	}{{name: "foo"}}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(x *testing.T) {
			x.Parallel()
			fmt.Println(tc.name)
		})
	}
}

func TestFunctionSuccessfulNoRangeTest(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
	}{{name: "foo"}, {name: "bar"}}

	t.Run(testCases[0].name, func(t *testing.T) {
		t.Parallel()
		fmt.Println(testCases[0].name)
	})
	t.Run(testCases[1].name, func(t *testing.T) {
		t.Parallel()
		fmt.Println(testCases[1].name)
	})
}

// Basic test cases - missing parallel calls
func TestFunctionMissingCallToParallel(t *testing.T) {} // want "Function TestFunctionMissingCallToParallel missing the call to method parallel"

// Range-based test cases
func TestFunctionRangeMissingCallToParallel(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
	}{{name: "foo"}}

	// this range loop should be okay as it does not have test Run
	for _, tc := range testCases {
		fmt.Println(tc.name)
	}

	for _, tc := range testCases { // want "Range statement for test TestFunctionRangeMissingCallToParallel missing the call to method parallel in test Run"
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
		})
	}
}

func TestFunctionMissingCallToParallelAndRangeNotUsingRangeValueInTDotRun(t *testing.T) { // want "Function TestFunctionMissingCallToParallelAndRangeNotUsingRangeValueInTDotRun missing the call to method parallel"
	testCases := []struct {
		name string
	}{{name: "foo"}}

	for _, tc := range testCases { // want "Range statement for test TestFunctionMissingCallToParallelAndRangeNotUsingRangeValueInTDotRun missing the call to method parallel in test Run"
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
	for _, tc := range testCases { // want "Range statement for test TestFunctionRangeNotUsingRangeValueInTDotRun does not reinitialise the variable tc"
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

// Multiple test run cases
func TestFunctionTwoTestRunMissingCallToParallel(t *testing.T) {
	t.Parallel()
	t.Run("1", func(t *testing.T) { // want "Function TestFunctionTwoTestRunMissingCallToParallel missing the call to method parallel in the test run"
		fmt.Println("1")
	})
	t.Run("2", func(t *testing.T) { // want "Function TestFunctionTwoTestRunMissingCallToParallel missing the call to method parallel in the test run"
		fmt.Println("2")
	})
}

func TestFunctionFirstOneTestRunMissingCallToParallel(t *testing.T) {
	t.Parallel()
	t.Run("1", func(t *testing.T) { // want "Function TestFunctionFirstOneTestRunMissingCallToParallel missing the call to method parallel in the test run"
		fmt.Println("1")
	})
	t.Run("2", func(t *testing.T) {
		t.Parallel()
		fmt.Println("2")
	})
}

func TestFunctionSecondOneTestRunMissingCallToParallel(t *testing.T) {
	t.Parallel()
	t.Run("1", func(x *testing.T) {
		x.Parallel()
		fmt.Println("1")
	})
	t.Run("2", func(t *testing.T) { // want "Function TestFunctionSecondOneTestRunMissingCallToParallel missing the call to method parallel in the test run"
		fmt.Println("2")
	})
}

// Mock testing.T to verify linter doesn't flag non-testing.T methods
type notATest int

func (notATest) Run(args ...interface{}) {}
func (notATest) Parallel()               {}
func (notATest) Setenv(_, _ string)      {}

func TestFunctionWithRunLookalike(t *testing.T) {
	t.Parallel()
	var other notATest
	// These aren't t.Run, so they shouldn't be flagged as Run invocations missing calls to Parallel.
	other.Run(1, 1)
	other.Run(2, 2)
}

func TestFunctionWithParallelLookalike(t *testing.T) { // want "Function TestFunctionWithParallelLookalike missing the call to method parallel"
	var other notATest
	// This isn't t.Parallel, so it doesn't qualify as a call to Parallel.
	other.Parallel()
}

// Test cases with different parameter names
func TestFunctionWithOtherTestingVar(q *testing.T) {
	q.Parallel()
}

// Setenv-related test cases
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

func TestFunctionSetenvChildrenCanBeParallel(t *testing.T) {
	// unable to call t.Parallel with t.Setenv
	t.Setenv("foo", "bar")
	t.Run("1", func(t *testing.T) { // want "Function TestFunctionSetenvChildrenCanBeParallel missing the call to method parallel in the test run"
		fmt.Println("1")
	})
	t.Run("2", func(t *testing.T) { // want "Function TestFunctionSetenvChildrenCanBeParallel missing the call to method parallel in the test run"
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
	t.Run("2", func(t *testing.T) { // want "Function TestFunctionRunWithSetenvSibling missing the call to method parallel in the test run"
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

// Edge cases and error handling
func TestParalleltestEdgeCases(t *testing.T) {
	t.Parallel()
	t.Run("edge_cases", func(t *testing.T) {
		t.Parallel()

		// Test nil test case
		t.Run("nil_test", func(t *testing.T) {
			t.Parallel()
			// Test handling of nil test case
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic for nil test case")
					}
				}()
				var nilTest *testing.T
				nilTest.Parallel()
			}()
		})

		// Test empty test case
		t.Run("empty_test", func(t *testing.T) {
			t.Parallel()
			// Test handling of empty test case
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic for empty test case")
					}
				}()
				var emptyTest testing.T
				emptyTest.Parallel()
			}()
		})

		// Test invalid test case
		t.Run("invalid_test", func(t *testing.T) {
			t.Parallel()
			// Test handling of invalid test case
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic for invalid test case")
					}
				}()
				var invalidTest interface{}
				invalidTest.(*testing.T).Parallel()
			}()
		})
	})
}

// Helper function test cases
func TestFunctionCallToParallelWhereTestContextIsAFunction(t *testing.T) {
	t.Parallel()
	t.Run("1", foo) // want "Function TestFunctionCallToParallelWhereTestContextIsAFunction missing the call to method parallel in the test run"
	t.Run("2", bar)
}

func foo(t *testing.T) {
	fmt.Println("1")
}

func bar(t *testing.T) {
	t.Parallel()
	fmt.Println("2")
}

// Nested test run cases with helper functions
func TestNestedTestRunsWithHelpers(t *testing.T) {
	t.Parallel()
	t.Run("outer", func(t *testing.T) {
		t.Parallel()
		t.Run("inner1", nestedHelper1) // want "Function TestNestedTestRunsWithHelpers missing the call to method parallel in the test run"
		t.Run("inner2", nestedHelper2) // want "Function TestNestedTestRunsWithHelpers missing the call to method parallel in the test run"
	})
}

func nestedHelper1(t *testing.T) {
	fmt.Println("nested1")
}

func nestedHelper2(t *testing.T) {
	fmt.Println("nested2")
}

// Range-based test cases with helper functions
func TestRangeHelperWithDifferentParamNames(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
	}{
		{name: "case1"},
		{name: "case2"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			t.Run("sub1", rangeHelperWithCustomParam) // want "Function TestRangeHelperWithDifferentParamNames missing the call to method parallel in the test run"
			t.Run("sub2", rangeHelperWithAnotherParam) // want "Function TestRangeHelperWithDifferentParamNames missing the call to method parallel in the test run"
		})
	}
}

func rangeHelperWithCustomParam(testT *testing.T) {
	fmt.Println("range custom")
}

func rangeHelperWithAnotherParam(t *testing.T) {
	fmt.Println("range another")
}
