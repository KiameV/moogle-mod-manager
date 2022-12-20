package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/mods"
	"io"
	"net/http"
)

const (
	getTagUrl        = "https://api.github.com/repos/%s/%s/releases/latest"
	listDownloadsUrl = "https://api.github.com/repos/%s/%s/releases/tags/%s"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func LatestReleaseFromMod(mod *mods.Mod) (tag string, err error) {
	if mod == nil || mod.ModKind.GitHub == nil {
		return "", errors.New("mod is nil or not a github mod")
	}
	return LatestRelease(mod.ModKind.GitHub.Owner, mod.ModKind.GitHub.Repo)
}

func LatestRelease(owner, repo string) (tag string, err error) {
	var (
		resp    *http.Response
		body    []byte
		release Release
		// Replace with the owner and repository of the desired GitHub repository
	)

	// Send a GET request to the GitHub REST API to retrieve the latest release
	resp, err = http.Get(fmt.Sprintf(getTagUrl, owner, repo))
	if err != nil {
		// Handle error
		return
	}
	defer func() { _ = resp.Body.Close() }()

	// Read the response body and unmarshal it into a Release struct
	if body, err = io.ReadAll(resp.Body); err != nil {
		// Handle error
		return
	}

	if err = json.Unmarshal(body, &release); err != nil {
		// Handle error
		return
	}

	return release.TagName, nil
}

type (
	Download struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	}
	DlRelease struct {
		Assets []Download `json:"assets"`
	}
)

func ListDownloads(owner, repo, tag string) (downloads []Download, err error) {
	// Send a GET request to the GitHub REST API to retrieve the specified release
	var (
		resp    *http.Response
		release DlRelease
		body    []byte
	)

	if resp, err = http.Get(fmt.Sprintf(listDownloadsUrl, owner, repo, tag)); err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Read the response body and unmarshal it into a Release struct
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	// Return the list of assets
	return release.Assets, nil
}
