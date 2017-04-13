package dht

import (
	"bytes"
	"fmt"
	"kademlia/internal/dht/internal/routing"
	"log"
	"net/rpc"
	"sort"
)

type DHT struct {
	Rt      *routing.RoutingTable
	ID      []byte
	Storage map[string]string // Key Value Store
}

func createDHT(id []byte) *DHT {
	dht := &DHT{
		Rt:      routing.NewRoutingTable(id),
		ID:      id,
		Storage: make(map[string]string),
	}
	return dht
}

func (d *DHT) remoteLookup(n *routing.Node, target []byte) []*routing.Node {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", n.Address, port))
	if err != nil {
		log.Println("Error in remote lookup: ", err)
		return nil
	}
	defer client.Close()
	var reply []*routing.Node
	err = client.Call("DHT.KClosestRPC", &target, &reply)
	if err != nil {
		log.Println("Error in calling RPC: ", err)
		return nil
	}
	return reply
}

func (d *DHT) KClosestRPC(target *[]byte, reply *[]*routing.Node) error {
	r := self.Rt.GetKClosest(*target)
	*reply = r.Nodes
	return nil
}

func (d *DHT) Lookup(target []byte) []*routing.Node {
	kclosest := d.Rt.GetKClosest(target).Nodes
	if len(kclosest) == 0 {
		return nil
	}
	closest := kclosest[0]

	seen := make(map[string]bool)
	// K closest nodes on this node.
	// Now have to query everyone.
	// Now alpha = 1 so query the first node.
	var shortlist routing.NeighborList
	shortlist.ID = target
	shortlist.Nodes = kclosest // refine this after every iteration.
	numresponses := 0
	i := 0
	for (numresponses) < routing.KSize && i < (shortlist.Len()) {
		if seen[string(shortlist.Nodes[i].ID)] {
			i++
			continue
		}

		seen[string(shortlist.Nodes[i].ID)] = true
		kclosest_r1 := d.remoteLookup(shortlist.Nodes[i], target)
		if kclosest_r1 == nil {
			log.Printf("Node %v has failed\n", shortlist.Nodes[i].ID)
			for k, v := range self.Rt.Buckets {
				for j, n := range v.Bucket {
					if n != nil {
						if string(n.ID) == string(shortlist.Nodes[i].ID) {
							log.Printf("Deleted node %v from the system\n", n.ID)
							self.Rt.Buckets[k].Bucket = append(self.Rt.Buckets[k].Bucket[:j], self.Rt.Buckets[k].Bucket[j+1:]...)
						}
					}
				}
			}
			if shortlist.Nodes != nil {
				shortlist.Nodes = append(shortlist.Nodes[:i], shortlist.Nodes[i+1:]...)
			}
			i++
			continue
		}
		i++
		//check for null
		numresponses++
		for _, v := range kclosest_r1 {
			if !seen[string(v.ID)] {
				shortlist.Nodes = append(shortlist.Nodes, v)
			}
		}
		//sort.Sort(shortlist) // now update it with the new kclosest nodes.

		//kclosest = shortlist.Nodes[:ksize]
		// We will get a list of k closest nodes according to
		// the closest node.
		// Now check if we found a closer Node to the current closest or not.
		if bytes.Compare(closest.ID, kclosest_r1[0].ID) == 0 {
			// found the best one
			//kclosest has the final result now.
			sort.Sort(shortlist)
			if len(kclosest) >= routing.KSize {
				kclosest = shortlist.Nodes[:routing.KSize]
			}
			return kclosest
		}
		// update closest node
		if routing.CalculateDistance(target, kclosest[0].ID).Cmp(routing.CalculateDistance(target, closest.ID)) == -1 {
			closest = kclosest[0] // update closest
		}
	}
	sort.Sort(shortlist)
	var m map[string]bool = make(map[string]bool)
	var flist []*routing.Node
	for _, n := range shortlist.Nodes {
		if !(m[string(n.ID)]) {
			flist = append(flist, n)
			m[string(n.ID)] = true
		}
	}
	if len(flist) < routing.KSize {
		return flist
	}
	kclosest = flist[:routing.KSize]
	return kclosest
}
