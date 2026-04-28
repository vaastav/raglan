package specrt

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

type ParseSpecPointVisitor struct {
	fset  *token.FileSet
	curFn string
	Pts   []*CompileTimeSpecPoint[any]
	Fns   map[string]bool
}

func (v *ParseSpecPointVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch x := n.(type) {
	case *ast.FuncDecl:
		// Start of function scope!
		name := x.Name.Name
		v.curFn = name
	case *ast.CallExpr:
		switch callee := x.Fun.(type) {
		case *ast.Ident:
			callee_name := callee.Name
			if strings.HasPrefix(callee_name, "iridescent") {
				// Only parse Iridescent function calls
				pieces := strings.Split(callee_name, "_")
				pt_irid_type := pieces[2]
				pt_go_type := pieces[3]
				arg := x.Args[0]
				arg_lit := arg.(*ast.BasicLit)
				arg_name := arg_lit.Value
				arg_name = arg_name[1 : len(arg_name)-1]
				values := []any{}
				if pt_irid_type == "general" {
					if len(x.Args) > 2 {
						for i := 2; i < len(x.Args); i = i + 1 {
							arg := x.Args[i]
							arg_lit := arg.(*ast.BasicLit)
							arg_val := arg_lit.Value
							values = append(values, arg_val)
						}
					}
				} else if pt_irid_type == "range" {
					start_arg := x.Args[1]
					sa_lit := start_arg.(*ast.BasicLit)
					sa_val := sa_lit.Value
					start, _ := strconv.ParseInt(sa_val, 10, 64)
					end_arg := x.Args[2]
					ea_lit := end_arg.(*ast.BasicLit)
					ea_val := ea_lit.Value
					end, _ := strconv.ParseInt(ea_val, 10, 64)
					for i := start; i <= end; i++ {
						values = append(values, i)
					}
				}
				pt := NewCompileTimeSpecPoint[any](arg_name, values)
				pt.ParentFn = v.curFn
				pt.IridType = pt_irid_type
				pt.GoType = pt_go_type
				v.Pts = append(v.Pts, pt)
				v.Fns[v.curFn] = true
			}
		default:
		}
	default:
	}

	return v
}

func parseOriginalModule(filename string) ([]*CompileTimeSpecPoint[any], error) {
	var points []*CompileTimeSpecPoint[any]
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return points, err
	}
	v := &ParseSpecPointVisitor{fset: fset, Fns: make(map[string]bool)}
	ast.Walk(v, file)
	points = v.Pts
	return points, nil
}
