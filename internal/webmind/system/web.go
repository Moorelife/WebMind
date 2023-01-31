package system

import "net"

// Web defines the data required to define a Web system.
type Web struct {
	Address net.Addr `json:"address"` // the address of the node.
}

func NewWeb(address net.Addr) *Node {
	node := Node{Address: address}
	return &node
}
