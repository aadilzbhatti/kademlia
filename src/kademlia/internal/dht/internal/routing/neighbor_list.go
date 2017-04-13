package routing

type NeighborList struct {
	Nodes []*Node
	ID    []byte // will be the key ID.
	// Implement len, swap and less functions to get sorting functionality
}

func (s NeighborList) Len() int {
	return len(s.Nodes)
}
func (s NeighborList) Swap(i, j int) {
	s.Nodes[i], s.Nodes[j] = s.Nodes[j], s.Nodes[i]
}

func (s NeighborList) Less(i, j int) bool {
	dist1 := CalculateDistance(s.Nodes[i].ID, s.ID)
	dist2 := CalculateDistance(s.Nodes[j].ID, s.ID)

	if dist1.Cmp(dist2) == 1 {
		return false
	}
	return true
}
