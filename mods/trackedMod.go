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

func (m *TrackedMod) GetModID() ModID {
	return m.Mod.ID
}

func (m *TrackedMod) GetDirSuffix() string {
	switch m.Mod.ModKind.Kind {
	case Hosted:
		return filepath.Join(util.CreateFileName(string(m.GetModID())), util.CreateFileName(m.Mod.Version))
	case Nexus:
		return filepath.Join("nexus", util.CreateFileName(string(m.GetModID())))
	case CurseForge:
		return filepath.Join("cf", util.CreateFileName(string(m.GetModID())))
	}
	panic(fmt.Sprintf("unknown kind %v", m.Mod.ModKind.Kind))
}

func (m *TrackedMod) GetMod() *Mod {
	return m.Mod
}

func (m *TrackedMod) Save() error {
	return util.SaveToFile(m.MoogleModFile, m.Mod)
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
