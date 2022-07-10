package state

import "github.com/kiamev/pr-modsync/config"

type GUI byte

const (
	None GUI = iota
	LocalMods
	RemoteMods
	Config
)

var (
	Game      config.Game
	CurrentUI GUI
	Errors    []error
)
