package secrets

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/util"
	"path/filepath"
)

const secretsFile = "secrets.json"

var secret = &sec{}

type (
	Key byte
	sec struct {
		NexusApiKey string `json:"nexusApiKey"`
		CfApiKey    string `json:"cfApiKey"`
	}
)

const (
	_ Key = iota
	NexusApiKey
	CfApiKey
)

func Initialize() {
	_ = util.LoadFromFile(filepath.Join(config.PWD, secretsFile), secret)
}

func Save() error {
	return util.SaveToFile(filepath.Join(config.PWD, secretsFile), secret)
}

func Get(k Key) (v string) {
	if k == NexusApiKey {
		v = secret.NexusApiKey
	} else if k == CfApiKey {
		v = secret.CfApiKey
	}
	return
}

func Set(k Key, v string) {
	if k == NexusApiKey {
		v = secret.NexusApiKey
	} else if k == CfApiKey {
		v = secret.CfApiKey
	}
}
