package complexity

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestComplexity(t *testing.T) {
	testcases := []struct {
		name       string
		code       string
		complexity int
	}{
		{
			name: "simple function",
			code: `
package main
func Double(n int) int {
	return n * 2
}`,
			complexity: 1,
		},
		{
			name: "if statement",
			code: `
package main
func Double(n int) int {
	if n%2 == 0 {
		return 0
	}
	return n
}`,
			complexity: 2,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			a := GetAST(t, testcase.code)
			c := Count(a)
			if c != testcase.complexity {
				t.Errorf("got=%d, want=%d", c, testcase.complexity)
			}
		})
	}
}

func GetAST(t *testing.T, code string) ast.Node {
	t.Helper()
	fset := token.NewFileSet()
	n, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		t.Log(err)
		return nil
	}
	return n
}
