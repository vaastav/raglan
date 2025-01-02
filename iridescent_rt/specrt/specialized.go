package specrt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

func setupSpecializedModule(filename string, global_fns map[string]bool, spec_points []*CompileTimeSpecPoint[any]) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return "", err
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if _, ok2 := global_fns[fn.Name.Name]; ok2 {
				var newStatements []ast.Stmt
				for _, stmt := range fn.Body.List {
					switch s := stmt.(type) {
					case *ast.AssignStmt:
						var keepStmt = true
						for _, rhs := range s.Rhs {
							if callExpr, ok := rhs.(*ast.CallExpr); ok {
								if ident, ok := callExpr.Fun.(*ast.Ident); ok && strings.HasPrefix(ident.Name, "iridescent") {
									keepStmt = false
									break
								}
							}
						}
						if keepStmt {
							newStatements = append(newStatements, stmt)
						}
					default:
						newStatements = append(newStatements, stmt)
					}
				}
				fn.Body.List = newStatements
			}
		}
		return true
	})

	for _, pt := range spec_points {
		if pt.IsSpecialized {
			ast.Inspect(file, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == pt.ParentFn {
					ast.Inspect(fn.Body, func(n ast.Node) bool {
						if ident, ok := n.(*ast.Ident); ok {
							if ident.Name == pt.Name {
								ident.Name = fmt.Sprintf("%v", pt.Current)
							}
						}
						return true
					})
				}
				return true
			})
		}
	}
	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if _, ok2 := global_fns[fn.Name.Name]; ok2 {
				fn.Name = ast.NewIdent(fn.Name.Name + "_Specialized")
			}
		}
		return true
	})

	new_file := strings.ReplaceAll(filename, ".go", "_specialized.go")
	f, err := os.Create(new_file)
	defer f.Close()
	err = printer.Fprint(f, fset, file)
	return new_file, nil
}
