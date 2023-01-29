package trinity

import "github.com/Moorelife/WebMind/internal/localnode"

type Trinity struct {
	LocalNodes []localnode.LocalNode //`json:"localnodes"`
}

func NewTrinity(address string, ports []int) *Trinity {
	if len(ports) != 3 {
		panic("trinity should have three ports")
	}
	var creation = Trinity{
		LocalNodes: make([]localnode.LocalNode, 3),
	}
	for key, value := range ports {
		creation.LocalNodes[key] = *localnode.NewLocalNode(address, value)
	}
	return &creation
}
