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

func (d *DHT) remoteLookup(n *Node, target []byte) []*Node {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", n.Address, port))
	if err != nil {
		log.Println("Error in remote lookup: ", err)
		return nil
	}
	defer client.Close()
	var reply []*Node
	err = client.Call("DHT.KClosestRPC", &target, &reply)
	if err != nil {
		log.Println("Error in calling RPC: ", err)
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
	if len(kclosest) == 0 {
		return nil
	}
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
		if kclosest_r1 == nil {
			log.Printf("Node %v has failed\n", shortlist.nodes[i].ID)
			log.Printf("%v -> shortlist\n", shortlist.nodes)
			for i, v := range self.Rt.buckets {
				fmt.Printf("%d, %v\n", i, v)
				for j, n := range v.bucket {
					fmt.Printf("%d, %v\n", j, n)
          if n != nil {
            if string(n.ID) == string(shortlist.nodes[i].ID) {
  						self.Rt.buckets[i].bucket = append(self.Rt.buckets[i].bucket[:j], self.Rt.buckets[i].bucket[j + 1:]...) // DELETE NODE FROM BUCKET IF DEAD
  					}
          }
				}
			}
			if shortlist.nodes != nil {
				for _, v := range shortlist.nodes {
					fmt.Printf("%v, ", v.ID)
				}
				fmt.Println("")
				shortlist.nodes = append(shortlist.nodes[:i], shortlist.nodes[i+1:]...)
				for _, v := range shortlist.nodes {
					fmt.Printf("%v, ", v.ID)
				}
				fmt.Println("")
			}
			i++
			continue
		}
		i++
		//check for null
		numresponses++
		for _, v := range kclosest_r1 {
			if !seen[string(v.ID)] {
				shortlist.nodes = append(shortlist.nodes, v)
			}
		}
		//sort.Sort(shortlist) // now update it with the new kclosest nodes.

		//kclosest = shortlist.nodes[:ksize]
		// We will get a list of k closest nodes according to
		// the closest node.
		// Now check if we found a closer Node to the current closest or not.
		if bytes.Compare(closest.ID, kclosest_r1[0].ID) == 0 {
			// found the best one
			//kclosest has the final result now.
			sort.Sort(shortlist)
			if len(kclosest) >= ksize {
				kclosest = shortlist.nodes[:ksize]
			}
			return kclosest
		}
		// update closest node
		if calculateDistance(target, kclosest[0].ID).Cmp(calculateDistance(target, closest.ID)) == -1 {
			closest = kclosest[0] // update closest
		}
	}
	sort.Sort(shortlist)
	var m map[string]bool = make(map[string]bool)
	var flist []*Node
	for _, n := range shortlist.nodes {
		if !(m[string(n.ID)]) {
			flist = append(flist, n)
			m[string(n.ID)] = true
		}
	}
	if len(flist) < ksize {
		return flist
	}
	kclosest = flist[:ksize]
	return kclosest
}
