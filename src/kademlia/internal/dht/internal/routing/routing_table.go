package routing

import (
	"sort"
)

const IDLength = 4

type RoutingTable struct {
	ID      []byte
	Buckets [IDLength * 8]*Kbucket
}

func NewRoutingTable(id []byte) *RoutingTable {
	var b [IDLength * 8]*Kbucket
	for i := 0; i < IDLength*8; i++ {
		b[i] = NewBucket(KSize)
	}
	rt := &RoutingTable{
		ID:      id,
		Buckets: b,
	}
	return rt
}

func (rt *RoutingTable) Insert(n *Node) {
	bucketIndex := rt.findBucketIndex(n.ID)
	rt.Buckets[bucketIndex].addNode(n)
}

func (rt *RoutingTable) findBucketIndex(target []byte) int {
	hashedValue := hash(string(target), 10)
	tgt := []byte(string(hashedValue))
	bucketIndex := getConflictingBit(tgt, rt.ID)
	return bucketIndex
}

func (rt *RoutingTable) GetKClosest(target []byte) *NeighborList {
	bucketIndex := rt.findBucketIndex(target)
	var closest NeighborList
	closest.ID = target
	for i := bucketIndex; i >= 0; i-- {
		// keep going till you find routing.KSize nodes.
		for j := 0; j < len(rt.Buckets[i].Bucket); j++ {
			if closest.Len() >= KSize {
				sort.Sort(closest)
				return &closest
			}
			if rt.Buckets[i].Bucket[j] != nil {
				closest.Nodes = append(closest.Nodes, rt.Buckets[i].Bucket[j]) // adding node
			}
		}

	}
	// Check Right subtrees now
	for i := bucketIndex + 1; i < IDLength*8; i++ {
		// keep going till you find routing.KSize nodes.
		for j := 0; j < len(rt.Buckets[i].Bucket); j++ {
			if closest.Len() >= KSize {

				sort.Sort(closest)
				return &closest
			}
			if rt.Buckets[i].Bucket[j] != nil {
				closest.Nodes = append(closest.Nodes, rt.Buckets[i].Bucket[j]) // adding node

			}
		}

	}

	sort.Sort(closest)
	return &closest
}
