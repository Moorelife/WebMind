package oneness

import (
	"github.com/Moorelife/WebMind/internal/localnode"
	"github.com/Moorelife/WebMind/internal/trinity"
)

type Oneness struct {
	Monitor      localnode.LocalNode //`json:"monitor`
	LocalTrinity trinity.Trinity     //`json:"local_trinity"`
}

func NewOneness(monitor localnode.LocalNode, trio []int) *Oneness {
	var creation = Oneness{
		Monitor: monitor,
	}
	for key, value := range trio {
		if value == 0 {
			trio[key] = *monitor.LocalPort + 1 + key
		}
	}
	creation.LocalTrinity = *trinity.NewTrinity(monitor.LocalAddress, trio)

	return &creation
}
