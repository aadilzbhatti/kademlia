package main

import (
  "net"
)

type JoinArgs struct {
  Id int
  hostname string
  port int
  Ip net.IP
}

type FindArgs struct {
  Key string
  PrevClosestDistance int
}

type SetArgs struct {

}

type OwnerArgs struct {

}

func (n *Node) Join(ja *JoinArgs, newNode *string, reply *string) error {
  // populate my buckets
  id := ja.Id
  bucket := getConflictingBit(id, self.Id) - 1
  entry := TableEntry{id, ja.Port, ja.Hostname}
  for _, v := range self.Table[bucket] {
    if v.Id == id {
      *reply = "ACK"
    }
  }
  lock.Lock()
  self.Table[bucket] = append(self.Table[bucket], entry)
  lock.Unlock()
  *reply = "ACK"

  // send a message to the other nodes
  if newNode != "" {
    for _, v := range self.Table {
      for _, b := range v {
        client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:8085", b.Hostname))
        if err != nil {
          log.Fatal("Error in dialing:", err)
          return
        }
        divCall := client.Go("Node.Join", ja, "", nil)
        replyCall := <-divCall.Done
      }
    }
  }
  // replicate keys TODO
}

func (n *Node) Find(fa *FindArgs, reply *string) error {
    // query alpha nodes in closest bucket
    // if
}
