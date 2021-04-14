package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/antonholmquist/jason"
	"github.com/mongodb-forks/digest"
	"go.mongodb.org/ops-manager/opsmngr"
)

type Payload struct {
	Name  string `json:"name"`
	Orgid string `json:"orgId"`
}

func getPublicAndPrivateKey() (string, string) {
	data, _ := ioutil.ReadFile("config.json")
	v, _ := jason.NewObjectFromBytes(data)
	pub, _ := v.GetString("public")
	priv, _ := v.GetString("private")
	return pub, priv
}
func getOrdId() string {
	data, _ := ioutil.ReadFile("config.json")
	v, _ := jason.NewObjectFromBytes(data)
	ordId, _ := v.GetString("orgId")
	return ordId
}

func getCloudClient() *opsmngr.Client {
	pub, priv := getPublicAndPrivateKey()
	t := digest.NewTransport(pub, priv)

	tc, err := t.Client()
	if err != nil {
		panic(err.Error())
	}

	clientops := opsmngr.SetBaseURL("https://cloud-qa.mongodb.com/" + opsmngr.APIPublicV1Path)
	client, err := opsmngr.New(tc, clientops)
	if err != nil {
		panic(err.Error())
	}
	return client

}

func createProject(name string) {
	client := getCloudClient()
	proj := opsmngr.Project{
		Name:  name,
		OrgID: getOrdId(),
	}
	p, _, err := client.Projects.Create(context.TODO(), &proj)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	// write project ID to a tmp file "cloud.txt"
	w := []byte(p.ID)
	err = ioutil.WriteFile("./cloud.txt", w, 0644)
	if err != nil {
		fmt.Printf("error: failed to write project ID to file: %v", err)
		return
	}
	fmt.Println("successfully created project mdb")
}

func updateAutomationAgent() {
	var tmp opsmngr.AutomationConfig
	json.Unmarshal([]byte(processesJSON), &tmp)

	// read project/groupID from file.
	r, err := ioutil.ReadFile("./cloud.txt")
	if err != nil {
		panic(err.Error())
	}

	groupId := string(r)
	client := getCloudClient()

	// Get automationConfig first
	automationConfig, _, err := client.Automation.GetConfig(context.TODO(), groupId)
	if err != nil {
		fmt.Printf("error: failed to get automationConfig: %v", err)
	}

	// edit the automationConfig now
	automationConfig.Processes = tmp.Processes
	automationConfig.ReplicaSets = tmp.ReplicaSets
	automationConfig.MonitoringVersions = tmp.MonitoringVersions
	automationConfig.BackupVersions = tmp.BackupVersions
	automationConfig.AgentVersion = tmp.AgentVersion
	automationConfig.Options = tmp.Options

	// Send POST request now
	_, err = client.Automation.UpdateConfig(context.TODO(), groupId, automationConfig)
	if err != nil {
		fmt.Printf("error: failed to update automationConfig: %v\n", err)
		return
	}
	fmt.Println("successfully updated automationConfig")
}

// to be used by agent and put in the file cloud.txt
func createAgentKey() {
	return
}
