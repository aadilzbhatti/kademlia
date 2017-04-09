package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func getHostName() string {
	out, err := exec.Command("hostname").Output()
	if err != nil {
		fmt.Println("Failed to obtain hostname")
	}
	return strings.Trim(string(out), "\n")
}
