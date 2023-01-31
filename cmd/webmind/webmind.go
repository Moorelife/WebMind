package main

import (
	"net"

	"github.com/Moorelife/WebMind/internal/webmind/system/node"
)

func main() {
	address := net.IP{192, 168, 2, 111}
	localNode := node.NewNode(address)
	// localAddress := "192.168.2.111"
	// monitor := localnode.NewLocalNode(localAddress, "11000")
	// redundantnode := NewRedundantNode(*monitor, []string{"11001", "11002", "11003"})
	//
	// jsonText, err := json.Marshal(redundantnode)
	// if err != nil {
	//	log.Println("Marshal failed!")
	// }
	// log.Printf("%s", jsonText)
	//
	// localNode := localnode.NewLocalNode(localAddress, "11000")
	// localNode.SetupLogging()
	// localNode.RetrievePublicAddress()
	// localNode.CreateAndRetrievePeerList()
	// localNode.SendPeerAddRequests()
	// localNode.StartSendingKeepAlive()
	// localNode.SetupExitHandler()
	// localNode.HandleRequests()
}
