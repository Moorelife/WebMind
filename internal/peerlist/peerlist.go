package peerlist

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Moorelife/WebMind/internal/trace"
)

// Peer holds all information about a network peer.
type Peer struct {
	addressPort string
}

// PeerList holds all information for a collection of peers.
type PeerList struct {
	peers map[string]*Peer
}

var peerList = PeerList{}

// NewPeer creates a new Peer with the given host and port.
func NewPeer(hostPort string) *Peer {
	result := Peer{hostPort}
	return &result
}

// Get returns the peer list as a sorted slice.
func Get(hostPort string) []Peer {
	trace.Entered("PeerList::GetPeerList endpoint")
	defer trace.Exited("PeerList::GetPeerList endpoint")

	url := fmt.Sprintf("http://%v/peer/list", hostPort)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to get peer list from %v (ERROR: %v)\n", hostPort, err)
	}

	defer resp.Body.Close()

	rawlist, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read body for peer list from %v (ERROR: %v)\n", hostPort, err)
	}

	fmt.Printf("%#s\n", string(rawlist))

	var result []Peer
	result = append(result, *NewPeer(hostPort))
	return result
}

// Add adds a new Peer to the PeerList
func Add(hostPort string) {
	trace.Entered("PeerList:Add")
	defer trace.Exited("PeerList:Add")
	if peerList.peers == nil {
		peerList.peers = make(map[string]*Peer)
	}
	peerList.peers[hostPort] = NewPeer(hostPort)
}

// Remove removes a Peer from the PeerList
func Remove(hostPort string) {
	trace.Entered("PeerList:Remove")
	defer trace.Exited("PeerList:Remove")
	delete(peerList.peers, hostPort)
}

func HandlePeerList(w http.ResponseWriter, r *http.Request) {
	trace.Entered("PeerList::HandlePeerList endpoint")
	defer trace.Exited("PeerList::HandlePeerList endpoint")

	var result []string
	for _, peer := range peerList.peers {
		result = append(result, string(peer.addressPort))
	}

	j, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error converting peerlist to JSON: %v", err)
	}

	w.Write(j)
}

func HandlePeerAdd(w http.ResponseWriter, r *http.Request) {
	trace.Entered("PeerList::HandlePeerAdd endpoint")
	defer trace.Exited("PeerList::HandlePeerAdd endpoint")

	Add(r.Host)
}
