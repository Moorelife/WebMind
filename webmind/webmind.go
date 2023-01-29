package main

import (
	"github.com/Moorelife/WebMind/internal/localnode"
)

// WebMind in its current state is JUST A LEARNING EXPERIMENT,
// and as such can not be expected to be fit for any given purpose.
// Please understand that you use the program at your own risk!!!

func main() {
	localNode := localnode.NewLocalNode()
	localNode.SetupLogging()
	localNode.RetrievePublicAddress()
	localNode.CreateAndRetrievePeerList()
	localNode.SendPeerAddRequests()
	localNode.StartSendingKeepAlive()
	localNode.SetupExitHandler()
	localNode.HandleRequests()
}
