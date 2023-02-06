package system

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Moorelife/WebMind/foundation"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Struct and Constructor ============================================

// Node defines the data required to define a node system.
type Node struct {
	Address net.TCPAddr `json:"Address"` // the address of the node.

	server http.Server
	wg     sync.WaitGroup
	ctime  time.Time

	// temporary stuff
	OtherPort int
}

// NewNode creates a new Node structure and returns a pointer to it
func NewNode(address net.TCPAddr) *Node {
	node := Node{Address: address, wg: sync.WaitGroup{}, ctime: time.Now()}
	return &node
}

// Core functionality ================================================

func (n *Node) Start() {
	http.HandleFunc("/", n.HandleRoot)
	http.HandleFunc("/shutdown", n.HandleShutdown)
	http.HandleFunc("/startup", n.HandleStartup)
	http.HandleFunc("/status", n.HandleStatus)

	n.wg.Add(1)
	n.server = http.Server{Addr: n.Address.String(), Handler: nil}

	go func() {
		if err := n.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error in ListenAndServe(): %v", err)
		}
	}()
	n.wg.Wait()
}

// WebHandler endpoints ==============================================

func (n *Node) HandleRoot(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	foundation.PrintRequest(r)
	log.Printf("Handling /")
	fmt.Fprintf(w, "Illegal request: %s", r.RequestURI)
}

func (n *Node) HandleShutdown(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Fprintf(w, "Shutting down webserver!")
	ctx, err := context.WithTimeout(context.Background(), 1*time.Second)
	if err != nil {

	}
	n.server.Shutdown(ctx)
	log.Printf("Handling /shutdown")
	n.wg.Done()
}

func (n *Node) HandleStartup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("Handling /startup")
	n.OtherPort = n.getPortFromRequest(r)
	fmt.Fprintf(w, "Starting up new node!")
	foundation.StartNode(strconv.Itoa(n.OtherPort))
}

func (n *Node) getPortFromRequest(r *http.Request) int {
	parts := strings.Split(r.RequestURI, "?")
	if len(parts) < 2 {
		return n.OtherPort
	}
	port, err := strconv.Atoi(parts[1])
	if err == nil {
		port = clamp(port, 0, 65535)
	}
	return port
}

func clamp(n, min, max int) int {
	if n < min {
		return min
	}
	
	if n > max {
		return max
	}

	return n
}

func (n *Node) HandleStatus(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("Handling /status")
	fmt.Fprintf(w, "Node up and running! (%v)", n.ctime.Round(0))
}

// Utility functions =================================================

// ToJSON converts the Node struct to indented JSON.
func (n *Node) ToJSON() string {
	b, err := json.Marshal(n)
	if err != nil {
		panic(fmt.Sprintf("Marshal failed: %v", err))
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Indent failed: %v", err))
	}

	return fmt.Sprintf("%s", &out)
}
