package main

import (
	"fmt"
	"log"

	"github.com/Moorelife/WebMind/internal/webmind"
)

// WebMind in its current state is JUST A LEARNING EXPERIMENT,
// and as such can not be expected to be fit for any given purpose.
// Please understand that you use the program at your own risk!!!

func main() {
	ctx := webmind.ParseArgsToContext()
	ctx = webmind.SetupLogging(ctx)

	ctx, err := webmind.BuildPublicAddress(ctx)
	if err != nil {
		log.Printf("error retrieving public address: %v", err)
	}

	webmind.CreateAndRetrievePeerList(ctx)

	webmind.SetupExitHandler(ctx)
	webmind.HandleRequests(fmt.Sprintf("%s", ctx.Value("port")))
}
