package main

import (
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/vaastav/iridescent/apps/userservice/wiring/specs"
)

func main() {
	name := "DistributedUsers"
	cmdbuilder.MakeAndExecute(
		name,
		specs.Default,
	)
}
