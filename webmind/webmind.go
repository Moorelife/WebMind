package main

import (
	"flag"
	"fmt"
	"github.com/Moorelife/WebMind/internal/ip"
	"github.com/Moorelife/WebMind/internal/trace"
	"log"
	"net/http"
	"sort"

	"github.com/Moorelife/WebMind/pkg/peerlist"
)

func main() {
	originServer := flag.String("origin", "", "origin server address")
	webPort := flag.Int("port", 7777, "https server port number")
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)
	log.Printf("WebMind started on port %v\n", *webPort)

	address, err := ip.GetPublicIP()
	if err != nil {
		log.Printf("GetPublicIP failed: %v", err)
	}
	selfAddrPort := fmt.Sprintf("%v:%v", address, *webPort)
	log.Println(selfAddrPort)
	peerlist.Add(selfAddrPort)
	if *originServer != "" {
		peerlist.Get(*originServer)
	}

	handleRequests(*webPort)
}

func handleRequests(port int) {

	// basic endpoints
	http.HandleFunc("/", serverRoot)
	http.HandleFunc("/trace/on", trace.HandleTraceOn)
	http.HandleFunc("/trace/off", trace.HandleTraceOff)

	// peerlist endpoints
	http.HandleFunc("/peer/list", peerlist.HandlePeerList)
	http.HandleFunc("/peer/add", peerlist.HandlePeerAdd)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// basic operations endpoints
func serverRoot(w http.ResponseWriter, r *http.Request) {
	trace.Entered("serverRoot endpoint")
	defer trace.Exited("serverRoot endpoint")

	printRequest(r)

	defer r.Body.Close()

	fmt.Fprintf(w, "WebMind up and running!")
}

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
