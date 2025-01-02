package tput

import (
	"log/slog"

	"github.com/blueprint-uservices/blueprint/blueprint/pkg/blueprint"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/pointer"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/golang"
)

func AddTputMeter(spec wiring.WiringSpec, serviceName string) {
	serverWrapper := serviceName + ".server.tputmeter"

	ptr := pointer.GetPointer(spec, serviceName)
	if ptr == nil {
		slog.Error("Unable to add tput meter to " + serviceName + " as it is not a pointer")
		return
	}

	serverNext := ptr.AddDstModifier(spec, serverWrapper)

	// Define the server wrapper
	spec.Define(serverWrapper, &TputServerWrapper{}, func(ns wiring.Namespace) (ir.IRNode, error) {
		var server golang.Service
		if err := ns.Get(serverNext, &server); err != nil {
			return nil, blueprint.Errorf("Tputer %s expected %s to be a golang.Service, but encountered %s", serverWrapper, serverNext, err)
		}

		return newTputServerWrapper(serverWrapper, server)
	})
}
