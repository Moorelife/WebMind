package main

import (
	"encoding/json"
	"github.com/Moorelife/WebMind/internal/localnode"
	"log"
)

// WebMind in its current state is JUST A LEARNING EXPERIMENT,
// and as such can not be expected to be fit for any given purpose.
// Please understand that you use the program at your own risk!!!

type Oneness struct {
	Monitor      localnode.LocalNode `json:"monitor`
	LocalTrinity Trinity             `json:"local_trinity"`
}

type Trinity struct {
	LocalNodes []localnode.LocalNode `json:"localnodes"`
}

func NewOneness(monitor localnode.LocalNode, trinity []string) *Oneness {
	var creation = Oneness{
		Monitor: monitor,
	}
	creation.LocalTrinity = *NewTrinity(monitor.LocalAddress, trinity)

	return &creation
}

func NewTrinity(address string, ports []string) *Trinity {
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

func main() {

	localAddress := "192.168.2.111"
	monitor := localnode.NewLocalNode(localAddress, "11000")
	oneness := NewOneness(*monitor, []string{"11001", "11002", "11003"})

	jsonText, err := json.Marshal(oneness)
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
