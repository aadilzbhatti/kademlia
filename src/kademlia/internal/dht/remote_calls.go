package dht

import (
	"fmt"
	"kademlia/internal/dht/internal/routing"
	"log"
	"net/rpc"
	"sort"
)

/**
 * RPC argument structs
 */
type JoinArgs struct {
	ID       []byte
	Hostname string
	Port     int
	NewNode  string
}

type FindArgs struct {
	Target []byte
	Node   routing.Node
}

type SetArgs struct {
	KVP KV
}

func (d *DHT) Join(ja *JoinArgs, reply *routing.Node) error {
	if string(ja.ID) == string(self.ID) {
		return nil
	}

	// populate my buckets
	n := routing.Node{ja.ID, ja.Hostname, ja.Port}
	myself := routing.Node{self.ID, fmt.Sprintf("sp17-cs425-g26-0%d.cs.illinois.edu", self.ID[0]), port}
	*reply = myself

	// send a message to the other nodes
	if ja.NewNode != "" {
		self.Rt.Insert(&n)
		kClosest := self.Lookup(ja.ID)
		for _, n := range kClosest {
			client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", n.Address, port))
			if err != nil {
				log.Println("Error in dial: ", err)
				return err
			}
			var reply routing.Node
			err = client.Call("DHT.Join", ja, &reply)
			if err != nil {
				log.Println("Error in join: ", err)
				return err
			}
			client.Close()
		}
	}
	return nil
}

func (d *DHT) Set(sa *SetArgs, reply *string) error {
	// find the k closest Nodes which have the key
	kClosest := self.Lookup(sa.KVP.Key)
	for _, n := range kClosest {
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", n.Address, port))
		if err != nil {
			log.Println("Error in dial: ", err)
			return err
		}
		var reply string
		err = client.Call("DHT.StoreKVP", sa, &reply)
		if err != nil {
			log.Println("Error in calling store: ", err)
		}
		client.Close()
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
	nodes := self.Lookup(*target)
	for _, v := range nodes {
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", v.Address, port))
		if err != nil {
			log.Println("Error in find RPC: ", err)
			continue
		}
		key := string(*target)
		err = client.Call("DHT.GetKVP", &key, &reply)
		if err != nil {
			log.Println("Error in getting the key: ", err)
			client.Close()
			continue
		}
		if reply.Value != nil {
			client.Close()
			return nil
		}
		client.Close()
	}
	return fmt.Errorf("No such key in system")
}

func (d *DHT) GetKVP(key *string, reply *KV) error {
	if value, ok := self.Storage[*key]; ok {
		*reply = KV{[]byte(*key), []byte(value)}
	} else {
		*reply = KV{[]byte(*key), nil}
	}
	return nil
}

func (d *DHT) Owners(key *[]byte, reply *[]*routing.Node) error {
	// find Nodes with given key
	*reply = self.Lookup(*key)
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
