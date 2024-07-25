package main

import (
	"fmt"

	"github.com/namnv2496/go-consistent-hashing/internal/logic"
)

func main() {
	// Create a new nodes instance
	n := logic.NewNodes()

	fmt.Println("================================== Add Nodes ==========================================")
	// Add nodes
	if err := n.AddNode("192.168.1.1"); err != nil {
		fmt.Println(err)
	}
	if err := n.AddNode("192.168.1.2"); err != nil {
		fmt.Println(err)
	}
	if err := n.AddNode("192.168.1.4"); err != nil {
		fmt.Println(err)
	}

	n.PrintNodes()
	fmt.Println("================================== Add keys ==========================================")
	// Route a key
	n.Route("my-key", 1)
	n.Route("my-key1", "phúc long")
	value := struct {
		name string
		age  int
	}{
		name: "nguyễn văn a",
		age:  15,
	}
	n.Route("my-key2", value)
	n.Route("my-key3", "highland")
	fmt.Println("================================= Print values in each node =========================================")

	nodeValue, _ := n.GetValueInNode("192.168.1.1")
	fmt.Println("value in server 192.168.1.1: ", nodeValue)
	nodeValue, _ = n.GetValueInNode("192.168.1.2")
	fmt.Println("value in server 192.168.1.2: ", nodeValue)
	nodeValue, _ = n.GetValueInNode("192.168.1.4")
	fmt.Println("value in server 192.168.1.4: ", nodeValue)
	fmt.Println("=================================== Remove a node 192.168.1.4 ====================================")
	// Remove a node
	if err := n.RemoveNode("192.168.1.4"); err != nil {
		fmt.Println(err)
	}
	fmt.Println("========================== add  key again after remove ===================================")
	fmt.Println("========== Previous key ''my-key1' route to '192.168.1.4' ==========")
	fmt.Println("========== After remove node key ''my-key1' route to '192.168.1.1' ==========")
	// Route the key again after removing a node
	n.Route("my-key5", "The coffee house")

	nodeValue, _ = n.GetValueInNode("192.168.1.1")
	fmt.Println("value in server 192.168.1.1: ", nodeValue)
	nodeValue, _ = n.GetValueInNode("192.168.1.2")
	fmt.Println("value in server 192.168.1.2: ", nodeValue)
	fmt.Println("============================================================================")

	_, err := n.GetValueInNode("192.168.1.2")
	if err != nil {
		fmt.Println("The ip 192.168.1.2 is not exist")
	}
}
