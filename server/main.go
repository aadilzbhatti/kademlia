package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

var barrier sync.WaitGroup
var nodeId int
var nodes = make([]net.Conn, 10)
var clients = make([]net.Conn, 10)
var lock = &sync.Mutex{}
var self Node

func main() {
	fmt.Println("Starting!")
	var host string = "sp17-cs425-g26-0%d.cs.illinois.edu"
	hostname := getHostName()
	fmt.Println(hostname)
	barrier.Add(2)
	for i := 1; i < 10; i++ {
		otherHost := fmt.Sprintf(host, i)
		if otherHost != hostname {
			if i < 4 {
				fmt.Printf("Connecting to node %d\n", i)
				go connectToHost(otherHost)
			}
		} else {
			nodeId = i
		}
	}
	barrier.Wait()
	conn, _ := net.Listen("tcp", hostname)
	self = initializeNode(nodeId, 10, 8080, hostname)
	setupJoinRPC()
	// set up RPC stuff
	barrier.Add(1)
	go handleConn(conn)
	barrier.Wait()
}

func connectToHost(host string) {
	for {
		conn, _ := net.Dial("tcp", host)
		if conn == nil {
			continue
		} else {
			nodes = append(nodes, conn)
			break
		}
	}
	defer barrier.Done()
}

func handleConn(ln net.Listener) {
	for {
		conn, _ := ln.Accept()
		lock.Lock()
		clients = append(clients, conn)
		lock.Unlock()
		// go handleClient(conn)
	}
}

func setupJoinRPC() {
	node := new(Node)
	rpc.Register(node)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":8085")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
