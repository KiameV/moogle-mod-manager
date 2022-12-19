package curseforge

import (
	"encoding/json"
	"errors"
	"fmt"
	converter "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/kiamev/moogle-mod-manager/config"
	u "github.com/kiamev/moogle-mod-manager/discover/remote/util"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	getModsByGameID        = "https://api.curseforge.com/v1/mods/search?gameId=%d"
	GetModsByGameIDAndName = "https://api.curseforge.com/v1/mods/search?gameId=%d&slug=%s"
	getModDataByModID      = "https://api.curseforge.com/v1/mods/%d"
	getModFilesByModID     = "https://api.curseforge.com/v1/mods/%d/files"
	getModDescByModID      = "https://api.curseforge.com/v1/mods/%d/description"
)

type client struct {
	compiler u.ModCompiler
}

func NewClient(compiler u.ModCompiler) *client {
	c := &client{compiler: compiler}
	compiler.SetFinder(c)
	return c
}

func IsCurseforge(url string) bool {
	return strings.Contains(url, "curseforge.com")
}

func (c *client) GetFromMod(in *mods.Mod) (found bool, mod *mods.Mod, err error) {
	if len(in.Games) == 0 {
		err = fmt.Errorf("no games found for mod %s", in.Name)
		return
	}
	var id uint64
	if id, err = in.ModIdAsNumber(); err != nil {
		err = fmt.Errorf("could not parse mod id %s for %s", in.ModID, in.Name)
		return
	}
	return c.get(fmt.Sprintf(getModDataByModID, id))
}

func (c *client) GetFromID(_ config.GameDef, id int) (found bool, mod *mods.Mod, err error) {
	return c.get(fmt.Sprintf(getModDataByModID, id))
}

func (c *client) GetFromUrl(url string) (found bool, mod *mods.Mod, err error) {
	// www.curseforge.com/final-fantasy-vi/mods/gau-rage-descriptions-extended-magna-roader-fix/files/4073052
	url = strings.Replace(url, "https://", "", 1)
	url = strings.Replace(url, "http://", "", 1)
	sp := strings.Split(url, "/")
	if len(sp) < 4 {
		err = errors.New("invalid url")
		return
	}
	var game config.GameDef
	if game, err = config.GameDefFromCfPath(config.CfPath(sp[1])); err != nil {
		return
	}
	slug := sp[3]
	return c.get(fmt.Sprintf(GetModsByGameIDAndName, game.Remote().CurseForge.ID, slug))
}

func (c *client) get(url string) (found bool, mod *mods.Mod, err error) {
	var (
		b      []byte
		result cfModResponse
		desc   string
		dls    fileParent
	)
	if b, err = sendRequest(url); err != nil {
		return
	}
	if err = json.Unmarshal(b, &result); err != nil {
		return
	}

	if dls, err = getDownloads(result.Data); err != nil {
		return
	}

	if desc, err = getDescription(result.Data); err != nil {
		return
	}

	return toMod(result.Data, desc, dls.Files)
}

func (c *client) GetNewestMods(game config.GameDef, lastID int) (result []*mods.Mod, err error) {
	var (
		b       []byte
		dls     fileParent
		mod     *mods.Mod
		desc    string
		include bool
	)
	if b, err = sendRequest(fmt.Sprintf(getModsByGameID, game.Remote().CurseForge.ID)); err != nil {
		return
	}
	var data struct {
		Mods []cfMod `json:"data"`
	}
	if err = json.Unmarshal(b, &data); err != nil {
		return
	}

	result = make([]*mods.Mod, 0, len(data.Mods))
	for _, m := range data.Mods {
		if m.ModID > lastID {
			if dls, err = getDownloads(m); err != nil {
				return
			}
			if desc, err = getDescription(m); err != nil {
				return
			}
			if include, mod, err = toMod(m, desc, dls.Files); err != nil {
				return
			} else if !include {
				continue
			}
			result = append(result, mod)
		}
	}
	return
}

func getDownloads(m cfMod) (dls fileParent, err error) {
	var (
		b   []byte
		url = fmt.Sprintf(getModFilesByModID, m.ModID)
	)
	if b, err = sendRequest(url); err != nil {
		return
	}
	err = json.Unmarshal(b, &dls)
	return
}

func getDescription(m cfMod) (s string, err error) {
	var (
		b   []byte
		d   description
		url = fmt.Sprintf(getModDescByModID, m.ModID)
	)
	if b, err = sendRequest(url); err != nil {
		return
	}
	err = json.Unmarshal(b, &d)
	s = httpToMarkdown(d.Data)
	return
}

func httpToMarkdown(s string) (result string) {
	var (
		c   = converter.NewConverter("", true, nil)
		err error
	)
	s = removeFont(s)
	if result, err = c.ConvertString(s); err != nil {
		result = s
	}
	result = strings.ReplaceAll(result, "<br />", "\n")
	result = strings.ReplaceAll(result, "\\\\_", "_")
	result = strings.ReplaceAll(result, "\\_", "_")
	return
}

func sendRequest(url string) (response []byte, err error) {
	var (
		apiKey = config.GetSecrets().CfApiKey
		req    *http.Request
		resp   *http.Response
	)
	if apiKey == "" {
		err = errors.New("no CurseForge Api Key set. Please go to File->Secrets")
		return
	}
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		err = fmt.Errorf("failed to create request to validate user with nexus %s: %v", url, err)
		return
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-key", apiKey)

	if resp, err = (&http.Client{}).Do(req); err != nil {
		err = fmt.Errorf("failed to make request to %s: %v", url, err)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	code := resp.StatusCode
	if code < 200 && code >= 300 {
		err = fmt.Errorf("received code [%d] from call to [%s]", code, url)
		return
	}

	if response, err = io.ReadAll(resp.Body); err != nil {
		err = fmt.Errorf("failed to read response's body for %s: %v", url, err)
	}
	return
}

func toMod(m cfMod, desc string, dls []CfFile) (include bool, mod *mods.Mod, err error) {
	var (
		modID = fmt.Sprintf("%d", m.ModID)
		game  config.GameDef
	)
	if game, err = config.GameDefFromCfID(m.Game); err != nil {
		return
	}
	mod = mods.NewMod(&mods.ModDef{
		ModID:        mods.NewModID(mods.CurseForge, modID),
		Name:         mods.ModName(m.Name),
		Version:      m.Version(),
		Description:  desc,
		ReleaseDate:  m.CreatedTime.Format("Jan 2, 2006"),
		ReleaseNotes: "",
		Link:         m.Links.WebsiteUrl,
		Preview: &mods.Preview{
			Url:   &m.Logo.Url,
			Local: nil,
		},
		ModKind: mods.ModKind{
			Kind: mods.CurseForge,
		},
		Games: []*mods.Game{{
			ID:       game.ID(),
			Versions: nil,
		}},
		Downloadables:  make([]*mods.Download, len(dls)),
		DonationLinks:  nil,
		AlwaysDownload: nil,
		Configurations: nil,
	})
	if len(m.Category) > 0 {
		if mod.Category, err = m.Category[0].toCategory(); err != nil {
			return
		}
	}
	if len(m.Author) > 0 {
		mod.Author = m.Author[0].Name
		mod.AuthorLink = m.Author[0].Url
	}

	var choices []*mods.Choice
	for i, d := range dls {
		mod.Downloadables[i] = &mods.Download{
			Name:    d.Name,
			Version: d.Version(),
			CurseForge: &mods.CurseForgeDownloadable{
				RemoteDownloadable: mods.RemoteDownloadable{
					FileID:   d.FileID,
					FileName: d.Name,
				},
				Url: d.DownloadUrl,
			},
		}
		dlf := &mods.DownloadFiles{
			DownloadName: d.Name,
			Dirs: []*mods.ModDir{
				{
					From:      string(game.BaseDir()),
					To:        string(game.BaseDir()),
					Recursive: true,
				},
			},
		}
		choices = append(choices, &mods.Choice{
			Name:                  d.Name,
			Description:           "",
			Preview:               nil,
			DownloadFiles:         dlf,
			NextConfigurationName: nil,
		})
	}

	include = true
	if len(choices) > 1 {
		mod.Configurations = []*mods.Configuration{
			{
				Name:        "Choose preference",
				Description: "",
				Preview:     nil,
				Root:        true,
				Choices:     choices,
			},
		}
	} else if len(choices) == 1 {
		mod.AlwaysDownload = append(mod.AlwaysDownload, choices[0].DownloadFiles)
	} else {
		include = false
	}
	return
}

func removeFont(s string) string {
	var i, j int
	for i = 0; i < len(s)-10; i++ {
		if s[i] == '[' && s[i+1] == 'f' && s[i+2] == 'o' && s[i+3] == 'n' && s[i+4] == 't' && s[i+5] == '=' {
			for j = i; j < len(s) && s[j] != ']'; j++ {
			}
			s = s[:i] + s[j+1:]
		}
	}
	s = strings.ReplaceAll(s, "[/font]", "")
	return s
}

func (c *client) Folder(game config.GameDef) string {
	return filepath.Join(config.PWD, "remote", string(game.ID()), string(mods.CurseForge))
}

func (c *client) GetMods(game config.GameDef) (result []*mods.Mod, err error) {
	if game == nil {
		return nil, errors.New("GetMods called with a nil game")
	}
	dir := c.Folder(game)
	_ = os.MkdirAll(dir, 0777)
	if err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "mod.json" || d.Name() == "mod.xml" {
			m := &mods.Mod{}
			if err = util.LoadFromFile(path, m); err != nil {
				return err
			}
			result = append(result, m)
		}
		return nil
	}); err != nil {
		return
	}
	return c.compiler.AppendNewMods(c.Folder(game), game, result)
}
