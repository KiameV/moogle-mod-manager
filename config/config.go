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

	windowsRegLookup = "Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\Steam App "
	idI              = "1173770"
	idII             = "1173780"
	idIII            = "1173790"
	idIV             = "1173800"
	idV              = "1173810"
	idVI             = "1173820"
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
		c.DirI = c.getGameDirFromRegistry(idI)
		c.DirII = c.getGameDirFromRegistry(idII)
		c.DirIII = c.getGameDirFromRegistry(idIII)
		c.DirIV = c.getGameDirFromRegistry(idIV)
		c.DirV = c.getGameDirFromRegistry(idV)
		c.DirVI = c.getGameDirFromRegistry(idVI)
		c.Theme = DarkThemeColor
	}
	c.setDefaults()
	return nil
}

func (c *Configs) getGameDirFromRegistry(gameId string) (dir string) {
	//only poke into registry for Windows, there's probably a similar method for Mac/Linux
	if runtime.GOOS == "windows" {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, windowsRegLookup+gameId, registry.QUERY_VALUE)
		if err != nil {
			return
		}
		if dir, _, err = key.GetStringValue("InstallLocation"); err != nil {
			dir = ""
		}
	}
	return

}

func (c *Configs) Save() (err error) {
	var (
		b []byte
		f *os.File
	)
	c.setDefaults()
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

func (c *Configs) setDefaults() {
	if c.ModsDir == "" {
		c.ModsDir = filepath.Join(PWD, "mods")
	}
	if c.ImgCacheDir == "" {
		c.ImgCacheDir = filepath.Join(PWD, "imgCache")
	}
	if c.DownloadDir == "" {
		c.DownloadDir = filepath.Join(PWD, "downloads")
	}
	if c.BackupDir == "" {
		c.BackupDir = filepath.Join(PWD, "backups")
	}
}
