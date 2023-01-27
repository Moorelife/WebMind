package trace

import (
	"fmt"
	"log"
	"net/http"
)

var trace = false // if true, logs enters and exits

func Entered(name string) {
	if trace {
		log.Printf("%v entered\n", name)
	}
}

func Exited(name string) {
	if trace {
		log.Printf("%v exited\n", name)
	}
}

func HandleTraceOn(w http.ResponseWriter, r *http.Request) {
	Entered("traceOn endpoint")
	defer Exited("traceOn endpoint")

	trace = true

	fmt.Fprintf(w, "Tracing activated")
}

func HandleTraceOff(w http.ResponseWriter, r *http.Request) {
	Entered("traceOff endpoint")
	defer Exited("traceOff endpoint")

	trace = false
	fmt.Fprintf(w, "Tracing deactivated")
}
