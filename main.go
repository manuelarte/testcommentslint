package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/manuelarte/testcommentslint/analyzer"
)

func main() {
	singlechecker.Main(analyzer.New())
}
