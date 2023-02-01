package main

import (
	"log"
	"net"

	"github.com/Moorelife/WebMind/internal/webmind/system"
)

func main() {
	address := net.TCPAddr{
		IP:   []byte{192, 168, 2, 111},
		Port: 14285,
	}
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
