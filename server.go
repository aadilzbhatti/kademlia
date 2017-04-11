package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

var barrier sync.WaitGroup
var nodeId int
var nodes = make([]net.Conn, 10)
var clients = make([]net.Conn, 10)
var lock = &sync.Mutex{}
var self Node
var port int = 3000
var myhost string

func startServer() {
  // port = getPort()

	// set up node ID
	fmt.Println("Starting!")
	host := "sp17-cs425-g26-0%d.cs.illinois.edu"
	myhost = getHostName()
	fmt.Println(myhost)
	bucket := make([][]TableEntry, 10)
	for i := 1; i < 10; i++ {
		otherHost := fmt.Sprintf(host, i)
		fmt.Println(myhost)
		if otherHost == myhost {
			nodeId = i
			break
		}
	}
	self := initializeNode(nodeId, 10, port, myhost)
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
		}
	}

	// handle connections
	barrier.Add(1)
	go handleSelf()
	barrier.Wait()
}

func handleSelf() {
	for {
		// periodically update k closest nodes for each key with KVPs (replicas)
	}
	defer barrier.Done()
}

func setupRPC() {
	node := new(Node)
	rpc.Register(node)
	l, e := net.Listen("tcp", ":3000")
	if e != nil {
		log.Fatal("Join listen error: ", e)
	}

	go rpc.Accept(l)
}

func makeJoinCall(self Node, host string) error {
	fmt.Printf("%v is self\n", self)
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
		return err
	}

	ja := JoinArgs{self.Id, self.Address, self.Port, "NEWNODE"}
	var reply string
	divCall := client.Go("Node.Join", &ja, &reply, nil)
	replyCall := <-divCall.Done
	log.Printf("Node %d joined the system: %s\n", self.Id, replyCall.Reply)
	return nil
}
