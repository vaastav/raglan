package specrt

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
)

func setupTrampolineModule(filename string, global_fns map[string]bool) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return "", err
	}

	var trampoline_fns []*ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			// Only trampoline functions that are global functions
			if _, ok2 := global_fns[fn.Name.Name]; ok2 {
				var args []string
				var rets []string
				for _, arg := range fn.Type.Params.List {
					args = append(args, arg.Names[0].Name)
				}
				for i := 0; i < len(fn.Type.Results.List)-1; i = i + 1 {
					rets = append(rets, fmt.Sprintf("ret%d", i))
				}

				arg_string := strings.Join(args, ",")
				ret_string := strings.Join(rets, ",")
				executor := newTemplateExecutor()
				tempargs := templateArgs{
					ArgString: arg_string,
					RetString: ret_string,
					Name:      fn.Name.Name,
				}
				new_body, err := executor.exec(fn.Name.Name, tempargs)
				if err != nil {
					log.Fatal(err)
				}
				new_body_fset := token.NewFileSet()
				new_body_node, err := parser.ParseFile(new_body_fset, "trampoline.go", new_body, 0)
				if err != nil {
					log.Println(new_body)
					log.Fatal(err)
				}
				var blkStmt *ast.BlockStmt
				ast.Inspect(new_body_node, func(n ast.Node) bool {
					if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "dummy" {
						blkStmt = fn.Body // Extract the block statements
						return false      // Stop traversing
					}
					return true
				})

				trampoline_fn := &ast.FuncDecl{
					Doc:  fn.Doc,
					Name: ast.NewIdent(fn.Name.Name + "_Trampoline"),
					Type: fn.Type,
					Body: blkStmt,
				}
				trampoline_fns = append(trampoline_fns, trampoline_fn)
			}
		}
		return true
	})
	var decls []ast.Decl
	file.Decls = decls
	for _, decl := range trampoline_fns {
		file.Decls = append(file.Decls, decl)
	}

	new_file := strings.ReplaceAll(filename, ".go", "_trampoline.go")
	f, err := os.Create(new_file)
	defer f.Close()
	err = printer.Fprint(f, fset, file)
	return new_file, nil
}

type templateExecutor struct {
	Funcs template.FuncMap
}

type templateArgs struct {
	ArgString string
	RetString string
	Name      string
}

func newTemplateExecutor() *templateExecutor {
	e := &templateExecutor{
		Funcs: template.FuncMap{},
	}
	return e
}

func (e *templateExecutor) exec(name string, args templateArgs) (string, error) {
	t, err := template.New(name).Funcs(e.Funcs).Parse(trampoline_template)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, args)
	return buf.String(), err
}

var trampoline_template = `
package foo

func dummy() {
{{if .RetString}}
{{.RetString}}, err := {{.Name}}_Specialized({{.ArgString}})
{{else}}
err := {{.Name}}_Specialized({{.ArgString}})
{{end}}
if err != nil {
	return {{.Name}}_Original({{.ArgString}})
}
{{if .RetString}}
return {{.RetString}}, err
{{else}}
return err
{{end}}
}
`
