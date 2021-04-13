package main

import (
	"context"
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

func createProject(name string) {
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

	proj := opsmngr.Project{
		Name:  name,
		OrgID: getOrdId(),
	}
	_, _, err = client.Projects.Create(context.TODO(), &proj)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

}

func updateAutomationAgent() {
	// get the json
}

// to be used by agent and put in the secret
func createAgentKey() {
	return
}
