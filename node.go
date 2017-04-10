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

func (n *Node) findKey(key string) (string, error) {
	for i := range n.Keys {
		if n.Keys[i].Key == key {
			return n.Keys[i].Value, nil
		}
	}
	return "", fmt.Errorf("Could not find key in node")
}

func (n *Node) storeKVP(KVP KV) {
	n.Keys = append(n.Keys, KVP)
}
