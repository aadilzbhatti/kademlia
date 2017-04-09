package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"net"
)

var alpha int = 3

type TableEntry struct {
	Id      int
	Port    int
	Address string
}

type KV struct {
	key   string
	value string
}

type Node struct {
	Table   [][]TableEntry
	Id      int
	Port    int
	Address string
	Keys    []KV
}

// func find(key string, size int, start Node) {
// 	hashVal := hash(key, size)
// 	bucket := getConflictingBit(start.Id, hashVal) - 1
// 	nodes := start.Table[bucket]
// }

func initializeRoutingTable(id int, numNodes int) [][]TableEntry {
  k := int(math.Ceil(math.Log2(float64(numNodes))))
  buckets := make([][]TableEntry, k)
  return buckets
}

func initializeNode(id int, numNodes int, port int, address string) Node {
  routingTable := initializeRoutingTable(id, numNodes)
	keys := make([]KV, 5)
  n := Node{routingTable, id, port, address}
  return n
}

func hash(key string, size int) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % size
}

func findKey(key string, this Node) (string, error) {
	for i := range this.Keys {
		if this.Keys[i].key == key {
			return this.Keys[i].value, nil
		}
	}
	return "", fmt.Errorf("Could not find key in node")
}

func getConflictingBit(id1 int, id2 int) int {
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
	return int(numBits)
}

func distance(id1 int, id2 int) int {
	return id1 ^ id2
}
