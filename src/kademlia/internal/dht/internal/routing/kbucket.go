package routing

type Node struct {
	ID      []byte
	Address string
	Port    int
}

type Kbucket struct {
	Size    int
	Bucket  []*Node
	seenMap map[string]bool
}

const KSize = 3

// Contains the implementation of kbuckets and the table itself.
func NewBucket(size int) *Kbucket {
	return &Kbucket{
		Size:    size,
		Bucket:  make([]*Node, KSize),
		seenMap: make(map[string]bool),
	}
}

func (k *Kbucket) addNode(n *Node) {
	// check if already exists
	// if it exists move to tail of the list
	exists := k.seenMap[string(n.ID[0])]
	if exists {
		// // move to the end.
	} else {
		if len(k.Bucket) == k.Size {
			// pinging stuff
			k.Bucket = k.Bucket[1:]
			k.Bucket = append(k.Bucket, n)
		} else {
			k.Bucket = append(k.Bucket, n)
		}
		k.seenMap[string(n.ID[0])] = true
	}
}
