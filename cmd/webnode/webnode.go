package main

import (
	"flag"
	"github.com/Moorelife/WebMind/internal/webmind1/localnode"
	"log"
	"time"
)

// WebNode in its current state is JUST A LEARNING EXPERIMENT,
// and as such can not be expected to be fit for any given purpose.
// Please understand that you use the program at your own risk!!!

func AddCommandFlagsToNode(l *localnode.LocalNode) {
	localPort := flag.Int("port", 7777, "http server port number")
	localAddress := flag.String("address", "localhost", "address of the node")
	flag.Parse()

	l.LocalPort = *localPort
	l.LocalAddress = *localAddress
}

func main() {
	localNode := localnode.NewLocalNode("localhost", 7777)

	AddCommandFlagsToNode(localNode)

	localNode.SetupLogging()
	// localNode.RetrievePublicAddress()
	// localNode.CreateAndRetrievePeerList()
	// localNode.SendPeerAddRequests()
	// localNode.StartSendingKeepAlive()
	// localNode.SetupExitHandler()

	for i := 10; i > 0; i-- {
		log.Printf("Countdown: %v\n", i)
		time.Sleep(1 * time.Second)
	}
	log.Printf("web interface running at %s:%d", localNode.LocalAddress, localNode.LocalPort)
	localNode.HandleRequests()
}
