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
	GameDef struct {
		ID                 GameID            `json:"id"`
		Name               GameName          `json:"name"`
		SteamID            SteamID           `json:"steamID"`
		Versions           []Version         `json:"versions"`
		BaseDir            BaseDir           `json:"baseDir"`
		Remote             Remote            `json:"remote"`
		DefaultInstallType InstallType       `json:"defaultInstallType"`
		LogoPath           string            `json:"-"`
		Logo               fyne.CanvasObject `json:"-"`
		InstallDir         string            `json:"-"`
	}
)

const (
	StreamingAssetsDir = "StreamingAssets"

	//Bundles  InstallType = "Bundles"
	// DLL Patcher https://discord.com/channels/371784427162042368/518331294858608650/863930606446182420
	//DllPatch   InstallType = "DllPatch"
	Archive InstallType = "Archive"
	Move    InstallType = "Move"
)

func GameDefs() []GameDef {
	return gameDefs
}

func GameDefFromID(id GameID) (GameDef, error) {
	for _, g := range gameDefs {
		if g.ID == id {
			return g, nil
		}
	}
	return GameDef{}, fmt.Errorf("game with ModID [%s] not found", id)
}

func GameDefFromNexusID(id NexusGameID) (GameDef, error) {
	for _, g := range gameDefs {
		if g.Remote.Nexus.ID == id {
			return g, nil
		}
	}
	return GameDef{}, fmt.Errorf("game with Nexus ModID [%d] not found", id)
}

func GameDefFromCfID(id CfGameID) (GameDef, error) {
	for _, g := range gameDefs {
		if g.Remote.CurseForge.ID == id {
			return g, nil
		}
	}
	return GameDef{}, fmt.Errorf("game with CurseForge ModID [%d] not found", id)
}

func GameDefFromCfPath(path CfPath) (GameDef, error) {
	for _, g := range gameDefs {
		if g.Remote.CurseForge.Path == path {
			return g, nil
		}
	}
	return GameDef{}, fmt.Errorf("game with CurseForge path [%s] not found", path)
}

func GameDefFromName(name GameName) (GameDef, error) {
	for _, g := range gameDefs {
		if g.Name == name {
			return g, nil
		}
	}
	return GameDef{}, fmt.Errorf("game with name [%s] not found", name)
}

func Initialize(dirs []string) (err error) {
	var (
		des    []os.DirEntry
		dirMap = make(map[GameID]string)
		logo   string
		def    string
	)
	for _, dir := range dirs {
		if err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if d.IsDir() {
				if err != nil {
					return err
				}
				if des, err = os.ReadDir(filepath.Join(path, d.Name())); err != nil {
					return err
				}

				for _, de := range des {
					if de.Name() == "game.json" {
						def = filepath.Join(path, de.Name())
					} else if strings.HasPrefix(de.Name(), "logo.") {
						logo = filepath.Join(path, de.Name())
					}
					if def != "" {
						var game GameDef
						if err = util.LoadFromFile(def, &game); err != nil {
							return err
						}
						game.LogoPath = logo
						gameDefs = append(gameDefs, game)
					}
				}
			}
			return nil
		}); err != nil {
			break
		}
	}

	if err = util.LoadFromFile(filepath.Join(PWD, "gameDirs.json"), &dirMap); err != nil {
		err = nil
	} else {
		for _, game := range gameDefs {
			if dir, ok := dirMap[game.ID]; ok {
				game.InstallDir = dir
			}
		}
	}

	if configs.FirstTime {
		for i, game := range gameDefs {
			if s := game.GetSteamDirFromRegistry(); s != "" {
				gameDefs[i].InstallDir = s
			}
		}
	}
	return
}

func (g *GameDef) GetSteamDirFromRegistry() (dir string) {
	//only poke into registry for Windows, there's probably a similar method for Mac/Linux
	if runtime.GOOS == "windows" {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, fmt.Sprintf("%s%s", windowsRegLookup, g.SteamID), registry.QUERY_VALUE)
		if err != nil {
			return
		}
		if dir, _, err = key.GetStringValue("InstallLocation"); err != nil {
			dir = ""
		}
	}
	return
}

func (g *GameDef) SetLogo(logo fyne.CanvasObject) {
	g.Logo = logo
}
