package main

import (
	"log"

	"golang.design/x/clipboard"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func main() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
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
			dumpClipboard()
		case <-hkLoad.Keydown():
			log.Printf("hotkey: %v is down\n", hkLoad)
			loadClipboard()
		}
	}
}

var cb []byte

func dumpClipboard() {
	cb = clipboard.Read(clipboard.FmtText)
}

func loadClipboard() {
	clipboard.Write(clipboard.FmtText, cb)
}
