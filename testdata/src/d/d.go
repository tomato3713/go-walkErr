package d

import "fmt"

var AError = fmt.Errorf("A")
var BError = fmt.Errorf("B")

func afunc() error { // want "return errors: d.AError"
	return AError
}

func bfunc(b int) error { // want "return errors: d.AError, d.BError"
	if b == 0 {
		return afunc()
	}
	return BError
}
