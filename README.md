# paralleltest

[![Test](https://github.com/kunwardeep/paralleltest/actions/workflows/test.yml/badge.svg)](https://github.com/kunwardeep/paralleltest/actions/workflows/test.yml)

The Go linter `paralleltest` checks that the t.Parallel gets called for the test method and for the range of test cases within the test.

## Installation

```sh
go install github.com/kunwardeep/paralleltest@latest
```

## Usage

```sh
paralleltest ./...
```

A few options can be activated by flag:

* `-i`: Ignore missing calls to `t.Parallel` and only report incorrect uses of it.
* `-ignoremissingsubtests`: Require that top-level tests specify `t.Parallel`, but don't require it in subtests (`t.Run(...)`).
* `-ignoreloopVar`: Ignore loop variable detection.

## Development

### Prerequisites

- Go 1.23.0 or later
- Make

### Local Development

1. Clone the repository:
```sh
git clone https://github.com/kunwardeep/paralleltest.git
cd paralleltest
```

2. Install development tools:
```sh
make install_devtools
```

3. Install dependencies:
```sh
make ensure_deps
```

4. Run tests:
```sh
make test
```

5. Run linter:
```sh
make lint
```

To fix linting issues automatically:
```sh
make lint_fix
```

### CI/CD

The project uses GitHub Actions for continuous integration. The workflow includes:

- Running tests with race condition detection
- Running golangci-lint for code quality checks

The workflow runs on:
- Pull requests
- Pushes to the main branch

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

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
