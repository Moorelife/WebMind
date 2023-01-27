package peerlist

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// The tests of the WebMind program assume that the origin node is running, isolated from all other instances
// of WebMind. We start the origin node at the following address:
var originNode = "192.168.2.111:7777"

func TestNewPeer(t *testing.T) {
	type args struct {
		hostPort string
	}
	tests := []struct {
		name string
		args args
		want *Peer
	}{
		{
			name: "NewPeerGoodFlow",
			args: args{hostPort: originNode},
			want: &Peer{
				addressPort: originNode,
				lastSeen:    time.Now(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPeer(tt.args.hostPort); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPeer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerList_RemoteAddToAll(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: originNode, lastSeen: time.Now()}
	Peers = PeerList{Users: map[string]*Peer{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}}
	type args struct {
		ownAddress string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "RemoteAdd",
			args: args{
				ownAddress: originNode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Peers.RemoteAddToAll(tt.args.ownAddress)
		})
	}
}

func TestPeerList_RemoteAdd(t *testing.T) {
	type args struct {
		addressToAdd string
		sendTo       string
	}
	tests := []struct {
		name string
		p    PeerList
		args args
	}{
		{
			name: "RemoteAdd",
			p:    PeerList{Users: map[string]*Peer{originNode: {originNode, time.Now()}}},
			args: args{
				addressToAdd: originNode,
				sendTo:       originNode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.RemoteAdd(tt.args.addressToAdd, tt.args.sendTo)
		})
	}
}

func TestPeerList_RemoteDeleteToAll(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: originNode, lastSeen: time.Now()}
	Peers = PeerList{Users: map[string]*Peer{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}}
	type args struct {
		exceptAddress string
	}
	tests := []struct {
		name string
		p    PeerList
		args args
	}{
		{
			name: "HandleLocalRemoteGetGoodFlow",
			args: args{exceptAddress: originNode},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.RemoteDeleteToAll(tt.args.exceptAddress)
		})
	}
}

func TestPeerList_RemoteDelete(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: originNode, lastSeen: time.Now()}
	peers := PeerList{Users: map[string]*Peer{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}}
	type args struct {
		selfAddress string
		sendTo      string
	}
	tests := []struct {
		name string
		p    PeerList
		args args
	}{
		{
			name: "HandleLocalRemoteGetGoodFlow",
			p:    peers,
			args: args{selfAddress: originNode, sendTo: "192.168.2.222:65535"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.RemoteDelete(tt.args.selfAddress, tt.args.sendTo)
		})
	}
}

func TestPeerList_RemoteGet(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: originNode, lastSeen: time.Now()}
	Peers = PeerList{Users: map[string]*Peer{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}}
	type args struct {
		hostPort string
	}
	tests := []struct {
		name string
		args args
		want []Peer
	}{
		{
			name: "HandleLocalRemoteGetGoodFlow",
			args: args{hostPort: originNode},
			want: Peers.LocalGet(originNode),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Peers.RemoteGet(tt.args.hostPort); len(got) < 1 {
				t.Errorf("len(RemoteGet()) < 1:  %v", len(got))
			}
		})
	}
}

func TestPeerList_RemoteGetNilBody(t *testing.T) {
	var Peer1 = Peer{addressPort: originNode, lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: "86.89.186.20:14285", lastSeen: time.Now()}
	var Peer3 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	Peers = PeerList{Users: map[string]*Peer{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2, Peer3.addressPort: &Peer3}}
	type args struct {
		hostPort string
	}
	tests := []struct {
		name string
		args args
		want []Peer
	}{
		{
			name: "HandleLocalRemoteGetNilBody",
			args: args{hostPort: "localhost:65535"},
			want: []Peer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Peers.RemoteGet(tt.args.hostPort); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerList_LocalAdd(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: originNode, lastSeen: time.Now()}
	peers := PeerList{Users: map[string]*Peer{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}}
	type args struct {
		hostPort string
	}
	tests := []struct {
		name string
		p    PeerList
		args args
	}{
		{name: "HandleLocalDeleteGoodFlow", p: peers, args: args{hostPort: "localhost:14285"}},
		{name: "HandleLocalDeleteNilPeers", p: PeerList{}, args: args{hostPort: originNode}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.LocalAdd(tt.args.hostPort)
		})
	}
}

func TestPeerList_LocalDelete(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: originNode, lastSeen: time.Now()}
	Peers = PeerList{Users: map[string]*Peer{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}}
	type args struct {
		hostPort string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "HandleLocalDeleteGoodFlow", args: args{hostPort: "192.168.2.222:65535"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Peers.LocalDelete(tt.args.hostPort)
		})
	}
}

func TestPeerList_CleanPeerList(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now().Add(-20 * time.Second)}
	var Peer2 = Peer{addressPort: originNode, lastSeen: time.Now()}
	Peers = PeerList{Users: map[string]*Peer{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}}
	type args struct {
		exceptAddress string
	}
	tests := []struct {
		name string
		p    PeerList
		args args
	}{
		{name: "HandleLocalDeleteGoodFlow", args: args{exceptAddress: "localhost:7777"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Peers.CleanPeerList(tt.args.exceptAddress)
		})
	}
}

func TestPeerList_HandleKeepAliveNoPeers(t *testing.T) {
	Peers = PeerList{}
	var reader io.Reader
	requestURL := fmt.Sprintf("http://%s/peer/HandleKeepAlive?%s", originNode, originNode)
	request, _ := http.NewRequest(http.MethodGet, requestURL, reader)
	request.RequestURI = requestURL
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "HandleHandleKeepAliveGoodFlow", args: args{w: httptest.NewRecorder(), r: request}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Peers.HandleKeepAlive(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlePeerAdd(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: originNode, lastSeen: time.Now()}
	Peers = PeerList{Users: map[string]*Peer{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}}
	var reader io.Reader
	requestURL := fmt.Sprintf("http://%s/peer/add?%s", originNode, originNode)
	request, _ := http.NewRequest(http.MethodGet, requestURL, reader)
	request.RequestURI = requestURL
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "HandlePeerAddGoodFlow", args: args{w: httptest.NewRecorder(), r: request}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandlePeerAdd(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlePeerDelete(t *testing.T) {
	var reader io.Reader
	requestURL := fmt.Sprintf("http://%s/peer/delete?%s", originNode, originNode)
	request, _ := http.NewRequest(http.MethodGet, requestURL, reader)
	request.RequestURI = requestURL
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "HandlePeerDeleteGoodFlow", args: args{w: httptest.NewRecorder(), r: request}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandlePeerDelete(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlePeerList(t *testing.T) {
	var reader io.Reader
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/peer/peer/list", originNode), reader)
	request.RequestURI = fmt.Sprintf("http://%s/peer/peer/list", originNode)
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "HandlePeerListGoodFlow", args: args{w: httptest.NewRecorder(), r: request}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandlePeerList(tt.args.w, tt.args.r)
		})
	}
}

func TestPeerList_logLocalList(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: originNode, lastSeen: time.Now()}
	peers := PeerList{Users: map[string]*Peer{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}}
	tests := []struct {
		name string
		p    PeerList
		cnt  bool
	}{
		{"logLocalListGoodFlow", peers, false},
		{"logLocalListCountOnly", peers, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.logLocalList(tt.cnt)
		})
	}
}

func TestPeerList_HandleKeepAlive(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: originNode, lastSeen: time.Now()}
	Peers = PeerList{Users: map[string]*Peer{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}}
	var reader io.Reader
	requestURL := fmt.Sprintf("http://%s/peer/delete?%s", originNode, originNode)
	request, _ := http.NewRequest(http.MethodGet, requestURL, reader)
	request.RequestURI = requestURL
	requestURL = fmt.Sprintf("http://%s/peer/delete?", originNode)
	requestNok, _ := http.NewRequest(http.MethodGet, requestURL, reader)
	requestNok.RequestURI = requestURL
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "HandleHandleKeepAliveGoodFlow", args: args{w: httptest.NewRecorder(), r: request}},
		{name: "HandleHandleKeepAliveNok", args: args{w: httptest.NewRecorder(), r: requestNok}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Peers.HandleKeepAlive(tt.args.w, tt.args.r)
		})
	}
}
