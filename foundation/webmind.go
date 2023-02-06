// Package webmind is the core code of WebMind 2.0.

package foundation

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
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

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func StartNode(port string) {
	cmd := exec.Command("app\\services\\node\\startnode.cmd", port)
	// log.Printf("CMD: %#v", cmd)
	cmd.Start()
	log.Printf("Node started on port %s", port)
}

func PrintRequest(r *http.Request) {
	log.Printf("-   Remote address: %v", r.RemoteAddr)
	log.Printf("-   Request URI %v", r.RequestURI)
	log.Printf("-   Address: %v", r.Host)
	log.Printf("-   Method: %v", r.Method)
}
