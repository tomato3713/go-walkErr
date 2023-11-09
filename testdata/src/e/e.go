package e

import (
	"e/sube"
	"fmt"
)

var AError = fmt.Errorf("A")

func afunc() error { // want "return errors: e.AError"
	return AError
}

func bfunc(b int) error { // want "return errors: e.AError, sube.AError"
	if b == 0 {
		return afunc()
	}
	return sube.AFunc()
}
