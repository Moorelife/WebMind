package peerlist

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Moorelife/WebMind/internal/trace"
)

// Peer holds all information about a network peer.
type Peer struct {
	addressPort string
	lastSeen    time.Time
	// public key?
}

// PeerList holds all information for a collection of peers.
type PeerList map[string]*Peer

// Peers contains the list of Peer objects for this instance of the client.
var Peers = PeerList{}

// NewPeer creates a new Peer with the given host and port.
func NewPeer(hostPort string) *Peer {
	result := Peer{hostPort, time.Now()}
	return &result
}

// Remote functions deal with sending requests to other peers regarding the peer list.

// RemoteAddToAll sends an Add peer request to all members on the peer list.
func (p *PeerList) RemoteAddToAll(ownAddress string) {
	trace.Entered("PeerList:RemoteAddToAll")
	defer trace.Exited("PeerList:RemoteAddToAll")

	for other := range Peers {
		if other != ownAddress {
			p.RemoteAdd(ownAddress, other)
		}
	}
}

// RemoteAdd sends a request to get this host added to the remote specified.
func (p *PeerList) RemoteAdd(addressToAdd, sendTo string) {
	trace.Entered("PeerList:RemoteAdd")
	defer trace.Exited("PeerList:RemoteAdd")

	url := fmt.Sprintf("http://%v/peer/add?%v", sendTo, addressToAdd)
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send peer list add request to %v (ERROR: %v)\n", sendTo, err)
	}
	log.Printf("Peerlist length after Add: %#v", len(p.LocalGet(addressToAdd)))
}

// RemoteDeleteToAll sends an Add peer request to all members on the peer list.
func (p *PeerList) RemoteDeleteToAll(exceptAddress string) {
	trace.Entered("PeerList:RemoteDeleteToAll")
	defer trace.Exited("PeerList:RemoteDeleteToAll")

	for other := range Peers {
		if other != exceptAddress {
			p.RemoteDelete(exceptAddress, other)
		}
	}
}

// RemoteDelete sends a request to get this host added to the remote specified.
func (p *PeerList) RemoteDelete(selfAddress, sendTo string) {
	trace.Entered("PeerList:RemoteDelete")
	defer trace.Exited("PeerList:RemoteDelete")

	url := fmt.Sprintf("http://%v/peer/delete?%v", sendTo, selfAddress)
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send peer list delete request to %v (ERROR: %v)\n", sendTo, err)
	}
	log.Printf("Peerlist length after Add: %#v", len(p.LocalGet(selfAddress)))
}

// RemoteGet returns the peer list as a sorted slice.
func (p *PeerList) RemoteGet(hostPort string) []Peer {
	trace.Entered("PeerList:RemoteGet")
	defer trace.Exited("PeerList:RemoteGet")

	url := fmt.Sprintf("http://%v/peer/list", hostPort)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to get peer list from %v (ERROR: %v)\n", hostPort, err)
	}

	defer resp.Body.Close()

	rawlist, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read body for peer list from %v (ERROR: %v)\n", hostPort, err)
	}

	peers := make([]string, 1)
	err = json.Unmarshal(rawlist, &peers)
	if err != nil {
		log.Printf("error converting peerlist from JSON: %v", err)
	}

	for _, peerAddr := range peers {
		Peers[peerAddr] = NewPeer(peerAddr)
	}

	return p.LocalGet(hostPort)
}

// Local functions deal with operations on the local peer list.

// LocalGet returns the peer list as a sorted slice.
func (p *PeerList) LocalGet(hostPort string) []Peer {
	trace.Entered("PeerList:LocalGet")
	defer trace.Exited("PeerList:LocalGet")

	log.Printf("Peerlist: %#v", Peers)

	var result []Peer
	for _, peer := range Peers {
		result = append(result, *peer)
	}

	return result
}

// LocalAdd adds a new Peer to the PeerList
func (p *PeerList) LocalAdd(hostPort string) {
	trace.Entered("PeerList:LocalAdd")
	defer trace.Exited("PeerList:LocalAdd")
	if Peers == nil {
		Peers = make(map[string]*Peer)
	}

	Peers[hostPort] = NewPeer(hostPort)
	log.Printf("AFTER ADD: %#v", p.LocalGet(hostPort))
}

// LocalDelete removes a Peer from the PeerList
func (p *PeerList) LocalDelete(hostPort string) {
	trace.Entered("PeerList:LocalDelete")
	defer trace.Exited("PeerList:LocalDelete")
	delete(Peers, hostPort)
	log.Printf("AFTER DELETE: %#v", p.LocalGet(hostPort))
}

// CleanPeerList removes entries that have not been seen in the last KeepAlive cycle.
func (p *PeerList) CleanPeerList(exceptAddress string) {
	for _, peer := range Peers {
		if peer.addressPort != exceptAddress &&
			peer.lastSeen.Before(time.Now().Add(-11*time.Second)) {
			p.LocalDelete(peer.addressPort)
		}
	}
}

func (p *PeerList) KeepAlive(w http.ResponseWriter, r *http.Request) {
	trace.Entered("WebMind:Internal:KeepAlive")
	//defer trace.Exited("WebMind:Internal:KeepAlive")
	//defer r.Body.Close()
	sender := strings.Split(r.RequestURI, "?")

	// if received, mark peer as still alive at this time.
	if len(sender) <= 1 {
		return
	}
	log.Printf("keepalive from %v", sender[1])
	peer := Peers[sender[1]]
	if peer == nil {
		return
	}
	peer.lastSeen = time.Now()

	fmt.Fprintf(w, "I'm still here...")
}

// Webserver handler functions.

// HandlePeerAdd deals with peer addition requests coming in over the network.
func HandlePeerAdd(w http.ResponseWriter, r *http.Request) {
	trace.Entered("PeerList:HandlePeerAdd endpoint")
	defer trace.Exited("PeerList:HandlePeerAdd endpoint")
	parts := strings.Split(r.RequestURI, "?")
	if len(parts) > 1 {
		log.Printf("Adding host %v to the peerlist\n", parts[1])
		Peers.LocalAdd(parts[1])
	}
}

// HandlePeerDelete deals with peer removal requests coming in over the network.
func HandlePeerDelete(w http.ResponseWriter, r *http.Request) {
	trace.Entered("PeerList:HandlePeerDelete endpoint")
	defer trace.Exited("PeerList:HandlePeerDelete endpoint")
	parts := strings.Split(r.RequestURI, "?")
	if len(parts) > 1 {
		log.Printf("Removing host %v from the peerlist\n", parts[1])
		Peers.LocalDelete(parts[1])
	}
}

// HandlePeerList deals with peer list requests coming in over the network.
func HandlePeerList(w http.ResponseWriter, r *http.Request) {
	trace.Entered("PeerList:HandlePeerList endpoint")
	defer trace.Exited("PeerList:HandlePeerList endpoint")

	var result []string
	for _, peer := range Peers {
		result = append(result, string(peer.addressPort))
	}

	j, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error converting peerlist to JSON: %v", err)
	}

	w.Write(j)
}
