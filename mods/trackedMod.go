package mods

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"path/filepath"
	"strings"
)

const moogleModName = "mod.moogle"

type (
	TrackedMod interface {
		ID() ModID
		Kind() Kind
		SubKind() SubKind
		Mod() *Mod
		SetMod(m *Mod)
		Enable()
		Enabled() bool
		EnabledPtr() *bool
		Disable()
		Save() error
		DisplayName() string
		DisplayNamePtr() *string
		SetDisplayName(name string)
		UpdatedMod() *Mod
		SetUpdatedMod(m *Mod)
		MoogleModFile() string
		InstallType(game config.GameDef) config.InstallType
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

func (m *TrackedModConc) InstallType(game config.GameDef) config.InstallType {
	return m.Mod().InstallType(game)
}

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

func (m *TrackedModConc) EnabledPtr() *bool {
	return &m.IsEnabled
}

func (m *TrackedModConc) Disable() {
	m.IsEnabled = false
}

func NewTrackerMod(mod *Mod, game config.GameDef) TrackedMod {
	tm := &TrackedModConc{
		IsEnabled: false,
		Mod_:      mod,
	}
	tm.MoogleModFile_ = filepath.Join(config.Get().GetModsFullPath(game), tm.ID().AsDir(), moogleModName)
	return tm
}

func (id ModID) AsDir() string {
	return strings.ReplaceAll(string(id), ".", "_")
}

func (m *TrackedModConc) ID() ModID {
	return m.Mod_.ID()
}

func (m *TrackedModConc) Kind() Kind {
	return m.Mod_.Kind()
}

func (m *TrackedModConc) SubKind() SubKind {
	return m.Mod_.SubKind()
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
