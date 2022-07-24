package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
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

type Configs struct {
	FirstTime   bool   `json:"firstTime"`
	WindowX     int    `json:"width"`
	WindowY     int    `json:"height"`
	DirI        string `json:"dir1"`
	DirII       string `json:"dir2"`
	DirIII      string `json:"dir3"`
	DirIV       string `json:"dir4"`
	DirV        string `json:"dir5"`
	DirVI       string `json:"dir6"`
	DownloadDir string `json:"downloadDir"`
	BackupDir   string `json:"backupDir"`
}

func Get() *Configs {
	return configs
}

func (c *Configs) SetGameDir(dir string, game Game) {
	switch game {
	case I:
		c.DirI = dir
	case II:
		c.DirII = dir
	case III:
		c.DirIII = dir
	case IV:
		c.DirIV = dir
	case V:
		c.DirV = dir
	case VI:
		c.DirVI = dir
	}
}

func (c *Configs) GetModDir(game Game) (dir string) {
	switch game {
	case I:
		dir = c.DirI
		if dir == "" {
			dir = "I"
		}
	case II:
		dir = c.DirII
		if dir == "" {
			dir = "II"
		}
	case III:
		dir = c.DirIII
		if dir == "" {
			dir = "III"
		}
	case IV:
		dir = c.DirIV
		if dir == "" {
			dir = "IV"
		}
	case V:
		dir = c.DirV
		if dir == "" {
			dir = "V"
		}
	case VI:
		dir = c.DirVI
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

func (c *Configs) Initialize() (err error) {
	var b []byte
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
