package main

import (
	"bytes"
	"sort"
)

type DHT struct {
	Rt      *RoutingTable
	ID      []byte
	Storage map[string]string // Key Value Store
}

func createDHT(id []byte) {
	dht := &DHT{
		Rt:      NewRoutingTable(id),
		ID:      id,
		Storage: make(map[string]string),
	}
}

func (d *DHT) remoteLookup(n *node, target []byte) {
	// make rpc call
}

func (d *DHT) lookup(target []byte) {
	kclosest := d.Rt.getKClosest(target).nodes
	closest := kclosest[0]

	var seen map[string]bool
	// K closest nodes on this node.
	// Now have to query everyone.
	// Now alpha = 1 so query the first node.
	var shortlist neighborList
	shortlist.ID = target
	shortlist.nodes = kclosest // refine this after every iteration.
	numresponses := 0
	i := 0
	for (numresponses) < ksize && i < len(shortlist) {
		if seen[string(shortlist[i].ID)] {
			i++
			continue
		}

		seen[string(shortlist[i].ID)] = true
		i++
		kclosest_r1 := d.remoteLookup(shortlist[i], target)
		//check for null
		numresponses++
		shortlist.nodes = append(shortlist.nodes, kclosest_r1...)
		//sort.Sort(shortlist) // now update it with the new kclosest nodes.

		//kclosest = shortlist.nodes[:ksize]
		// We will get a list of k closest nodes according to
		// the closest node.
		// Now check if we found a closer node to the current closest or not.
		if bytes.Compare(closest.ID, kclosest_r1[0] == 0) {
			// found the best one
			//kclosest has the final result now.
			sort.Sort(shortlist)
			kclosest = shortlist.nodes[:ksize]
			return kclosest
		}
		// update closest node
		if calculateDistance(target, kclosest[0]) < calculateDistance(target, closest) {
			closest = kclosest[0] // update closest
		}
	}
	sort.Sort(shortlist)
	kclosest = shortlist.nodes[:ksize]
	return kclosest
}
