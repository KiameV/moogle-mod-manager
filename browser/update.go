package browser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	Version = "v0.6.5"

	tagUrl = `https://api.github.com/repos/KiameV/ffprModManager/tags`
	relUrl = `https://github.com/KiameV/ffprModManager/releases/%s`
)

type (
	comparer struct {
		major, minor, patch uint64
		pre                 int64
	}
	tag struct {
		Name string `json:"name"`
	}
)

func newComparer(v string) *comparer {
	v = strings.TrimPrefix(v, "v")
	var (
		c   = &comparer{}
		sp  = strings.Split(v, "_pre")
		err error
	)
	if len(sp) > 1 {
		if c.pre, err = strconv.ParseInt(sp[1], 10, 64); err != nil {
			c.pre = -1
		}
		v = sp[0]
	}
	sp = strings.Split(v, ".")
	if c.major, err = strconv.ParseUint(sp[0], 10, 64); err != nil {
		return nil
	}
	if c.minor, err = strconv.ParseUint(sp[1], 10, 64); err != nil {
		return nil
	}
	if c.patch, err = strconv.ParseUint(sp[2], 10, 64); err != nil {
		return nil
	}
	return c
}

func (c comparer) isGreaterThan(cv *comparer) bool {
	if cv == nil {
		return true
	}
	if cv.pre != 0 && c.pre == 0 {
		return true
	}
	if c.pre != 0 {
		if c.major == cv.major && c.minor == cv.minor && c.patch == cv.patch {
			if cv.pre == 0 {
				return false
			}
			return c.pre > cv.pre
		}
	}
	return c.major > cv.major ||
		c.minor > cv.minor ||
		c.patch > cv.patch
}

func CheckForUpdate() (hasNewer bool, version string, err error) {
	var (
		r              *http.Response
		b              []byte
		tags           []tag
		highestVersion = newComparer(Version)
	)
	if r, err = http.Get(tagUrl); err != nil {
		return
	}
	if b, err = io.ReadAll(r.Body); err != nil {
		return
	}
	if err = json.Unmarshal(b, &tags); err != nil {
		return
	}

	for _, t := range tags {
		if strings.Contains(t.Name, ".") {
			i := newComparer(t.Name)
			if !highestVersion.isGreaterThan(i) {
				hasNewer = true
				version = t.Name
				highestVersion = i
			}
		}
	}
	return
}

func Update(tag string) (err error) {
	url := fmt.Sprintf(relUrl, tag)
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return
}
