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
		name, err = getName(url)
		buf       *bytes.Buffer
		out       *os.File
		file      = path.Join(toDir, name)
	)
	if err != nil {
		return "", err
	}

	if _, err = os.Stat(file); err == nil {
		// Already downloaded
		return file, nil
	}

	if buf, err = download(url); err != nil {
		return "", err
	}

	if err = os.MkdirAll(toDir, 0777); err != nil {
		return "", err
	}

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
	buf, err := download(url)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func DownloadAsBytes(url string) ([]byte, error) {
	buf, err := download(url)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func download(url string) (buf *bytes.Buffer, err error) {
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return
	}
	defer func() { resp.Body.Close() }()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("failed to download the mod's source at %s", url)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)

	return
}

func getName(url string) (name string, err error) {
	sp := strings.Split(url, "/")
	if len(sp) == 0 {
		err = fmt.Errorf("invalid url: %s", url)
		return
	}
	name = sp[len(sp)-1]
	if i := strings.Index(name, "?"); i >= 0 {
		name = name[:i]
	}
	return
}
