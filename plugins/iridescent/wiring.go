package iridescent

import (
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/goproc"
)

func AddIridescent(spec wiring.WiringSpec, procName string, exploration_duration string, monitoring_period string, strategy string, specialization_file string) {
	irid_name := procName + ".iridescent_runtime"

	// Initialize the Iridescent Runtime IR Node.

	spec.Define(irid_name, &IridescentProcNode{}, func(ns wiring.Namespace) (ir.IRNode, error) {
		return newIridescentRT(exploration_duration, monitoring_period, strategy, specialization_file)
	})

	goproc.AddToProcess(spec, procName, irid_name)
}
