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

func TestUserServiceReorder(t *testing.T) {
	rifp := NewReorderIfPass()
	src := `package main

func GetUserInfo(ctx context.Context, u *workflow.UserServiceImpl, username string) (workflow.Info, error) {
	iridescent_instr_range_int("userinfo", 0, 5)
	var info workflow.Info
	//reorder:if userinfo
	if find_user(ctx, u.NaDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 0)
		return info, nil
	} else if find_user(ctx, u.EuDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 1)
		return info, nil
	} else if find_user(ctx, u.ApDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 2)
		return info, nil
	} else if find_user(ctx, u.SaDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 3)
		return info, nil
	} else if find_user(ctx, u.AfDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 4)
		return info, nil
	} else if find_user(ctx, u.OcDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 5)
		return info, nil
	}
	return workflow.Info{}, errors.New("User does not exist!")
}
`

	reordered_src := `package main

func GetUserInfo(ctx context.Context, u *workflow.UserServiceImpl, username string) (workflow.Info, error) {
	iridescent_instr_range_int("userinfo", 0, 5)
	var info workflow.Info
	//reorder:if userinfo
	if find_user(ctx, u.OcDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 5)
		return info, nil
	} else if find_user(ctx, u.AfDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 4)
		return info, nil
	} else if find_user(ctx, u.SaDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 3)
		return info, nil
	} else if find_user(ctx, u.ApDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 2)
		return info, nil
	} else if find_user(ctx, u.EuDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 1)
		return info, nil
	} else if find_user(ctx, u.NaDB, username, &info) {
		autotune.GetRuntime().SpecRT.Instrument("userinfo", 0)
		return info, nil
	}

	return workflow.Info{}, errors.New("User does not exist!")
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	require.NoError(t, err)
	// Fully reverse the branch order
	rifp.SetOrder("userinfo", []int{5, 4, 3, 2, 1, 0})

	file = rifp.Modify(fset, file)
	var out bytes.Buffer
	printer.Fprint(&out, fset, file)

	log.Println(out.String())

	require.Equal(t, reordered_src, out.String())
}
