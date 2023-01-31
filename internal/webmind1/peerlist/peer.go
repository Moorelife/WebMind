package peerlist

import (
	"log"
	"time"
)

// Peer holds all information about a network peer.
type Peer struct {
	addressPort string
	lastSeen    time.Time
	// public key?
}

// NewPeer creates a new Peer with the given host and port.
func NewPeer(hostPort string) *Peer {
	result := Peer{hostPort, time.Now()}
	return &result
}

func (p *Peer) is(localAddress string) bool {
	return p.addressPort == localAddress
}

func (p *Peer) timedOut() bool {
	return p.lastSeen.Before(time.Now().Add(-2 * KeepAliveInterval))
}

func (p *Peer) refresh() {
	p.lastSeen = time.Now().Round(0)
}

func (p *Peer) log() {
	log.Printf("%v  	%v", p.addressPort, p.lastSeen.Round(0))
}
