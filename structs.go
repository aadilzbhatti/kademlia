package main

/**
 * Node structs
 */
type TableEntry struct {
	Id      int
	Port    int
	Address string
}

type KV struct {
	Key   string
	Value string
}

type Node struct {
	Table   [][]TableEntry
	Id      int
	Port    int
	Address string
	Keys    []KV
}

/**
 * RPC argument structs
 */
type JoinArgs struct {
	Id       int
	Hostname string
	Port     int
	NewNode  string
  HostNode *Node
}

type FindArgs struct {
	Key                 string
	PrevClosestDistance float64
}

type FindReply struct {
	KVP     KV
	Closest []Node
}

type SetArgs struct {
	KVP KV
}

type OwnerArgs struct {
	Key string
}

type ListLocalArgs struct{}
