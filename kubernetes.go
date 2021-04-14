package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

// context can be either "kind-cluster-a" or "kind-cluster-b"
func getClient(context string) *kubernetes.Clientset {
	home := os.Getenv("HOME")
	kubeconfig := filepath.Join(home, ".kube", "config")

	config, _ := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
	clientset, _ := kubernetes.NewForConfig(config)
	return clientset
}

// creates mdb namespace in cluster
func createNamespace(c *kubernetes.Clientset) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "mdb",
			Labels: map[string]string{"istio-injection": "enabled"},
		},
	}
	_, err := c.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error: failed to create namespace: %v\n", err)
		return
	}

	fmt.Println("successfully created namespace")
}

// creates the service in "mdb" namespace of the cluster
func createService(c *kubernetes.Clientset, name string) {
	decode := scheme.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(svc), nil, nil)
	if err != nil {
		fmt.Printf("error decoding: %v\n", err)
	}

	switch obj.(type) {
	case *v1.Service:
		svc := obj.(*v1.Service)
		svc.ObjectMeta.Labels = map[string]string{"app": name + "-svc"}
		svc.ObjectMeta.Namespace = "mdb"
		svc.ObjectMeta.Name = name + "-svc"
		svc.Spec.Selector = map[string]string{"app": name}
		_, err := c.CoreV1().Services("mdb").Create(context.TODO(), svc, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("error: failed to create service: %v\n", err)
			return
		}
		fmt.Printf("successfully created service: %s\n", name+"-svc")
	default:
		panic("unkown service type")
	}
}

func createSecret(c *kubernetes.Clientset) {
	// read file with stored secret
	r, err := ioutil.ReadFile("./agent.txt")
	if err != nil {
		panic(err.Error())
	}

	secretObj := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "agent-secret",
		},
		Data: map[string][]byte{"agentApiKey": r},
	}

	_, err = c.CoreV1().Secrets("mdb").Create(context.TODO(), secretObj, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error: failed to create secret: %v", err)
		return
	}
	fmt.Println("successfully created secret object")
}

func createPod(c *kubernetes.Clientset, name string) {
	decode := scheme.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(pod), nil, nil)
	if err != nil {
		fmt.Printf("error decoding: %v\n", err)
	}

	switch obj.(type) {
	case *v1.Pod:
		podObj := obj.(*v1.Pod)

		// Start editing things now
		podObj.ObjectMeta.Name = name
		podObj.ObjectMeta.Labels = map[string]string{"app": name}

		arr := podObj.Spec.Containers[0].Env
		arr = append(arr, v1.EnvVar{Name: "GROUP_ID", Value: getGroupId()})
		podObj.Spec.Containers[0].Env = arr

		podObj.Spec.Hostname = name

		_, err := c.CoreV1().Pods("mdb").Create(context.TODO(), podObj, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("error: failed to create pod: %v\n", err)
			return
		}

		fmt.Printf("successfully created pod: %s\n", name)
	default:
		panic("unkown pod type")

	}
}

func deployMongoDBRS() {
	// i. get client for cluster-a / cluster-b
	ca := getClient("kind-cluster-a")
	cb := getClient("kind-cluster-b")

	// // ii. create namespace in cluster-a / cluster-b
	createNamespace(ca)
	createNamespace(cb)

	// iii. create services 1 service in clustera and 2 services in cluster2
	createService(ca, "my-replica-set-0")
	createService(cb, "my-replica-set-1")
	createService(cb, "my-replica-set-2")

	// iv: create the secret object in each both cluster
	createSecret(ca)
	createSecret(cb)

	// v. create the pods 1 pod in clustera and 2 pods in cluster2
	createPod(ca, "my-replica-set-0")
	createPod(cb, "my-replica-set-1")
	createPod(cb, "my-replica-set-2")
}

// also cleanup cloudQA ??
func deleteNamespace() {

}
