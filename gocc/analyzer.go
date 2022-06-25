package gocc

import (
	"go/ast"

	"github.com/mi11km/go-workspace/gocc/complexity"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "gcc",
	Doc:      "checks cyclomatic complexity",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	inspt := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	inspt.Preorder(nodeFilter, func(node ast.Node) {
		count := complexity.Count(node)
		if count >= 10 {
			fd := node.(*ast.FuncDecl)
			pass.Reportf(node.Pos(), "function %s complexity=%d", fd.Name.Name, count)
		}
	})
	return nil, nil
}
