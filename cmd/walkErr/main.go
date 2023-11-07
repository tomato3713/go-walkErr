package main

import (
	"github.com/tomato3713/walkErr"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(walkErr.Analyzer)
}
