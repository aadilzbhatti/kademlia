package main

import (
	"fmt"
	"log"
	"math"
	"net/rpc"
)

func (n *Node) Join(ja *JoinArgs, reply *string) error {
  log.Println("In join ", ja.Id)
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
	if ja.NewNode != "" {
		na := ja
		na.NewNode = ""
		for _, v := range self.Table { // in reality we'd send this message to our k-closest, not all
			for _, b := range v {
				client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", b.Address, port))
				if err != nil {
					log.Fatal("Error in dialing:", err)
					return err
				}
				var remoteReply string
				divCall := client.Go("Node.Join", &na, &remoteReply, nil)
				replyCall := <-divCall.Done
				fmt.Println(replyCall)
			}
		}
	}

	// replicate keys TODO
	return nil
}

func (n *Node) Find(fa *FindArgs, reply *FindReply) error {
	// if our distance(id, hashed key) = 0
	// check ourselves for the key
	// reply FOUND if we found it
	// if not found, query alpha nodes in closest bucket (found by getBucket)
	// once found, reply KV to original node

	return nil
}

func (n *Node) Set(sa *SetArgs, reply *string) error {
	// find the node which has the key (via Find) (in reality we'd call FIND on the k-closest, not all)
	for _, v := range self.Table {
		for _, b := range v {
			client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", b.Address, port))
			if err != nil {
				log.Fatal("Error in dialing:", err)
				return err
			}
			fa := FindArgs{sa.KVP.Key, math.Inf(1)}
			var fr *FindReply
			divCall := client.Go("Node.Find", &fa, fr, nil)
			replyCall := <-divCall.Done
			fmt.Println(replyCall)

			// if we have found k-closest nodes, we store the KV-pair in them (thus storing in k-closest for replication)
			if fr != nil {
				for _, n := range fr.Closest {
					n.storeKVP(fr.KVP)
				}
				break
			}

		}
	}

	// reply ACK to original node
	*reply = "ACK"
	return nil
}

func (n *Node) Owners(oa *OwnerArgs, reply *[]Node) error {
	// find node with given key
	for _, v := range self.Table {
		for _, b := range v {
			client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", b.Address, port))
			if err != nil {
				log.Fatal("Error in dialing:", err)
				return err
			}

			fa := FindArgs{oa.Key, math.Inf(1)}
			var fr *FindReply
			divCall := client.Go("Node.Find", &fa, fr, nil)
			replyCall := <-divCall.Done
			fmt.Println(replyCall)

			// if we have found k-closest nodes, we reply with those nodes
			if fr != nil {
				*reply = fr.Closest
				return nil
			}
		}
	}

	return nil
}

func (n *Node) ListLocal(ll *ListLocalArgs, reply *[]KV) error {
	// reply with all keys in our node
	*reply = self.Keys

	return nil
}
