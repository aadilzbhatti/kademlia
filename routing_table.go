package dht

import (
	"math/big"
	"net"
	"sort"
)

// Contains the implementation of kbuckets and the table itself.

type node struct {
	ID   []byte
	IP   net.IP
	Port int
}

type Kbucket struct {
	Size   int
	bucket []*node
}

func NewBucket(size int) {
	return &Kbucket{
		Size:   size,
		bucket: make([]*node, 0),
	}
}

func (k *Kbucket) addNode(n *Node) {
	// check if already exists
	// if it exists move to tail of the list
	pos, exists := checkNodeExists(n) // should return pos and bool
	if exists {
		// move to the end.
		k.bucket = append(k.bucket[:pos], k.bucket[pos+1:])
		append(k.bucket, n)
	} else {
		if len(k.bucket) == k.Size {
			// pinging stuff
		} else {
			append(k.bucket, n)
		}
	}
}

func (k *Kbucket) checkNodeExists(n *Node) {
	for i := range k.bucket {
		if k.bucket[i].ID == n.ID {
			return i, true
		}
	}
	return -1, false
}

type RoutingTable struct {
	ID      []byte
	buckets [IDLength * 8]*Kbucket
}

func NewRoutingTable(id []byte) *RoutingTable {
	l := IDLength * 8
	rt := &RoutingTable{
		ID:      id,
		buckets: [l]*Kbucket{},
	}
	return rt
}

func (rt *RoutingTable) insert(n *node) {
	bucketIndex := rt.findBucketIndex(n.ID)
	rt.buckets[bucketIndex].add(n)
}

func (rt *RoutingTable) findBucketIndex(target []byte) int {
	bucketIndex := getConflictingBit(target, rt.ID)
}

func (rt *RoutingTable) getKClosest(target []byte) *neighborList {
	bucketIndex := rt.findBucketIndex(target)
	var closest neighborList
	closest.ID = target
	for i := bucketIndex; i >= 0; i-- {
		// keep going till you find ksize nodes.
		for j := 0; j < len(rt.buckets[i]); j++ {
			if closest.Len() >= ksize {
				return sort.Sort(closest)
			}
			closest.nodes = append(closest.nodes, rt.buckets[i][j]) // adding node
		}

	}
	// Check Right subtrees now
	for i := bucketIndex + 1; i < IDLength*8; i++ {
		// keep going till you find ksize nodes.
		for j := 0; j < len(rt.buckets[i]); j++ {
			if closest.Len() >= ksize {
				return sort.Sort(closest)
			}
			closest.nodes = append(closest.nodes, rt.buckets[i][j]) // adding node
		}

	}

	return sort.Sort(closest)
}

func getConflictingBit(id1, id2 []byte) int {
	for i := 0; i < len(id1); i++ {
		res := id1[i] ^ id2[i]

		// This is a byte. Need to get bit position.
		for j := 0; j < 8; j++ {
			if hasBit(res, uint(7-j)) {
				return 160 - (8*i + j) - 1
			}

		}
	}
	return 0 // same id
}

//http://stackoverflow.com/questions/23192262/how-would-you-set-and-clear-a-single-bit-in-go

func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

type neighborList struct {
	nodes []*nodes
	ID    []byte // will be the key ID.
	// Implement len, swap and less functions to get sorting functionality
}

func (s *neighborList) Len() int {
	return len(s.nodes)
}
func (s *neighborList) Swap(i, j int) {
	s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
}

func (s *neighborList) Less(i, j) bool {
	dist1 := calculateDistance(s.nodes[i], s.ID)
	dist2 := calculateDistance(s.nodes[j], s.ID)

	if dist1.Cmp(dist2) == 1 {
		return false
	}
	return true
}

func calculateDistance(id1, id2 []byte) *big.Int {
	a := new(big.Int).setBytes(id1)
	b := new(big.Int).setBytes(id2)
	dist := new(big.Int).Xor(a, b)
	return dist
}
