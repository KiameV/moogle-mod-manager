package model

import (
	"fyne.io/fyne/v2/data/binding"
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
	Enabled       bool   `json:"Enabled"`
	MoogleModFile string `json:"MoogleModFile"`
	//Installed     []*InstalledDownload `json:"Installed"`
	Mod         *mods.Mod      `json:"-"`
	NameBinding binding.String `json:"-"`
	UpdatedMod  *mods.Mod      `json:"-"`
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
	k := m.Mod.ModKind
	if k.Kind == mods.Hosted {
		return filepath.Join(util.CreateFileName(m.GetModID()), util.CreateFileName(k.Hosted.Version))
	}
	return filepath.Join(util.CreateFileName(m.GetModID()))
}

func (m *TrackedMod) GetMod() *mods.Mod {
	return m.Mod
}

type InstalledDownload struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
}

func NewInstalledDownload(name, version string) *InstalledDownload {
	return &InstalledDownload{
		Name:    name,
		Version: version,
	}
}
