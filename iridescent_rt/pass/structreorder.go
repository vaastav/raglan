package pass

import (
	"go/ast"
	"go/token"
	"strings"
)

type ReorderStructPass struct {
	Orders map[string][]int
}

func NewReorderStructPass() *ReorderStructPass {
	rstrp := &ReorderStructPass{Orders: make(map[string][]int)}
	return rstrp
}

// --------------------
// Parse: //reorder:struct name_location
// --------------------
func (pass *ReorderStructPass) getReorderID(cg *ast.CommentGroup) (string, bool) {
	if cg == nil {
		return "", false
	}

	for _, c := range cg.List {
		text := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
		if strings.HasPrefix(text, "reorder:struct") {
			parts := strings.Fields(text)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), true // e.g. "name_location"
			}
		}
	}
	return "", false
}

func (pass *ReorderStructPass) reorderStruct(st *ast.StructType, id string) {
	if st.Fields == nil || len(st.Fields.List) == 0 {
		return
	}
	order := pass.GetOrder(id)

	fields := st.Fields.List

	reordered := make([]*ast.Field, len(fields))
	for i, idx := range order {
		if idx < 0 || idx >= len(fields) {
			break
		}
		reordered[i] = fields[idx]
	}
	st.Fields.List = reordered
}

func (pass *ReorderStructPass) GetOrder(id string) []int {
	return pass.Orders[id]
}

func (pass *ReorderStructPass) SetOrder(id string, order []int) {
	pass.Orders[id] = order
}

func (pass *ReorderStructPass) Modify(fset *token.FileSet, file *ast.File) *ast.File {
	cmap := ast.NewCommentMap(fset, file, file.Comments)

	ast.Inspect(file, func(n ast.Node) bool {
		// ASSUMPTION: Only 1 type declaration per gen decl
		// This doesn't support declarations of the following pattern:
		// //reorder:struct struct_a
		// type (
		//   A struct {...}
		//   B struct {...}
		//)
		gen, ok := n.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			return true
		}

		// Check if this struct has been marked for reordering
		groups := cmap[gen]

		var id string
		var found bool

		for _, cg := range groups {
			if id, found = pass.getReorderID(cg); found {
				break
			}
		}

		if !found {
			return true
		}

		ts, ok := gen.Specs[0].(*ast.TypeSpec)
		if !ok {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		pass.reorderStruct(st, id)

		return false
	})

	return file
}
