package pass

import (
	"bytes"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReorderIfPass(t *testing.T) {
	rifp := NewReorderIfPass()
	src := `package main

func reorderable_if() {
	//reorder:if branch_local
	if b > 0 {
		println("b")
	} else if a > 0 {
		println("a")
	} else if c > 0 {
		println("c")
	} else {
		println("default")
	}
}
`

	reordered_src := `package main

func reorderable_if() {
	//reorder:if branch_local
	if c > 0 {
		println("c")
	} else if a > 0 {
		println("a")
	} else if b > 0 {
		println("b")
	} else {
		println("default")
	}
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	require.NoError(t, err)
	// Fully reverse the branch order
	rifp.SetOrder("branch_local", []int{2, 1, 0})

	file = rifp.Modify(fset, file)
	var out bytes.Buffer
	printer.Fprint(&out, fset, file)

	log.Println(out.String())

	require.Equal(t, reordered_src, out.String())
}
