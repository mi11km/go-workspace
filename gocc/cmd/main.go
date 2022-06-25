package main

import (
	"github.com/mi11km/go-workspace/gocc"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(gocc.Analyzer)
}
