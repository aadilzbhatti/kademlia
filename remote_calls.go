package main

import (
	"fmt"
	"log"
	"net/rpc"
	"sort"
)

func (d *DHT) Join(ja *JoinArgs, reply *Node) error {
	if string(ja.ID) == string(self.ID) {
		return nil
	}

	// populate my buckets
	n := Node{ja.ID, ja.Hostname, ja.Port}
	myself := Node{self.ID, fmt.Sprintf("sp17-cs425-g26-0%d.cs.illinois.edu", self.ID[0]), port}
	*reply = myself

	// send a message to the other nodes
	if ja.NewNode != "" {
		self.Rt.insert(&n)
		kClosest := self.lookup(ja.ID)
		for _, n := range kClosest {
			client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", n.Address, port))
			defer client.Close()
			if err != nil {
				log.Println("Error in dial: ", err)
				return err
			}
			var reply Node
			err = client.Call("DHT.Join", ja, &reply)
			if err != nil {
				log.Println("Error in join: ", err)
				return err
			}
		}
	}
	return nil
}

func (d *DHT) Set(sa *SetArgs, reply *string) error {
	// find the k closest Nodes which have the key
	kClosest := self.lookup(sa.KVP.Key)
	for _, n := range kClosest {
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", n.Address, port))
		if err != nil {
			log.Println("Error in dial: ", err)
			return err
		}
		defer client.Close()
		var reply string
		err = client.Call("DHT.StoreKVP", sa, &reply)
    if err != nil {
      log.Println("Error in calling store: ", err)
    }
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

func (d *DHT) Find(target *[]byte, reply *KV) error {
	nodes := self.lookup(*target)
	fmt.Println("FIND")
	for _, v := range nodes {
		fmt.Printf("%v, ", v.ID)
	}
	fmt.Println("")
	for _, v := range nodes {
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", v.Address, port))
		defer client.Close()
		if err != nil {
			log.Println("Error in find RPC: ", err)
			continue
		}
		key := string(*target)
		err = client.Call("DHT.GetKVP", &key, &reply)
		if err != nil {
			log.Println("Error in getting the key: ", err)
			continue
		}
		if reply.Value != nil {
			return nil
		}
		break
	}
	*reply = KV{*target, nil}
	return nil
}

func (d *DHT) GetKVP(key *string, reply *KV) error {
	value := self.Storage[*key]
	*reply = KV{[]byte(*key), []byte(value)}
	return nil
}

func (d *DHT) Owners(key *[]byte, reply *[]*Node) error {
	// find Nodes with given key
	*reply = self.lookup(*key)
	return nil
}

func (d *DHT) ListLocal(args *string, reply *[]KV) error {
	list := make([]KV, 0)
	keys := make([]string, 0)
	for k, _ := range self.Storage {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		list = append(list, KV{[]byte(k), []byte(self.Storage[k])})
	}

	// reply with all keys in our Node
	*reply = list
	return nil
}
