package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"os/exec"
	"strings"
  "net"
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
 * Returns the bucket that node with id=id2
 * would be in for the routing table of node
 * with id=id1
 *
 * @type {int}
 */
func getBucket(id1 int, id2 int) int {
	if id1 == 0 {
		return 1
	}
	maxId := math.Max(float64(id1), float64(id2))
	numBits := uint(math.Ceil(math.Log2(maxId)))
	for i := numBits; i > 0; i-- {
		msb_1 := id1 & (1 << i)
		msb_2 := id2 & (1 << i)
		if msb_1 != msb_2 {
			return int(numBits) - int(i)
		}
	}
	return int(numBits) - 1
}

/**
 * Returns the distance between nodes with
 * id=id1 and id=id2
 */
func distance(id1 int, id2 int) int {
	return id1 ^ id2
}
