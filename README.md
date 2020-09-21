# go-printf-func-name

The Go linter `go-printf-func-name` checks that printf-like functions are named with `f` at the end.

For example, `myLog` should be named `myLogf` by Go convention:

```go
package main

import "log"

func myLog(format string, args ...interface{}) {
	const prefix = "[my] "
	log.Printf(prefix + format, args...)
}
```
