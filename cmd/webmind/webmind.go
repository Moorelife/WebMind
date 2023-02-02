package main

import (
	"github.com/Moorelife/WebMind/internal/webmind"
	"github.com/Moorelife/WebMind/internal/webmind/system"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
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
	node.Start()
	log.Printf("Ending program, restarting")
	Phoenix()
}

func Phoenix() {
	cmd := exec.Command("B:\\webmind.exe")
	cmd.Start()
	log.Printf("Phoenix has risen!!!")
	time.Sleep(5 * time.Second)
	os.Exit(1)
}

// Saved Experiments =================================================

func constructAndPrintStructs(address net.TCPAddr) {
	localNode := system.NewNode(address)
	log.Println(localNode.ToJSON())
}
