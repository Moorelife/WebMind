package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
)

// Hive defines the data required to define a Hive system.
type Hive struct {
	Address net.TCPAddr `json:"address"` // the address of the Hive.
}

// NewHive creates a new Hive structure and returns a pointer to it
func NewHive(address net.TCPAddr) *Hive {
	Hive := Hive{Address: address}
	return &Hive
}

// ToJSON converts the Hive struct to indented JSON.
func (n *Hive) ToJSON() string {
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
