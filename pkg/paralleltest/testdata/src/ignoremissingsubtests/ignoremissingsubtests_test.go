package t

import (
	"fmt"
	"testing"
)

func NoATestFunction()                                  {}
func TestingFunctionLooksLikeATestButIsNotWithParam()   {}
func TestingFunctionLooksLikeATestButIsWithParam(i int) {}
func AbcFunctionSuccessful(t *testing.T)                {}

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

func TestFunctionMissingCallToParallel(t *testing.T) {} // want "Function TestFunctionMissingCallToParallel missing the call to method parallel"
func TestFunctionRangeMissingCallToParallel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
	}{{name: "foo"}}

	// this range loop should be okay as it does not have test Run
	for _, tc := range testCases {
		fmt.Println(tc.name)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
		})
	}
}

func TestFunctionMissingCallToParallelAndRangeNotUsingRangeValueInTDotRun(t *testing.T) { // want "Function TestFunctionMissingCallToParallelAndRangeNotUsingRangeValueInTDotRun missing the call to method parallel"
	testCases := []struct {
		name string
	}{{name: "foo"}}

	for _, tc := range testCases {
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

func TestFunctionTwoTestRunMissingCallToParallel(t *testing.T) {
	t.Parallel()

	t.Run("1", func(t *testing.T) {
		fmt.Println("1")
	})
	t.Run("2", func(t *testing.T) {
		fmt.Println("2")
	})
}

func TestFunctionFirstOneTestRunMissingCallToParallel(t *testing.T) {
	t.Parallel()

	t.Run("1", func(t *testing.T) {
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
	t.Run("2", func(t *testing.T) {
		fmt.Println("2")
	})
}

type notATest int

func (notATest) Run(args ...interface{}) {}
func (notATest) Parallel()               {}

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

func TestFunctionWithOtherTestingVar(q *testing.T) {
	q.Parallel()
}
