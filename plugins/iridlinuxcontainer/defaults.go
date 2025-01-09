package iridlinuxcontainer

import (
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/blueprint/ioutil"
	"github.com/blueprint-uservices/blueprint/blueprint/pkg/ir"
	"github.com/blueprint-uservices/blueprint/plugins/linux"
)

func RegisterAsDefaultBuilder() {
	ir.RegisterDefaultNamespace[linux.Process]("linux", buildDefaultLinuxWorkspace)
}

func buildDefaultLinuxWorkspace(outputDir string, nodes []ir.IRNode) error {
	ctr := newIridLinuxContainerNode("linux")
	ctr.Nodes = nodes
	ctrDir, err := ioutil.CreateNodeDir(outputDir, "linux")
	if err != nil {
		return err
	}
	return ctr.GenerateArtifacts(ctrDir)
}
