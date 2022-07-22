package model

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"path/filepath"
)

func NewTrackerMod(game config.Game, mod *mods.Mod) *TrackerMod {
	return &TrackerMod{
		Enabled: false,
		Mod:     mod,
		Dir:     filepath.Join(config.GetModDir(game), mod.ID),
	}
}

type TrackerMod struct {
	Enabled bool      `json:"Enabled"`
	Dir     string    `json:"Dir"`
	Mod     *mods.Mod `json:"-"`
}

func (m *TrackerMod) IsEnabled() bool {
	return m.Enabled
}

func (m *TrackerMod) SetIsEnabled(isEnabled bool) {
	m.Enabled = isEnabled
}

func (m *TrackerMod) Toggle() bool {
	m.Enabled = !m.Enabled
	return m.Enabled
}

func (m *TrackerMod) GetModID() string {
	return m.Mod.ID
}

func (m *TrackerMod) GetDir() string {
	return m.Dir
}

func (m *TrackerMod) GetMod() *mods.Mod {
	return m.Mod
}
