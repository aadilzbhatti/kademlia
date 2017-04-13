package routing

import (
	"hash/fnv"
	"math/big"
)

//http://stackoverflow.com/questions/23192262/how-would-you-set-and-clear-a-single-bit-in-go

func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func CalculateDistance(id1, id2 []byte) *big.Int {
	a := new(big.Int).SetBytes(id1)
	b := new(big.Int).SetBytes(id2)
	dist := new(big.Int).Xor(a, b)
	return dist
}

/**
 * Returns the hash mod the number nodes in
 * the system
 *
 * @type {int}
 */
func hash(key string, size int) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % size
}

/**
 * Returns the bucket that Nodewith id=id2
 * would be in for the routing table of node
 * with id=id1
 *
 * @type {int}
 */
func getConflictingBit(id1, id2 []byte) int {
	for i := 0; i < len(id1); i++ {
		res := id1[i] ^ id2[i]

		// This is a byte. Need to get bit position.
		for j := 0; j < 8; j++ {
			if hasBit(int(res), uint(7-j)) {
				return 32 - (8*i + j) - 1
			}

		}
	}
	return 0 // same id
}
