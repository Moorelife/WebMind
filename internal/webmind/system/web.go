package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
)

// Web defines the data required to define a Web system.
type Web struct {
	Address net.TCPAddr `json:"address"` // the address of the Web.
}

// NewWeb creates a new Web structure and returns a pointer to it
func NewWeb(address net.TCPAddr) *Web {
	Web := Web{Address: address}
	return &Web
}

// ToJSON converts the Web struct to indented JSON.
func (n *Web) ToJSON() string {
	b, err := json.Marshal(n)
	if err != nil {
		panic(fmt.Sprintf("Marshal failed: %v", err))
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, " ", "  ")
	if err != nil {
		panic(fmt.Sprintf("Indent failed: %v", err))
	}

	return fmt.Sprintf("%s", &out)
}
