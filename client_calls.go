package main

import (
  "net"
  "net/rpc"
  "fmt"
  "log"
)

func clientSet(key string, value string) error {
  client, err := net.Dial("tcp", fmt.Sprintf("%s:8080", hostname))
  if err != nil {
    log.Fatal("Could not connect to server:", err)
    return err
  }
  KVP := KV{key, value}
  sa := SetArgs{KVP}
  var reply string
  divCall := client.Go("Node.Set", &sa, &reply)
  replyCall := <-divCall.Done
  if replyCall.Error != nil {
    return replyCall.Error
  }
  fmt.Printf("Successfully SET %s=%s\n", key, value)
  return nil
}
