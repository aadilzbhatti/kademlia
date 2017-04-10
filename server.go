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

func startServer() {
	// set up node ID
	fmt.Println("Starting!")
	host := "sp17-cs425-g26-0%d.cs.illinois.edu"
	hostname := getHostName()
	fmt.Println(hostname)
	bucket := make([][]TableEntry, 10)
	for i := 1; i < 10; i++ {
		otherHost := fmt.Sprintf(host, i)
		if otherHost == hostname {
			nodeId = i
			break
		}
	}
	self := initializeNode(nodeId, 10, 8080, hostname)
	self.Table = bucket
	fmt.Println(self)

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

	for {
		l, e := net.Listen("tcp", fmt.Sprintf("%s:8080", hostname))
		if e != nil {
			log.Fatal("Join listen error:", e)
		}

		go rpc.Accept(l)
	}
}

func makeJoinCall(self Node, host string) error {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:8080", host))
	if err != nil {
		log.Fatal("Erorr in dialing:", err)
		return err
	}

	ja := JoinArgs{self.Id, self.Address, self.Port, "NEWNODE"}
	var reply string
	divCall := client.Go("Node.Join", ja, &reply, nil)
	replyCall := <-divCall.Done
	fmt.Println(replyCall)
	return nil
}
