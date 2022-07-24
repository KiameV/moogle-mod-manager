package config

import (
	"encoding/json"
	"os"
	"path"
)

const file = "modsync.config"

var config ConfigData

const (
	WindowWidth  = 820
	WindowHeight = 800
)

var (
	PWD string
)

type ConfigData struct {
	WindowX   int    `json:"width"`
	WindowY   int    `json:"height"`
	DirI      string `json:"dir1"`
	DirII     string `json:"dir2"`
	DirIII    string `json:"dir3"`
	DirIV     string `json:"dir4"`
	DirV      string `json:"dir5"`
	DirVI     string `json:"dir6"`
	ModDir    string `json:"mod-dir"`
	BackupDir string `json:"backup-dir"`
}

func Get() *ConfigData {
	return &config
}

func init() {
	var (
		b   []byte
		err error
	)
	if PWD, err = os.Getwd(); err != nil {
		PWD = "."
	}
	if b, err = os.ReadFile(path.Join(PWD, file)); err == nil {
		_ = json.Unmarshal(b, &config)
	}
}

func Save() {
	if f, e1 := os.Create(path.Join(PWD, file)); e1 == nil {
		if config.WindowX == 0 {
			config.WindowX = WindowWidth
		}
		if config.WindowY == 0 {
			config.WindowY = WindowHeight
		}
		b, err := json.Marshal(&config)
		if err == nil {
			os.WriteFile(path.Join(PWD, file), b, 644)
		}
		_, _ = f.Write(b)
	}
}

func (c *ConfigData) SetGameDir(dir string, game Game) {
	switch game {
	case I:
		Get().DirI = dir
	case II:
		Get().DirII = dir
	case III:
		Get().DirIII = dir
	case IV:
		Get().DirIV = dir
	case V:
		Get().DirV = dir
	case VI:
		Get().DirVI = dir
	}
}

func GetModDir(game Game) (dir string) {
	switch game {
	case I:
		dir = Get().DirI
		if dir == "" {
			dir = "I"
		}
	case II:
		dir = Get().DirII
		if dir == "" {
			dir = "II"
		}
	case III:
		dir = Get().DirIII
		if dir == "" {
			dir = "III"
		}
	case IV:
		dir = Get().DirIV
		if dir == "" {
			dir = "IV"
		}
	case V:
		dir = Get().DirV
		if dir == "" {
			dir = "V"
		}
	case VI:
		dir = Get().DirVI
		if dir == "" {
			dir = "VI"
		}
	}
	dir = path.Join(PWD, "mods", dir)
	return
}

func GetBackupDir(game Game) (s string) {
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
	s = path.Join(PWD, "backup", s)
	return
}
