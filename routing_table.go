package main

import (
	"math/big"
	"sort"
	"fmt"
)

const ksize = 3
const IDLength = 4

var seenMap map[string]bool = make(map[string]bool)

// Contains the implementation of kbuckets and the table itself.
func NewBucket(size int) *Kbucket {
	return &Kbucket{
		Size:   size,
		bucket: make([]*Node, ksize),
	}
}

func (k *Kbucket) addNode(n *Node) {
	// check if already exists
	// if it exists move to tail of the list
	exists := seenMap[string(n.ID)]
	//pos, exists := k.checkNodeExists(n) // should return pos and bool
	if exists {
		// // move to the end.
	} else {
		fmt.Printf("Adding %v\n", n)
		if len(k.bucket) == k.Size {
			// pinging stuff
			k.bucket = k.bucket[1:]
			k.bucket = append(k.bucket, n)
		} else {
			seenMap[string(n.ID)] = true
			k.bucket = append(k.bucket, n)
		}
	}
}

func (k *Kbucket) checkNodeExists(n *Node) (int, bool) {
	// for i := range k.bucket {
	// 	if string(k.bucket[i].ID) == string(n.ID) {
	// 		return i, true
	// 	}
	// }
	return -1, false
}

type RoutingTable struct {
	ID      []byte
	buckets [IDLength * 8]*Kbucket
}

func NewRoutingTable(id []byte) *RoutingTable {
	var b [IDLength * 8]*Kbucket
	for i := 0; i < IDLength*8; i++ {
		b[i] = NewBucket(ksize)
	}
	rt := &RoutingTable{
		ID:      id,
		buckets: b,
	}
	return rt
}

func (rt *RoutingTable) insert(n *Node) {
	bucketIndex := rt.findBucketIndex(n.ID)
	rt.buckets[bucketIndex].addNode(n)
}

func (rt *RoutingTable) findBucketIndex(target []byte) int {
	hashedValue := hash(string(target), 10)
	tgt := []byte(string(hashedValue))
	bucketIndex := getConflictingBit(tgt, rt.ID)
	return bucketIndex
}

func (rt *RoutingTable) getKClosest(target []byte) *neighborList {
	bucketIndex := rt.findBucketIndex(target)
	var closest neighborList
	closest.ID = target
	for i := bucketIndex; i >= 0; i-- {
		// keep going till you find ksize nodes.
		for j := 0; j < len(rt.buckets[i].bucket); j++ {
			if closest.Len() >= ksize {
				sort.Sort(closest)
				return &closest
			}
			if rt.buckets[i].bucket[j] != nil {
				closest.nodes = append(closest.nodes, rt.buckets[i].bucket[j]) // adding node
			}
		}

	}
	// Check Right subtrees now
	for i := bucketIndex + 1; i < IDLength*8; i++ {
		// keep going till you find ksize nodes.
		for j := 0; j < len(rt.buckets[i].bucket); j++ {
			if closest.Len() >= ksize {

				sort.Sort(closest)
				return &closest
			}
			if rt.buckets[i].bucket[j] != nil {
				closest.nodes = append(closest.nodes, rt.buckets[i].bucket[j]) // adding node

			}
		}

	}

	sort.Sort(closest)
	return &closest
}

//http://stackoverflow.com/questions/23192262/how-would-you-set-and-clear-a-single-bit-in-go

func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

type neighborList struct {
	nodes []*Node
	ID    []byte // will be the key ID.
	// Implement len, swap and less functions to get sorting functionality
}

func (s neighborList) Len() int {
	return len(s.nodes)
}
func (s neighborList) Swap(i, j int) {
	s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
}

func (s neighborList) Less(i, j int) bool {
	dist1 := calculateDistance(s.nodes[i].ID, s.ID)
	dist2 := calculateDistance(s.nodes[j].ID, s.ID)

	if dist1.Cmp(dist2) == 1 {
		return false
	}
	return true
}

func calculateDistance(id1, id2 []byte) *big.Int {
	a := new(big.Int).SetBytes(id1)
	b := new(big.Int).SetBytes(id2)
	dist := new(big.Int).Xor(a, b)
	return dist
}
