package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Moorelife/WebMind/internal/redundantnode"
	"log"
	"os"
)

func main() {
	configName := flag.String("config", "redundantnode.json", "filename to load the configuration into.")
	flag.Parse()

	data := readConfigFile(*configName)
	redundantNode := createRedundantNodeFromJSON(data)
	logConfiguration(data)

	log.Printf("%v", redundantNode)

	redundantNode.CreateNodes()

}

func readConfigFile(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Reading configuration failed: %v", err))
	}
	return data
}

func createRedundantNodeFromJSON(data []byte) *redundantnode.RedundantNode {
	redundantNode := redundantnode.RedundantNode{}
	err := json.Unmarshal(data, &redundantNode)
	if err != nil {
		panic(fmt.Sprintf("Unmarshal failed: %v", err))
	}
	return &redundantNode
}

func logConfiguration(data []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Indent failed: %v", err))
	}
	log.Printf("%s", out)
}
