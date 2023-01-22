package peerlist

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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

var Peers = PeerList{}

// NewPeer creates a new Peer with the given host and port.
func NewPeer(hostPort string) *Peer {
	result := Peer{hostPort}
	return &result
}

// Remote functions deal with sending requests to other peers regarding the peer list.

// RemoteAddToAll sends an Add peer request to all members on the peer list.
func RemoteAddToAll(ownAddress string) {
	trace.Entered("PeerList:RemoteAddToAll")
	defer trace.Exited("PeerList:RemoteAddToAll")

	for other := range Peers.peers {
		if other != ownAddress {
			RemoteAdd(ownAddress, other)
		}
	}
}

// RemoteAdd sends a request to get this host added to the remote specified.
func RemoteAdd(toAddress, sendTo string) {
	trace.Entered("PeerList:RemoteAdd")
	defer trace.Exited("PeerList:RemoteAdd")

	url := fmt.Sprintf("http://%v/peer/add?%v", sendTo, toAddress)
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send peer list add request to %v (ERROR: %v)\n", sendTo, err)
	}
}

// RemoteDeleteToAll sends an Add peer request to all members on the peer list.
func RemoteDeleteToAll(exceptAddress string) {
	trace.Entered("PeerList:RemoteDeleteToAll")
	defer trace.Exited("PeerList:RemoteDeleteToAll")

	for other := range Peers.peers {
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
}

// RemoteGet returns the peer list as a sorted slice.
func RemoteGet(hostPort string) []Peer {
	trace.Entered("PeerList:RemoteGet")
	defer trace.Exited("PeerList:RemoteGet")

	url := fmt.Sprintf("http://%v/peer/list", hostPort)

	log.Printf(">>>> Sending request: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to get peer list from %v (ERROR: %v)\n", hostPort, err)
	}

	log.Println(">>>> Defer body.Close()")

	defer resp.Body.Close()

	log.Println(">>>> reading body")

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
		Peers.peers[peerAddr] = NewPeer(peerAddr)
	}

	fmt.Printf(">>>> %#s\n", string(rawlist))

	return LocalGet(hostPort)
}

// Local functions deal with operations on the local peer list.

// LocalGet returns the peer list as a sorted slice.
func LocalGet(hostPort string) []Peer {
	trace.Entered("PeerList:LocalGet")
	defer trace.Exited("PeerList:LocalGet")

	var result []Peer
	result = append(result, *NewPeer(hostPort))
	return result
}

// LocalAdd adds a new Peer to the PeerList
func LocalAdd(hostPort string) {
	trace.Entered("PeerList:LocalAdd")
	defer trace.Exited("PeerList:LocalAdd")
	if Peers.peers == nil {
		Peers.peers = make(map[string]*Peer)
	}
	Peers.peers[hostPort] = NewPeer(hostPort)
}

// LocalDelete removes a Peer from the PeerList
func LocalDelete(hostPort string) {
	trace.Entered("PeerList:LocalDelete")
	defer trace.Exited("PeerList:LocalDelete")
	delete(Peers.peers, hostPort)
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
	for _, peer := range Peers.peers {
		result = append(result, string(peer.addressPort))
	}

	j, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error converting peerlist to JSON: %v", err)
	}

	w.Write(j)
}
