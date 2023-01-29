package localnode

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"

	"github.com/Moorelife/WebMind/internal/ip"
	"github.com/Moorelife/WebMind/internal/peerlist"
	"github.com/Moorelife/WebMind/internal/trace"
)

type LocalNode struct {
	OriginsFile  *string  // command line: filename of the file specifying the origin node addresses.
	LocalAddress string   // calculated: combination of public IP and port number of the local node.
	LocalPort    *string  // command line: port number for this instance of the program.
	FullList     *bool    // command line: specifies if full peer list is logged instead of just count.
	Trace        *bool    // command line: specifies if call trace log statements are executed.
	LogFile      *os.File // logging: maintains the file handle of the log file, to close on exit.
}

func NewLocalNode() *LocalNode {
	private := LocalNode{}
	private.parseArgs()
	private.LocalAddress = ip.GetPublicIP() + ":" + *private.LocalPort
	return &private
}

func (l *LocalNode) parseArgs() {
	l.OriginsFile = flag.String("origins", "./origins.json", "origins list file")
	l.LocalPort = flag.String("port", "7777", "http server port number")
	l.FullList = flag.Bool("full", true, "switch to list peers, not just list count")
	l.Trace = flag.Bool("trace", true, "switch to activate call tracing")
	flag.Parse()
}

// SetupLogging sets up logging and stores logging related arguments in the LocalNode struct if needed.
func (l *LocalNode) SetupLogging() {
	trace.Entered("WebMind:Internal:SetupLogging")
	defer trace.Exited("WebMind:Internal:SetupLogging")

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	logFileName := fmt.Sprintf("..\\..\\logs\\%v.log", *l.LocalPort)

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.Printf("WebMind started on port %v", *l.LocalPort)
}

// RetrievePublicAddress retrieves the public address and places it in the LocalNode struct.
func (l *LocalNode) RetrievePublicAddress() {
	trace.Entered("WebMind:Internal:RetrievePublicAddress")
	defer trace.Exited("WebMind:Internal:RetrievePublicAddress")

	address := ip.GetPublicIP()
	address = fmt.Sprintf("%v:%v", strings.Trim(address, " "), *l.LocalPort)
	err := os.Setenv("WEBMIND_LOCALPEER", address)
	if err != nil {
		log.Printf("Could not set WEBMIND_LOCALPEER environment variable")
	}

}

func (l *LocalNode) CreateAndRetrievePeerList() {
	trace.Entered("WebMind:Internal:CreateAndRetrievePeerList")
	defer trace.Exited("WebMind:Internal:CreateAndRetrievePeerList")

	peerlist.Peers.LocalAdd(l.LocalAddress)

	origins, err := os.ReadFile(*l.OriginsFile)
	if err != nil || len(origins) == 0 {
		log.Printf("os.ReadFile() returned no data: %v", err)
	}
	originList := make([]string, 1)

	err = json.Unmarshal(origins, &originList)
	if err != nil {
		log.Printf("Unmarshalling originsList failed: %v", err)
		return
	}

	for _, origin := range originList {
		go peerlist.Peers.RemoteGet(origin)
	}
}

// SendPeerAddRequests sends a peer add request to each system in the peer list.
func (l *LocalNode) SendPeerAddRequests() {
	trace.Entered("WebMind:Internal:SendPeerAddRequests")
	defer trace.Exited("WebMind:Internal:SendPeerAddRequests")
	log.Printf("PEERLIST: %v", peerlist.Peers)
	peerlist.Peers.RemoteAddToAll(l.LocalAddress)
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
			peerlist.Peers.LocalDelete(l.LocalAddress)
			peerlist.Peers.RemoteDeleteToAll(l.LocalAddress)

			l.LogFile.Close()
			os.Exit(0)
		}
	}()
}

func (l *LocalNode) HandleRequests() {
	http.HandleFunc("/", HandleServerRootRequests)
	http.HandleFunc("/peer/add", peerlist.HandlePeerAdd)
	http.HandleFunc("/peer/list", peerlist.HandlePeerList)
	http.HandleFunc("/peer/delete", peerlist.HandlePeerDelete)
	http.HandleFunc("/peer/keepalive", peerlist.Peers.HandleKeepAlive)

	log.Fatal(http.ListenAndServe(l.LocalAddress, nil))
}

// StartSendingKeepAlive starts a go routine that sends a /keepalive request to all peers every two seconds.
func (l *LocalNode) StartSendingKeepAlive() {
	trace.Entered("WebMind:Internal:StartSendingKeepAlive")
	defer trace.Exited("WebMind:Internal:StartSendingKeepAlive")

	go func() {
		for true {
			log.Printf("SendKeepAlives still running")
			for key, _ := range peerlist.Peers.Users {
				if key != l.LocalAddress {
					url := fmt.Sprintf("http://%v/peer/keepalive?%v", key, l.LocalAddress)
					log.Printf("sending keepalive to %v", key)
					_, err := http.Get(url)
					if err != nil {
						log.Printf("Failed to send keepalive to %v", key)
					}
				}
			}
			peerlist.Peers.CleanPeerList(l.LocalAddress)
			peerlist.Peers.LogLocalList(*l.FullList)
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
