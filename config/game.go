package config

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/kiamev/moogle-mod-manager/util"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var gameDefs []GameDef

type (
	InstallType string
	GameID      string
	GameName    string
	SteamID     string
	BaseDir     string
	NexusGameID int
	NexusPath   string
	CfGameID    int
	CfPath      string
	VersionID   string
	Nexus       struct {
		ID   NexusGameID `json:"id"`
		Path NexusPath   `json:"path"`
	}
	CurseForge struct {
		ID   CfGameID `json:"id"`
		Path CfPath   `json:"path"`
	}
	Remote struct {
		Nexus      Nexus      `json:"nexus"`
		CurseForge CurseForge `json:"curseforge"`
	}
	SteamVersion struct {
		Build    uint   `json:"build"`
		Manifest uint64 `json:"manifest"`
	}
	Version struct {
		Version VersionID     `json:"version"`
		Steam   *SteamVersion `json:"steam,omitempty"`
	}
	gameDef struct {
		ID_                 GameID            `json:"id"`
		Name_               GameName          `json:"name"`
		SteamID_            SteamID           `json:"steamID"`
		Versions_           []Version         `json:"versions"`
		BaseDir_            BaseDir           `json:"baseDir"`
		Remote_             Remote            `json:"remote"`
		DefaultInstallType_ InstallType       `json:"defaultInstallType"`
		LogoPath_           string            `json:"-"`
		Logo_               fyne.CanvasObject `json:"-"`
		InstallDir_         string            `json:"-"`
	}
	GameDef interface {
		ID() GameID
		Name() GameName
		SteamID() SteamID
		Versions() []Version
		BaseDir() BaseDir
		Remote() Remote
		DefaultInstallType() InstallType
		LogoPath() string
		SetLogoPath(path string)
		Logo() fyne.CanvasObject
		SetLogo(logo fyne.CanvasObject)
		SteamDirFromRegistry() string
	}
)

func (g *gameDef) ID() GameID {
	return g.ID_
}

func (g *gameDef) Name() GameName {
	return g.Name_
}

func (g *gameDef) SteamID() SteamID {
	return g.SteamID_
}

func (g *gameDef) Versions() []Version {
	return g.Versions_
}

func (g *gameDef) BaseDir() BaseDir {
	return g.BaseDir_
}

func (g *gameDef) Remote() Remote {
	return g.Remote_
}

func (g *gameDef) DefaultInstallType() InstallType {
	return g.DefaultInstallType_
}

func (g *gameDef) LogoPath() string {
	return g.LogoPath_
}

func (g *gameDef) SetLogoPath(path string) {
	g.LogoPath_ = path
}

func (g *gameDef) Logo() fyne.CanvasObject {
	return g.Logo_
}

func (g *gameDef) InstallDir() string {
	return g.InstallDir_
}

func (g *gameDef) InstallDirPtr() *string {
	return &g.InstallDir_
}

func (g *gameDef) SetInstallDir(dir string) {
	g.InstallDir_ = dir
}

const (
	StreamingAssetsDir = "StreamingAssets"

	//Bundles  InstallType = "Bundles"
	// DLL Patcher https://discord.com/channels/371784427162042368/518331294858608650/863930606446182420
	//DllPatch   InstallType = "DllPatch"
	MoveToArchive       InstallType = "MoveToArchive"
	Move                InstallType = "Move"
	ImmediateDecompress InstallType = "ImmediateDecompress"
)

func GameDefs() []GameDef {
	return gameDefs
}

func GameDefFromID(id GameID) (GameDef, error) {
	for _, g := range gameDefs {
		if g.ID() == id {
			return g, nil
		}
	}
	return nil, fmt.Errorf("game with ModID [%s] not found", id)
}

func GameDefFromNexusID(id NexusGameID) (GameDef, error) {
	for _, g := range gameDefs {
		if g.Remote().Nexus.ID == id {
			return g, nil
		}
	}
	return nil, fmt.Errorf("game with Nexus ModID [%d] not found", id)
}

func GameDefFromNexusPath(path NexusPath) (GameDef, error) {
	for _, g := range gameDefs {
		if g.Remote().Nexus.Path == path {
			return g, nil
		}
	}
	return nil, fmt.Errorf("game with Nexus Path [%s] not found", path)
}

func GameDefFromCfID(id CfGameID) (GameDef, error) {
	for _, g := range gameDefs {
		if g.Remote().CurseForge.ID == id {
			return g, nil
		}
	}
	return nil, fmt.Errorf("game with CurseForge ModID [%d] not found", id)
}

func GameDefFromCfPath(path CfPath) (GameDef, error) {
	for _, g := range gameDefs {
		if g.Remote().CurseForge.Path == path {
			return g, nil
		}
	}
	return nil, fmt.Errorf("game with CurseForge path [%s] not found", path)
}

func GameDefFromName(name GameName) (GameDef, error) {
	for _, g := range gameDefs {
		if g.Name() == name {
			return g, nil
		}
	}
	return nil, fmt.Errorf("game with name [%s] not found", name)
}

func Initialize(dirs []string) (err error) {
	var des []os.DirEntry
	for _, dir := range dirs {
		if err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if d.IsDir() {
				if err != nil {
					return err
				}
				if strings.Contains(path, ".git") {
					return nil
				}
				if des, err = os.ReadDir(path); err != nil {
					return err
				}
				var (
					logo string
					def  string
				)
				for _, de := range des {
					if de.Name() == "game.json" {
						def = filepath.Join(path, de.Name())
					} else if strings.HasPrefix(de.Name(), "logo.") {
						logo = filepath.Join(path, de.Name())
					}
					if def != "" && logo != "" {
						break
					}
				}
				if def != "" {
					var game gameDef
					if err = util.LoadFromFile(def, &game); err != nil {
						return err
					}
					game.LogoPath_ = logo
					gameDefs = append(gameDefs, &game)
				}
			}
			return nil
		}); err != nil {
			return
		}
	}
	return
}

func (g *gameDef) SteamDirFromRegistry() (dir string) {
	//only poke into registry for Windows, there's probably a similar method for Mac/Linux
	if runtime.GOOS == "windows" {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, fmt.Sprintf("%s%s", windowsRegLookup, g.SteamID_), registry.QUERY_VALUE)
		if err != nil {
			return
		}
		if dir, _, err = key.GetStringValue("InstallLocation"); err != nil {
			dir = ""
		}
	}
	return
}

func (g *gameDef) SetLogo(logo fyne.CanvasObject) {
	g.Logo_ = logo
}

func GameIDs() []string {
	names := make([]string, len(gameDefs)+1)
	names[0] = ""
	for i, g := range gameDefs {
		names[i+1] = string(g.ID())
	}
	return names
}

func (t *InstallType) Is(i InstallType) bool {
	return t != nil && *t == i
}
