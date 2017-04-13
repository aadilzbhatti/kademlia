package dht

import (
	"fmt"
	"os/exec"
	"strings"
)

/**
 * Returns the hostname of the current node
 *
 * @type {string}
 */
func getHostName() string {
	out, err := exec.Command("hostname").Output()
	if err != nil {
		fmt.Println("Failed to obtain hostname")
	}
	return strings.TrimSpace(string(out))
}

/**
 * Returns the distance between nodes with
 * id=id1 and id=id2
 */
func distance(id1 int, id2 int) int {
	return id1 ^ id2
}
