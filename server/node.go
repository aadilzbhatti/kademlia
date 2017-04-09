package main

import (
	"fmt"
	"math"
)

var alpha int = 3

func initializeRoutingTable(id int, numNodes int) [][]TableEntry {
	k := int(math.Ceil(math.Log2(float64(numNodes))))
	buckets := make([][]TableEntry, k)
	return buckets
}

func initializeNode(id int, numNodes int, port int, address string) Node {
	routingTable := initializeRoutingTable(id, numNodes)
	keys := make([]KV, 5)
	n := Node{routingTable, id, port, address, keys}
	return n
}

func findKey(key string, this Node) (string, error) {
	for i := range this.Keys {
		if this.Keys[i].key == key {
			return this.Keys[i].value, nil
		}
	}
	return "", fmt.Errorf("Could not find key in node")
}
