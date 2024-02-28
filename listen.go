package sharedclipboard

import (
	"log"

	"github.com/trueaniki/gopeers"
	"golang.design/x/clipboard"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func listen(peer *gopeers.Peer) func() {
	return func() {
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

		var dump []byte
		go func() {
			for msg := range peer.ReadChan {
				dump = msg
			}
		}()

		for {
			select {
			case <-hkDump.Keydown():
				log.Printf("hotkey: %v is down\n", hkDump)
				peer.WriteChan <- clipboard.Read(clipboard.FmtText)
			case <-hkLoad.Keydown():
				log.Printf("hotkey: %v is down\n", hkLoad)
				if dump != nil {
					clipboard.Write(clipboard.FmtText, dump)
				}
			}
		}
	}
}

func Listen(peer *gopeers.Peer) {
	mainthread.Init(listen(peer))
}
