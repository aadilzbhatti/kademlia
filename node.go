package main

import (
  "math"
  "fmt"
  "net"
)

type TableEntry struct {
  Id int
  Port int
  Address net.IP
}

type RoutingTable struct {
  buckets [][]TableEntry
}

type Node struct {
  Table RoutingTable
  Id int
  Port int
  Address net.IP
}

// func initializeRoutingTable(id int, numNodes int) RoutingTable {
//   k := int(math.Ceil(math.Log2(float64(numNodes))))
//   buckets := make([][]int, k)
//   table := RoutingTable{buckets}
//   return table
// }
//
// func initializeNode(id int, numNodes int) Node {
//   routingTable := initializeRoutingTable(id, numNodes)
//   n := Node{routingTable, id}
//   return n
// }

func getConflictingBit(id1 int, id2 int) int {
  if id1 == 0 {
    return 1
  }
  maxId := math.Max(float64(id1), float64(id2))
  numBits := uint(math.Ceil(math.Log2(maxId)))
  for i := numBits; i > 0; i-- {
    msb_1 := id1 & (1 << i)
    msb_2 := id2 & (1 << i)
    if msb_1 != msb_2 {
      return int(numBits) - int(i)
    }
  }
  return int(numBits)
}

func distance(id1 int, id2 int) int {
  return id1 ^ id2
}

func find(key string) {
  
}

func main() {
  b := getConflictingBit(6, 0)
  d := distance(6, 0)
  fmt.Println(b, d)
}
