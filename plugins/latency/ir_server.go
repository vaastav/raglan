package latency

import (
	"fmt"
	"strconv"

	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/service"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/plugins/golang"
	"github.com/blueprint-uservices/blueprint/plugins/golang/gocode"
)

type PercentileLatencyServerWrapper struct {
	golang.Service
	golang.GeneratesFuncs
	golang.Instantiable

	InstanceName string
	Wrapped      golang.Service

	outputPackage string
	percentile    float64
}

func (node *PercentileLatencyServerWrapper) ImplementsGolangNode() {}

func (node *PercentileLatencyServerWrapper) Name() string {
	return node.InstanceName
}

func (node *PercentileLatencyServerWrapper) String() string {
	return node.Name() + " = PercentileLatencyServerWrapper(" + node.Wrapped.Name() + "," + strconv.FormatFloat(node.percentile, 'f', 2, 64) + ")"
}

func (node *PercentileLatencyServerWrapper) AddInterfaces(builder golang.ModuleBuilder) error {
	return node.Wrapped.AddInterfaces(builder)
}

func newPercentileLatencyServerWrapper(name string, server golang.Service, percentile float64) (*PercentileLatencyServerWrapper, error) {
	node := &PercentileLatencyServerWrapper{}
	node.InstanceName = name
	node.Wrapped = server
	node.outputPackage = "latency"
	node.percentile = percentile

	return node, nil
}

func (node *PercentileLatencyServerWrapper) GetInterface(ctx ir.BuildContext) (service.ServiceInterface, error) {
	return node.Wrapped.GetInterface(ctx)
}

func (node *PercentileLatencyServerWrapper) GenerateFuncs(builder golang.ModuleBuilder) error {
	service, err := golang.GetGoInterface(builder, node.Wrapped)
	if err != nil {
		return err
	}

	err = generateServerHandler(builder, service, node.outputPackage, node.percentile)
	if err != nil {
		return err
	}

	return nil
}

func (node *PercentileLatencyServerWrapper) AddInstantiation(builder golang.NamespaceBuilder) error {
	if builder.Visited(node.InstanceName) {
		return nil
	}

	iface, err := golang.GetGoInterface(builder, node.Wrapped)
	if err != nil {
		return err
	}

	constructor := &gocode.Constructor{
		Package: builder.Module().Info().Name + "/" + node.outputPackage,
		Func: gocode.Func{
			Name: fmt.Sprintf("New_%v_PercentileLatencyHandler", iface.BaseName),
			Arguments: []gocode.Variable{
				{Name: "ctx", Type: &gocode.UserType{Package: "context", Name: "Context"}},
				{Name: "service", Type: iface},
			},
		},
	}
	return builder.DeclareConstructor(node.InstanceName, constructor, []ir.IRNode{node.Wrapped})
}
