package main

import "fmt"

var AError = fmt.Errorf("A")

func afunc() error { // want "return errors: main.AError"
	return AError
}
