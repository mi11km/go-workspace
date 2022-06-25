package complexity

import "go/ast"

func Count(node ast.Node) int {
	if node == nil {
		return 0
	}
	count := 1
	ast.Inspect(node, func(node ast.Node) bool {
		switch node.(type) {
		case *ast.IfStmt:
			count++
		case *ast.ForStmt:
			count++
		}
		return true
	})
	return count
}
