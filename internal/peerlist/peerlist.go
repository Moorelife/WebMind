package peerlist

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
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
type PeerList struct {
	Users map[string]*Peer
	rw    sync.RWMutex
}

// Peers contains the list of Peer objects for this instance of the client.
var Peers = PeerList{}

// NewPeer creates a new Peer with the given host and port.
func NewPeer(hostPort string) *Peer {
	result := Peer{hostPort, time.Now()}
	return &result
}

func init() {
	Peers.Users = make(map[string]*Peer, 1)
}

// Remote functions deal with sending requests to other peers regarding the peer list.

// RemoteAddToAll sends an Add peer request to all members on the peer list.
func (p *PeerList) RemoteAddToAll(ownAddress string) {
	trace.Entered("PeerList:RemoteAddToAll")
	defer trace.Exited("PeerList:RemoteAddToAll")

	for other := range Peers.Users {
		if other != ownAddress {
			p.RemoteAdd(ownAddress, other)
		}
	}
	p.logLocalList()
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
	p.logLocalList()
}

// RemoteDeleteToAll sends an Add peer request to all members on the peer list.
func (p *PeerList) RemoteDeleteToAll(exceptAddress string) {
	trace.Entered("PeerList:RemoteDeleteToAll")
	defer trace.Exited("PeerList:RemoteDeleteToAll")

	for other := range Peers.Users {
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
	p.logLocalList()
}

// RemoteGet returns the peer list as a sorted slice.
func (p *PeerList) RemoteGet(hostPort string) []Peer {
	trace.Entered("PeerList:RemoteGet")
	defer trace.Exited("PeerList:RemoteGet")

	url := fmt.Sprintf("http://%v/peer/list", hostPort)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to get peer list from %v, retrying once (ERROR: %v)\n", hostPort, err)
		resp, err = http.Get(url)
	}

	if resp == nil || resp.Body == nil {
		return []Peer{}
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
		Peers.Users[peerAddr] = NewPeer(peerAddr)
	}
	p.logLocalList()

	return p.LocalGet(hostPort)
}

// Local functions deal with operations on the local peer list.

// LocalGet returns the peer list as a sorted slice.
func (p *PeerList) LocalGet(hostPort string) []Peer {
	trace.Entered("PeerList:LocalGet")
	defer trace.Exited("PeerList:LocalGet")

	log.Printf("Peerlist: %#v", Peers)

	var result []Peer
	for _, peer := range Peers.Users {
		result = append(result, *peer)
	}

	return result
}

// LocalAdd adds a new Peer to the PeerList
func (p *PeerList) LocalAdd(hostPort string) {
	trace.Entered("PeerList:LocalAdd")
	defer trace.Exited("PeerList:LocalAdd")

	Peers.rw.Lock()
	defer Peers.rw.Unlock()
	Peers.Users[hostPort] = NewPeer(hostPort)
	p.logLocalList()
}

// LocalDelete removes a Peer from the PeerList
func (p *PeerList) LocalDelete(hostPort string) {
	trace.Entered("PeerList:LocalDelete")
	defer trace.Exited("PeerList:LocalDelete")
	Peers.rw.Lock()
	defer Peers.rw.Unlock()
	delete(Peers.Users, hostPort)
	p.logLocalList()
}

// CleanPeerList removes entries that have not been seen in the last HandleKeepAlive cycle.
func (p *PeerList) CleanPeerList(exceptAddress string) {
	for _, peer := range Peers.Users {
		if peer.addressPort != exceptAddress &&
			peer.lastSeen.Before(time.Now().Add(-11*time.Second)) {
			p.LocalDelete(peer.addressPort)
		}
	}
}

var countOnly = true // if true, logs only the count, not the entries

func (p *PeerList) logLocalList() {
	if countOnly {
		log.Printf("PEERLIST COUNT: %v", len(Peers.Users))
		return
	}
	log.Print("----------------------------------------------------------------------------------")
	log.Print("Address 			LastSeen")
	for _, peer := range Peers.Users {
		log.Printf("%v  	%v", peer.addressPort, peer.lastSeen)
	}
	log.Print("----------------------------------------------------------------------------------")
}

// Webserver handler functions.

// HandlePeerAdd deals with peer addition requests coming in over the network.
func HandlePeerAdd(w http.ResponseWriter, r *http.Request) {
	trace.Entered("PeerList:HandlePeerAdd endpoint")
	defer trace.Exited("PeerList:HandlePeerAdd endpoint")
	parts := strings.Split(r.RequestURI, "?")
	if len(parts) > 1 {
		// log.Printf("Adding host %v to the peerlist\n", parts[1])
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
	Peers.rw.Lock()
	defer Peers.rw.Unlock()
	for _, peer := range Peers.Users {
		result = append(result, string(peer.addressPort))
	}

	j, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error converting peerlist to JSON: %v", err)
	}

	w.Write(j)
}

func (p *PeerList) HandleKeepAlive(w http.ResponseWriter, r *http.Request) {
	trace.Entered("WebMind:Internal:HandleKeepAlive")
	defer trace.Exited("WebMind:Internal:HandleKeepAlive")
	defer r.Body.Close()
	sender := strings.Split(r.RequestURI, "?")

	// if received, mark peer as still alive at this time.
	if len(sender) <= 1 {
		return
	}
	peer := Peers.Users[sender[1]]
	if peer == nil {
		return
	}
	peer.lastSeen = time.Now()
	log.Printf("HandleKeepAlive from %v", sender[1])
	p.RemoteAddToAll(sender[1])
	fmt.Fprintf(w, "I'm still here...")
}
