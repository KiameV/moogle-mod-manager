package authored

import (
	"encoding/json"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"os"
	"path"
)

const file = "authored.json"

var lookup = make(map[mods.ModID]string)

func Initialize() (err error) {
	f := path.Join(config.PWD, file)
	if _, err = os.Stat(f); err != nil {
		return nil
	}
	var b []byte
	if b, err = os.ReadFile(f); err != nil {
		return fmt.Errorf("failed to read %s: %v", file, err)
	}
	if err = json.Unmarshal(b, &lookup); err != nil {
		return fmt.Errorf("failed to read %s: %v", file, err)
	}
	return nil
}

func GetDir(modID mods.ModID) (dir string, found bool) {
	if modID != "" {
		dir, found = lookup[modID]
	}
	return
}

func SetDir(modID mods.ModID, dir string) (err error) {
	lookup[modID] = dir
	var (
		b []byte
		f *os.File
	)

	if b, err = json.MarshalIndent(&lookup, "", "\t"); err != nil {
		return fmt.Errorf("failed to prepare %s: %v", file, err)
	}

	if f, err = os.Create(path.Join(config.PWD, file)); err != nil {
		return fmt.Errorf("failed to create %s: %v", file, err)
	}
	defer func() { _ = f.Close() }()

	if _, err = f.Write(b); err != nil {
		return fmt.Errorf("failed to write %s: %v", file, err)
	}
	return nil
}
