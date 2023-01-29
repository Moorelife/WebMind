package trace

import (
	"log"
)

var trace = true

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
