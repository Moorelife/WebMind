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
	localNode := system.NewNode(address)
	log.Println(localNode.ToJSON())

	localWeb := system.NewWeb(address)
	log.Println(localWeb.ToJSON())

	localHive := system.NewHive(address)
	log.Println(localHive.ToJSON())
}
