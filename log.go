package main

import "fmt"

func log(str string, args ...interface{}) {
	if verbose {
		fmt.Printf(str, args...)
	}
}
