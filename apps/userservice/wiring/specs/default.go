package specs

import (
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/blueprint-uservices/blueprint/plugins/goproc"
	"github.com/blueprint-uservices/blueprint/plugins/http"
	"github.com/blueprint-uservices/blueprint/plugins/mongodb"
	"github.com/blueprint-uservices/blueprint/plugins/workflow"
	"github.com/blueprint-uservices/blueprint/plugins/workload"
	wf "github.com/vaastav/raglan/apps/userservice/workflow"
	"github.com/vaastav/raglan/apps/userservice/workload/workloadgen"
	"github.com/vaastav/raglan/plugins/iridescent"
	"github.com/vaastav/raglan/plugins/iridlinuxcontainer"
)

var Default = cmdbuilder.SpecOption{
	Name:        "default",
	Description: "Deploys in a iridescent container",
	Build:       makeDefaultSpec,
}

func applyHTTPDefaults(spec wiring.WiringSpec, serviceName, proc_name, ctrName string) string {
	http.Deploy(spec, serviceName)
	iridescent.AddIridescent(spec, proc_name, "20s", "2s", "linear", "/src/"+proc_name+"/workflow/guest/db.go")
	goproc.CreateProcess(spec, proc_name, serviceName)
	return iridlinuxcontainer.CreateContainer(spec, ctrName, proc_name)
}

func makeDefaultSpec(spec wiring.WiringSpec) ([]string, error) {
	var services []string

	nadb := mongodb.Container(spec, "nadb")
	eudb := mongodb.Container(spec, "eudb")
	apdb := mongodb.Container(spec, "apdb")
	sadb := mongodb.Container(spec, "sadb")
	afdb := mongodb.Container(spec, "afdb")
	ocdb := mongodb.Container(spec, "ocdb")

	user_service := workflow.Service[wf.UserService](spec, "user_service", nadb, eudb, apdb, sadb, afdb, ocdb)

	ctr := applyHTTPDefaults(spec, user_service, "user_proc", "user_ctr")
	services = append(services, ctr)

	wlgen := workload.Generator[workloadgen.MultiWorkload](spec, "wlgen", user_service)
	services = append(services, wlgen)

	return services, nil
}
