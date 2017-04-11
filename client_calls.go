package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func clientSet(key string, value string) error {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		log.Fatal("Could not connect to server:", err)
		return err
	}
	KVP := KV{[]byte(key), []byte(value)}
	sa := SetArgs{KVP}
	var reply string
	divCall := client.Go("DHT.Set", &sa, &reply, nil)
	replyCall := <-divCall.Done
	if replyCall.Error != nil {
		return replyCall.Error
	}
	log.Printf("Successfully SET %s=%s\n", key, value)
	return nil
}

func clientGet(key string) error {
  client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
  if err != nil {
    log.Fatal("Could not connect to server:", err)
    return err
  }
	target := []byte(key)
	var reply KV
	err = client.Call("DHT.Find", &target, &reply)
	if err != nil {
		log.Fatal("Error in find: ", err)
		return err
	}
	if reply.Value != nil {
		log.Printf("FOUND %s=%s\n", key, string(reply.Value))
	} else {
		log.Printf("Could not find key %s\n", key)
	}
	return nil
}

func clientOwners(key string) error {
	client, _ := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	target := []byte(key)
	var reply []node
	err := client.Call("DHT.Owners", &target, &reply)
	if err != nil {
		log.Fatal("Error in owners: ", err)
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
}

func clientListLocal() error {
	client, _ := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	var reply []KV
	err := client.Call("DHT.ListLocal", nil, &reply)
	if reply != nil {
		log.Printf("LISTING ALL KEYS AT NODE %d\n", self.ID)
		for _, v := range reply {
			log.Printf("%s=%s\n", string(v.Key), string(v.Value))
		}
	} else {
		log.Printf("No keys located at this node.")
	}
}

func clientBatch() error {
	
}
