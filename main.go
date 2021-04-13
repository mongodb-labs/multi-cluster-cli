// This CLI enables you to deploy MongoDB pods accross multiple clusters
// It assumes the clusters are already setup prior to using this. At some point
// we would like to create the cluster as well maybe??
// Create project: ./mccli --op project --name mdb
package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

type Object string

const (
	Project  Object = "project"
	AgentKey        = "agentKey"
	MongoDB         = "mongo"
)

// cluster names are currently hardcoded
// "cluster-a" and "cluster-b", maybe tweak it later
func main() {
	var op, name string
	var setupKind bool

	flag.StringVar(&op, "op", "", "operation to perform")
	flag.StringVar(&name, "name", "", "name of the project to create")
	flag.BoolVar(&setupKind, "setkind", false, "flag to setup kind cluster")
	flag.Parse()

	if setupKind {
		setUpKubernetesClusters()
	}

	switch Object(op) {
	case Project:
		// assert if name is passed as well
		// read the Orgid from yaml
		createProject(name)
		updateAutomationAgent()
	case AgentKey:
		createAgentKey()
	case MongoDB:
		deployMongoDBRS()
	default:
		fmt.Println("Enter a valid option...")
	}
}
