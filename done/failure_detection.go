package main

import (
	"encoding/gob"
	"fmt"
	set "github.com/deckarep/golang-set"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Failure struct {
	Id            int
	Counter       int
	TimeSincePing int64
	Failed        bool
}

var servers = make([]net.Conn, 0, 0)
var mutex = &sync.Mutex{}
var seenMessages = set.NewSet()
var wg sync.WaitGroup
var nodeNumber int
var numNodes int
var failureNode int

var hostname string
var messageHost string

var T_gossip int = 500
var T_fail int = 10000
var T_cleanup int = 10000
var failureHost string
var failureServers []Failure
var failureMutex = &sync.Mutex{}
var failureConns []net.UDPConn

type failureCallback func(nodeNumber int)

func startFailureDetection(numNodes int, hostname string, callback failureCallback) {
	fmt.Println("NumNodes is ", numNodes)
	fmt.Println("Launching failure server...")
	fmt.Println(time.Now())

	gob.Register(&Failure{})
	gob.Register(&net.UDPConn{})
	failureHost = fmt.Sprintf("%s:8084", hostname)

	failureServers = make([]Failure, numNodes)
	failureConns = make([]net.UDPConn, numNodes)
	connectToFailureCluster(numNodes)
	fmt.Println("Connected to failure servers")

	wg.Add(1)
	go failureDetection(callback)
	wg.Wait()
	fmt.Println("All connections closed. Shutting down server...")
}

func failureDetection(callback failureCallback) {
	fmt.Println("Starting FD")
	time.Sleep(10 * time.Second)
	numNodes = len(failureServers)
	for {
		failureServers[failureNode-1].Failed = false // our Nodehas not failed if we are here
		gossipIdx := rand.Intn(numNodes)             // pick a random node
		if gossipIdx == failureNode-1 {
			continue
		}
		failureMutex.Lock()
		gossipMember := failureServers[gossipIdx]
		failureMutex.Unlock()
		if gossipMember.Failed {
			continue
		}
		if gossipMember.Id == 0 && !gossipMember.Failed {
			failureMutex.Lock()
			failureServers[gossipIdx].Failed = true
			failureMutex.Unlock()
		}

		time.Sleep(time.Duration(T_gossip) * time.Millisecond)
		failureMutex.Lock()
		fmt.Printf("%v - Failures are %v\n", (time.Now().UnixNano() / int64(time.Millisecond)), failureServers)
		failureServers[failureNode-1].Counter++
		failureServers[failureNode-1].TimeSincePing = time.Now().Unix()
		failureMutex.Unlock()

		go writeBuf(gossipMember, callback)
		go readBuf(failureServers[failureNode-1], callback)
	}
	defer wg.Done()
}

func writeBuf(gossipMember Failure, callback failureCallback) {
	if gossipMember.Id == 0 {
		return
	}
	// fmt.Printf("PINGING %d\n", gossipMember.Id)
	failureMutex.Lock()
	conn := failureConns[gossipMember.Id-1]
	failureMutex.Unlock()
	enc := gob.NewEncoder(&conn)
	failureMutex.Lock()
	err := enc.Encode(failureServers)
	if err != nil {
		failureServers[gossipMember.Id-1].Failed = true
		fmt.Printf("%v - Node%d has failed\n", (time.Now().UnixNano() / int64(time.Millisecond)), gossipMember.Id)
		msg := fmt.Sprintf("Node%d has failed\n", gossipMember.Id)
		handleFailedNode(gossipMember.Id, callback)
		writeToServers(msg)
	}
	failureMutex.Unlock()
}

func readBuf(self Failure, callback failureCallback) {
	failureMutex.Lock()
	conn := failureConns[self.Id-1]
	failureMutex.Unlock()
	dec := gob.NewDecoder(&conn)
	var newFailures []Failure
	err := dec.Decode(&newFailures)
	if err != nil {
	}
	for i := range failureServers {
		diff := time.Now().Unix() - failureServers[i].TimeSincePing
		if diff > int64(T_fail) && !failureServers[i].Failed {
			failureServers[i].Failed = true
			fmt.Printf("%v - Node%d has failed\n", (time.Now().UnixNano() / int64(time.Millisecond)), failureServers[i].Id)
			msg := fmt.Sprintf("Node%d has failed\n", failureServers[i].Id)
			writeToServers(msg)
			go handleFailedNode(i, callback)
		}
	}
	failureMutex.Lock()
	failureServers = mergeFailures(newFailures, failureServers)
	failureMutex.Unlock()
}

func handleFailedNode(Nodeint, callback failureCallback) {
	time.Sleep(time.Duration(T_cleanup) * time.Millisecond)
	callback(node)
	fmt.Println("Removed node", node)
}

func mergeFailures(newFailures []Failure, failures []Failure) []Failure {
	if len(newFailures) == 0 {
		return failures
	}
	ret := make([]Failure, numNodes)
	for i := 0; i < numNodes; i++ {
		if newFailures[i].Id == 0 {
			ret[i] = newFailures[i]
			ret[i].Failed = true
		}
		if !failures[i].Failed && !newFailures[i].Failed {
			if newFailures[i].Counter > failures[i].Counter {
				ret[i] = newFailures[i]
			} else {
				ret[i] = failures[i]
			}
		} else {
			ret[i] = failures[i]
			ret[i].Failed = true
		}
	}
	return ret
}

func writeToServers(message string) {
	data := message + "\n"
	mutex.Lock()
	for i := range servers {
		servers[i].Write([]byte(data))
	}
	mutex.Unlock()
}

func connectToFailureCluster(numNodes int) {
	wg.Add(numNodes)
	for i := 1; i <= numNodes; i++ {
		serverName := fmt.Sprintf("sp17-cs425-g26-0%d.cs.illinois.edu:8084", i)
		if serverName != failureHost {
			go connectToFailureServer(serverName, hostname, i, false)
		} else {
			go connectToFailureServer(serverName, hostname, i, true)
			failureNode = i
		}
	}
	wg.Wait()
}

func setUpSelf(nodeNum int, serverName string) Failure {
	ServerAddr, err := net.ResolveUDPAddr("udp", serverName)
	if err != nil {
		fmt.Println("Error in resolve: ", err)
	}
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		fmt.Println("Error in listen: ", err)
	}
	failureConns[nodeNum-1] = *ServerConn
	return Failure{nodeNum, 0, time.Now().Unix(), false}
}

func connectToFailureServer(Nodestring, hostname string, nodeNum int, self bool) {
	if self {
		failureMutex.Lock()
		ret := setUpSelf(nodeNum, node)
		failureServers[nodeNum-1] = ret
		failureMutex.Unlock()
	} else {
		for {
			ServerAddr, err := net.ResolveUDPAddr("udp", node)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			LocalAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:0", hostname))
			if err != nil {
				fmt.Println("Error: ", err)
			}
			conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
			if conn == nil {
				continue
			} else {
				if err != nil {
					fmt.Println("Error: ", err)
				}
				newMember := Failure{nodeNum, 0, time.Now().Unix(), false}
				failureMutex.Lock()
				failureConns[nodeNum-1] = *conn
				failureServers[nodeNum-1] = newMember
				failureMutex.Unlock()
				break
			}
		}
	}
	defer wg.Done()
}
