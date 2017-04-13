package main

/**
 * Nodestructs
 */

type Node struct {
	ID      []byte
	Address string
	Port    int
}

type Kbucket struct {
	Size   int
	bucket []*Node
	seenMap map[string]bool
}

type KV struct {
	Key   []byte
	Value []byte
}

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
	Node   Node
}

type SetArgs struct {
	KVP KV
}
