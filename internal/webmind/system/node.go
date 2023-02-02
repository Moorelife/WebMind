package system

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Moorelife/WebMind/internal/webmind"
	"log"
	"net"
	"net/http"
	"strconv"
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
	log.Printf("Creating node.WaitGroup")
	node := Node{Address: address, wg: sync.WaitGroup{}, ctime: time.Now()}
	return &node
}

// Core functionality ================================================

func (n *Node) Start() {
	http.HandleFunc("/", n.HandleRoot)
	http.HandleFunc("/kill", n.HandleKill)
	http.HandleFunc("/spawn", n.HandleSpawn)
	http.HandleFunc("/heartbeat", n.HandleHeartbeat)

	log.Printf("Creating node.server")
	n.wg.Add(1)
	n.server = http.Server{Addr: n.Address.String(), Handler: nil}
	log.Printf("Created node.server")

	go func() {
		log.Printf("Entering server.ListenAndServe()")
		if err := n.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error in ListenAndServe(): %v", err)
		}
	}()
	n.wg.Wait()
}

// WebHandler endpoints ==============================================

func (n *Node) HandleRoot(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	webmind.PrintRequest(r)
	log.Printf("Handling /")
	fmt.Fprintf(w, "Node up and running! (%v)", time.Now().Sub(n.ctime))
}

func (n *Node) HandleKill(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Fprintf(w, "Killing webserver!")
	ctx, err := context.WithTimeout(context.Background(), 1*time.Second)
	if err != nil {

	}
	n.server.Shutdown(ctx)
	log.Printf("Handling /kill")
	n.wg.Done()
}

func (n *Node) HandleSpawn(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("Handling /spawn")
	fmt.Fprintf(w, "Spawning new node!")
	webmind.Phoenix(strconv.Itoa(n.OtherPort))
}

func (n *Node) HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("Handling /heartbeat")
	fmt.Fprintf(w, "Heartbeat answered!")
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
