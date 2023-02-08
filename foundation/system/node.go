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

// Node defines the data required to set up a node system.
type Node struct {
	Address net.TCPAddr `json:"Address"` // the address of the node.
	Source  net.TCPAddr `json:"Source"`  // the address of the node that started this.

	server http.Server
	wg     sync.WaitGroup
	ctime  time.Time
}

// NewNode creates a new Node structure and returns a pointer to it
func NewNode(source net.TCPAddr, address net.TCPAddr) *Node {
	node := Node{
		Address: address,
		Source:  source,
		wg:      sync.WaitGroup{},
		ctime:   time.Now(),
	}
	return &node
}

// Core functionality ================================================

func (n *Node) Start() {
	http.HandleFunc("/", n.HandleRoot)
	http.HandleFunc("/shutdown", n.HandleShutdown)
	http.HandleFunc("/startup", n.HandleStartup)
	http.HandleFunc("/status", n.HandleStatus)
	http.HandleFunc("/sync", n.HandleSync)

	n.startSyncCheck()
	n.startWebServer()
	n.wg.Wait()
}

func (n *Node) startSyncCheck() {
	go func() {
		for true {
			time.Sleep(5 * time.Second)
			url := fmt.Sprintf("http://%s/sync?%d", n.Source.String(), n.Address.Port)
			response, err := http.Get(url)
			if err != nil {
				log.Printf("sync failed, restarting peer...")
				foundation.StartNode(strconv.Itoa(n.Source.Port), strconv.Itoa(n.Address.Port))
			}
			if response != nil {
				defer response.Body.Close()
				if response.StatusCode != http.StatusOK {
					log.Printf("illegal response received")
				}
			}
		}
	}()
}

func (n *Node) startWebServer() {
	n.wg.Add(1)
	n.server = http.Server{Addr: n.Address.String(), Handler: nil}

	go func() {
		if err := n.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error in ListenAndServe(): %v", err)
			n.wg.Done()
		}
	}()
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
	fmt.Fprintf(w, "Shutting down node %s", n.Address.String())
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	n.server.Shutdown(ctx)
	log.Printf("Handling /shutdown")
	n.wg.Done()
}

func (n *Node) HandleStartup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("Handling /startup")
	n.Source.Port = n.getPortFromRequest(r)
	fmt.Fprintf(w, "Starting up new node!")
	foundation.StartNode(strconv.Itoa(n.Source.Port), strconv.Itoa(n.Address.Port))
}

func (n *Node) HandleSync(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("Handling /sync")
	n.Source.Port = n.getPortFromRequest(r)
	fmt.Fprintf(w, "%s", n.Address)
}

func (n *Node) getPortFromRequest(r *http.Request) int {
	parts := strings.Split(r.RequestURI, "?")
	if len(parts) < 2 {
		return n.Source.Port
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
