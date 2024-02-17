package main

import (
	"flag"

	"github.com/trueaniki/gopeers"
	"github.com/trueaniki/shared-clipboard"
	"golang.design/x/clipboard"
)

var (
	network  string
	confPath string
)

func main() {
	flag.StringVar(&network, "network", "", "network to scan in CIDR format")
	flag.StringVar(&network, "net", "", "network to scan in CIDR format")

	flag.StringVar(&confPath, "conf", "", "path to config file")
	flag.StringVar(&confPath, "c", "", "path to config file")

	flag.Parse()
	if network == "" {
		panic("network flag is required")
	}

	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	locals := gopeers.PingSweep(network)
	peer := gopeers.NewPeer(locals)
	peer.Start()

	sharedclipboard.Start(peer)
}
