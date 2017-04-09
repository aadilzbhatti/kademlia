package main

import (
	"net"
)

/**
 * Node structs
 */
type TableEntry struct {
	Id      int
	Port    int
	Address string
}

type KV struct {
	key   string
	value string
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
	Ip       net.IP
}

type FindArgs struct {
	Key                 string
	PrevClosestDistance int
}

type SetArgs struct {
	KVP KV
}

type OwnerArgs struct {
	Key string
}

type ListLocalArgs struct{}
