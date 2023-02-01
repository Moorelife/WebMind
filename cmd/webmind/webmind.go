package main

import (
	"github.com/Moorelife/WebMind/internal/webmind"
	"log"
	"net"

	"github.com/Moorelife/WebMind/internal/webmind/system"
)

func main() {
	address := net.TCPAddr{
		IP:   []byte{0, 0, 0, 0}, // accept any connection
		Port: 14285,
	}
	webmind.SetupLogging(address.String())
	constructAndPrintStructs(address)
	log.Printf("Starting Web interface at: %s", address.String())
	node := system.NewNode(address)
	node.Start()
}

// Saved Experiments =================================================

func constructAndPrintStructs(address net.TCPAddr) {
	localNode := system.NewNode(address)
	log.Println(localNode.ToJSON())
}
