package util

import (
	"encoding/base64"
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
