package main

import (
	"flag"
	"github.com/Moorelife/WebMind/foundation/system"
	"net"
)

func main() {
	port1 := flag.Int("p1", 14285, "port number for the 1st node.")
	port2 := flag.Int("p2", 42851, "port number for the 2nd node.")
	port3 := flag.Int("p3", 28514, "port number for the 3rd node.")
	flag.Parse()

	addresses := []*system.Node{
		system.NewNode(net.TCPAddr{IP: []byte{192, 168, 2, 111}, Port: *port1}),
		system.NewNode(net.TCPAddr{IP: []byte{192, 168, 2, 111}, Port: *port2}),
		system.NewNode(net.TCPAddr{IP: []byte{192, 168, 2, 111}, Port: *port3}),
	}

	web := system.NewWeb(addresses)
	web.Start()
}
