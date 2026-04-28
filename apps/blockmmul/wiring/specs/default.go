package specs

import (
	"strings"

	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/blueprint-uservices/blueprint/plugins/goproc"
	"github.com/blueprint-uservices/blueprint/plugins/healthchecker"
	"github.com/blueprint-uservices/blueprint/plugins/http"
	"github.com/blueprint-uservices/blueprint/plugins/workflow"
	"github.com/vaastav/raglan/apps/blockmmul/workflow/services"
	"github.com/vaastav/raglan/plugins/iridescent"
	"github.com/vaastav/raglan/plugins/iridlinuxcontainer"
)

var Docker = cmdbuilder.SpecOption{
	Name:        "docker",
	Description: "Deploys each service in a separate container with http, uses mongodb as NoSQL database backends, and applies a number of modifiers",
	Build:       makeDockerSpec,
}

func makeDockerSpec(spec wiring.WiringSpec) ([]string, error) {
	applyDockerDefaults := func(spec wiring.WiringSpec, serviceName string) string {
		proc_name := goproc.CreateProcess(spec, strings.ReplaceAll(serviceName, "_service", "_proc"))
		iridescent.AddIridescent(spec, proc_name, "20s", "2s", "linear", "/src/"+proc_name+"/workflow/guest/mmul.go")
		healthchecker.AddHealthCheckAPI(spec, serviceName)
		http.Deploy(spec, serviceName)
		goproc.AddToProcess(spec, proc_name, serviceName)
		return iridlinuxcontainer.Deploy(spec, serviceName)
	}

	mmul_service := workflow.Service[services.MatrixMulService](spec, "mmul_service")

	ctr := applyDockerDefaults(spec, mmul_service)

	return []string{ctr}, nil
}
