package main

import (
	"fmt"
	"hash/fnv"
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
 * Returns the hash mod the number nodes in
 * the system
 *
 * @type {int}
 */
func hash(key string, size int) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % size
}

/**
 * Returns the bucket that Nodewith id=id2
 * would be in for the routing table of node
 * with id=id1
 *
 * @type {int}
 */
func getConflictingBit(id1, id2 []byte) int {
	for i := 0; i < len(id1); i++ {
		res := id1[i] ^ id2[i]

		// This is a byte. Need to get bit position.
		for j := 0; j < 8; j++ {
			if hasBit(int(res), uint(7-j)) {
				return 32 - (8*i + j) - 1
			}

		}
	}
	return 0 // same id
}

/**
 * Returns the distance between nodes with
 * id=id1 and id=id2
 */
func distance(id1 int, id2 int) int {
	return id1 ^ id2
}
