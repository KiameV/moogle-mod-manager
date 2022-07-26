package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/windows/registry"
)

const configsFile = "configs.json"

var configs = &Configs{}

const (
	WindowWidth  = 1000
	WindowHeight = 850
)

var (
	PWD string
)

type ThemeColor byte

const (
	DarkThemeColor ThemeColor = iota
	LightThemeColor
)

type Configs struct {
	FirstTime   bool       `json:"firstTime"`
	WindowX     int        `json:"width"`
	WindowY     int        `json:"height"`
	DirI        string     `json:"dir1"`
	DirII       string     `json:"dir2"`
	DirIII      string     `json:"dir3"`
	DirIV       string     `json:"dir4"`
	DirV        string     `json:"dir5"`
	DirVI       string     `json:"dir6"`
	ModsDir     string     `json:"modDir"`
	ImgCacheDir string     `json:"imgCacheDir"`
	DownloadDir string     `json:"downloadDir"`
	BackupDir   string     `json:"backupDir"`
	Theme       ThemeColor `json:"theme"`
}

func Get() *Configs {
	return configs
}

func (c *Configs) GetModsFullPath(game Game) string {
	return filepath.Join(c.ModsDir, c.GetGameDirSuffix(game))
}

func (c *Configs) GetDownloadFullPath(game Game) string {
	return filepath.Join(c.DownloadDir, c.GetGameDirSuffix(game))
}

func (c *Configs) GetBackupFullPath(game Game) string {
	return filepath.Join(c.BackupDir, c.GetGameDirSuffix(game))
}

func (c *Configs) GetGameDirSuffix(game Game) (s string) {
	switch game {
	case I:
		s = "I"
	case II:
		s = "II"
	case III:
		s = "III"
	case IV:
		s = "IV"
	case V:
		s = "V"
	case VI:
		s = "VI"
	}
	return
}

func TryGetGameDirFromRegistry(gameId string) string {
	path := ""                     //initialize the path
	if runtime.GOOS == "windows" { //only poke into registry for Windows, there's probably a similar method for Mac/Linux
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, "Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\Steam App "+gameId, registry.QUERY_VALUE)
		if err != nil {
			return path //since path is initialized to "", we just assume the user does not own this game
		}
		path, valtype, err := key.GetStringValue("InstallLocation")
		if err != nil {
			fmt.Errorf("Failed to get path for Game ID %v: ValType: %q Err: %v", gameId, valtype, err)
		}
		return path
	}
	return path

}

func (c *Configs) Initialize() (err error) {
	var b []byte
	if PWD, err = os.Getwd(); err != nil {
		PWD = "."
	}

	p := filepath.Join(PWD, configsFile)
	if _, err = os.Stat(p); err == nil {
		if b, err = os.ReadFile(p); err != nil {
			return fmt.Errorf("failed to read configs file: %v", err)
		}
		if err = json.Unmarshal(b, c); err != nil {
			return fmt.Errorf("failed to unmarshal configs file: %v", err)
		}
	} else {
		c.FirstTime = true
		c.DirI = TryGetGameDirFromRegistry("1173770")
		c.DirII = TryGetGameDirFromRegistry("1173780")
		c.DirIII = TryGetGameDirFromRegistry("1173790")
		c.DirIV = TryGetGameDirFromRegistry("1173800")
		c.DirV = TryGetGameDirFromRegistry("1173810")
		c.DirVI = TryGetGameDirFromRegistry("1173820")
		c.ModsDir = filepath.Join(PWD, "mods")
		c.ImgCacheDir = filepath.Join(PWD, "imgCache")
		c.DownloadDir = filepath.Join(PWD, "downloads")
		c.BackupDir = filepath.Join(PWD, "backups")
		c.Theme = DarkThemeColor
	}
	return nil
}

func (c *Configs) Save() (err error) {
	var (
		b []byte
		f *os.File
	)
	if b, err = json.MarshalIndent(c, "", "\t"); err != nil {
		return fmt.Errorf("failed to marshal configs: %v", err)
	}
	if f, err = os.Create(filepath.Join(PWD, configsFile)); err != nil {
		return fmt.Errorf("failed to create configs file: %v", err)
	}
	if _, err = f.Write(b); err != nil {
		return fmt.Errorf("failed to write configs file: %v", err)
	}
	return
}
