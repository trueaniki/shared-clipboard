package main

import (
	sharedclipboard "shared-clipboard"

	"github.com/trueaniki/gopeers"
	"golang.design/x/clipboard"
)

func main() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	locals := gopeers.PingSweep("192.168.100.0/24")
	peer := gopeers.NewPeer(locals)
	peer.Start()

	sharedclipboard.Start(peer)
}
