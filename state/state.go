package state

import (
	"sort"
	"sync"
	"time"

	"github.com/zerocruft/capacitor"
)

var (
	nodes map[string]trackedNode
	mutex sync.Mutex
)

func init() {
	nodes = map[string]trackedNode{}
	mutex = sync.Mutex{}
}

type trackedNode struct {
	Node        capacitor.FluxNode
	LastPing    time.Time
	Connections int
}

func AddNode(node capacitor.FluxNode, numberOfConnections int) {
	mutex.Lock()
	defer mutex.Unlock()
	nodes[node.ClientEndpoint] = trackedNode{
		Node:        node,
		LastPing:    time.Now(),
		Connections: numberOfConnections,
	}
}

func RemoveNode(key string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(nodes, key)
}

func ToNodeSlice() []capacitor.FluxNode {
	mutex.Lock()
	defer mutex.Unlock()
	peers := []capacitor.FluxNode{}

	for _, v := range nodes {
		peers = append(peers, v.Node)
	}
	return peers
}

func CopyOfNodes() []trackedNode {
	nodesCopy := []trackedNode{}
	mutex.Lock()
	defer mutex.Unlock()
	for _, x := range nodes {
		nodesCopy = append(nodesCopy, x)
	}
	return nodesCopy
}

func GetNodeWithLightestLoad() capacitor.FluxNode {
	x := CopyOfNodes()
	if len(x) == 0 {
		return capacitor.FluxNode{}
	}

	sort.Sort(ByAmountOfConnections(x))
	return x[0].Node
}

type ByAmountOfConnections []trackedNode

func (n ByAmountOfConnections) Len() int {
	return len(n)
}
func (s ByAmountOfConnections) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByAmountOfConnections) Less(i, j int) bool {
	return s[i].Connections < s[j].Connections
}
