package peerlist

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

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
			args: args{hostPort: "localhost:14285"},
			want: &Peer{
				addressPort: "localhost:14285",
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
	type args struct {
		ownAddress string
	}
	tests := []struct {
		name string
		p    PeerList
		args args
	}{
		{
			name: "RemoteAdd",
			p:    PeerList{"1.2.3.4:5": {"1.2.3.4:5", time.Now()}},
			args: args{
				ownAddress: "1.2.3.4:5",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.RemoteAddToAll(tt.args.ownAddress)
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
			p:    PeerList{"1.2.3.4:5": {"1.2.3.4:5", time.Now()}},
			args: args{
				addressToAdd: "1.2.3.4:5",
				sendTo:       "6.7.8.9:10",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.RemoteAdd(tt.args.addressToAdd, tt.args.sendTo)
		})
	}
}

func TestPeerList_LocalAdd(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: "192.168.2.111:14285", lastSeen: time.Now()}
	var peerList = PeerList{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}
	type args struct {
		hostPort string
	}
	tests := []struct {
		name string
		p    PeerList
		args args
	}{
		{name: "HandleLocalDeleteGoodFlow", p: peerList, args: args{hostPort: "localhost:14285"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.LocalAdd(tt.args.hostPort)
		})
	}
}

func TestPeerList_LocalDelete(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: "192.168.2.111:14285", lastSeen: time.Now()}
	Peers = PeerList{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}
	type args struct {
		hostPort string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "HandleLocalDeleteGoodFlow", args: args{hostPort: "localhost:14285"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Peers.LocalDelete(tt.args.hostPort)
		})
	}
}

func TestPeerList_CleanPeerList(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now().Add(-20 * time.Second)}
	var Peer2 = Peer{addressPort: "192.168.2.111:14285", lastSeen: time.Now()}
	Peers = PeerList{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}
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

func TestPeerList_KeepAlive(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: "1.2.3.4:5", lastSeen: time.Now()}
	Peers = PeerList{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}
	var reader io.Reader
	request, err := http.NewRequest(http.MethodGet, "http://localhost:14285/peer/keepalive?1.2.3.4:5", reader)
	request.RequestURI = "http://localhost:14285/peer/keepalive?1.2.3.4:5"
	requestNok, err := http.NewRequest(http.MethodGet, "http://localhost:14285/peer/keepalive", reader)
	requestNok.RequestURI = "http://localhost:14285/peer/keepalive"
	if err != nil {
		t.Fatal("TEST")
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "HandleKeepAliveGoodFlow", args: args{w: httptest.NewRecorder(), r: request}},
		{name: "HandleKeepAliveNok", args: args{w: httptest.NewRecorder(), r: requestNok}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Peers.KeepAlive(tt.args.w, tt.args.r)
		})
	}
}

func TestPeerList_KeepAliveNoPeers(t *testing.T) {
	Peers = nil
	var reader io.Reader
	request, err := http.NewRequest(http.MethodGet, "http://localhost:14285/peer/keepalive?1.2.3.4:5", reader)
	request.RequestURI = "http://localhost:14285/peer/keepalive?1.2.3.4:5"
	if err != nil {
		t.Fatal("TEST")
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "HandleKeepAliveGoodFlow", args: args{w: httptest.NewRecorder(), r: request}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Peers.KeepAlive(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlePeerAdd(t *testing.T) {
	var Peer1 = Peer{addressPort: "localHost:14285", lastSeen: time.Now()}
	var Peer2 = Peer{addressPort: "192.168.2.111:7777", lastSeen: time.Now()}
	Peers = PeerList{Peer1.addressPort: &Peer1, Peer2.addressPort: &Peer2}
	var reader io.Reader
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:14285/peer/add?1.2.3.4:5", reader)
	request.RequestURI = "http://localhost:14285/peer/add?1.2.3.4:5"
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
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:14285/peer/delete?1.2.3.4:5", reader)
	request.RequestURI = "http://localhost:14285/peer/delete?1.2.3.4:5"
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
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:14285/peer/list", reader)
	request.RequestURI = "http://localhost:14285/peer/peer/list"
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
