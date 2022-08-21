package mods

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/util"
	"path/filepath"
)

const moogleModName = "mod.moogle"

func NewTrackerMod(mod *Mod, game config.Game) *TrackedMod {
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
	Mod         *Mod   `json:"-"`
	UpdatedMod  *Mod   `json:"-"`
	DisplayName string `json:"-"`
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
	if k.Kind == Hosted {
		return filepath.Join(util.CreateFileName(m.GetModID()), util.CreateFileName(m.Mod.Version))
	}
	return filepath.Join(util.CreateFileName(m.GetModID()))
}

func (m *TrackedMod) GetMod() *Mod {
	return m.Mod
}

func (m *TrackedMod) GetBranchName() string {
	return fmt.Sprintf("%s_%s", m.Mod.ID, m.Mod.Version)
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
