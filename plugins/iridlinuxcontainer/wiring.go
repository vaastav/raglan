package iridlinuxcontainer

import (
	"strings"

	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/namespaceutil"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/wiring"
	"github.com/blueprint-uservices/blueprint/plugins/linux"
)

func AddToContainer(spec wiring.WiringSpec, containerName, childName string) {
	namespaceutil.AddNodeTo[Container](spec, containerName, childName)
}

// Deploy can be used by wiring specs to deploy a process-level service in a linux container.
//
// Adds a modifier to the service that, during compilation, will create the linux container if
// not already created.
//
// The name of the container created is determined by attempting to replace a "_service" suffix
// with "_ctr", or adding "_ctr" if serviceName doesn't end with "_service", e.g.
//
//	user_service => user_ctr
//	user => user_ctr
//	user_srv => user_srv_ctr
//
// After calling [Deploy], serviceName will be a container-level service.
//
// Returns the name of the container
func Deploy(spec wiring.WiringSpec, serviceName string) string {
	servicePrefix, _ := strings.CutSuffix(serviceName, "_service")
	ctrName := servicePrefix + "_ctr"
	CreateContainer(spec, ctrName, serviceName)
	return ctrName
}

func CreateContainer(spec wiring.WiringSpec, containerName string, children ...string) string {
	// If any children were provided in this call, add them to the process via a property
	for _, childName := range children {
		AddToContainer(spec, containerName, childName)
	}

	// A linux container node is simply a namespace that accumulates linux process nodes
	spec.Define(containerName, &Container{}, func(namespace wiring.Namespace) (ir.IRNode, error) {
		ctr := newIridLinuxContainerNode(containerName)
		_, err := namespaceutil.InstantiateNamespace(namespace, &linuxContainerNamespace{ctr})
		return ctr, err
	})

	return containerName
}

// A [wiring.NamespaceHandler] used to build golang process nodes
type linuxContainerNamespace struct {
	*Container
}

// Implements [wiring.NamespaceHandler]
func (ctr *Container) Accepts(nodeType any) bool {
	_, isLinuxProcess := nodeType.(linux.Process)
	return isLinuxProcess
}

// Implements [wiring.NamespaceHandler]
func (ctr *Container) AddEdge(name string, edge ir.IRNode) error {
	ctr.Edges = append(ctr.Edges, edge)
	return nil
}

// Implements [wiring.NamespaceHandler]
func (ctr *Container) AddNode(name string, node ir.IRNode) error {
	ctr.Nodes = append(ctr.Nodes, node)
	return nil
}
