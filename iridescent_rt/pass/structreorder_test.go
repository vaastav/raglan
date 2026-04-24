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

func TestReorderStructPass(t *testing.T) {
	rstrp := NewReorderStructPass()
	src := `package main
	
//reorder:struct struct_local
type ReorderableStruct struct {
	Field0 int
	Field1 string
	Field2 []int
	Field3 bool
	Field4 []UserDefined
}`

	reordered_src := `package main

//reorder:struct struct_local
type ReorderableStruct struct {
	Field4	[]UserDefined
	Field3	bool
	Field2	[]int
	Field1	string
	Field0	int
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	require.NoError(t, err)
	// Fully reverse the struct field order
	rstrp.SetOrder("struct_local", []int{4, 3, 2, 1, 0})
	file = rstrp.Modify(fset, file)
	var out bytes.Buffer
	printer.Fprint(&out, fset, file)

	log.Println(out.String())

	require.Equal(t, reordered_src, out.String())
}
