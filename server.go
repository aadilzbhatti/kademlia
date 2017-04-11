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
var nodeId int
var nodes = make([]net.Conn, 10)
var clients = make([]net.Conn, 10)
var lock = &sync.Mutex{}
var self Node
var port int = 3000
var myhost string
var T int = 45

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
	bucket := make([][]TableEntry, 10)
	for i := 1; i < 10; i++ {
		otherHost := fmt.Sprintf(host, i)
		if otherHost == myhost {
			nodeId = i
			break
		}
	}
	self = initializeNode(nodeId, 10, port, myhost)
	self.Table = bucket

	// set up RPCs
	go setupRPC()

	// add nodes {1, 2, 3} \ nodeID to buckets
	for i := 1; i < 4; i++ {
		if nodeId != i {
			err := makeJoinCall(self, fmt.Sprintf(host, i))
			if err != nil {
				log.Fatal("Failed to join node:", i)
			}
			log.Println("Joined in!")
		}
	}

	// continuously republish keys
	barrier.Add(1)
	go handleSelf()
	barrier.Wait()
}

/**
 * Republishes keys every T time units
 */
func republishKeys() {
	for {
    time.Sleep(T * time.Second)
		// periodically update k closest nodes for each key with KVPs (replicas)
	}
	defer barrier.Done()
}

/**
 * Sets up the RPC channel for this node
 */
func setupRPC() {
	node := new(Node)
	rpc.Register(node)
	l, e := net.Listen("tcp", ":3000")
	if e != nil {
		log.Fatal("Join listen error: ", e)
	}

	go rpc.Accept(l)
}

/**
 * Wrapper for a node to join a node at a hostname
 */
func makeJoinCall(self Node, host string) error {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
		return err
	}

	// make the RPC
	ja := JoinArgs{self.Id, self.Address, self.Port, "NEWNODE"}
	var reply Node
	divCall := client.Go("Node.Join", &ja, &reply, nil)
	replyCall := <-divCall.Done
	log.Println(replyCall)

	// insert the new guy into my bucket
	bucket := getBucket(reply.Id, self.Id)
	entry := TableEntry{reply.Id, reply.Port, reply.Hostname}
	lock.Lock()
	self.Table[bucket] = append(self.Table[bucket], entry)

	return nil
}
