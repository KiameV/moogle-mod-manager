package browser

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func Download(url, toDir string) (string, error) {
	var (
		buf, name, err = download(url)
		out            *os.File
		file           = path.Join(toDir, name)
	)
	// Create the file
	if out, err = os.Create(file); err != nil {
		return "", err
	}
	defer func() { _ = out.Close() }()

	// Write the body to file
	_, err = io.Copy(out, buf)
	return file, err
}

func DownloadAsString(url string) (string, error) {
	buf, _, err := download(url)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func DownloadAsBytes(url string) ([]byte, error) {
	buf, _, err := download(url)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func download(url string) (buf *bytes.Buffer, name string, err error) {
	var (
		resp *http.Response
		sp   = strings.Split(url, "/")
	)
	if resp, err = http.Get(url); err != nil {
		return
	}
	defer func() { resp.Body.Close() }()

	if len(sp) == 0 {
		err = fmt.Errorf("invalid url: %s", url)
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("failed to download the mod's source at %s", url)
		return
	}
	defer resp.Body.Close()

	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	name = sp[len(sp)-1]
	return
}
