// Package webmind is the core code of WebMind 2.0.

package webmind

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

// WebMind contains programming that the system and process "classes" require.
type WebMind struct {
}

// General Utility functions =========================================

// SetupLogging sets up logging and stores logging related arguments in the LocalNode struct if needed.
func SetupLogging(program string) {
	saneName := strings.Replace(program, ".", "_", -1)
	saneName = strings.Replace(saneName, ":", "_", -1)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	logFileName := fmt.Sprintf("%v.log", saneName)

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func PrintRequest(r *http.Request) {
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
