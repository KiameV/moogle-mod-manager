package config

import (
	"errors"
	"fyne.io/fyne/v2"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
	"strings"
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

	//idChronoCross    = "1133760"
	// TODO BoF
)

type DirKind byte

const (
	ModsDirKind DirKind = iota
	DownloadDirKind
	BackupDirKind
	GameDirKind
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

func (c *Configs) GetModsFullPath(game GameDef) string {
	return filepath.Join(c.ModsDir, string(game.ID))
}

func (c *Configs) GetDownloadFullPathForUtility() string {
	return filepath.Join(c.DownloadDir, "utility")
}

func (c *Configs) GetDownloadFullPathForGame(game GameDef) string {
	return filepath.Join(c.DownloadDir, string(game.ID))
}

func (c *Configs) GetBackupFullPath(game GameDef) string {
	return filepath.Join(c.BackupDir, string(game.ID))
}

func (c *Configs) AddDir(game GameDef, dirKind DirKind, from string) (string, error) {
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

func (c *Configs) GetDir(game GameDef, dirKind DirKind) (dir string, err error) {
	switch dirKind {
	case ModsDirKind:
		dir = c.GetModsFullPath(game)
	case DownloadDirKind:
		dir = c.GetDownloadFullPathForGame(game)
	case BackupDirKind:
		dir = c.GetBackupFullPath(game)
	case GameDirKind:
		dir = game.InstallDir
	default:
		err = errors.New("unknown dir kind")
	}
	return
}

func (c *Configs) RemoveDir(game GameDef, dirKind DirKind, from string) (string, error) {
	dir, err := c.GetDir(game, dirKind)
	if err != nil {
		return "", err
	}
	dir = strings.ReplaceAll(dir, "\\", "/")
	from = strings.ReplaceAll(from, "\\", "/")
	from = strings.TrimPrefix(from, dir)
	return strings.TrimPrefix(from, "/"), nil
}

func (c *Configs) RemoveGameDir(game GameDef, to string) string {
	dir := game.InstallDir
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
	return nil
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
