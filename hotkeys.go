package sharedclipboard

import (
	"errors"
	"strings"

	parsehotkeys "github.com/trueaniki/go-parse-hotkeys"
	"golang.design/x/hotkey"
)

type Hotkeys struct {
	HKDump *hotkey.Hotkey
	HKLoad *hotkey.Hotkey
}

func ParseHotkeys(definition string) (*Hotkeys, error) {
	hk := &Hotkeys{}
	defs := strings.Split(definition, "\n")
	for _, def := range defs {
		if def == "" {
			continue
		}
		if strings.HasPrefix(def, "#") {
			continue
		}
		if strings.HasPrefix(strings.ToLower(strings.Trim(def, " ")), "HKDump") {
			hks := strings.Split(def, "=")[1]
			hkdump, err := parsehotkeys.Parse(hks, "+")
			if err != nil {
				return nil, err
			}
			hk.HKDump = hkdump
		}
		if strings.HasPrefix(strings.ToLower(strings.Trim(def, " ")), "HKLoad") {
			hks := strings.Split(def, "=")[1]
			hkload, err := parsehotkeys.Parse(hks, "+")
			if err != nil {
				return nil, err
			}
			hk.HKLoad = hkload
		}
	}

	if hk.HKDump == nil {
		return nil, errors.New("HKDump not found")
	}
	if hk.HKLoad == nil {
		return nil, errors.New("HKLoad not found")
	}

	return hk, nil
}
