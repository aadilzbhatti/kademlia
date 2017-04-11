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
	log.Printf("FOUND %s=%s\n", key, string(reply.Value))
	return nil
}
