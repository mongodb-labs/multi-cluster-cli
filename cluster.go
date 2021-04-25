package main

import (
	"bufio"
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
	cmd := exec.Command("/bin/sh", "./scripts/set-up-cillium.sh")

	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()
	cmd.Start()
	cmd.Wait()

	// verify the installation
	deadline := time.Now().Add(4 * time.Minute)
	for {
		out, _ := exec.Command("/bin/sh", "./scripts/check-cilium.sh").Output()
		if strings.Count(string(out), "6/6 reachable") == 6 {
			fmt.Println("Cillium installed successfully")
			break
		}
		if time.Now().After(deadline) {
			fmt.Println("error: Gave up waiting for cillium to be installed")
		}
		fmt.Println("waiting for cilium to be installed ....")
		time.Sleep(10 * time.Second)
	}
}

// TODO : install istio https://istio.io/latest/docs/setup/install/multicluster/multi-primary/
func setUpIstio() {
	cmd := exec.Command("/bin/sh", "./scripts/install-istio.sh")
	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()
	cmd.Start()
	cmd.Wait()
}

// run "docker system prune -a" ? as pre-requesite
func setUpKubernetesClusters() {
	// check if kind clusters exist if yes, skip cluster setup.
	// specify a flag to delte kind clusters explicitly
	setupKindClustersWithCillium()
	// hack wait for 5 second
	time.Sleep(5 * time.Second)
	setUpIstio()
}
