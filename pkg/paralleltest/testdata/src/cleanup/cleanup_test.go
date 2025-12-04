package cleanup

import (
	"fmt"
	"os"
	"testing"
)

// Test with t.Parallel and defer - should report an issue when checkcleanup flag is enabled
func TestWithParallelAndDefer(t *testing.T) {
	t.Parallel()

	tempFile := "test.tmp"
	f, _ := os.Create(tempFile)
	defer os.Remove(tempFile) // want "Function TestWithParallelAndDefer uses defer with t.Parallel, use t.Cleanup instead to ensure cleanup runs after parallel subtests complete"
	defer f.Close()           // want "Function TestWithParallelAndDefer uses defer with t.Parallel, use t.Cleanup instead to ensure cleanup runs after parallel subtests complete"

	t.Run("subtest", func(t *testing.T) {
		t.Parallel()
		fmt.Fprintf(f, "test data\n")
	})
}

// Test with t.Parallel and t.Cleanup - should be fine
func TestWithParallelAndCleanup(t *testing.T) {
	t.Parallel()

	tempFile := "test.tmp"
	f, _ := os.Create(tempFile)
	t.Cleanup(func() {
		f.Close()
		os.Remove(tempFile)
	})

	t.Run("subtest", func(t *testing.T) {
		t.Parallel()
		fmt.Fprintf(f, "test data\n")
	})
}

// Test without t.Parallel but with defer - should only report missing parallel, not defer issue
func TestWithoutParallelButWithDefer(t *testing.T) { // want "Function TestWithoutParallelButWithDefer missing the call to method parallel"
	tempFile := "test.tmp"
	f, _ := os.Create(tempFile)
	defer os.Remove(tempFile)
	defer f.Close()

	fmt.Fprintf(f, "test data\n")
}

// Test with t.Parallel but no defer - should be fine
func TestWithParallelButNoDefer(t *testing.T) {
	t.Parallel()

	fmt.Println("test")
}

// Test with Setenv (can't parallel) but has defer - should be fine (Setenv prevents parallel)
func TestWithSetenvAndDefer(t *testing.T) {
	t.Setenv("TEST_VAR", "value")
	defer fmt.Println("cleanup")

	fmt.Println("test")
}

func TestWithMultipleDefersAndParallelButNoSubtests(t *testing.T) {
	t.Parallel()

	defer fmt.Println("cleanup 1")
	defer fmt.Println("cleanup 2")
	defer fmt.Println("cleanup 3")

	fmt.Println("test")
}

// Test demonstrating the issue: defer runs before subtests complete with t.Parallel
func TestDemonstratingProblem(t *testing.T) {
	t.Parallel()

	counter := 0
	defer func() { // want "Function TestDemonstratingProblem uses defer with t.Parallel, use t.Cleanup instead to ensure cleanup runs after parallel subtests complete\n"
		// This runs immediately when the test function returns,
		// BEFORE subtests complete!
		fmt.Printf("Counter value in defer: %d\n", counter)
	}()

	t.Run("subtest1", func(t *testing.T) {
		t.Parallel()
		counter++
	})

	t.Run("subtest2", func(t *testing.T) {
		t.Parallel()
		counter++
	})
	// Function returns here, defer runs, but subtests are still running!
}

// Test showing correct usage with t.Cleanup
func TestCorrectUsageWithCleanup(t *testing.T) {
	t.Parallel()

	counter := 0
	t.Cleanup(func() {
		// This runs AFTER all subtests complete
		fmt.Printf("Counter value in cleanup: %d\n", counter)
	})

	t.Run("subtest1", func(t *testing.T) {
		t.Parallel()
		counter++
	})

	t.Run("subtest2", func(t *testing.T) {
		t.Parallel()
		counter++
	})
	// t.Cleanup runs after all subtests finish
}

func TestNestedDefer(t *testing.T) {
	t.Parallel()
	t.Run("1", func(t *testing.T) {
		t.Parallel()
		defer fmt.Println("cleanup 1") // want "Function literal uses defer with t.Parallel, use t.Cleanup instead to ensure cleanup runs after parallel subtests complete\n"
		t.Run("2", func(t *testing.T) {
			t.Parallel()
			defer fmt.Println("cleanup 2") // okay if there's no nesting.
		})
	})
}

func helperDefer(t *testing.T) {
	defer fmt.Println("cleanup 2")
}

func TestDeferOkayInHelper(t *testing.T) {
	t.Parallel()
	helperDefer(t)
	t.Run("1", func(t *testing.T) {
		t.Parallel()
	})
}
