package main

import (
	"flag"
	"github.com/Moorelife/WebMind/foundation"
	"github.com/Moorelife/WebMind/foundation/system"
	"log"
	"net"
)

func main() {
	port := flag.Int("port", 14285, "port number for the node.")
	flag.Parse()

	address := net.TCPAddr{
		IP:   []byte{192, 168, 2, 111}, // accept any connection
		Port: *port,
	}
	foundation.SetupLogging(address.String())

	log.Printf("Starting Web interface at %s", address.String())

	node := system.NewNode(address)
	node.Start()

	log.Printf("Ending program =================================")
}
