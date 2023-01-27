package webmind

import (
	"context"
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

// ParseArgsToContext parses all command line arguments and adds them to a context.
func ParseArgsToContext() context.Context {
	trace.Entered("WebMind:Internal:ParseArgsToContext")
	defer trace.Exited("WebMind:Internal:ParseArgsToContext")

	originsFile := flag.String("origins", "./origins.json", "origins list file")
	webPort := flag.String("port", "7777", "http server port number")
	flag.Parse()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "origins", *originsFile)
	ctx = context.WithValue(ctx, "port", *webPort)

	log.Printf("origin: %v", ctx.Value("origins"))
	log.Printf("port: %v", ctx.Value("port"))

	return ctx
}

// SetupLogging sets up logging and stores logging related arguments in the context if needed.
func SetupLogging(ctx context.Context) (context.Context, *os.File) {
	trace.Entered("WebMind:Internal:SetupLogging")
	defer trace.Exited("WebMind:Internal:SetupLogging")

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	logFileName := fmt.Sprintf("..\\..\\logs\\%v.log", ctx.Value("port"))

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.Printf("WebMind started on port %v\n", ctx.Value("port"))
	return ctx, logFile
}

// RetrievePublicAddress retrieves the public address and places it in the context.
// It returns an error if the public address resolver cannot process the request.
func RetrievePublicAddress(ctx context.Context) context.Context {
	trace.Entered("WebMind:Internal:RetrievePublicAddress")
	defer trace.Exited("WebMind:Internal:RetrievePublicAddress")

	address := ip.GetPublicIP()
	address = fmt.Sprintf("%v:%v", strings.Trim(address, " "), ctx.Value("port"))
	ctx = context.WithValue(ctx, "selfAddress", address)

	log.Printf("selfAddress: %v", ctx.Value("selfAddress"))

	return ctx
}

func CreateAndRetrievePeerList(ctx context.Context) {
	trace.Entered("WebMind:Internal:CreateAndRetrievePeerList")
	defer trace.Exited("WebMind:Internal:CreateAndRetrievePeerList")
	peerlist.Peers.LocalAdd(fmt.Sprintf("%s", ctx.Value("selfAddress")))

	origins, err := os.ReadFile(fmt.Sprintf("%s", ctx.Value("origins")))
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
func SendPeerAddRequests(ctx context.Context) {
	trace.Entered("WebMind:Internal:SendPeerAddRequests")
	defer trace.Exited("WebMind:Internal:SendPeerAddRequests")
	log.Printf("PEERLIST: %v", peerlist.Peers)
	peerlist.Peers.RemoteAddToAll(fmt.Sprintf("%s", ctx.Value("selfAddress")))
}

// SetupExitHandler catches the Ctrl-C signal and executes any needed cleanup.
func SetupExitHandler(ctx context.Context, logFile *os.File) {
	trace.Entered("WebMind:Internal:SetupExitHandler")
	defer trace.Exited("WebMind:Internal:SetupExitHandler")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("***** Ctrl-C pressed: %v *****\n", sig)
			peerlist.Peers.LocalDelete(fmt.Sprintf("%s", ctx.Value("selfAddress")))
			peerlist.Peers.RemoteDeleteToAll(fmt.Sprintf("%s", ctx.Value("selfAddress")))

			logFile.Close()
			os.Exit(0)
		}
	}()
}

func HandleRequests(port string) {
	trace.Entered("WebMind:Internal:HandleRequests")
	defer trace.Exited("WebMind:Internal:HandleRequests")

	// basic endpoints
	http.HandleFunc("/", HandleServerRootRequests)

	http.HandleFunc("/trace/on", trace.HandleTraceOn)
	http.HandleFunc("/trace/off", trace.HandleTraceOff)

	// peerlist endpoints
	http.HandleFunc("/peer/add", peerlist.HandlePeerAdd)
	http.HandleFunc("/peer/list", peerlist.HandlePeerList)
	http.HandleFunc("/peer/delete", peerlist.HandlePeerDelete)
	http.HandleFunc("/peer/keepalive", peerlist.Peers.HandleKeepAlive)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// basic operations endpoints
func HandleServerRootRequests(w http.ResponseWriter, r *http.Request) {
	trace.Entered("WebMind:Internal:HandleServerRootRequests")
	defer trace.Exited("WebMind:Internal:HandleServerRootRequests")

	printRequest(r)

	defer r.Body.Close()

	fmt.Fprintf(w, "WebMind up and running!")
}

// basic operations endpoints
func printRequest(r *http.Request) {
	trace.Entered("WebMind:Internal:printRequest")
	defer trace.Exited("WebMind:Internal:printRequest")

	log.Printf("-   Remote address: %v", r.RemoteAddr)
	log.Printf("-   Request URI %v", r.RequestURI)
	log.Printf("-   Address: %v", r.Host)
	log.Printf("-   Method: %v", r.Method)
	log.Printf("-   ContentLength: %v", r.ContentLength)
	log.Println("-   Header:")
	printHeaderMap(r.Header)
}

func printHeaderMap(header http.Header) {
	trace.Entered("WebMind:Internal:printHeaderMap")
	defer trace.Exited("WebMind:Internal:printHeaderMap")

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

// StartSendingKeepAlive starts a go routine that sends a /keepalive request to all peers every two seconds.
func StartSendingKeepAlive(ctx context.Context) {
	trace.Entered("WebMind:Internal:StartSendingKeepAlive")
	defer trace.Exited("WebMind:Internal:StartSendingKeepAlive")

	go func() {
		self := fmt.Sprintf("%v", ctx.Value("selfAddress"))
		for true {
			log.Printf("SendKeepAlives still running")
			for key, _ := range peerlist.Peers.Users {
				if key != self {
					url := fmt.Sprintf("http://%v/peer/keepalive?%v", key, self)
					log.Printf("sending keepalive to %v", key)
					_, err := http.Get(url)
					if err != nil {
						log.Printf("Failed to send keepalive to %v", key)
					}
				}
			}
			peerlist.Peers.CleanPeerList(fmt.Sprintf("%s", ctx.Value("selfAddress")))
			peerlist.Peers.LogLocalList(peerlist.CountOnly)
			time.Sleep(peerlist.KeepAliveInterval)
		}
		log.Print("***************************************************")
		log.Print("***** KEEPALIVE GO FUNCTION IS TERMINATING!!! *****")
		log.Print("***************************************************")
	}()
}
