package main

import (
	"flag"
	"github.com/Moorelife/WebMind/internal/ip"
	"github.com/Moorelife/WebMind/internal/localnode"
	"github.com/Moorelife/WebMind/internal/peerlist"
)

// WebNode in its current state is JUST A LEARNING EXPERIMENT,
// and as such can not be expected to be fit for any given purpose.
// Please understand that you use the program at your own risk!!!

func AddCommandFlagsToNode(l *localnode.LocalNode) {
	l.OriginsFile = flag.String("origins", "./webmind.json", "origins list file")
	l.LocalPort = flag.String("port", "7777", "http server port number")
	l.FullList = flag.Bool("full", true, "switch to list peers, not just list count")
	l.Trace = flag.Bool("trace", true, "switch to activate call tracing")

	l.LocalAddress = ip.GetPublicIP() + ":" + *l.LocalPort
	l.Peers = peerlist.NewPeerList()
	flag.Parse()
}

func main() {
	localNode := localnode.NewLocalNode()

	AddCommandFlagsToNode(localNode)

	localNode.SetupLogging()
	localNode.RetrievePublicAddress()
	localNode.CreateAndRetrievePeerList()
	localNode.SendPeerAddRequests()
	localNode.StartSendingKeepAlive()
	localNode.SetupExitHandler()
	localNode.HandleRequests()
}
