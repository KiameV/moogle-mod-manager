package nexus

import (
	"encoding/json"
	"errors"
	"fmt"
	converter "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/frustra/bbcode"
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
	nexusApiModUrl         = "https://api.nexusmods.com/v1/games/%s/mods/%s.json"
	nexusApiModDlUrl       = "https://api.nexusmods.com/v1/games/%s/mods/%s/files.json%s"
	nexusApiModDlUrlSuffix = "?category=main,update,optional,miscellaneous"
	nexusUrl               = "https://www.nexusmods.com/%s/mods/%d"
	nexusApiNewestModsUrl  = "https://api.nexusmods.com/v1/games/%s/mods/latest_added.json"

	//nexusUsersApiUrl = "https://users.nexusmods.com/oauth/token"

	// NexusFileDownload file_id, NexusGameID
	NexusFileDownload = "https://www.nexusmods.com/Core/Libs/Common/Widgets/DownloadPopUp?id=%d&game_id=%v"
)

type client struct {
	compiler u.ModCompiler
}

func NewClient(compiler u.ModCompiler) *client {
	c := &client{compiler: compiler}
	compiler.SetFinder(c)
	return c
}

func IsNexus(url string) bool {
	return strings.Contains(url, "nexusmods.com")
}

func (c *client) GetFromMod(in *mods.Mod) (found bool, mod *mods.Mod, err error) {
	if len(in.Games) == 0 {
		err = fmt.Errorf("no games found for mod %s", in.Name)
		return
	}
	var (
		id   uint64
		game config.GameDef
	)
	if game, err = config.GameDefFromID(in.Games[0].ID); err != nil {
		return
	}
	if id, err = in.ModIdAsNumber(); err != nil {
		err = fmt.Errorf("could not parse mod id %s for %s", in.ModID, in.Name)
		return
	}
	return c.GetFromUrl(fmt.Sprintf(nexusUrl, game.Remote().Nexus.Path, id))
}

func (c *client) GetFromID(game config.GameDef, id int) (found bool, mod *mods.Mod, err error) {
	return c.GetFromUrl(fmt.Sprintf(nexusUrl, game.Remote().Nexus.Path, id))
}

func (c *client) GetFromUrl(url string) (found bool, mod *mods.Mod, err error) {
	var (
		sp    = strings.Split(url, "/")
		path  string
		modID string
		b     []byte
		nMod  nexusMod
		nDls  fileParent
	)
	for i, s := range sp {
		if s == "mods" {
			if i > 0 && i < len(sp)-1 {
				path = sp[i-1]
				modID = strings.Split(sp[i+1], "?")[0]
				break
			}
		}
	}
	if path == "" || modID == "" {
		err = fmt.Errorf("could not get GameDef and Mod ModID from %s", url)
		return
	}
	if b, err = sendRequest(fmt.Sprintf(nexusApiModUrl, path, modID)); err != nil {
		return
	}
	if err = json.Unmarshal(b, &nMod); err != nil {
		return
	}
	if nMod.Name == "" {
		err = errors.New("no mod found for " + modID)
		return
	}
	if nDls, err = getDownloads(config.NexusPath(path), modID); err != nil {
		return
	}
	return toMod(nMod, nDls.Files)
}

func (c *client) GetNewestMods(game config.GameDef, lastID int) (result []*mods.Mod, err error) {
	var (
		b     []byte
		path  = game.Remote().Nexus.Path
		nDls  fileParent
		mod   *mods.Mod
		found bool
	)
	if b, err = sendRequest(fmt.Sprintf(nexusApiNewestModsUrl, path)); err != nil {
		return
	}
	var nMods []nexusMod
	if err = json.Unmarshal(b, &nMods); err != nil {
		return
	}

	result = make([]*mods.Mod, 0, len(nMods))
	for _, nMod := range nMods {
		if nMod.ModID > lastID {
			if nDls, err = getDownloads(path, fmt.Sprintf("%d", nMod.ModID)); err != nil {
				return
			}
			if found, mod, err = toMod(nMod, nDls.Files); err != nil {
				return
			} else if found {
				continue
			}
			result = append(result, mod)
		}
	}
	return
}

func getDownloads(path config.NexusPath, modID string) (nDls fileParent, err error) {
	var (
		b   []byte
		url = fmt.Sprintf(nexusApiModDlUrl, path, modID, nexusApiModDlUrlSuffix)
	)
	if b, err = sendRequest(url); err != nil {
		return
	}
	err = json.Unmarshal(b, &nDls)
	return
}

func sendRequest(url string) (response []byte, err error) {
	var (
		apiKey = config.GetSecrets().NexusApiKey
		req    *http.Request
		resp   *http.Response
	)
	if apiKey == "" {
		err = errors.New("no Nexus Api Key set. Please go to File->Secrets")
		return
	}
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		err = fmt.Errorf("failed to create request to validate user with nexus %s: %v", url, err)
		return
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("apikey", apiKey)

	if resp, err = (&http.Client{}).Do(req); err != nil {
		err = fmt.Errorf("failed to make request to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

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

func toMod(n nexusMod, dls []NexusFile) (include bool, mod *mods.Mod, err error) {
	var (
		modID = fmt.Sprintf("%d", n.ModID)
		game  config.GameDef
	)
	if game, err = config.GameDefFromNexusPath(n.GamePath); err != nil {
		return
	}
	mod = mods.NewMod(&mods.ModDef{
		ModID:        mods.NewModID(mods.Nexus, modID),
		Name:         mods.ModName(n.Name),
		Version:      n.Version,
		Author:       n.Author,
		AuthorLink:   n.AuthorLink,
		Category:     "",
		ReleaseDate:  n.CreatedTime.Format("Jan 2, 2006"),
		ReleaseNotes: "",
		Link:         fmt.Sprintf(nexusUrl, n.GamePath, n.ModID),
		Preview: &mods.Preview{
			Url:   &n.PictureUrl,
			Local: nil,
			//Size:  nil,
		},
		ModKind: mods.ModKind{
			Kind: mods.Nexus,
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
	compiler := bbcode.NewCompiler(true, true)
	c := converter.NewConverter("", true, nil)
	cd := compiler.Compile(n.Description)
	cd = removeFont(cd)
	if mod.Description, err = c.ConvertString(cd); err != nil {
		mod.Description = n.Description
		err = nil
	}
	mod.Description = strings.ReplaceAll(mod.Description, "<br />", "\n")
	mod.Description = strings.ReplaceAll(mod.Description, "\\\\_", "_")
	mod.Description = strings.ReplaceAll(mod.Description, "\\_", "_")

	var choices []*mods.Choice
	for i, d := range dls {
		mod.Downloadables[i] = &mods.Download{
			Name:    d.Name,
			Version: d.Version,
			Nexus: &mods.RemoteDownloadable{
				FileID:   d.FileID,
				FileName: d.FileName,
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
			Description:           d.Description,
			Preview:               nil,
			DownloadFiles:         dlf,
			NextConfigurationName: nil,
		})
	}

	include = true
	if len(choices) > 1 {
		mod.Configurations = []*mods.Configuration{
			{
				Name:        "Choose",
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
	return filepath.Join(config.PWD, "remote", string(game.ID()), string(mods.Nexus))
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

/*
function convert() {
  var left = document.getElementById('left_ta');
  var right = document.getElementById('right_ta');

  var left_value = left.value;

  //preprocessing for tf2toolbox BBCode
  if(left_value.search(/TF2Toolbox/gmi) != -1) {
    left_value = left_value
    .replace(/(\(List generated at .+?\[\/URL\]\))((?:.|\n)+)/gmi, '$2\n\n\n$1') //Move TF2Toolbox link to bottom
    .replace('(List generated at', '(List generated from')
    .replace(/[^\S\n]+\(List/gmi,'(List')
    .replace(/\[b\]\[u\](.+?)\[\/u\]\[\/b\]/gmi,'[b]$1[/b]\n') //Fix double emphasized titles
    .replace(/(\n)\[\*\]\[b\](.+?)\[\/b\]/gmi, '$1\[\*\] $2');
  }

  //general BBcode conversion
  left_value = left_value
    .replace(/\[b\]((?:.|\n)+?)\[\/b\]/gmi, '**$1**') //bold; replace [b] $1 [/b] with ** $1 **
    .replace(/\[\i\]((?:.|\n)+?)\[\/\i\]/gmi, '*$1*')  //italics; replace [i] $1 [/u] with * $1 *
    .replace(/\[\u\]((?:.|\n)+?)\[\/\u\]/gmi, '$1')  //remove underline;
    .replace(/\[s\]((?:.|\n)+?)\[\/s\]/gmi, '~~ $1~~') //strikethrough; replace [s] $1 [/s] with ~~ $1 ~~
    .replace(/\[center\]((?:.|\n)+?)\[\/center\]/gmi, '$1') //remove center;
    .replace(/\[quote\=.+?\]((?:.|\n)+?)\[\/quote\]/gmi, '$1') //remove [quote=] tags
    .replace(/\[size\=.+?\]((?:.|\n)+?)\[\/size\]/gmi, '## $1') //Size [size=] tags
    .replace(/\[color\=.+?\]((?:.|\n)+?)\[\/color\]/gmi, '$1') //remove [color] tags
    .replace(/\[list\=1\]((?:.|\n)+?)\[\/list\]/gmi, function (match, p1, offset, string) {return p1.replace(/\[\*\]/gmi, '1. ');})
    .replace(/(\n)\[\*\]/gmi, '$1* ') //lists; replcae lists with + unordered lists.
    .replace(/\[\/*list\]/gmi, '')
    .replace(/\[img\]((?:.|\n)+?)\[\/img\]/gmi,'![$1]($1)')
    .replace(/\[url=(.+?)\]((?:.|\n)+?)\[\/url\]/gmi,'[$2]($1)')
    .replace(/\[code\](.*?)\[\/code\]/gmi, '`$1`')
    .replace(/\[code\]((?:.|\n)+?)\[\/code\]/gmi, function (match, p1, offset, string) {return p1.replace(/^/gmi, '    ');})
    .replace(/\[php\](.*?)\[\/php\]/gmi, '`$1`')
    .replace(/\[php\]((?:.|\n)+?)\[\/php\]/gmi, function (match, p1, offset, string) {return p1.replace(/^/gmi, '    ');})
    .replace(/\[pawn\](.*?)\[\/pawn\]/gmi, '`$1`')
    .replace(/\[pawn\]((?:.|\n)+?)\[\/pawn\]/gmi, function (match, p1, offset, string) {return p1.replace(/^/gmi, '    ');});

  //post processing for tf2toolbox BBCode
  if(left_value.search(/TF2Toolbox/gmi) != -1) {
    left_value = left_value
    .replace('/bbcode_lookup.php))', '/bbcode_lookup.php) and converted to /r/tf2trade ready Markdown by Dum\'s [converter](http://jondum.github.com/BBCode-To-Markdown-Converter/)).') //add a linkback
    .replace(/\*\*.+?\*\*[\s]+?None[\s]{2}/gmi, ''); //remove empty sections

  }

  right.value = left_value;

}
*/
