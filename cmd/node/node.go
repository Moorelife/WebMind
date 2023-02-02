package main

import (
	"flag"
	"github.com/Moorelife/WebMind/internal/webmind"
	"github.com/Moorelife/WebMind/internal/webmind/system"
	"log"
	"net"
)

func main() {
	port := flag.Int("port", 14285, "port number for the monitor node.")
	flag.Parse()

	address := net.TCPAddr{
		IP:   []byte{192, 168, 2, 111}, // accept any connection
		Port: *port,
	}
	webmind.SetupLogging(address.String())

	constructAndPrintStructs(address)
	log.Printf("Starting Web interface at %s", address.String())

	node := system.NewNode(address)
	node.OtherPort = 11000
	node.Start()

	log.Printf("Ending program =================================")
}

// Saved Experiments =================================================

func constructAndPrintStructs(address net.TCPAddr) {
	localNode := system.NewNode(address)
	log.Println(localNode.ToJSON())
}
