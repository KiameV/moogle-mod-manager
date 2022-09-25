package config

import (
	"errors"
	"fyne.io/fyne/v2"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
	idChronoCross    = "1133760"
	// TODO BoF
)

type DirKind byte

const (
	ModsDirKind DirKind = iota
	DownloadDirKind
	BackupDirKind
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
	DirChrCrs   string     `json:"chrCrs"`
	DirBofIII   string     `json:"bof3"`
	DirBofIV    string     `json:"bof4"`
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

func (c *Configs) Size() fyne.Size {
	size := fyne.NewSize(WindowWidth, WindowHeight)
	if x := configs.WindowX; x != 0 {
		size.Width = float32(x)
	}
	if y := configs.WindowY; y != 0 {
		size.Height = float32(y)
	}
	return size
}

func (c *Configs) GetModsFullPath(game Game) string {
	return filepath.Join(c.ModsDir, c.GetGameDirSuffix(game))
}

func (c *Configs) GetDownloadFullPathForUtility() string {
	return filepath.Join(c.DownloadDir, "utility")
}

func (c *Configs) GetDownloadFullPathForGame(game Game) string {
	return filepath.Join(c.DownloadDir, c.GetGameDirSuffix(game))
}

func (c *Configs) GetBackupFullPath(game Game) string {
	return filepath.Join(c.BackupDir, c.GetGameDirSuffix(game))
}

func (c *Configs) AddDir(game Game, dirKind DirKind, from string) (string, error) {
	dir, err := c.GetDir(game, dirKind)
	if err != nil {
		return "", err
	}
	dir = strings.ReplaceAll(dir, "\\", "/")
	from = strings.ReplaceAll(from, "\\", "/")
	if strings.HasPrefix(from, dir) {
		return from, nil
	}
	return filepath.Join(dir, from), nil
}

func (c *Configs) GetDir(game Game, dirKind DirKind) (dir string, err error) {
	switch dirKind {
	case ModsDirKind:
		dir = c.GetModsFullPath(game)
	case DownloadDirKind:
		dir = c.GetDownloadFullPathForGame(game)
	case BackupDirKind:
		dir = c.GetBackupFullPath(game)
	default:
		err = errors.New("unknown dir kind")
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
	case ChronoCross:
		s = c.DirChrCrs
	case BofIII:
		s = c.DirBofIII
	case BofIV:
		s = c.DirBofIV
	}
	return
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
	case ChronoCross:
		s = "chronocross"
	case BofIII:
		s = "bofIII"
	case BofIV:
		s = "bofIV"
	case Utility:
		s = "utility"
	}
	return
}

func (c *Configs) RemoveGameDir(game Game, to string) string {
	dir := c.GetGameDir(game)
	dir = strings.ReplaceAll(dir, "\\", "/")
	to = strings.ReplaceAll(to, "\\", "/")
	to = strings.TrimPrefix(to, dir)
	return strings.TrimPrefix(to, "/")
}

func (c *Configs) Initialize() (err error) {
	if PWD, err = os.Getwd(); err != nil {
		PWD = "."
	}
	if err = util.LoadFromFile(filepath.Join(PWD, configsFile), c); err != nil {
		c.FirstTime = true
		c.Theme = DarkThemeColor
	}
	c.setDefaults()

	if c.DirI == "" {
		c.DirI = c.getGameDirFromRegistry(idI)
	}
	if c.DirII == "" {
		c.DirII = c.getGameDirFromRegistry(idII)
	}
	if c.DirIII == "" {
		c.DirIII = c.getGameDirFromRegistry(idIII)
	}
	if c.DirIV == "" {
		c.DirIV = c.getGameDirFromRegistry(idIV)
	}
	if c.DirV == "" {
		c.DirV = c.getGameDirFromRegistry(idV)
	}
	if c.DirVI == "" {
		c.DirVI = c.getGameDirFromRegistry(idVI)
	}
	if c.DirChrCrs == "" {
		c.DirChrCrs = c.getGameDirFromRegistry(idChronoCross)
	}
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
