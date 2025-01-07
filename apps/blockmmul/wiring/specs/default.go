package specs

import (
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/blueprint-uservices/blueprint/plugins/goproc"
	"github.com/blueprint-uservices/blueprint/plugins/healthchecker"
	"github.com/blueprint-uservices/blueprint/plugins/http"
	"github.com/blueprint-uservices/blueprint/plugins/linuxcontainer"
	"github.com/blueprint-uservices/blueprint/plugins/workflow"
	"github.com/vaastav/iridescent/apps/blockmmul/workflow/services"
	"github.com/vaastav/iridescent/plugins/iridescent"
)

var Docker = cmdbuilder.SpecOption{
	Name:        "docker",
	Description: "Deploys each service in a separate container with http, uses mongodb as NoSQL database backends, and applies a number of modifiers",
	Build:       makeDockerSpec,
}

func makeDockerSpec(spec wiring.WiringSpec) ([]string, error) {
	iridescent.AddIridescent(spec, "mmul_proc", "20s", "2s", "linear", "../workflow/guest/mmul.go")
	applyDockerDefaults := func(spec wiring.WiringSpec, serviceName string) string {
		healthchecker.AddHealthCheckAPI(spec, serviceName)
		http.Deploy(spec, serviceName)
		goproc.Deploy(spec, serviceName)
		return linuxcontainer.Deploy(spec, serviceName)
	}

	mmul_service := workflow.Service[services.MatrixMulService](spec, "mmul_service")

	ctr := applyDockerDefaults(spec, mmul_service)

	return []string{ctr}, nil
}
