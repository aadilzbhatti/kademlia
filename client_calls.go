package main

import (
	"fmt"
	"log"
	"net/rpc"
  "os"
	"bufio"
	"strings"
	"regexp"
)

func clientSet(key string, value string) error {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		log.Println("Could not connect to server:", err)
		return err
	}
  defer client.Close()
	KVP := KV{[]byte(key), []byte(value)}
	sa := SetArgs{KVP}
	var reply string
	err = client.Call("DHT.Set", &sa, &reply)
	if err != nil {
    log.Println("Error in set: ", err)
		return err
	}
	log.Printf("Successfully SET %s=%s\n", key, value)
	return nil
}

func clientGet(key string) error {
  client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
  if err != nil {
    log.Println("Could not connect to server:", err)
    return err
  }
	target := []byte(key)
	var reply KV
	err = client.Call("DHT.Find", &target, &reply)
	if err != nil {
		log.Println("Error in find: ", err)
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
	var reply []Node
	err := client.Call("DHT.Owners", &target, &reply)
	if err != nil {
		log.Println("Error in owners: ", err)
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
	return nil
}

func clientListLocal() error {
	client, _ := rpc.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	var reply []KV
	var args string
	err := client.Call("DHT.ListLocal", &args, &reply)
	if err != nil {
		log.Println("Error in list local: ", err)
		return err
	}
	if reply != nil {
		log.Printf("LISTING ALL KEYS AT NODE %d\n", self.ID)
		for _, v := range reply {
			log.Printf("%s=%s\n", string(v.Key), string(v.Value))
		}
	} else {
		log.Printf("No keys located at this node.")
	}
	return nil
}

func clientBatch(fname string) error {
	r, _ := regexp.Compile("(GET) (.*)|(SET) (.*) (.*)|(LIST_LOCAL)|(OWNERS) (.*)|(BATCH (.*))")
  if file, err := os.Open(fname); err == nil {
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      line := strings.TrimSpace(scanner.Text())
			if r.MatchString(line) {
				res := r.FindStringSubmatch(line)
				for i := range res {
					if i > 0 && res[i] != "" {
						runCommand(res, i)
						break
					}
				}
			}
    }
  } else {
		return err
	}
	return nil
}
