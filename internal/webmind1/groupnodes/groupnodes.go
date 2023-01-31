package groupnodes

import (
	"github.com/Moorelife/WebMind/internal/webmind1/localnode"
)

type GroupNodes struct {
	LocalNodes []localnode.LocalNode //`json:"localnodes"`
}

func NewGroupNodes(address string, ports []int) *GroupNodes {
	if len(ports) != 3 {
		panic("groupnodes should have three ports")
	}
	var creation = GroupNodes{
		LocalNodes: make([]localnode.LocalNode, 3),
	}
	for key, value := range ports {
		creation.LocalNodes[key] = *localnode.NewLocalNode(address, value)
	}
	return &creation
}
