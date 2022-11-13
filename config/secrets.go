package config

import (
	"github.com/kiamev/moogle-mod-manager/util"
	"path/filepath"
)

const secretsFile = "secrets.json"

var secret = &Secret{}

type Secret struct {
	NexusApiKey string `json:"nexusApiKey"`
	CfApiKey    string `json:"cfApiKey"`
}

func GetSecrets() *Secret {
	return secret
}

func (s *Secret) Initialize() {
	_ = util.LoadFromFile(filepath.Join(PWD, secretsFile), s)
}

func (s *Secret) Save() error {
	return util.SaveToFile(filepath.Join(PWD, secretsFile), s)
}
