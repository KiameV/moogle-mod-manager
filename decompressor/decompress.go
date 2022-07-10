package decompressor

import (
	"path/filepath"
)

func Unzip(src, dest string) (err error) {
	if ext := filepath.Ext(src); ext == ".7z" {
		return handle7zip(src, dest)
	}

	return
}
