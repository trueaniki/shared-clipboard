package main

import (
	"log"

	"github.com/trueaniki/gopeers"
	"golang.design/x/clipboard"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

var (
	peer *gopeers.Peer
)

func main() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	locals := gopeers.PingSweep("192.168.100.0/24")
	peer = gopeers.NewPeer(locals)
	peer.Start()

	mainthread.Init(listenHotkeys)
}

func listenHotkeys() {
	hkDump := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyA)
	hkLoad := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyD)

	err := hkDump.Register()
	if err != nil {
		panic(err)
	}

	err = hkLoad.Register()
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-hkDump.Keydown():
			log.Printf("hotkey: %v is down\n", hkDump)
			peer.WriteChan <- clipboard.Read(clipboard.FmtText)
		case <-hkLoad.Keydown():
			log.Printf("hotkey: %v is down\n", hkLoad)
			clipboard.Write(clipboard.FmtText, <-peer.ReadChan)
		}
	}
}
