package pass

import (
	"go/ast"
	"go/token"
)

type SpecPass interface {
	Modify(fset *token.FileSet, f *ast.File) *ast.File
}
