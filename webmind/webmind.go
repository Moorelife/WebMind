package main

import (
	"fmt"
	"github.com/Moorelife/WebMind/internal/webmind"
)

// WebMind in its current state is JUST A LEARNING EXPERIMENT,
// and as such can not be expected to be fit for any given purpose.
// Please understand that you use the program at your own risk!!!

func main() {
	ctx := webmind.ParseArgsToContext()
	ctx = webmind.SetupLogging(ctx)
	ctx = webmind.RetrievePublicAddress(ctx)

	webmind.CreateAndRetrievePeerList(ctx)
	webmind.SendPeerAddRequests(ctx)
	webmind.StartSendingKeepAlive(ctx)
	webmind.SetupExitHandler(ctx)
	webmind.HandleRequests(fmt.Sprintf("%s", ctx.Value("port")))
}
