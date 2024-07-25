package logic

import (
	"fmt"
	"sort"

	"github.com/namnv2496/go-consistent-hashing/internal/model"
	"github.com/namnv2496/go-consistent-hashing/internal/utils"
)

type Nodes interface {
	AddNode(ip string) error
	Route(key string, value any) (string, error)
	RemoveNode(ip string) error
	PrintNodes()
	GetValueInNode(ip string) (map[int]any, error)
}

type nodes struct {
	Server     map[int]*model.Node
	TotalNodes int
}

var keys []int

// NewNodes creates a new nodes instance
func NewNodes() Nodes {
	return &nodes{
		Server:     make(map[int]*model.Node),
		TotalNodes: 0,
	}
}

// AddNode adds a node to the server
func (n *nodes) AddNode(ip string) error {
	hash := utils.HashFunction(ip)
	if _, exists := n.Server[hash]; exists {
		return fmt.Errorf("ip already exists")
	}
	n.Server[hash] = &model.Node{
		Ip:       ip,
		Position: hash,
		Value:    make(map[int]any),
	}
	n.TotalNodes++
	keys = n.sortedKeys()
	return nil
}

// Route finds the appropriate server for a given key
func (n *nodes) Route(key string, value any) (string, error) {
	if n.TotalNodes == 0 {
		return "", fmt.Errorf("no nodes available")
	}
	hash := utils.HashFunction(key)
	nextNode := n.findNextNode(hash)
	if nextNode != nil {
		nextNode.Value[hash] = value
		fmt.Printf("Key %s is routed hash %d to server: %s\n", key, hash, nextNode.Ip)
		return nextNode.Ip, nil
	}

	// If no appropriate node is found, wrap around to the first node
	if n.TotalNodes == 0 {
		return "", fmt.Errorf("no node was found")
	}
	fmt.Printf("Route key %s, hash: %d to begin node: %s \n", key, hash, n.Server[keys[0]].Ip)
	n.Server[keys[0]].Value[hash] = value
	return n.Server[keys[0]].Ip, nil
}

// RemoveNode removes a node from the server
func (n *nodes) RemoveNode(ip string) error {
	hash := utils.HashFunction(ip)
	if _, exists := n.Server[hash]; !exists {
		return fmt.Errorf("ip does not exists")
	}
	var result []int
	for _, k := range keys {
		if k != hash {
			result = append(result, k)
		}
	}
	keys = result
	n.distributeKeysOnRemovedNode(hash)
	delete(n.Server, hash)
	n.TotalNodes--
	return nil
}

func (n *nodes) PrintNodes() {

	for i, k := range keys {
		fmt.Println("Node ", i+1, " Ip = ", n.Server[k].Ip, " position = ", n.Server[k].Position)
	}
}

func (n *nodes) GetValueInNode(ip string) (map[int]any, error) {
	for _, node := range n.Server {
		if node.Ip == ip {
			return node.Value, nil
		}
	}
	return nil, fmt.Errorf("ip is not exist %s", ip)
}

func (n *nodes) distributeKeysOnRemovedNode(hash int) error {
	nodeValues := n.Server[hash].Value
	nextNode := n.findNextNode(hash)

	fmt.Printf("distributeKeysOnRemovedNode %s to next node %s", n.Server[hash].Ip, nextNode.Ip)
	fmt.Println("Vlaue: ", nodeValues)
	for key, value := range nodeValues {
		nextNode.Value[key] = value
	}
	return nil
}

// sortedKeys returns the sorted keys of the server map
func (n *nodes) sortedKeys() []int {
	keys := make([]int, 0, len(n.Server))
	for k := range n.Server {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// find next node
func (n *nodes) findNextNode(hash int) *model.Node {
	// keys := n.sortedKeys()
	for _, k := range keys {
		if hash <= k {
			return n.Server[k]
		}
	}
	return nil
}
