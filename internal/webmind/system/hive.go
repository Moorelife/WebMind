package system

import "net"

// Hive defines the data required to define a Hive system.
type Hive struct {
	Address net.Addr `json:"address"` // the address of the node.
}

func NewHive(address net.Addr) *Node {
	node := Node{Address: address}
	return &node
}
