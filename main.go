// This CLI enables you to deploy MongoDB pods accross multiple clusters
// It assumes the clusters are already setup prior to using this. At some point
// we would like to create the cluster as well maybe??
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	flag "github.com/spf13/pflag"
)

type Object string

const (
	Project  Object = "project"
	APIKeys         = "apiKeys"
	AgentKey        = "agentKey"
	MongoDB         = "mongo"
)

type ProjectData struct {
	Name  string `json:"name"`
	OrgId string `json:orgId`
}

func createProject(orgId, name string) {
	fmt.Println("attempting to create project ...")
	data := ProjectData{
		Name:  name,
		OrgId: orgId,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		// do something
	}

	body := bytes.NewReader(dataBytes)

	req, err := http.NewRequest("POST", "https://cloud.mongodb.com/api/public/v1.0/groups?pretty=true", body)
	if err != nil {

	}
	// TODO: read the yaml file and populate this part
	req.SetBasicAuth("ZVKHFSEB", "5a148efd-61aa-465b-994e-61dcc7a8f77d")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error in making project request: %v", err)
	}
	defer resp.Body.Close()

	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// do something about it
	}
	fmt.Println(string(rb))
}

func createAPIKey() error {
	return nil
}

// to be used by agent and put in the secret
func createAgentKey() error {
	return nil
}

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
		createProject("5f847622c73f2019e8cc5917", name)
	case APIKeys:
		createAPIKey()
	case AgentKey:
		createAgentKey()
	case MongoDB:
		deployMongoDBRS()
	default:
		fmt.Println("Enter a valid option...")
	}
}
