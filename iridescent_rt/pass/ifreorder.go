package pass

import (
	"go/ast"
	"go/token"
	"log"
	"strings"
)

type ReorderIfPass struct {
	Orders map[string][]int
}

func NewReorderIfPass() *ReorderIfPass {
	rifp := &ReorderIfPass{Orders: make(map[string][]int)}
	return rifp
}

type ifBranch struct {
	cond ast.Expr
	body *ast.BlockStmt
}

// --------------------
// Parse: //reorder:if name_location
// --------------------
func (pass *ReorderIfPass) getReorderID(cg *ast.CommentGroup) (string, bool) {
	if cg == nil {
		return "", false
	}

	for _, c := range cg.List {
		text := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
		if strings.HasPrefix(text, "reorder:if") {
			parts := strings.Fields(text)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), true // e.g. "name_location"
			}
		}
	}
	return "", false
}

func (pass *ReorderIfPass) collectIfChain(n *ast.IfStmt) ([]*ifBranch, ast.Stmt) {
	var branches []*ifBranch
	current := n

	for current != nil {
		branches = append(branches, &ifBranch{cond: current.Cond, body: current.Body})

		if next, ok := current.Else.(*ast.IfStmt); ok {
			current = next
		} else {
			return branches, current.Else
		}
	}
	return branches, nil
}

func (pass *ReorderIfPass) buildIfChain(n *ast.IfStmt, branches []*ifBranch, elseStmt ast.Stmt) {
	curr := n
	for i := 0; i < len(branches); i++ {
		curr.Cond = branches[i].cond
		curr.Body = branches[i].body

		if i == len(branches)-1 {
			curr.Else = elseStmt
			break
		}

		nextIf, ok := curr.Else.(*ast.IfStmt)
		if !ok {
			nextIf = &ast.IfStmt{}
			curr.Else = nextIf
		}
		curr = nextIf
	}
}

func (pass *ReorderIfPass) reorderIfChain(n *ast.IfStmt, id string) {
	branches, elseStmt := pass.collectIfChain(n)
	order := pass.GetOrder(id)

	newBranches := make([]*ifBranch, len(branches))

	for i, idx := range order {
		if idx < 0 || idx >= len(branches) {
			// Invalid situation just return default
			log.Println("Invalid if-else order")
		}
		// Reorder the branch
		newBranches[i] = branches[idx]
	}

	pass.buildIfChain(n, newBranches, elseStmt)
}

func (pass *ReorderIfPass) GetOrder(id string) []int {
	return pass.Orders[id]
}

func (pass *ReorderIfPass) SetOrder(id string, order []int) {
	pass.Orders[id] = order
}

func (pass *ReorderIfPass) Modify(fset *token.FileSet, file *ast.File) *ast.File {
	cmap := ast.NewCommentMap(fset, file, file.Comments)

	ast.Inspect(file, func(n ast.Node) bool {
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}

		// Check if this has been marked for reordering
		groups := cmap[ifStmt]
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

		pass.reorderIfChain(ifStmt, id)

		return false
	})
	return file
}
