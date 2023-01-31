package main

import (
	"encoding/json"
	"github.com/Moorelife/WebMind/internal/webmind1/localnode"
	"log"
)

// WebMind in its current state is JUST A LEARNING EXPERIMENT,
// and as such can not be expected to be fit for any given purpose.
// Please understand that you use the program at your own risk!!!

type RedundantNode struct {
	Monitor         localnode.LocalNode `json:"monitor`
	LocalGroupNodes GroupNodes          `json:"local_groupnodes"`
}

type GroupNodes struct {
	LocalNodes []localnode.LocalNode `json:"localnodes"`
}

func NewRedundantNode(monitor localnode.LocalNode, groupnodes []string) *RedundantNode {
	var creation = RedundantNode{
		Monitor: monitor,
	}
	creation.LocalGroupNodes = *NewGroupNodes(monitor.LocalAddress, groupnodes)

	return &creation
}

func NewGroupNodes(address string, ports []string) *GroupNodes {
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

func main() {

	localAddress := "192.168.2.111"
	monitor := localnode.NewLocalNode(localAddress, "11000")
	redundantnode := NewRedundantNode(*monitor, []string{"11001", "11002", "11003"})

	jsonText, err := json.Marshal(redundantnode)
	if err != nil {
		log.Println("Marshal failed!")
	}
	log.Printf("%s", jsonText)

	localNode := localnode.NewLocalNode(localAddress, "11000")
	localNode.SetupLogging()
	localNode.RetrievePublicAddress()
	localNode.CreateAndRetrievePeerList()
	localNode.SendPeerAddRequests()
	localNode.StartSendingKeepAlive()
	localNode.SetupExitHandler()
	localNode.HandleRequests()
}
