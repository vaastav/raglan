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

func isIridescentAnnotationCallExpr(expr *ast.CallExpr) bool {
	if ident, ok := expr.Fun.(*ast.Ident); ok && strings.HasPrefix(ident.Name, "iridescent") {
		return true
	}
	return false
}

func (srt *SpecializationRuntime) setupSpecializedModule(filename string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if _, ok2 := srt.GlobalFns[fn.Name.Name]; ok2 {
				var newStatements []ast.Stmt
				for _, stmt := range fn.Body.List {
					switch s := stmt.(type) {
					case *ast.AssignStmt:
						var keepStmt = true
						for _, rhs := range s.Rhs {
							if callExpr, ok := rhs.(*ast.CallExpr); ok {
								if isIridescentAnnotationCallExpr(callExpr) {
									keepStmt = false
								}
							}
						}
						if keepStmt {
							newStatements = append(newStatements, stmt)
						}
					case *ast.ExprStmt:
						var keepStmt = true
						if callExpr, ok := s.X.(*ast.CallExpr); ok {
							if isIridescentAnnotationCallExpr(callExpr) {
								keepStmt = false
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

	for _, pt := range srt.Pts {
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
			if _, ok2 := srt.GlobalFns[fn.Name.Name]; ok2 {
				fn.Name = ast.NewIdent(fn.Name.Name + "_Specialized")
			}
		}
		return true
	})

	// Apply custom specialization passes
	for _, pass := range srt.Passes {
		file = pass.Modify(fset, file)
	}

	new_file := strings.ReplaceAll(filename, ".go", "_specialized.go")
	f, err := os.Create(new_file)
	defer f.Close()
	err = printer.Fprint(f, fset, file)
	return new_file, nil
}
