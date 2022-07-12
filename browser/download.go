package browser

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Download(url, toDir string) (string, error) {
	var (
		resp, err = http.Get(url)
		out       *os.File
		sp        = strings.Split(url, "/")
		file      string
	)
	if err != nil {
		return "", err
	}
	defer func() { resp.Body.Close() }()

	if len(sp) == 0 {
		return "", fmt.Errorf("invalid url: %s", url)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to download the mod's source at %s, url")
	}

	file = filepath.Join(toDir, sp[len(sp)-1])
	// Create the file
	if out, err = os.Create(file); err != nil {
		return "", err
	}
	defer func() { _ = out.Close() }()

	// Write the body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", err
	}
	return file, nil
}
