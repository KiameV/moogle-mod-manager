//go:build !windows

package config

func (g *gameDef) SteamDirFromRegistry() string {
	return ""
}
