package config

import (
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/windows/registry"
)

const configsFile = "configs.json"

var (
	PWD     string
	configs = &Configs{}
)

const (
	WindowWidth  = 1200
	WindowHeight = 850

	windowsRegLookup = "Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\Steam App "
	idI              = "1173770"
	idII             = "1173780"
	idIII            = "1173790"
	idIV             = "1173800"
	idV              = "1173810"
	idVI             = "1173820"
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
	DefaultGame string     `json:"openTo"`
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

func (c *Configs) GetGameDir(game Game) (s string) {
	switch game {
	case I:
		s = c.DirI
	case II:
		s = c.DirII
	case III:
		s = c.DirIII
	case IV:
		s = c.DirIV
	case V:
		s = c.DirV
	case VI:
		s = c.DirVI
	}
	return
}

func (c *Configs) Initialize() (err error) {
	if PWD, err = os.Getwd(); err != nil {
		PWD = "."
	}
	if err = util.LoadFromFile(filepath.Join(PWD, configsFile), c); err != nil {
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
	c.setDefaults()
	return util.SaveToFile(filepath.Join(PWD, configsFile), c)
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
