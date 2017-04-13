package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

var barrier sync.WaitGroup
var nodeId []byte
var lock = &sync.Mutex{}
var self DHT
var port int = 3000
var hostname string

/**
 * Starts up the server allowing for nodes to join the
 * distributed hash table
 */
func startServer() {

	// set up NodeID
	fmt.Println("Starting!")
	host := "sp17-cs425-g26-0%d.cs.illinois.edu"
	hostname = getHostName()
	for i := 1; i < 10; i++ {
		otherHost := fmt.Sprintf(host, i)
		if otherHost == hostname {
			nodeId = []byte(string(i))
			break
		}
	}
	self = *(createDHT(nodeId))

	// set up RPCs
	go setupRPC()

	// add nodes {1, 2, 3} \ nodeID to buckets
	for i := 1; i < 10; i++ {
		if string(nodeId) != string(i) {
			go makeJoinCall(self, fmt.Sprintf(host, i), i)
		}
	}

	// continuously republish keys
	barrier.Add(1)
	go republishKeys()
	barrier.Wait()
}

/**
 * Republishes keys every T time units
 */
func republishKeys() {
  client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	defer client.Close()
  if err != nil {
		log.Printf("Error in republish dial: ", err)
  }

	for {
		time.Sleep(20 * time.Second)
    log.Printf("Republishing at node %d\n", nodeId)
		// periodically update k closest nodes for each key with KVPs (replicas)

    for k, v := range self.Storage {
      sa := SetArgs{KV{[]byte(k), []byte(v)}}
			var reply string
			err = client.Call("DHT.Set", &sa, &reply)
			if err != nil {
				log.Printf("Failed to republish on node %d\n", nodeId)
			}
    }
	}
	defer barrier.Done()
}

/**
 * Sets up the RPC channel for this node
 */
func setupRPC() {
	dht := new(DHT)
	rpc.Register(dht)

	l, e := net.Listen("tcp", ":3000")
	if e != nil {
		log.Println("Join listen error: ", e)
	}

	go rpc.Accept(l)
}

/**
 * Wrapper for a Node to join a Node at a hostname
 */
func makeJoinCall(self DHT, host string, id int) {
	log.Printf("Node %v is trying to join Node %v\n", self.ID, id)
  for {
    client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
    if err != nil {
      continue
    }

    // make the RPC
    ja := JoinArgs{self.ID, hostname, port, "NEWNODE"}
    var reply Node
    err = client.Call("DHT.Join", &ja, &reply)
		if err != nil {
			log.Println("Error in initial join: ", err)
			break
		}

		log.Printf("Node %v has joined Node %v\n", self.ID, id)
		client.Close()
		break
  }
}
