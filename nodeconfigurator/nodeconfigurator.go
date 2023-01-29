package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/Moorelife/WebMind/internal/localnode"
	"github.com/Moorelife/WebMind/internal/oneness"
	"log"
	"os"
)

// NodeConfigurator makes it easy to create a configuration for a redundant local node configuration, where three
// instances interact locally to provide redundancy, and where a fourth one will keep them operating smoothly.
// All four nodes run on the same IPv4 address, each with distinct ports. These can be consecutive if only the
// monitor port is specified, making the redundant ports the three next port numbers. See main() for parameters.

func main() {
	address := flag.String("address", "localhost", "IPv4 address of the monitor and the local nodes.")
	monitorPort := flag.Int("port", 11000, "port number for the monitor node.")
	t1Port := flag.Int("t1", 0, "port number for the trinity node 1, 0 means monitorPort + 1.")
	t2Port := flag.Int("t2", 0, "port number for the trinity node 2, 0 means monitorPort + 2.")
	t3Port := flag.Int("t3", 0, "port number for the trinity node 3, 0 means monitorPort + 3.")
	configFilename := flag.String("file", "webmind.json", "filename to write the configuration into.")
	flag.Parse()

	monitor := localnode.NewLocalNode(*address, *monitorPort)
	redundant := oneness.NewOneness(*monitor, []int{*t1Port, *t2Port, *t3Port})

	jsonText, err := json.Marshal(redundant)
	if err != nil {
		log.Printf("Marshal failed: %v", err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, jsonText, "", "\t")
	if err != nil {
		log.Printf("Indent failed: %v", err)
	}

	file, err := os.OpenFile(*configFilename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 644)
	_, err = out.WriteTo(file)
	if err != nil {
		log.Printf("Write to file failed: %v", err)
	}
}
