package main

import "fmt"

var AError = fmt.Errorf("A")
var BError = fmt.Errorf("B")

func afunc() error { // want "return errors: main.AError"
	return AError
}

func bfunc(b int) error { // want "return errors: main.AError, main.BError"
	if b == 0 {
		return afunc()
	}
	return BError
}
