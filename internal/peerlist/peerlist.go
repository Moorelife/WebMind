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

var KeepAliveInterval = 60 * time.Second

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

func init() {
	Peers.Users = make(map[string]*Peer, 1)
}

// NewPeer creates a new Peer with the given host and port.
func NewPeer(hostPort string) *Peer {
	result := Peer{hostPort, time.Now()}
	return &result
}

// Remote functions deal with sending requests to other peers regarding the peer list.

// RemoteAddToAll sends an Add peer request to all members on the peer list.
func (p *PeerList) RemoteAddToAll(localAddress string) {
	trace.Entered("PeerList:RemoteAddToAll")
	defer trace.Exited("PeerList:RemoteAddToAll")

	for other := range Peers.Users {
		if other != localAddress {
			p.RemoteAdd(localAddress, other)
		}
	}
}

// RemoteAdd sends a request to get this host added to the remote specified.
func (p *PeerList) RemoteAdd(localAddress, sendTo string) {
	trace.Entered("PeerList:RemoteAdd")
	defer trace.Exited("PeerList:RemoteAdd")

	url := fmt.Sprintf("http://%v/peer/add?%v", sendTo, localAddress)
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to add to peer %v\n", sendTo)
		p.LocalDelete(sendTo)
	}
}

// RemoteDeleteToAll sends an Add peer request to all members on the peer list.
func (p *PeerList) RemoteDeleteToAll(localAddress string) {
	trace.Entered("PeerList:RemoteDeleteToAll")
	defer trace.Exited("PeerList:RemoteDeleteToAll")

	for other := range Peers.Users {
		if other != localAddress {
			p.RemoteDelete(localAddress, other)
		}
	}
}

// RemoteDelete sends a request to get this host removed from the remote specified.
func (p *PeerList) RemoteDelete(localAddress, sendTo string) {
	trace.Entered("PeerList:RemoteDelete")
	defer trace.Exited("PeerList:RemoteDelete")

	url := fmt.Sprintf("http://%v/peer/delete?%v", sendTo, localAddress)
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send peer list delete request to %v (ERROR: %v)\n", sendTo, err)
	}
}

// RemoteGet returns the peer list as a sorted slice.
func (p *PeerList) RemoteGet(localAddress string) []Peer {
	url := fmt.Sprintf("http://%v/peer/list", localAddress)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to get peer list from %v, retrying once (ERROR: %v)\n", localAddress, err)
		resp, err = http.Get(url)
		if err != nil {
			log.Printf("failed to get peer list from %v, skipping: %v)\n", localAddress, err)
		}
	}

	if resp == nil || resp.Body == nil {
		return []Peer{}
	}
	defer resp.Body.Close()

	rawlist, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read body for peer list from %v (ERROR: %v)\n", localAddress, err)
	}

	peers := make([]string, 1)
	err = json.Unmarshal(rawlist, &peers)
	if err != nil {
		log.Printf("error converting peerlist from JSON: %v", err)
	}

	for _, peerAddr := range peers {
		p.LocalAdd(peerAddr)
	}

	return p.LocalGet()
}

// Local functions deal with operations on the local peer list.

// LocalGet returns the peer list as a sorted slice.
func (p *PeerList) LocalGet() []Peer {
	var result []Peer
	for _, peer := range Peers.Users {
		result = append(result, *peer)
	}

	return result
}

// LocalAdd adds a new Peer to the PeerList
func (p *PeerList) LocalAdd(addressToAdd string) {
	Peers.rw.Lock()
	defer Peers.rw.Unlock()
	if Peers.Users[addressToAdd] == nil {
		log.Printf("Adding host %v to the peerlist\n", addressToAdd)
		Peers.Users[addressToAdd] = NewPeer(addressToAdd)
	}
}

// LocalDelete removes a Peer from the PeerList
func (p *PeerList) LocalDelete(addressToDelete string) {
	if len(Peers.Users) > 2 {
		p.rw.Lock()
		defer p.rw.Unlock()
		delete(p.Users, addressToDelete)
		log.Printf("Removing host %v from the peerlist\n", addressToDelete)
	} else {
		p.RemoteGet(addressToDelete)
	}
}

// CleanPeerList removes entries that have not been seen in the last HandleKeepAlive cycle.
// It refuses to delete the local address, or the last remaining remote address.
func (p *PeerList) CleanPeerList(localAddress string) {
	for _, peer := range Peers.Users {
		if peer.addressPort != localAddress &&
			peer.lastSeen.Before(time.Now().Add(-KeepAliveInterval)) {
			p.LocalDelete(peer.addressPort)
		}
	}
}

var CountOnly = false // if true, logs only the count, not the entries

func (p *PeerList) LogLocalList(fullList bool) {
	if !fullList {
		log.Printf("PEERLIST COUNT: %v", len(Peers.Users))
		return
	}
	log.Print("----------------------------------------------------------------------------------")
	log.Printf("Address 			LastSeen         COUNT: %v", len(Peers.Users))
	for _, peer := range Peers.Users {
		log.Printf("%v  	%v", peer.addressPort, peer.lastSeen.Round(0))
	}
	log.Print("----------------------------------------------------------------------------------")
}

// Webserver handler functions.

// HandlePeerAdd deals with peer addition requests coming in over the network.
func HandlePeerAdd(w http.ResponseWriter, r *http.Request) {
	log.Printf("PeerList:HandlePeerAdd endpoint")

	parts := strings.Split(r.RequestURI, "?")
	if len(parts) > 1 {
		Peers.LocalAdd(parts[1])
	}
}

// HandlePeerDelete deals with peer removal requests coming in over the network.
func HandlePeerDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("PeerList:HandlePeerDelete endpoint")

	parts := strings.Split(r.RequestURI, "?")
	if len(parts) > 1 {
		Peers.LocalDelete(parts[1])
	}
}

// HandlePeerList deals with peer list requests coming in over the network.
func HandlePeerList(w http.ResponseWriter, r *http.Request) {
	log.Println("PeerList:HandlePeerList endpoint")

	var result []string
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
	log.Println("WebMind:Internal:HandleKeepAlive")

	defer r.Body.Close()
	sender := strings.Split(r.RequestURI, "?")

	if len(sender) <= 1 {
		return
	}
	peer := Peers.Users[sender[1]]
	if peer == nil {
		peer = NewPeer(sender[1])
	}
	peer.lastSeen = time.Now().Round(0)
	// log.Printf("HandleKeepAlive from %v", sender[1])
	p.LocalAdd(sender[1])
	// p.RemoteAddToAll(sender[1])
	fmt.Fprintf(w, "I'm still here...")
}
