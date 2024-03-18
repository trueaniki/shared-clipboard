package sharedclipboard

import (
	"log"

	"github.com/trueaniki/gopeers"
	"golang.design/x/clipboard"
	"golang.design/x/hotkey/mainthread"
)

func listen(peer *gopeers.Peer, hk *Hotkeys) func() {
	hkShare := hk.HKShare
	hkAdopt := hk.HKAdopt

	return func() {
		err := hkShare.Register()
		if err != nil {
			panic(err)
		}

		err = hkAdopt.Register()
		if err != nil {
			panic(err)
		}

		var dump []byte
		go func() {
			for msg := range peer.ReadChan {
				dump = msg
			}
		}()

		for {
			select {
			case <-hkShare.Keydown():
				log.Printf("hotkey: %v is down\n", hkShare)
				peer.WriteChan <- clipboard.Read(clipboard.FmtText)
			case <-hkAdopt.Keydown():
				log.Printf("hotkey: %v is down\n", hkAdopt)
				if dump != nil {
					clipboard.Write(clipboard.FmtText, dump)
				}
			}
		}
	}
}

func Listen(peer *gopeers.Peer, hk *Hotkeys) {
	mainthread.Init(listen(peer, hk))
}
