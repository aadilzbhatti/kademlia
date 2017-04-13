package dht

import (
	"fmt"
	"kademlia/internal/dht/internal/routing"
	"log"
	"net/rpc"
)

func Set(key string, value string) error {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		log.Println("Could not connect to server:", err)
		return err
	}
	defer client.Close()
	KVP := KV{[]byte(key), []byte(value)}
	sa := SetArgs{KVP}
	var reply string
	err = client.Call("DHT.Set", &sa, &reply)
	if err != nil {
		log.Println("Error in set: ", err)
		return err
	}
	log.Printf("Successfully SET %s=%s\n", key, value)
	return nil
}

func Get(key string) error {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		log.Println("Could not connect to server:", err)
		return err
	}
	defer client.Close()
	target := []byte(key)
	var reply KV
	err = client.Call("DHT.Find", &target, &reply)
	if err != nil {
		log.Println("Error in find: ", err)
		return err
	}
	if reply.Value != nil {
		log.Printf("FOUND %s=%s\n", key, string(reply.Value))
	} else {
		log.Printf("Could not find key %s\n", key)
	}
	return nil
}

func Owners(key string) error {
	client, _ := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	defer client.Close()
	target := []byte(key)
	var reply []routing.Node
	err := client.Call("DHT.Owners", &target, &reply)
	if err != nil {
		log.Println("Error in owners: ", err)
		return err
	}
	if reply != nil {
		log.Printf("LISTING OWNERS OF KEY %s:\n", key)
		for _, v := range reply {
			log.Printf("Node %d\n", v.ID)
		}
	} else {
		log.Printf("No owners for this key.")
	}
	return nil
}

func ListLocal() error {
	client, _ := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	defer client.Close()
	var reply []KV
	var args string
	err := client.Call("DHT.ListLocal", &args, &reply)
	if err != nil {
		log.Println("Error in list local: ", err)
		return err
	}
	if reply != nil {
		log.Printf("LISTING ALL KEYS AT NODE %d\n", self.ID)
		for _, v := range reply {
			log.Printf("%s=%s\n", string(v.Key), string(v.Value))
		}
	} else {
		log.Printf("No keys located at this node.")
	}
	return nil
}
