package system

import "net"

// Web defines the data required to define a Web system.
type Node struct {
	Address net.Addr `json:"address"` // the address of the node.
}

func NewNode(address net.Addr) *Node {
	node := Node{Address: address}
	return &node
}
