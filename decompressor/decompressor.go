package decompressor

import (
	"fmt"
	"path"
)

type Decompressor interface {
	DecompressTo(dest string) error
}

func NewDecompressor(src string) (Decompressor, error) {
	switch path.Ext(src) {
	case ".7z":
		return new7zDecompressor(src), nil
	case ".zip", ".rar":
		return newArchiveDecompressor(src), nil
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", path.Ext(src))
	}
}
