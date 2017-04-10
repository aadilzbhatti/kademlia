package main

import (
	"fmt"
	"hash/fnv"
	"math"
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
	return strings.Trim(string(out), "\n")
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

// func (n *Node) getkClosestContacts(target []byte, sender []byte) *sortedList {
// 	// Get all the neighbors - Find K closest ones
// 	// Since we have only 10 nodes it's okay to search all of them.
// 	// This is more computation than network cost so its okay.
// 	bucketIndex := getBucket(target, sender)
// 	// Collect other buckets
// 	sortedNodes := &sortedList{}
// 	for i := 0; i < len(n.Table[bucketIndex]); i++ {
// 		if n.Table[bucketIndex][i].ID == n.ID {
// 			continue
// 		}
// 		sortedNodes.nodes = append(sortedNodes.nodes, rt.buckets[bucketIndex][i])
// 	}
//
// 	for i := bucketIndex; i >= 0; i++ {
// 		for j := 0; j < len(rt.buckets[i]); j++ {
// 			if rt.buckets[bucketIndex][j].ID == rt.ID {
// 				continue
// 			}
// 			sortedNodes.nodes = append(sortedNodes.nodes, rt.buckets[i][j])
// 		}
// 	}
//
// 	for i := bucketIndex + 1; i < b; i++ {
// 		for j := 0; j < len(rt.buckets[i]); j++ {
// 			if rt.buckets[bucketIndex][j].ID == rt.ID {
// 				continue // dont add yourself to the list? Should we do this or not?
// 			}
// 			sortedNodes.nodes = append(sortedNodes.nodes, rt.buckets[i][j])
// 		}
// 	}
//
// 	sort.Sort(sortedNodes)
//
// 	return sortedNodes
//
// 	// Collected all nodes.
// 	// Break when we have 3 nodes.
// 	// Sort all the nodes in the order of distance
//
// }
//
// //return closest_nodes // probably need to sort by distance or something.
// func generateID() {
// 	ID := make([]byte, 20)
// 	_, err := rand.Read(ID)
// 	return result
// }
//
// func getConflictingBit(id1, id2 []byte) int {
// 	for i := 0; i < len(id1); i++ {
// 		res := id1[i] ^ id2[i]
//
// 		// This is a byte. Need to get bit position.
// 		for j := 0; j < 8; j++ {
// 			if hasBit(res, uint(7-j)) {
// 				return 160 - (8*i + j) - 1
// 			}
//
// 		}
// 	}
// 	return -1 // same id
// }
//
// //http://stackoverflow.com/questions/23192262/how-would-you-set-and-clear-a-single-bit-in-go
//
// func hasBit(n int, pos uint) bool {
// 	val := n & (1 << pos)
// 	return (val > 0)
// }
//
// type sortedList struct {
// 	nodes []*nodes
// 	ID    []byte // will be the local ID.
// 	// Implement len, swap and less functions to get sorting functionality
// }
//
// func (s *sortedList) Len() int {
// 	return len(s.nodes)
// }
// func (s *sortedList) Swap(i, j int) {
// 	s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
// }
//
// func (s *sortedList) Less(i, j) bool {
// 	dist1 := calculateDistance(s.nodes[i], s.ID)
// 	dist2 := calculateDistance(s.nodes[j], s.ID)
//
// 	if dist1.Cmp(dist2) == 1 {
// 		return false
// 	}
// 	return true
// }
//
// func calculateDistance(id1, id2 []byte) *big.Int {
// 	a := new(big.Int).setBytes(id1)
// 	b := new(big.Int).setBytes(id2)
// 	dist := new(big.Int).Xor(a, b)
// 	return dist
// }
