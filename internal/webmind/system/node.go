package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
)

// Node defines the data required to define a node system.
type Node struct {
	Address net.TCPAddr `json:"address"` // the address of the node.
}

// NewNode creates a new Node structure and returns a pointer to it
func NewNode(address net.TCPAddr) *Node {
	node := Node{Address: address}
	return &node
}

// ToJSON converts the Node struct to indented JSON.
func (n *Node) ToJSON() string {
	b, err := json.Marshal(n)
	if err != nil {
		panic(fmt.Sprintf("Marshal failed: %v", err))
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, " ", "  ")
	if err != nil {
		panic(fmt.Sprintf("Indent failed: %v", err))
	}

	return fmt.Sprintf("%s", &out)
}
