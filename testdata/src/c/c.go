package main

import "fmt"

var AError = fmt.Errorf("A")
var BError = fmt.Errorf("B")

func afunc(b int) error { // want "return errors: main.AError, main.BError"
	if b == 0 {
		return AError
		if b == 1 {
			return BError
		}
	}
	return BError
}
