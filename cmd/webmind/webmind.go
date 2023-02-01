package main

import (
	"github.com/Moorelife/WebMind/internal/webmind"
	"github.com/Moorelife/WebMind/internal/webmind/system"
	"log"
	"net"
)

func main() {
	address := net.TCPAddr{
		IP:   []byte{0, 0, 0, 0}, // accept any connection
		Port: 14285,
	}
	webmind.SetupLogging(address.String())

	constructAndPrintStructs(address)
	log.Print("=======================================================================")
	log.Printf("         Starting Web interface at: %s", address.String())
	log.Print("=======================================================================")
	node := system.NewNode(address)

	log.Printf("Entering node.Start()")
	node.Start()
}

// Saved Experiments =================================================

func constructAndPrintStructs(address net.TCPAddr) {
	localNode := system.NewNode(address)
	log.Println(localNode.ToJSON())
}
