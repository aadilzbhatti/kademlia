package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func (d *DHT) Join(ja *JoinArgs, reply *Node) error {
	if string(ja.ID) == string(self.ID) {
		return nil
	}
	log.Printf("Node%v is trying to join Node%v\n", ja.ID, self.ID)

	// populate my buckets
	n := Node{ja.ID, ja.Hostname, ja.Port}
	myself := Node{self.ID, fmt.Sprintf("sp17-cs425-g26-0%d.cs.illinois.edu", self.ID), port}
	self.Rt.insert(&n)
	*reply = myself

	// send a message to the other nodes
	if ja.NewNode != "" {
		kClosest := self.lookup(ja.ID)
		for _, n := range kClosest {
			client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", n.Address, port))
			if err != nil {
				log.Fatal("Error in dial: ", err)
				return err
			}
			var reply Node
			divCall := client.Go("DHT.Join", ja, &reply, nil)
			replyCall := <-divCall.Done
			log.Println(replyCall.Reply)
		}
	}
	log.Printf("Node%v has joined Node%v\n", ja.ID, self.ID)
	return nil
}

func (d *DHT) Set(sa *SetArgs, reply *string) error {
	// find the Nodewhich has the key (via Find)
	kClosest := self.lookup(sa.KVP.Key)
	for _, n := range kClosest {
		fmt.Printf("Id -> %v\n", n.ID)
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", n.Address, port))
		if err != nil {
			log.Fatal("Error in dial: ", err)
			return err
		}
		var reply string
		divCall := client.Go("DHT.StoreKVP", sa, &reply, nil)
		replyCall := <-divCall.Done
		log.Println(replyCall.Reply)
	}

	// reply ACK to original Node
	*reply = "ACK"
	return nil
}

func (d *DHT) StoreKVP(sa *SetArgs, reply *string) error {
	self.Storage[string(sa.KVP.Key)] = string(sa.KVP.Value)
	*reply = "ACK"
	return nil
}

// func (d *DHT) Owners(oa *OwnerArgs, reply *[]Node) error {
// 	// find Nodewith given key
// 	for _, v := range self.Table {
// 		for _, b := range v {
// 			client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", b.Address, port))
// 			if err != nil {
// 				log.Fatal("Error in dialing:", err)
// 				return err
// 			}
//
// 			fa := FindArgs{oa.Key, math.Inf(1)}
// 			var fr *FindReply
// 			divCall := client.Go("Node.Find", &fa, fr, nil)
// 			replyCall := <-divCall.Done
// 			fmt.Println(replyCall)
//
// 			// if we have found k-closest Nodes, we reply with those Nodes
// 			if fr != nil {
// 				*reply = fr.Closest
// 				return nil
// 			}
// 		}
// 	}
//
// 	return nil
// }
//
// func (d *DHT) ListLocal(ll *ListLocalArgs, reply *[]KV) error {
// 	// reply with all keys in our Node
// 	*reply = self.Keys
//
// 	return nil
// }
