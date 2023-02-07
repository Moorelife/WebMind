package main

import (
	"flag"
	"github.com/Moorelife/WebMind/foundation"
	"github.com/Moorelife/WebMind/foundation/system"
	"log"
	"net"
)

func main() {
	sourcePort := flag.Int("source", 14285, "port number for the source node.")
	nodePort := flag.Int("port", 11111, "port number for the new node.")
	flag.Parse()

	log.Printf("Source Port: %v", *sourcePort)
	log.Printf("Node Port:   %v", *nodePort)

	addr := []byte{192, 168, 2, 111}
	sourceAddress := net.TCPAddr{
		IP:   addr,
		Port: *sourcePort,
	}
	nodeAddress := net.TCPAddr{
		IP:   addr,
		Port: *nodePort,
	}
	foundation.SetupLogging(nodeAddress.String())

	log.Printf("Starting Web interface at %s", nodeAddress.String())

	node := system.NewNode(sourceAddress, nodeAddress)
	node.Start()

	log.Printf("Ending program =================================")
}
