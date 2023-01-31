package localnode

import (
	"encoding/json"
	"fmt"
	"github.com/Moorelife/WebMind/internal/webmind1/ip"
	"github.com/Moorelife/WebMind/internal/webmind1/peerlist"
	"github.com/Moorelife/WebMind/internal/webmind1/trace"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"time"
)

type LocalNode struct {
	LocalAddress string `json:"LocalAddress"`
	LocalPort    int    `json:"LocalPort"`

	OriginsFile string
	FullList    bool
	Trace       bool

	peers *peerlist.PeerList

	logFile *os.File
}

func NewLocalNode(address string, port int) *LocalNode {
	creation := LocalNode{
		LocalAddress: address,
		LocalPort:    port,
		peers:        peerlist.NewPeerList(),
	}

	return &creation
}

// SetupLogging sets up logging and stores logging related arguments in the LocalNode struct if needed.
func (l *LocalNode) SetupLogging() {
	trace.Entered("WebMind:Internal:SetupLogging")
	defer trace.Exited("WebMind:Internal:SetupLogging")

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	logFileName := fmt.Sprintf("%v.log", l.LocalPort)

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.Printf("started on port %v", l.LocalPort)
}

// RetrievePublicAddress retrieves the public address and places it in the LocalNode struct.
func (l *LocalNode) RetrievePublicAddress() {
	trace.Entered("WebMind:Internal:RetrievePublicAddress")
	defer trace.Exited("WebMind:Internal:RetrievePublicAddress")

	address := ip.GetPublicIP()
	address = fmt.Sprintf("%v:%v", strings.Trim(address, " "), l.LocalPort)
	err := os.Setenv("WEBMIND_LOCALPEER", address)
	if err != nil {
		log.Printf("Could not set WEBMIND_LOCALPEER environment variable")
	}

}

func (l *LocalNode) CreateAndRetrievePeerList() {
	trace.Entered("WebMind:Internal:CreateAndRetrievePeerList")
	defer trace.Exited("WebMind:Internal:CreateAndRetrievePeerList")

	l.peers.LocalAdd(l.LocalAddress)

	origins, err := os.ReadFile(l.OriginsFile)
	if err != nil || len(origins) == 0 {
		panic(fmt.Sprintf("os.ReadFile() returned no data: %v", err))
	}
	originList := make([]string, 1)

	err = json.Unmarshal(origins, &originList)
	if err != nil {
		panic(fmt.Sprintf("Unmarshalling originsList failed: %v", err))
		return
	}

	for _, origin := range originList {
		go l.peers.RemoteGet(origin)
	}
}

// SendPeerAddRequests sends a peer add request to each system in the peer list.
func (l *LocalNode) SendPeerAddRequests() {
	trace.Entered("WebMind:Internal:SendPeerAddRequests")
	defer trace.Exited("WebMind:Internal:SendPeerAddRequests")
	log.Printf("PEERLIST: %v", l)
	l.peers.RemoteAddToAll(l.LocalAddress)
}

// SetupExitHandler catches the Ctrl-C signal and executes any needed cleanup.
func (l *LocalNode) SetupExitHandler() {
	trace.Entered("WebMind:Internal:SetupExitHandler")
	defer trace.Exited("WebMind:Internal:SetupExitHandler")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("***** Ctrl-C pressed: %v *****\n", sig)
			l.peers.LocalDelete(l.LocalAddress)
			l.peers.RemoteDeleteToAll(l.LocalAddress)

			l.logFile.Close()
			os.Exit(0)
		}
	}()
}

func (l *LocalNode) HandleRequests() {
	http.HandleFunc("/", HandleServerRootRequests)
	http.HandleFunc("/peer/add", l.peers.HandlePeerAdd)
	http.HandleFunc("/peer/list", l.peers.HandlePeerList)
	http.HandleFunc("/peer/delete", l.peers.HandlePeerDelete)
	http.HandleFunc("/peer/keepalive", l.peers.HandleKeepAlive)

	localAddress := l.LocalAddress + ":" + strconv.Itoa(l.LocalPort)
	log.Fatal(http.ListenAndServe(localAddress, nil))
}

// StartSendingKeepAlive starts a go routine that sends a /keepalive request to all peers every two seconds.
func (l *LocalNode) StartSendingKeepAlive() {
	trace.Entered("WebMind:Internal:StartSendingKeepAlive")
	defer trace.Exited("WebMind:Internal:StartSendingKeepAlive")

	go func() {
		for true {
			log.Printf("SendKeepAlives still running")
			for key, _ := range l.peers.Users {
				if key != l.LocalAddress {
					url := fmt.Sprintf("http://%v/peer/keepalive?%v", key, l.LocalAddress)
					log.Printf("sending keepalive to %v", key)
					_, err := http.Get(url)
					if err != nil {
						log.Printf("Failed to send keepalive to %v", key)
					}
				}
			}
			l.peers.CleanPeerList(l.LocalAddress)
			l.peers.LogLocalList(l.FullList)
			time.Sleep(peerlist.KeepAliveInterval)
		}
		log.Print("***************************************************")
		log.Print("***** KEEPALIVE GO FUNCTION IS TERMINATING!!! *****")
		log.Print("***************************************************")
	}()
}

func HandleServerRootRequests(w http.ResponseWriter, r *http.Request) {
	log.Println("WebMind:Internal:HandleServerRootRequests")

	printRequest(r)
	defer r.Body.Close()

	fmt.Fprintf(w, "WebMind up and running!")
}

// basic operations endpoints
func printRequest(r *http.Request) {
	log.Printf("-   Remote address: %v", r.RemoteAddr)
	log.Printf("-   Request URI %v", r.RequestURI)
	log.Printf("-   Address: %v", r.Host)
	log.Printf("-   Method: %v", r.Method)
	log.Printf("-   ContentLength: %v", r.ContentLength)
	log.Println("-   Header:")
	printHeaderMap(r.Header)
}

func printHeaderMap(header http.Header) {
	type KeyValue struct {
		Key   string
		Value []string
	}

	s := make([]KeyValue, 0, len(header))

	for k, v := range header {
		s = append(s, KeyValue{k, v})
	}

	sort.SliceStable(s, func(i, j int) bool {
		return s[i].Key < s[j].Key
	})

	for _, v := range s {
		// if v.Key == "User-Agent" {
		log.Println("-      ", v.Key, ": ", v.Value)
		// }
	}
}
