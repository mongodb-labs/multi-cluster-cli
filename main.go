// This CLI enables you to deploy MongoDB pods accross multiple clusters
package main

import (
	"fmt"
	"time"

	flag "github.com/spf13/pflag"
)

type Object string

const (
	Project Object = "project"
	Cluster        = "cluster"
	MongoDB        = "mongo"
)

// cluster names are currently hardcoded
// "cluster-a" and "cluster-b", maybe tweak it later
func main() {
	var op string

	flag.StringVar(&op, "op", "", "operation to perform, available options: cluster, project and mongo")
	// TODO: maybe make this configurable?
	// flag.StringVar(&name, "name", "", "name of the project to create")
	flag.Parse()

	switch Object(op) {
	case Project:
		// assert if name is passed as well
		// read the Orgid from yaml
		createProject()
		// hack: at times patching automation agent right after project creation fails
		time.Sleep(5 * time.Second)
		updateAutomationAgent()

		time.Sleep(5 * time.Second)
		createAgentKey()
	case MongoDB:
		deployMongoDBRS()
	case Cluster:
		setUpKubernetesClusters()
	default:
		fmt.Println("Enter a valid option...")
	}
}
