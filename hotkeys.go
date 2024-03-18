package sharedclipboard

import (
	"errors"
	"strings"

	parsehotkeys "github.com/trueaniki/go-parse-hotkeys"
	"golang.design/x/hotkey"
)

type Hotkeys struct {
	HKShare *hotkey.Hotkey
	HKAdopt *hotkey.Hotkey
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
		if strings.HasPrefix(strings.ToLower(strings.Trim(def, " ")), "share") {
			hks := strings.Split(def, "=")[1]
			hkshare, err := parsehotkeys.Parse(hks, "+")
			if err != nil {
				return nil, err
			}
			hk.HKShare = hkshare
		}
		if strings.HasPrefix(strings.ToLower(strings.Trim(def, " ")), "adopt") {
			hks := strings.Split(def, "=")[1]
			hkadopt, err := parsehotkeys.Parse(hks, "+")
			if err != nil {
				return nil, err
			}
			hk.HKAdopt = hkadopt
		}
	}

	if hk.HKShare == nil {
		return nil, errors.New("hotkey 'Share' not found")
	}
	if hk.HKAdopt == nil {
		return nil, errors.New("hotkey 'Adopt' not found")
	}

	return hk, nil
}
