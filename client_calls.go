package main

import (
	"log"
	"fmt"
	"net/rpc"
	"math"
)

func clientSet(key string, value string) error {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		log.Fatal("Could not connect to server:", err)
		return err
	}
	KVP := KV{key, value}
	sa := SetArgs{KVP}
	var reply string
	divCall := client.Go("Node.Set", &sa, &reply, nil)
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
	sa := FindArgs{Key, math.Inf(1)}
  var reply FindReply
	divCall := client.Go("Node.Get", &sa, &reply, nil)
	replyCall := <-divCall.Done
	if replyCall.Error != nil {
		return replyCall.Error
	}
	log.Printf("FOUND %s=%s\n", key, reply.KVP.Value)
	return nil
}
