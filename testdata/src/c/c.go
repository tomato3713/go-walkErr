package c

import "fmt"

var AError = fmt.Errorf("A")
var BError = fmt.Errorf("B")

func afunc(b int) error { // want "return errors: c.AError, c.BError"
	if b == 0 {
		return AError
		if b == 1 {
			return BError
		}
	}
	return BError
}
