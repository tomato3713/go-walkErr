package b

import "fmt"

var AError = fmt.Errorf("A")
var BError = fmt.Errorf("B")

func afunc() (int, error) { return 0, AError } // want "return errors: b.AError"

func bfunc() int { return 0 }
