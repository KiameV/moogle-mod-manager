package mods

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/util"
	"path/filepath"
)

const moogleModName = "mod.moogle"

type (
	TrackedMod interface {
		ID() ModID
		Kind() Kind
		Mod() *Mod
		Enable()
		Enabled() bool
		Disable()
		Toggle() bool
		DirSuffix() string
		Save() error
	}
	trackedMod struct {
		IsEnabled     bool   `json:"Enabled"`
		MoogleModFile string `json:"MoogleModFile"`
		//Installed     []*InstalledDownload `json:"Installed"`
		m           *Mod   `json:"-"`
		UpdatedMod  *Mod   `json:"-"`
		DisplayName string `json:"-"`
	}
)

func (m *trackedMod) Enable() {
	m.IsEnabled = true
}

func (m *trackedMod) Disable() {
	m.IsEnabled = false
}

func NewTrackerMod(mod *Mod, game config.GameDef) TrackedMod {
	tm := &trackedMod{
		IsEnabled: false,
		m:         mod,
	}
	tm.MoogleModFile = filepath.Join(config.Get().GetModsFullPath(game), tm.DirSuffix(), moogleModName)
	return tm
}

func (m *trackedMod) ID() ModID {
	return m.ID()
}

func (m *trackedMod) Kind() Kind {
	return m.Kind()
}

func (m *trackedMod) Mod() *Mod {
	return m.Mod()
}

func (m *trackedMod) Enabled() bool {
	return m.IsEnabled
}

func (m *trackedMod) Toggle() bool {
	m.IsEnabled = !m.IsEnabled
	return m.IsEnabled
}

func (m *trackedMod) DirSuffix() string {
	switch m.Kind() {
	case Hosted:
		return filepath.Join(util.CreateFileName(string(m.ID())), util.CreateFileName(m.Mod().Version))
	case Nexus:
		return filepath.Join(util.CreateFileName(string(m.ID())))
	case CurseForge:
		return filepath.Join(util.CreateFileName(string(m.ID())))
	}
	panic(fmt.Sprintf("unknown kind %v", m.Kind()))
}

func (m *trackedMod) Save() error {
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
