package iridescent

import (
	"fmt"
	"log/slog"

	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/service"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/plugins/golang"
	"github.com/blueprint-uservices/blueprint/plugins/workflow/workflowspec"
	"github.com/vaastav/iridescent/iridescent_rt/autotune"
)

type IridescentProcNode struct {
	golang.Node
	golang.Instantiable

	InstanceName string
	Spec         *workflowspec.Service

	// Constructor args
	Dur          string
	Period       string
	Strategy     string
	SpecFilePath string
}

func (node *IridescentProcNode) Name() string {
	return node.InstanceName
}

func (node *IridescentProcNode) String() string {
	return node.Name() + " = IridescentProcNode()"
}

func newIridescentRT(name string, duration string, period string, strategy_name string, spec_file_path string) (*IridescentProcNode, error) {
	spec, err := workflowspec.GetService[*autotune.IridescentRT]()
	if err != nil {
		return nil, err
	}

	node := &IridescentProcNode{}
	node.InstanceName = name
	node.Spec = spec
	node.Dur = duration
	node.Period = period
	node.Strategy = strategy_name
	node.SpecFilePath = spec_file_path
	return node, nil
}

// Implements golang.Instantiable
func (node *IridescentProcNode) AddInstantiation(builder golang.NamespaceBuilder) error {
	if builder.Visited(node.InstanceName) {
		return nil
	}

	slog.Info(fmt.Sprintf("Instantiating IridescentProcNode %v in %v/%v", node.InstanceName, builder.Info().Package.PackageName, builder.Info().FileName))

	constructor := node.Spec.Constructor.AsConstructor()
	return builder.DeclareConstructor(node.InstanceName, constructor, []ir.IRNode{&ir.IRValue{Value: node.Dur}, &ir.IRValue{Value: node.Period}, &ir.IRValue{Value: node.Strategy}, &ir.IRValue{Value: node.SpecFilePath}})
}

// Implements golang.ProvidesModule
func (node *IridescentProcNode) AddToWorkspace(builder golang.WorkspaceBuilder) error {
	return node.Spec.AddToWorkspace(builder)
}

// Implements golang.ProvidesInterface
func (node *IridescentProcNode) AddInterfaces(builder golang.ModuleBuilder) error {
	return node.Spec.AddToModule(builder)
}

// Implements service.ServiceNode
func (node *IridescentProcNode) GetInterface(ctx ir.BuildContext) (service.ServiceInterface, error) {
	return node.Spec.Iface.ServiceInterface(ctx), nil
}

// Implements golang.Node
func (node *IridescentProcNode) ImplementsGolangNode() {}
