package p

func notPrintfFuncAtAll() {}

func funcWithEllipsis(args ...interface{}) {}

func printfLikeButWithStrings(format string, args ...string) {}

func printfLikeButWithBadFormat(format int, args ...interface{}) {}

func secondArgIsNotEllipsis(format string, arg int) {}

func printfLikeButWithExtraInterfaceMethods(format string, args ...interface {
	String() string
}) {
}

func prinfLikeFuncf(format string, args ...interface{}) {}

func prinfLikeFuncWithReturnValue(format string, args ...interface{}) string {
	return ""
}

func prinfLikeFuncWithAnotherFormatArgName(msg string, args ...interface{}) {}

func prinfLikeFunc(format string, args ...interface{}) {} // want "printf-like formatting function"

func prinfLikeFuncWithExtraArgs1(extraArg, format string, args ...interface{}) {} // want "printf-like formatting function"

func prinfLikeFuncWithExtraArgs2(extraArg int, format string, args ...interface{}) {} // want "printf-like formatting function"
