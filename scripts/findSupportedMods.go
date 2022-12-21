package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	u "net/url"
	"os"
	"strings"
)

const (
	// game path, mod id
	getFilesUrl = "https://api.nexusmods.com/v1/games/:1/mods/:2/files.json?category=main%2Cupdate%2Coptional%2Cmiscellaneous"
	// game id, mod id, file name.json
	getStructureUrl = "https://file-metadata.nexusmods.com/file/nexus-files-s3-meta/%s/%s/%s.json"
)

// https://file-metadata.nexusmods.com/file/nexus-files-s3-meta/4335/30/Lopsho's%20Classic%20SNES%20UI-30-1-0-1647718556.zip.json

type (
	filesResp struct {
		Files []files `json:"files"`
	}
	files struct {
		ID   int    `json:"file_id"`
		Name string `json:"file_name"`
	}
	FileType string
	dlFiles  struct {
		Path     string    `json:"path"`
		Name     string    `json:"name"`
		Type     FileType  `json:"type"`
		Children []dlFiles `json:"children"`
	}
	game struct {
		id   string
		path string
	}
)

const (
	Dir  FileType = "directory"
	File FileType = "file"
)

var games = []game{
	//{path: "finalfantasypixelremaster", id: "3934"},
	{path: "finalfantasy2pixelremaster", id: "3958"},
	{path: "finalfantasy3pixelremaster", id: "3942"},
	{path: "finalfantasy4pixelremaster", id: "4022"},
	{path: "finalfantasy5pixelremaster", id: "4137"},
	{path: "finalfantasy6pixelremaster", id: "4335"},
}

func GetFileContents(apiKey string) (string, error) {
	for z, g := range games {
		fo, _ := os.Create(fmt.Sprintf("%d-output.txt", z+2))
		for i := 1; i < 60; i++ {
			var fs filesResp
			url := strings.Replace(getFilesUrl, ":1", g.path, 1)
			url = strings.Replace(url, ":2", fmt.Sprintf("%d", i), 1)
			if err := sendRequest(apiKey, url, &fs); err != nil {
				continue
			}
			if len(fs.Files) == 1 {
				url = fmt.Sprintf(getStructureUrl, g.id, fmt.Sprintf("%d", i), u.QueryEscape(fs.Files[0].Name))
				var dlf dlFiles
				if err := sendRequest(apiKey, url, &dlf); err != nil {
					continue
				}
				var fs []string
				compileFiles(dlf, &fs)
				for _, f := range fs {
					if strings.HasPrefix(f, "FINAL FANTASY") ||
						strings.Contains(f, "StandaloneWindows") ||
						strings.Contains(f, "/StreamingAssets") ||
						strings.Contains(f, "/FINAL FANTASY") {
						var sb strings.Builder
						sb.WriteString(fmt.Sprintf("ID %d\n", i))
						for _, f := range fs {
							sb.WriteString("- ")
							sb.WriteString(f)
							sb.WriteString("\n")
						}
						fo.WriteString(sb.String())
						break
					}
				}
			}
		}
		fo.Close()
	}
	return "", nil
}

func sendRequest(apiKey, url string, to any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("apikey", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, to); err != nil {
		return err
	}
	return nil
}

func compileFiles(dlf dlFiles, files *[]string) {
	for _, f := range dlf.Children {
		if f.Type == Dir {
			compileFiles(f, files)
		} else if f.Type == File {
			*files = append(*files, f.Path)
		}
	}
}

func main() {
	GetFileContents(os.Getenv("apikey"))
}
