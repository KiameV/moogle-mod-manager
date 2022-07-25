package model

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/util"
	"path/filepath"
)

const moogleModName = "mod.moogle"

func NewTrackerMod(mod *mods.Mod, game config.Game) *TrackedMod {
	tm := &TrackedMod{
		Enabled: false,
		Mod:     mod,
	}
	tm.MoogleModFile = filepath.Join(config.Get().GetModsFullPath(game), tm.GetDirSuffix(), moogleModName)
	return tm
}

type TrackedMod struct {
	Enabled       bool      `json:"Enabled"`
	MoogleModFile string    `json:"MoogleModFile"`
	Mod           *mods.Mod `json:"-"`
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

func (m *TrackedMod) GetDirSuffix() string {
	return filepath.Join(util.CreateFileName(m.GetModID()), util.CreateFileName(m.Mod.Version))
}

func (m *TrackedMod) GetMod() *mods.Mod {
	return m.Mod
}
