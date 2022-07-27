package util

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func CreateFileName(s string) string {
	if reg, err := regexp.Compile("[^a-zA-Z0-9]+"); err == nil {
		s = strings.TrimSpace(reg.ReplaceAllString(s, ""))
		if len(s) > 0 {
			return s
		}
	}
	return strings.Trim(base64.StdEncoding.EncodeToString([]byte(s)), "=")
}

func LoadFromFile(file string, i interface{}) error {
	var (
		b   []byte
		ext = filepath.Ext(file)
		err error
	)
	if _, err = os.Stat(file); err != nil {
		return err
	}
	if b, err = os.ReadFile(file); err != nil {
		return fmt.Errorf("failed to read %s: %v", file, err)
	}
	switch ext {
	case ".json", ".moogle":
		err = json.Unmarshal(b, &i)
	case ".xml":
		err = xml.Unmarshal(b, &i)
	default:
		return fmt.Errorf("unknown file extension: %s", file)
	}
	if err = json.Unmarshal(b, i); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %v", file, err)
	}
	return nil
}

func SaveToFile(file string, i interface{}) (err error) {
	var (
		b []byte
		f *os.File
	)
	if b, err = json.MarshalIndent(i, "", "\t"); err != nil {
		return fmt.Errorf("failed to marshal %s: %v", file, err)
	}
	if f, err = os.Create(file); err != nil {
		return fmt.Errorf("failed to create %s: %v", file, err)
	}
	if _, err = f.Write(b); err != nil {
		return fmt.Errorf("failed to write %s: %v", file, err)
	}
	return
}
