package specrt

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

func setupOriginalModule(filename string, global_fns map[string]bool) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return "", err
	}
	// Modify the AST to delete function calls to iridescent instrumentation function
	// Currently, these are replaced by the identifier itself
	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if _, ok2 := global_fns[fn.Name.Name]; ok2 {
				fn.Name = ast.NewIdent(fn.Name.Name + "_Original")
			}
		}
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			// Iterate over the right-hand side of the assignment
			for i, rhsExpr := range assignStmt.Rhs {
				if callExpr, ok := rhsExpr.(*ast.CallExpr); ok {
					if ident, ok := callExpr.Fun.(*ast.Ident); ok && strings.HasPrefix(ident.Name, "iridescent") {
						if len(assignStmt.Lhs) > i {
							if lhsIdent, ok := assignStmt.Lhs[i].(*ast.Ident); ok {
								assignStmt.Rhs[i] = &ast.Ident{Name: lhsIdent.Name}
							}
						}
					}
				}
			}
		}
		return true
	})
	new_file := strings.ReplaceAll(filename, ".go", "_original.go")
	f, err := os.Create(new_file)
	defer f.Close()
	err = printer.Fprint(f, fset, file)
	return new_file, nil
}
