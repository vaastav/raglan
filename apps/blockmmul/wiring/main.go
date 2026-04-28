package main

import (
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/vaastav/raglan/apps/blockmmul/wiring/specs"
)

func main() {
	name := "MatMul"
	cmdbuilder.MakeAndExecute(
		name,
		specs.Docker,
	)
}
