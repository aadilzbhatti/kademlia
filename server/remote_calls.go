package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func (n *Node) Join(ja *JoinArgs, newNode *string, reply *string) error {
	// populate my buckets
	id := ja.Id
	bucket := getBucket(id, self.Id)
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
	if *newNode != "" {
		for _, v := range self.Table {
			for _, b := range v {
				client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:8085", b.Address))
				if err != nil {
					log.Fatal("Error in dialing:", err)
					return err
				}
				divCall := client.Go("Node.Join", ja, "", nil)
				replyCall := <-divCall.Done
				fmt.Println(replyCall)
			}
		}
	}

	// replicate keys TODO
	return nil
}

func (n *Node) Find(fa *FindArgs, reply *KV) error {
	// if our distance(id, hashed key) = 0
	// check ourselves for the key
	// reply FOUND if we found it
	// if not found, query alpha nodes in closest bucket (found by getBucket)
	// once found, reply KV to original node

	return nil
}

func (n *Node) Set(sa *SetArgs, reply *string) error {
	// find the node which has the key (via Find)
	// if it is ours, set K -> V
	// update replicas
	// reply ACK to original node

	return nil
}

func (n *Node) Owners(oa *OwnerArgs, reply *Node) error {
	// find node with given key
	// reply with that node

	return nil
}

func (n *Node) ListLocal(ll *ListLocalArgs, reply *[]KV) error {
	// reply with all keys in our node

	return nil
}
