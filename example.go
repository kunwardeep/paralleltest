package main

import "log"

func myLog(format string, args ...interface{}) {
	const prefix = "[my] "
	log.Printf(prefix+format, args...)
}

func Add(x, y int) int {
	return x + y
}
