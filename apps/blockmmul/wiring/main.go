package main

import (
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/vaastav/iridescent/apps/blockmmul/wiring/specs"
)

func main() {
	name := "MatMul"
	cmdbuilder.MakeAndExecute(
		name,
		specs.Docker,
	)
}
