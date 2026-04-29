package tput

import (
	"fmt"
	"reflect"

	"github.com/blueprint-uservices/blueprint/blueprint/pkg/coreplugins/service"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/plugins/golang"
	"github.com/blueprint-uservices/blueprint/plugins/golang/gocode"
)

type TputServerWrapper struct {
	golang.Service
	golang.GeneratesFuncs
	golang.Instantiable

	InstanceName string
	Wrapped      golang.Service

	outputPackage string
}

func (node *TputServerWrapper) ImplementsGolangNode() {}

func (node *TputServerWrapper) Name() string {
	return node.InstanceName
}

func (node *TputServerWrapper) String() string {
	return node.Name() + " = TputServerWrapper(" + node.Wrapped.Name() + ")"
}

func (node *TputServerWrapper) AddInterfaces(builder golang.ModuleBuilder) error {
	return node.Wrapped.AddInterfaces(builder)
}

func newTputServerWrapper(name string, server ir.IRNode) (*TputServerWrapper, error) {
	serverNode, is_callable := server.(golang.Service)
	if !is_callable {
		return nil, fmt.Errorf("Tput server wrapper requires %s to be a golang service but got %s", server.Name(), reflect.TypeOf(server).String())
	}

	node := &TputServerWrapper{}
	node.InstanceName = name
	node.Wrapped = serverNode
	node.outputPackage = "tput"

	return node, nil
}

func (node *TputServerWrapper) GetInterface(ctx ir.BuildContext) (service.ServiceInterface, error) {
	return node.Wrapped.GetInterface(ctx)
}

func (node *TputServerWrapper) GenerateFuncs(builder golang.ModuleBuilder) error {
	service, err := golang.GetGoInterface(builder, node.Wrapped)
	if err != nil {
		return err
	}
	err = generateServerHandler(builder, service, node.outputPackage)
	if err != nil {
		return err
	}
	return nil
}

func (node *TputServerWrapper) AddInstantiation(builder golang.NamespaceBuilder) error {
	// Only generate instantiation code for this instance once
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
			Name: fmt.Sprintf("New_%v_TputHandler", iface.BaseName),
			Arguments: []gocode.Variable{
				{Name: "ctx", Type: &gocode.UserType{Package: "context", Name: "Context"}},
				{Name: "service", Type: iface},
			},
		},
	}

	return builder.DeclareConstructor(node.InstanceName, constructor, []ir.IRNode{node.Wrapped})
}
