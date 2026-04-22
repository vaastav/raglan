package latency

import (
	"log/slog"

	"github.com/blueprint-uservices/blueprint/blueprint/pkg/blueprint"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/pointer"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/golang"
)

func AddLatencyMeter(spec wiring.WiringSpec, serviceName string, percentile float64) {
	serverWrapper := serviceName + ".server.latmeter"

	ptr := pointer.GetPointer(spec, serviceName)
	if ptr == nil {
		slog.Error("Unable to add latency meter to " + serviceName + " as it is not a pointer")
		return
	}

	serverNext := ptr.AddDstModifier(spec, serverWrapper)

	spec.Define(serverWrapper, &PercentileLatencyServerWrapper{}, func(ns wiring.Namespace) (ir.IRNode, error) {
		var server golang.Service

		if err := ns.Get(serverNext, &server); err != nil {
			return nil, blueprint.Errorf("Latency Meter %s expected %s to be a golang.Service, but encountered %s", serverWrapper, serverNext, err)
		}

		return newPercentileLatencyServerWrapper(serverWrapper, server, percentile)
	})
}
