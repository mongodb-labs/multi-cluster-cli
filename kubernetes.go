package main

import (
	"context"
	"fmt"
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
		fmt.Printf("Error: failed to create namespace: %v\n", err)
		return
	}

	fmt.Println("successfully created namespace")
}

// creates the service in "mdb" namespace of the cluster
func createService(c *kubernetes.Clientset, name string) {
	decode := scheme.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte("./yamls/svc.yaml"), nil, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	switch obj.(type) {
	case *v1.Service:
		svc := obj.(*v1.Service)
		svc.ObjectMeta.Labels = map[string]string{"app": name}
		svc.ObjectMeta.Namespace = "mdb"
		svc.ObjectMeta.Name = name
		svc.Spec.Selector = map[string]string{"app": name}
		_, err := c.CoreV1().Services("mdb").Create(context.TODO(), svc, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("failed to create service: %v\n", err)
			return
		}
		fmt.Printf("successfully created service: %s\n", name)
	default:
		panic("unkown service type")
	}
}

func createPod() {

}

func deployMongoDBRS() {
	// i. get client for cluster-a / cluster-b
	ca := getClient("kind-cluster-a")
	cb := getClient("kind-cluster-b")

	// ii. create namespace in cluster-a / cluster-b
	createNamespace(ca)
	createNamespace(cb)

	// iii. create services 1 service in clustera and 2 services in cluster2
	createService(ca, "my-replica-set-0-svc")
	createService(cb, "my-replica-set-1-svc")
	createService(cb, "my-replica-set-2-svc")

	// iv: create the secret object in each namespace
	// TODO
	// v. create the pods 1 pod in clustera and 2 pods in cluster2
	// TODO
}
