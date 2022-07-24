package model

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"path"
)

func NewTrackerMod(game config.Game, mod *mods.Mod) *TrackedMod {
	return &TrackedMod{
		Enabled: false,
		Mod:     mod,
		Dir:     path.Join(config.GetModDir(game), mod.ID),
	}
}

type TrackedMod struct {
	Enabled bool      `json:"Enabled"`
	Dir     string    `json:"Dir"`
	Mod     *mods.Mod `json:"-"`
}

func (m *TrackedMod) IsEnabled() bool {
	return m.Enabled
}

func (m *TrackedMod) SetIsEnabled(isEnabled bool) {
	m.Enabled = isEnabled
}

func (m *TrackedMod) Toggle() bool {
	m.Enabled = !m.Enabled
	return m.Enabled
}

func (m *TrackedMod) GetModID() string {
	return m.Mod.ID
}

func (m *TrackedMod) GetDir() string {
	return m.Dir
}

func (m *TrackedMod) GetMod() *mods.Mod {
	return m.Mod
}
