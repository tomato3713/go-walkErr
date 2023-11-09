package a

import "fmt"

var AError = fmt.Errorf("A")

func afunc() error { // want "return errors: a.AError"
	return AError
}
