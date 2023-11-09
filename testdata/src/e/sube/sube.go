package sube

import "fmt"

var AError = fmt.Errorf("A")

func AFunc() error { // want "return errors: sube.AError"
	return AError
}
