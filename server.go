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
var myhost string

/**
 * Starts up the server allowing for nodes to join the
 * distributed hash table
 */
func startServer() {

	// set up node ID
	fmt.Println("Starting!")
	host := "sp17-cs425-g26-0%d.cs.illinois.edu"
	myhost = getHostName()
	fmt.Println(myhost)
	for i := 1; i < 10; i++ {
		otherHost := fmt.Sprintf(host, i)
		if otherHost == myhost {
			nodeId = []byte(string(i))
			break
		}
	}
	self = *(createDHT(nodeId))

	// set up RPCs
	go setupRPC()

	// add nodes {1, 2, 3} \ nodeID to buckets
	for i := 1; i < 4; i++ {
		if string(nodeId) != string(i) {
			err := makeJoinCall(self, fmt.Sprintf(host, i))
			if err != nil {
				log.Fatal("Failed to join node:", i)
			}
			log.Println("Joined in!")
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
	for {
    time.Sleep(45 * time.Second)
		// periodically update k closest nodes for each key with KVPs (replicas)
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
		log.Fatal("Join listen error: ", e)
	}

	go rpc.Accept(l)
}

/**
 * Wrapper for a node to join a node at a hostname
 */
func makeJoinCall(self DHT, host string) error {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
		return err
	}

	// make the RPC
	ja := JoinArgs{self.ID, hostname, port, "NEWNODE"}
	var reply node
	divCall := client.Go("Node.Join", &ja, &reply, nil)
	replyCall := <-divCall.Done
	log.Println(replyCall)

	// insert the new guy into my bucket
	self.Rt.insert(&reply)

	return nil
}
