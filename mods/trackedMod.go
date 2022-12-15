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
		SetMod(m *Mod)
		Enable()
		Enabled() bool
		Disable()
		Toggle() bool
		DirSuffix() string
		Save() error
		DisplayName() string
		DisplayNamePtr() *string
		SetDisplayName(name string)
		UpdatedMod() *Mod
		SetUpdatedMod(m *Mod)
		MoogleModFile() string
	}
	// TrackedModConc is public for serialization purposes
	TrackedModConc struct {
		IsEnabled      bool   `json:"Enabled"`
		MoogleModFile_ string `json:"MoogleModFile"`
		//Installed     []*InstalledDownload `json:"Installed"`
		Mod_         *Mod   `json:"-"`
		UpdatedMod_  *Mod   `json:"-"`
		DisplayName_ string `json:"-"`
	}
)

func (m *TrackedModConc) DisplayNamePtr() *string {
	return &m.DisplayName_
}

func (m *TrackedModConc) SetDisplayName(name string) {
	m.DisplayName_ = name
}

func (m *TrackedModConc) SetUpdatedMod(updatedMod *Mod) {
	m.UpdatedMod_ = updatedMod
}

func (m *TrackedModConc) MoogleModFile() string {
	return m.MoogleModFile_
}

func (m *TrackedModConc) UpdatedMod() *Mod {
	return m.UpdatedMod_
}

func (m *TrackedModConc) DisplayName() string {
	return m.DisplayName_
}

func (m *TrackedModConc) Enable() {
	m.IsEnabled = true
}

func (m *TrackedModConc) Disable() {
	m.IsEnabled = false
}

func NewTrackerMod(mod *Mod, game config.GameDef) TrackedMod {
	tm := &TrackedModConc{
		IsEnabled: false,
		Mod_:      mod,
	}
	tm.MoogleModFile_ = filepath.Join(config.Get().GetModsFullPath(game), tm.DirSuffix(), moogleModName)
	return tm
}

func (m *TrackedModConc) ID() ModID {
	return m.Mod_.ID()
}

func (m *TrackedModConc) Kind() Kind {
	return m.Mod_.Kind()
}

func (m *TrackedModConc) Mod() *Mod {
	return m.Mod_
}

func (m *TrackedModConc) SetMod(mod *Mod) {
	m.Mod_ = mod
}

func (m *TrackedModConc) Enabled() bool {
	return m.IsEnabled
}

func (m *TrackedModConc) Toggle() bool {
	m.IsEnabled = !m.IsEnabled
	return m.IsEnabled
}

func (m *TrackedModConc) DirSuffix() string {
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

func (m *TrackedModConc) Save() error {
	return m.Mod_.Save(m.MoogleModFile_)
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
