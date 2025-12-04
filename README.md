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

You can configure paralleltest by using a .paralleltest.yaml file in either the local directory or in your home directory.
```
# IgnoreMissing check that missing calls to t.Parallel are not reported, default false
ignoreMissing: false
# IgnoreMissingSubtests check that missing calls to t.Parallel in subtests are not reported, default true
ignoreMissingSubtests: false
# CheckCleanup check that defer is not used with t.Parallel (use t.Cleanup instead), default false
checkClean: false
# ExtraSigs is a list of extra functions that cannot be used with t.Parallel, default []
extraSigs:
    - .CantBeParallel
    - .IgnoreParallel
```
With Go 1.22, we no longer need to check usage loop variables for `t.Parallel` calls.
## Development

### Prerequisites

- Go 1.24.0 or later
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
// Function literal missing the call to method parallel in the t.Run
```

### Using `defer` with `t.Parallel()` (requires `-checkcleanup` flag)

When `t.Parallel()` is called, the test function returns immediately, causing `defer` statements to execute before subtests complete. This can lead to cleanup happening too early, causing test failures or resource leaks.

```go
// bad - defer runs before subtests complete
func TestWithDeferAndParallel(t *testing.T) {
  t.Parallel()

  tempFile := "test.tmp"
  f, _ := os.Create(tempFile)
  defer f.Close()         // runs immediately when function returns
  defer os.Remove(tempFile)  // file deleted before subtests finish!

  t.Run("subtest1", func(t *testing.T) {
    t.Parallel()
    fmt.Fprintf(f, "data")  // may fail - file already closed!
  })
}

// good - t.Cleanup runs after all subtests complete
func TestWithCleanupAndParallel(t *testing.T) {
  t.Parallel()

  tempFile := "test.tmp"
  f, _ := os.Create(tempFile)
  t.Cleanup(func() {
    f.Close()
    os.Remove(tempFile)
  })

  t.Run("subtest1", func(t *testing.T) {
    t.Parallel()
    fmt.Fprintf(f, "data")  // works correctly
  })
}
// Error displayed (with -checkcleanup flag)
// Function TestWithDeferAndParallel uses defer with t.Parallel, use t.Cleanup instead to ensure cleanup runs after parallel subtests complete
```

**Note:** This check is disabled by default. Enable it with the `-checkcleanup` flag.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
