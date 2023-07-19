//go:build windows

package config

import (
	"fmt"
	"runtime"

	"golang.org/x/sys/windows/registry"
)

const (
	windowsRegLookup = "Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\Steam App "
)

func (g *gameDef) SteamDirFromRegistry() (dir string) {
	// only poke into registry for Windows, there's probably a similar method for Mac/Linux
	if runtime.GOOS == "windows" {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, fmt.Sprintf("%s%s", windowsRegLookup, g.SteamID_), registry.QUERY_VALUE)
		if err != nil {
			return
		}
		if dir, _, err = key.GetStringValue("InstallLocation"); err != nil {
			dir = ""
		}
	}
	return
}
