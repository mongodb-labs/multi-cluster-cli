package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func setupKindClustersWithCillium() {

	// run the bash script
	// i. It spins up two kind clusters
	// ii. Connects them with cillium
	// iii. verifies the kind installation with cillium
	// err := cmd.Run()

	// verify the installation
	deadline := time.Now().Add(4 * time.Minute)
	for {
		out, _ := exec.Command("/bin/sh", "./check-cilium.sh").Output()
		// fmt.Println(string(out))
		if strings.Count(string(out), "6/6 reachable") == 6 {
			fmt.Println("Cillium installed successfully")
			break
		}
		if time.Now().After(deadline) {
			fmt.Println("Error: Gave up waiting for cillium to be installed")
		}
		fmt.Println("waiting for cilium to be installed ....")
		time.Sleep(10 * time.Second)
	}
}

// TODO : install istio https://istio.io/latest/docs/setup/install/multicluster/multi-primary/
func setUpIstio() {
	fmt.Println("hello")
	cmd := exec.Command("/bin/sh", "./install-istio.sh", "`< cluster-contexts.txt`")
	fmt.Println(cmd)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// run "docker system prune -a" ? as pre-requesite
func setUpKubernetesClusters() {
	// check if kind clusters exist if yes, skip cluster setup.
	// specify a flag to delte kind clusters explicitly
	setupKindClustersWithCillium()
	// setUpIstio()
}
