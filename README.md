# paralleltest

[![Test](https://github.com/kunwardeep/paralleltest/actions/workflows/test.yml/badge.svg)](https://github.com/kunwardeep/paralleltest/actions/workflows/test.yml)

The Go linter `paralleltest` checks that the t.Parallel gets called for the test method and for the range of test cases within the test.

## Usage

```sh
paralleltest ./...
```

A few options can be activated by flag:

* `-i`: Ignore missing calls to `t.Parallel` and only report incorrect uses of it.
* `-ignoremissingsubtests`: Require that top-level tests specify `t.Parallel`, but don't require it in subtests (`t.Run(...)`).
* `-ignoreloopVar`: Ignore loop variable detection.

## Examples

### Missing `t.Parallel()` in the test method

```go
// bad
func TestFunctionMissingCallToParallel(t *testing.T) {
}

// good
func TestFunctionMissingCallToParallel(t *testing.T) {
  t.Parallel()
  // ^ call to t.Parallel()
}
// Error displayed
// Function TestFunctionMissingCallToParallel missing the call to method parallel
```

### Missing `t.Parallel()` in the range method

```go
// bad
func TestFunctionRangeMissingCallToParallel(t *testing.T) {
  t.Parallel()

  testCases := []struct {
    name string
  }{{name: "foo"}}

  for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
      fmt.Println(tc.name)
    })
  }
}

// good
func TestFunctionRangeMissingCallToParallel(t *testing.T) {
  t.Parallel()

  testCases := []struct {
    name string
  }{{name: "foo"}}

  for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
      t.Parallel()
      // ^ call to t.Parallel()
      fmt.Println(tc.name)
    })
  }
}
// Error displayed
// Range statement for test TestFunctionRangeMissingCallToParallel missing the call to method parallel in t.Run
```

### `t.Parallel()` is called in the range method but testcase variable not being used

```go
// bad
func TestFunctionRangeNotUsingRangeValueInTDotRun(t *testing.T) {
  t.Parallel()

  testCases := []struct {
    name string
  }{{name: "foo"}}
  for _, tc := range testCases {
    t.Run("this is a test name", func(t *testing.T) {
      // ^ call to tc.name missing
      t.Parallel()
      fmt.Println(tc.name)
    })
  }
}

// good
func TestFunctionRangeNotUsingRangeValueInTDotRun(t *testing.T) {
  t.Parallel()

  testCases := []struct {
    name string
  }{{name: "foo"}}
  for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
      // ^ call to tc.name
      t.Parallel()
      fmt.Println(tc.name)
    })
  }
}
// Error displayed
// Range statement for test TestFunctionRangeNotUsingRangeValueInTDotRun does not use range value in t.Run
```

### `t.Parallel()` is called in the range method and test case variable tc being used, but is not reinitialised (<a href="https://gist.github.com/kunwardeep/80c2e9f3d3256c894898bae82d9f75d0" target="_blank">More Info</a>)
```go
// bad
func TestFunctionRangeNotReInitialisingVariable(t *testing.T) {
  t.Parallel()

  testCases := []struct {
    name string
  }{{name: "foo"}}
  for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
      t.Parallel()
      fmt.Println(tc.name)
    })
  }
}

// good
func TestFunctionRangeNotReInitialisingVariable(t *testing.T) {
  t.Parallel()

  testCases := []struct {
    name string
  }{{name: "foo"}}
  for _, tc := range testCases {
    tc:=tc
    // ^ tc variable reinitialised
    t.Run(tc.name, func(t *testing.T) {
      t.Parallel()
      fmt.Println(tc.name)
    })
  }
}
// Error displayed
// Range statement for test TestFunctionRangeNotReInitialisingVariable does not reinitialise the variable tc
```
