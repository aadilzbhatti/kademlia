package main

/**
 * Node structs
 */

type Node struct {
	ID      []byte
	Address string
	Port    int
}

type KV struct {
	Key   []byte
	Value []byte
}
