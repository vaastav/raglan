package specrt

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"
)

func setupOriginalModule(filename string, global_fns map[string]bool) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	var newDecls []ast.Decl
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if _, ok2 := global_fns[fn.Name.Name]; ok2 {
				newDecls = append(newDecls, decl)
			}
		} else {
			newDecls = append(newDecls, decl)
		}
	}
	file.Decls = newDecls
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
					if isIridescentAnnotationCallExpr(callExpr) {
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
	astutil.Apply(file, func(c *astutil.Cursor) bool {
		node := c.Node()

		exprStmt, ok := node.(*ast.ExprStmt)
		if !ok {
			return true
		}

		call, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isIridescentAnnotationCallExpr(call) {
			c.Delete()
		}

		return true
	}, nil)
	new_file := strings.ReplaceAll(filename, ".go", "_original.go")

	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, file)
	if err != nil {
		return "", err
	}

	// Run imports logic
	out, err := imports.Process(new_file, buf.Bytes(), nil)
	if err != nil {
		return "", err
	}

	f, err := os.Create(new_file)
	defer f.Close()
	_, err = f.Write(out)
	return new_file, nil
}
