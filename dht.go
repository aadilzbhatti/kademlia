package main

import (
	"bytes"
	"fmt"
	"log"
	"net/rpc"
	"sort"
)

type DHT struct {
	Rt      *RoutingTable
	ID      []byte
	Storage map[string]string // Key Value Store
}

func createDHT(id []byte) *DHT {
	dht := &DHT{
		Rt:      NewRoutingTable(id),
		ID:      id,
		Storage: make(map[string]string),
	}
	return dht
}

// register dht
func (d *DHT) remoteLookup(n *Node, target []byte) []*Node {
	fmt.Printf("%s:%d\n", n.Address, port)
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", n.Address, port))
	if err != nil {
		log.Fatal("Error in remote lookup: ", err)
		return nil
	}
	var reply []*Node
	err = client.Call("DHT.KClosestRPC", &target, &reply)
	if err != nil {
		log.Fatal("Error in calling RPC: ", err)
		return nil
	}
	return reply
}

func (d *DHT) KClosestRPC(target *[]byte, reply *[]*Node) error {
	r := self.Rt.getKClosest(*target)
	*reply = r.nodes
	return nil
}

func (d *DHT) lookup(target []byte) []*Node {
	kclosest := d.Rt.getKClosest(target).nodes
	closest := kclosest[0]

	seen := make(map[string]bool)
	// K closest nodes on this node.
	// Now have to query everyone.
	// Now alpha = 1 so query the first node.
	var shortlist neighborList
	shortlist.ID = target
	shortlist.nodes = kclosest // refine this after every iteration.
	numresponses := 0
	i := 0
	for (numresponses) < ksize && i < (shortlist.Len()) {
		if seen[string(shortlist.nodes[i].ID)] {
			i++
			continue
		}

		seen[string(shortlist.nodes[i].ID)] = true
		kclosest_r1 := d.remoteLookup(shortlist.nodes[i], target)
		i++
		//check for null
		numresponses++
		shortlist.nodes = append(shortlist.nodes, kclosest_r1...)
		//sort.Sort(shortlist) // now update it with the new kclosest nodes.

		//kclosest = shortlist.nodes[:ksize]
		// We will get a list of k closest nodes according to
		// the closest node.
		// Now check if we found a closer Nodeto the current closest or not.
		if bytes.Compare(closest.ID, kclosest_r1[0].ID) == 0 {
			// found the best one
			//kclosest has the final result now.
			sort.Sort(shortlist)
			kclosest = shortlist.nodes[:ksize]
			return kclosest
		}
		// update closest node
		if calculateDistance(target, kclosest[0].ID).Cmp(calculateDistance(target, closest.ID)) == -1 {
			closest = kclosest[0] // update closest
		}
	}
	sort.Sort(shortlist)
	kclosest = shortlist.nodes[:ksize]
	return kclosest
}
