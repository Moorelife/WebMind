package webmind

import (
	"context"
	"flag"
	"fmt"
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

	originServer := flag.String("origin", "", "origin server address")
	webPort := flag.String("port", "7777", "https server port number")
	flag.Parse()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "origin", *originServer)
	ctx = context.WithValue(ctx, "port", *webPort)

	log.Printf("origin: %v", ctx.Value("origin"))
	log.Printf("port: %v", ctx.Value("port"))

	return ctx
}

// SetupLogging sets up logging and stores logging related arguments in the context if needed.
func SetupLogging(ctx context.Context) context.Context {
	trace.Entered("WebMind:Internal:SetupLogging")
	defer trace.Exited("WebMind:Internal:SetupLogging")

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)
	log.Printf("WebMind started on port %v\n", ctx.Value("port"))
	return ctx
}

// RetrievePublicAddress retrieves the public address and places it in the context.
// It returns an error if the public address resolver cannot process the request.
func RetrievePublicAddress(ctx context.Context) (context.Context, error) {
	trace.Entered("WebMind:Internal:RetrievePublicAddress")
	defer trace.Exited("WebMind:Internal:RetrievePublicAddress")

	address, err := ip.GetPublicIP()
	if err != nil {
		log.Printf("GetPublicIP failed: %v", err)
		return ctx, err
	}
	address = fmt.Sprintf("%v:%v", strings.Trim(address, " "), ctx.Value("port"))
	ctx = context.WithValue(ctx, "selfAddress", address)

	log.Printf("selfAddress: %v", ctx.Value("selfAddress"))

	return ctx, err
}

func CreateAndRetrievePeerList(ctx context.Context) {
	peerlist.LocalAdd(fmt.Sprintf("%s", ctx.Value("selfAddress")))
	if fmt.Sprintf("%s", ctx.Value("origin")) != "" {
		peerlist.RemoteGet(fmt.Sprintf("%s", ctx.Value("origin")))
	}
}

// SendPeerAddRequests sends a peer add request to each system in the peer list.
func SendPeerAddRequests(ctx context.Context) {
	trace.Entered("WebMind:Internal:SendPeerAddRequests")
	defer trace.Exited("WebMind:Internal:SendPeerAddRequests")
	log.Printf("PEERLIST: %v", peerlist.Peers)
	peerlist.RemoteAddToAll(fmt.Sprintf("%s", ctx.Value("selfAddress")))
}

// SetupExitHandler catches the Ctrl-C signal and executes any needed cleanup.
func SetupExitHandler(ctx context.Context) {
	trace.Entered("WebMind:Internal:SetupExitHandler")
	defer trace.Exited("WebMind:Internal:SetupExitHandler")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("***** Ctrl-C pressed: %v *****\n", sig)
			peerlist.LocalDelete(fmt.Sprintf("%s", ctx.Value("selfAddress")))
			peerlist.RemoteDeleteToAll(fmt.Sprintf("%s", ctx.Value("selfAddress")))

			os.Exit(0)
		}
	}()
}

func HandleRequests(port string) {
	trace.Entered("WebMind:Internal:HandleRequests")
	defer trace.Exited("WebMind:Internal:HandleRequests")

	// basic endpoints
	http.HandleFunc("/", serverRoot)
	http.HandleFunc("/keepalive", keepAlive)

	http.HandleFunc("/trace/on", trace.HandleTraceOn)
	http.HandleFunc("/trace/off", trace.HandleTraceOff)

	// peerlist endpoints
	http.HandleFunc("/peer/add", peerlist.HandlePeerAdd)
	http.HandleFunc("/peer/list", peerlist.HandlePeerList)
	http.HandleFunc("/peer/delete", peerlist.HandlePeerDelete)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// basic operations endpoints
func serverRoot(w http.ResponseWriter, r *http.Request) {
	trace.Entered("WebMind:Internal:serverRoot")
	defer trace.Exited("WebMind:Internal:serverRoot")

	printRequest(r)

	defer r.Body.Close()

	fmt.Fprintf(w, "WebMind up and running!")
}

// basic operations endpoints
func keepAlive(w http.ResponseWriter, r *http.Request) {
	trace.Entered("WebMind:Internal:keepAlive")
	defer trace.Exited("WebMind:Internal:keepAlive")
	defer r.Body.Close()
	sender := strings.Split(r.RequestURI, "?")
	log.Printf("keepalive from %v", sender[1])

	fmt.Fprintf(w, "I'm still here...")
}

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

// SendKeepAlive starts a go routine that sends a /keepalive request to all peers every two seconds.
func SendKeepAlive(ctx context.Context) {
	trace.Entered("WebMind:Internal:SendKeepAlives")
	defer trace.Exited("WebMind:Internal:SendKeepAlives")

	go func() {
		self := fmt.Sprintf("%v", ctx.Value("selfAddress"))
		for true {
			for key, peer := range peerlist.Peers {
				if key != self {
					url := fmt.Sprintf("http://%v/keepalive?%v", key, self)
					_, err := http.Get(url)
					if err != nil {
						log.Printf("Failed to send keepalive to %v (ERROR: %v)\n", peer, err)
					}
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()
}
