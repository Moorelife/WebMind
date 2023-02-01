package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Moorelife/WebMind/internal/webmind"
	"net"
	"net/http"
)

// Struct and Constructor ============================================

// Node defines the data required to define a node system.
type Node struct {
	Address net.TCPAddr `json:"Address"` // the address of the node.
}

// NewNode creates a new Node structure and returns a pointer to it
func NewNode(address net.TCPAddr) *Node {
	node := Node{Address: address}
	return &node
}

// Core functionality ================================================

func (n *Node) Start() {
	http.HandleFunc("/", HandleServerRootRequests)

	err := http.ListenAndServe(n.Address.String(), nil)
	if err != nil {
		panic(fmt.Sprintf("ListenAndServe ended: %v", err))
	}
}

// WebHandler endpoints ==============================================

func HandleServerRootRequests(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	webmind.PrintRequest(r)

	fmt.Fprintf(w, "Node up and running!")
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
