package system

import (
	"fmt"
	"github.com/Moorelife/WebMind/foundation"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Struct and Constructor ============================================

// Node defines the data required to set up a web system.
// A Web has a number of peers, that keep each other alive.
type Web struct {
	Nodes []*Node `json:"Nodes"` // the addresses of the nodes of the web.

	wg    sync.WaitGroup
	ctime time.Time
}

// NewNode creates a new Node structure and returns a pointer to it
func NewWeb(nodes []*Node) *Web {
	web := Web{Nodes: nodes}
	return &web
}

// Core functionality ================================================

func (w *Web) Start() {
	for _, node := range w.Nodes {
		foundation.StartNode(strconv.Itoa(node.Address.Port))
	}
}

// WebHandler endpoints ==============================================

func (w *Web) HandleRoot(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	foundation.PrintRequest(r)
	log.Printf("Handling /")
	fmt.Fprintf(rw, "Illegal request: %s", r.RequestURI)
}

func (w *Web) HandleSync(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("Handling /sync")
	panic("Implement me")
}
