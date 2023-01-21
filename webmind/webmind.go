package main

import (
	"fmt"
	"log"

	"github.com/Moorelife/WebMind/internal/peerlist"
	"github.com/Moorelife/WebMind/internal/webmind"
)

func main() {
	ctx := webmind.ParseArgsToContext()

	ctx = webmind.SetupLogging(ctx)

	ctx, err := webmind.BuildPublicAddress(ctx)
	if err != nil {
		log.Printf("error retrieving public address: %v", err)
	}

	webmind.SetupExitHandler(ctx)

	log.Println(fmt.Sprintf("%s", ctx.Value("selfAddress")))
	peerlist.Add(fmt.Sprintf("%s", ctx.Value("selfAddress")))
	if fmt.Sprintf("%s", ctx.Value("origin")) != "" {
		peerlist.Get(fmt.Sprintf("%s", ctx.Value("origin")))
	}

	webmind.HandleRequests(fmt.Sprintf("%s", ctx.Value("port")))
}
