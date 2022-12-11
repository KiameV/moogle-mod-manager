package curseforge

import (
	"encoding/json"
	"errors"
	"fmt"
	converter "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"io"
	"net/http"
	"strings"
)

type CfGame string
type CfGameID int

const (
	FFI   CfGame = "final-fantasy-i"
	FFII  CfGame = "final-fantasy-ii"
	FFIII CfGame = "final-fantasy-iii"
	FFIV  CfGame = "final-fantasy-iv"
	FFV   CfGame = "final-fantasy-v"
	FFVI  CfGame = "final-fantasy-vi"

	IdFFI   CfGameID = 5230
	IdFFII  CfGameID = 5001
	IdFFIII CfGameID = 5026
	IdFFIV  CfGameID = 4741
	IdFFV   CfGameID = 5021
	IdFFVI  CfGameID = 4773

	getModsByGameID        = "https://api.curseforge.com/v1/mods/search?gameId=%d"
	GetModsByGameIDAndName = "https://api.curseforge.com/v1/mods/search?gameId=%d&slug=%s"
	getModDataByModID      = "https://api.curseforge.com/v1/mods/%d"
	getModFilesByModID     = "https://api.curseforge.com/v1/mods/%d/files"
	getModDescByModID      = "https://api.curseforge.com/v1/mods/%d/description"
)

type CurseForgeClient struct{}

func IsCurseforge(url string) bool {
	return strings.Index(url, "curseforge.com") >= 0
}

func GameToCfGame(game config.Game) CfGame {
	switch game {
	case config.I:
		return FFI
	case config.II:
		return FFII
	case config.III:
		return FFIII
	case config.IV:
		return FFIV
	case config.V:
		return FFV
	default:
		return FFVI
	}
}

func CfGameIDToGame(id CfGameID) config.Game {
	switch id {
	case IdFFI:
		return config.I
	case IdFFII:
		return config.II
	case IdFFIII:
		return config.III
	case IdFFIV:
		return config.IV
	case IdFFV:
		return config.V
	default:
		return config.VI
	}
}

func CfGameToGameID(game CfGame) CfGameID {
	switch game {
	case FFI:
		return IdFFI
	case FFII:
		return IdFFII
	case FFIII:
		return IdFFIII
	case FFIV:
		return IdFFIV
	case FFV:
		return IdFFV
	default:
		return IdFFVI
	}
}

func CfGameToGame(game CfGame) config.Game {
	switch game {
	case FFI:
		return config.I
	case FFII:
		return config.II
	case FFIII:
		return config.III
	case FFIV:
		return config.IV
	case FFV:
		return config.V
	default:
		return config.VI

	}
}

func GameToID(game config.Game) CfGameID {
	switch game {
	case config.I:
		return IdFFI
	case config.II:
		return IdFFII
	case config.III:
		return IdFFIII
	case config.IV:
		return IdFFIV
	case config.V:
		return IdFFV
	default:
		return IdFFVI
	}
}

func (c *CurseForgeClient) GetFromMod(in *mods.Mod) (found bool, mod *mods.Mod, err error) {
	if len(in.Games) == 0 {
		err = errors.New("no games found for mod " + in.Name)
		return
	}
	var id uint64
	if id, err = in.ModIdAsNumber(); err != nil {
		err = fmt.Errorf("could not parse mod id %s for %s", in.ID, in.Name)
		return
	}
	return c.get(fmt.Sprintf(getModDataByModID, id))
}

func (c *CurseForgeClient) GetFromID(_ config.Game, id int) (found bool, mod *mods.Mod, err error) {
	return c.get(fmt.Sprintf(getModDataByModID, id))
}

func (c *CurseForgeClient) GetFromUrl(url string) (found bool, mod *mods.Mod, err error) {
	// www.curseforge.com/final-fantasy-vi/mods/gau-rage-descriptions-extended-magna-roader-fix/files/4073052
	url = strings.Replace(url, "https://", "", 1)
	url = strings.Replace(url, "http://", "", 1)
	sp := strings.Split(url, "/")
	if len(sp) < 4 {
		err = errors.New("invalid url")
		return
	}
	game := sp[1]
	slug := sp[3]
	return c.get(fmt.Sprintf(GetModsByGameIDAndName, CfGameToGameID(CfGame(game)), slug))
}

func (c *CurseForgeClient) get(url string) (found bool, mod *mods.Mod, err error) {
	var (
		b      []byte
		result cfMods
		desc   string
		dls    fileParent
	)
	if b, err = sendRequest(url); err != nil {
		return
	}
	if err = json.Unmarshal(b, &result); err != nil {
		return
	}
	if len(result.Data) == 0 {
		err = errors.New("no mod found for at " + url)
		return
	}

	if dls, err = getDownloads(result.Data[0]); err != nil {
		return
	}

	if desc, err = getDescription(result.Data[0]); err != nil {
		return
	}

	return toMod(result.Data[0], desc, dls.Files)
}

func (c *CurseForgeClient) GetNewestMods(game config.Game, lastID int) (result []*mods.Mod, err error) {
	var (
		b     []byte
		dls   fileParent
		mod   *mods.Mod
		desc  string
		found bool
	)
	if b, err = sendRequest(fmt.Sprintf(getModsByGameID, GameToID(game))); err != nil {
		return
	}
	var nMods []cfMod
	if err = json.Unmarshal(b, &nMods); err != nil {
		return
	}

	result = make([]*mods.Mod, 0, len(nMods))
	for _, m := range nMods {
		if m.ModID > lastID {
			if dls, err = getDownloads(m); err != nil {
				return
			}
			if desc, err = getDescription(m); err != nil {
				return
			}
			if found, mod, err = toMod(m, desc, dls.Files); err != nil {
				return
			} else if found {
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
	modID := fmt.Sprintf("%d", m.ModID)
	game := CfGameIDToGame(m.Game)
	mod = &mods.Mod{
		ID:           mods.NewModID(mods.CurseForge, modID),
		Name:         m.Name,
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
			Name:     config.GameToName(game),
			Versions: nil,
		}},
		Downloadables:  make([]*mods.Download, len(dls)),
		DonationLinks:  nil,
		AlwaysDownload: nil,
		Configurations: nil,
	}
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
					From:      string(mods.GameToInstallBaseDir(game)),
					To:        string(mods.GameToInstallBaseDir(game)),
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
