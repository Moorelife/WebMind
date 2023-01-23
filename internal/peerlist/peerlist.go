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
func RemoteAddToAll(ownAddress string) {
	trace.Entered("PeerList:RemoteAddToAll")
	defer trace.Exited("PeerList:RemoteAddToAll")

	for other := range Peers {
		if other != ownAddress {
			RemoteAdd(ownAddress, other)
		}
	}
}

// RemoteAdd sends a request to get this host added to the remote specified.
func RemoteAdd(addressToAdd, sendTo string) {
	trace.Entered("PeerList:RemoteAdd")
	defer trace.Exited("PeerList:RemoteAdd")

	url := fmt.Sprintf("http://%v/peer/add?%v", sendTo, addressToAdd)
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send peer list add request to %v (ERROR: %v)\n", sendTo, err)
	}
	log.Printf("Peerlist length after Add: %#v", len(LocalGet(addressToAdd)))
}

// RemoteDeleteToAll sends an Add peer request to all members on the peer list.
func RemoteDeleteToAll(exceptAddress string) {
	trace.Entered("PeerList:RemoteDeleteToAll")
	defer trace.Exited("PeerList:RemoteDeleteToAll")

	for other := range Peers {
		if other != exceptAddress {
			RemoteDelete(exceptAddress, other)
		}
	}
}

// RemoteDelete sends a request to get this host added to the remote specified.
func RemoteDelete(selfAddress, sendTo string) {
	trace.Entered("PeerList:RemoteDelete")
	defer trace.Exited("PeerList:RemoteDelete")

	url := fmt.Sprintf("http://%v/peer/delete?%v", sendTo, selfAddress)
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send peer list delete request to %v (ERROR: %v)\n", sendTo, err)
	}
	log.Printf("Peerlist length after Add: %#v", len(LocalGet(selfAddress)))
}

// RemoteGet returns the peer list as a sorted slice.
func RemoteGet(hostPort string) []Peer {
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

	return LocalGet(hostPort)
}

// Local functions deal with operations on the local peer list.

// LocalGet returns the peer list as a sorted slice.
func LocalGet(hostPort string) []Peer {
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
func LocalAdd(hostPort string) {
	trace.Entered("PeerList:LocalAdd")
	defer trace.Exited("PeerList:LocalAdd")
	if Peers == nil {
		Peers = make(map[string]*Peer)
	}

	Peers[hostPort] = NewPeer(hostPort)
	log.Printf("AFTER ADD: %#v", LocalGet(hostPort))
}

// LocalDelete removes a Peer from the PeerList
func LocalDelete(hostPort string) {
	trace.Entered("PeerList:LocalDelete")
	defer trace.Exited("PeerList:LocalDelete")
	delete(Peers, hostPort)
	log.Printf("AFTER DELETE: %#v", LocalGet(hostPort))
}

// CleanPeerList removes entries that have not been seen in the last KeepAlive cycle.
func CleanPeerList() {
	for _, peer := range Peers {
		if peer.lastSeen.Before(time.Now().Add(-20 * time.Second)) {
			LocalDelete(peer.addressPort)
		}
	}
}

func KeepAlive(w http.ResponseWriter, r *http.Request) {
	trace.Entered("WebMind:Internal:KeepAlive")
	defer trace.Exited("WebMind:Internal:KeepAlive")
	defer r.Body.Close()
	sender := strings.Split(r.RequestURI, "?")
	log.Printf("keepalive from %v", sender[1])

	// if received, mark peer as still alive at this time.
	peer := Peers[sender[1]]
	if peer == nil {
		return
	}
	peer.lastSeen = time.Now()
	log.Printf("keepAlive ")

	fmt.Fprintf(w, "I'm still here...")
}

// Webserver handler functions.

// HandlePeerAdd deals with peer addition requests coming in over the network.
func HandlePeerAdd(w http.ResponseWriter, r *http.Request) {
	trace.Entered("PeerList:HandlePeerAdd endpoint")
	defer trace.Exited("PeerList:HandlePeerAdd endpoint")
	parts := strings.Split(r.RequestURI, "?")
	log.Printf("Added host %v to the peerlist\n", parts[1])
	LocalAdd(parts[1])
}

// HandlePeerDelete deals with peer removal requests coming in over the network.
func HandlePeerDelete(w http.ResponseWriter, r *http.Request) {
	trace.Entered("PeerList:HandlePeerDelete endpoint")
	defer trace.Exited("PeerList:HandlePeerDelete endpoint")
	parts := strings.Split(r.RequestURI, "?")
	log.Printf("Removed host %v from the peerlist\n", parts[1])
	LocalDelete(parts[1])
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
