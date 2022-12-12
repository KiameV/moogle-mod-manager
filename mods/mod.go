package mods

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods/managed/cache"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/util"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type (
	SelectType string
	Category   string
	ModID      string
)

const (
	Auto   SelectType = "Auto"
	Select SelectType = "Select"
	Radio  SelectType = "Radio"
)

var SelectTypes = []string{string(Auto), string(Select), string(Radio)}

const (
	BattleScene        Category = "Battle Scene"
	EnemySprite        Category = "Enemy Sprite"
	GameOverhauls      Category = "Game Overhauls"
	Gameplay           Category = "Gameplay"
	General            Category = "General"
	Fonts              Category = "Fonts"
	PlayerNpcSprites   Category = "Player/NPC Sprite"
	ScriptText         Category = "Script/Text"
	Soundtrack         Category = "Soundtrack"
	TileSet            Category = "Tile Set"
	TitleScreen        Category = "Title Screen"
	UIGeneral          Category = "UI: General"
	UiMenuPortraits    Category = "UI: Menu Portraits"
	UiTextBoxPortraits Category = "UI: Textbox Portraits"
	UiWindowFrames     Category = "UI: Window Frames"
	Utility            Category = "Utility"
)

var Categories = []string{
	string(BattleScene),
	string(EnemySprite),
	string(GameOverhauls),
	string(Gameplay),
	string(PlayerNpcSprites),
	string(ScriptText),
	string(Soundtrack),
	string(TileSet),
	string(TitleScreen),
	string(UIGeneral),
	string(UiMenuPortraits),
	string(UiWindowFrames),
	string(UiTextBoxPortraits),
	string(Utility)}

type Mod struct {
	ID                  ModID             `json:"ID" xml:"ID"`
	Name                string            `json:"Name" xml:"Name"`
	Author              string            `json:"Author" xml:"Author"`
	AuthorLink          string            `json:"AuthorLink" xml:"AuthorLink"`
	ReleaseDate         string            `json:"ReleaseDate" xml:"ReleaseDate"`
	Category            Category          `json:"Category" xml:"Category"`
	Description         string            `json:"Description" xml:"Description"`
	ReleaseNotes        string            `json:"ReleaseNotes" xml:"ReleaseNotes"`
	Link                string            `json:"Link" xml:"Link"`
	Version             string            `json:"Version"`
	Preview             *Preview          `json:"Preview,omitempty" xml:"Preview,omitempty"`
	ModKind             ModKind           `json:"ModKind" xml:"ModKind"`
	ModCompatibility    *ModCompatibility `json:"Compatibility,omitempty" xml:"ModCompatibility,omitempty"`
	Downloadables       []*Download       `json:"Downloadable" xml:"Downloadables"`
	DonationLinks       []*DonationLink   `json:"DonationLink" xml:"DonationLinks"`
	Games               []*Game           `json:"Games" xml:"Games"`
	AlwaysDownload      []*DownloadFiles  `json:"AlwaysDownload,omitempty" xml:"AlwaysDownload,omitempty"`
	Configurations      []*Configuration  `json:"Configuration,omitempty" xml:"Configurations,omitempty"`
	ConfigSelectionType SelectType        `json:"ConfigSelectionType" xml:"ConfigSelectionType"`
	IsManuallyCreated   bool              `json:"IsManuallyCreated" xml:"IsManuallyCreated"`
}

func UniqueModID(game config.Game, modID ModID) string {
	return fmt.Sprintf("%d.%s", game, modID)
}

func (m *Mod) UniqueModID(game config.Game) string {
	return UniqueModID(game, NewModID(m.ModKind.Kind, string(m.ID)))
}

func (m *Mod) ModIdAsNumber() (uint64, error) {
	sp := strings.Split(string(m.ID), ".")
	return strconv.ParseUint(sp[len(sp)-1], 10, 64)
}

func (m *Mod) BranchName() string {
	return fmt.Sprintf("%s_%s", m.ID, m.Version)
}

type Preview struct {
	Url   *string       `json:"Url,omitempty" xml:"Url,omitempty"`
	Local *string       `json:"Local,omitempty" xml:"Local,omitempty"`
	img   *canvas.Image `json:"-" xml:"-"`
}

type Size struct {
	X int `json:"X" xml:"X"`
	Y int `json:"Y" xml:"Y"`
}

func (p *Preview) Get() *canvas.Image {
	if p == nil {
		return nil
	}
	if p.img == nil {
		var (
			r   fyne.Resource
			err error
		)
		if p.Local != nil {
			f := filepath.Join(state.GetBaseDir(), *p.Local)
			if _, err = os.Stat(f); err == nil {
				r, err = fyne.LoadResourceFromPath(f)
			}
		}
		if r == nil && p.Url != nil {
			if r, err = cache.GetImage(*p.Url); err != nil {
				r, err = fyne.LoadResourceFromURLString(*p.Url)
			}
		}
		if r == nil || err != nil {
			return nil
		}
		p.img = canvas.NewImageFromResource(r)
		p.img.SetMinSize(fyne.Size{Width: float32(300), Height: float32(300)})
		p.img.FillMode = canvas.ImageFillContain
	}
	return p.img
}

type ModCompatibility struct {
	Requires []*ModCompat `json:"Require" xml:"Requires"`
	Forbids  []*ModCompat `json:"Forbid" xml:"Forbids"`
	//OrderConstraints []ModCompat `json:"OrderConstraint"`
}

func (c *ModCompatibility) HasItems() bool {
	return c != nil && (len(c.Requires) > 0 || len(c.Forbids) > 0)
}

type InstallType string

const (
	//Bundles  InstallType = "Bundles"
	//Memoria  InstallType = "Memoria"
	//Magicite InstallType = "Magicite"
	//BepInEx  InstallType = "BepInEx"
	// DLL Patcher https://discord.com/channels/371784427162042368/518331294858608650/863930606446182420
	//DllPatch   InstallType = "DllPatch"
	Compressed InstallType = "Compressed"
)

var InstallTypes = []string{ /*string(Bundles), string(Memoria), string(Magicite), string(BepInEx), /*string(DllPatch),*/ string(Compressed)}

type Game struct {
	Name     config.GameName `json:"Name" xml:"Name"`
	Versions []string        `json:"Version,omitempty" xml:"GameVersions,omitempty"`
}

type DownloadFiles struct {
	DownloadName string `json:"DownloadName" xml:"DownloadName"`
	// IsInstallAll is used by nexus mods when a mod.xml is not used
	Files []*ModFile `json:"File,omitempty" xml:"Files,omitempty"`
	Dirs  []*ModDir  `json:"Dir,omitempty" xml:"Dirs,omitempty"`
}

func (f *DownloadFiles) IsEmpty() bool {
	return len(f.Files) == 0 && len(f.Dirs) == 0
}

type ModFile struct {
	From    string  `json:"From" xml:"From"`
	To      string  `json:"To" xml:"To"`
	Archive *string `json:"Archive,omitempty" xml:"Archive,omitempty"`
}

type ModDir struct {
	From      string  `json:"From" xml:"From"`
	To        string  `json:"To" xml:"To"`
	Recursive bool    `json:"Recursive" xml:"Recursive"`
	Archive   *string `json:"Archive,omitempty" xml:"Archive,omitempty"`
}

type Configuration struct {
	Name        string    `json:"Name" xml:"Name"`
	Description string    `json:"Description" xml:"Description"`
	Preview     *Preview  `json:"Preview,omitempty" xml:"Preview, omitempty"`
	Root        bool      `json:"Root" xml:"Root"`
	Choices     []*Choice `json:"Choice" xml:"Choices"`

	NextConfigurationName *string `json:"NextConfigurationName,omitempty" xml:"NextConfigurationName"`
}

type Choice struct {
	Name          string         `json:"Name" xml:"Name"`
	Description   string         `json:"Description" xml:"Description"`
	Preview       *Preview       `json:"Preview,omitempty" xml:"Preview,omitempty"`
	DownloadFiles *DownloadFiles `json:"DownloadFiles,omitempty" xml:"DownloadFiles,omitempty"`

	NextConfigurationName *string `json:"NextConfigurationName,omitempty" xml:"NextConfigurationName"`
}

type DonationLink struct {
	Name string `json:"Name" xml:"Name"`
	Link string `json:"Link" xml:"Link"`
}

func (m *Mod) Validate() string {
	sb := strings.Builder{}
	if m.ID == "" {
		sb.WriteString("ID is required\n")
	}
	if m.Name == "" {
		sb.WriteString("Name is required\n")
	}
	if m.Version == "" {
		sb.WriteString("Version is required\n")
	}
	if m.Author == "" {
		sb.WriteString("Author is required\n")
	}
	if m.ReleaseDate == "" {
		sb.WriteString("Release Date is required\n")
	}
	if m.Category == "" {
		//	sb.WriteString("Category is required\n")
	}
	if m.Description == "" {
		sb.WriteString("Description is required\n")
	}
	if m.Link == "" {
		sb.WriteString("Link is required\n")
	}

	/*if m.Preview != nil {
		if m.Preview.Size.X <= 50 || m.Preview.Size.Y <= 50 {
			sb.WriteString("Preview size must be greater than 50\n")
		}
	}*/

	kind := m.ModKind.Kind
	dlableNames := make(map[string]bool)
	if len(m.Downloadables) == 0 {
		sb.WriteString("Must have at least one Downloadables\n")
	}
	for _, d := range m.Downloadables {
		if d.Name == "" {
			sb.WriteString("Downloadables' name is required\n")
		}
		//if strings.Index(d.Name, " ") != -1 {
		//	sb.WriteString(fmt.Sprintf("Downloadables [%s]'s name cannot contain spaces\n"))
		//}
		if kind == Hosted {
			if d.Hosted == nil {
				sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Hosted is required\n", d.Name))
			} else {
				if len(d.Hosted.Sources) == 0 {
					sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Source is required\n", d.Name))
				}
				for _, s := range d.Hosted.Sources {
					u, err := url.Parse(s)
					if err != nil {
						sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Source [%s] is not a valid url: %v\n", d.Name, s, err))
					}
					i := strings.LastIndex(u.Path, "/")
					j := strings.LastIndex(u.Path, ".")
					if i == -1 || j == -1 {
						sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Source [%s] is not a valid url, maybe missing .zip/.rar/.7z\n", d.Name, s))
					}
					s = u.Path[i+1 : j]
					if d.Name != s {
						sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Source [%s] must be the same as the name after the extension is removed\n", d.Name, s))
					}
					dlableNames[s] = true
				}
			}
		} else if kind == Nexus {
			if d.Nexus == nil {
				sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Nexus is required\n", d.Name))
			}
			if d.Nexus.FileID <= 0 {
				sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Nexus FileID must be greater than 0\n", d.Name))
			}
			if d.Nexus.FileName == "" {
				sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Nexus FileName is required\n", d.Name))
			}
			dlableNames[d.Name] = true
		} else if kind == CurseForge {
			if d.CurseForge == nil {
				sb.WriteString(fmt.Sprintf("Downloadables [%s]'s CF is required\n", d.Name))
			}
			if d.CurseForge.FileID <= 0 {
				sb.WriteString(fmt.Sprintf("Downloadables [%s]'s CF FileID must be greater than 0\n", d.Name))
			}
			if d.CurseForge.FileName == "" {
				sb.WriteString(fmt.Sprintf("Downloadables [%s]'s CF FileName is required\n", d.Name))
			}
			dlableNames[d.Name] = true
		}
		//if d.InstallType == "" {
		//	sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Install Type is required\n", d.Name))
		//}
	}

	if len(m.AlwaysDownload) == 0 && len(m.Configurations) == 0 {
		sb.WriteString("One \"Always Download\", at least one \"Configuration\" or both are required\n")
	}

	for _, ad := range m.AlwaysDownload {
		if ad.IsEmpty() {
			sb.WriteString(fmt.Sprintf("AlwaysDownload [%s]' Must have at least one File or Dir specified\n", ad.DownloadName))
		}
		if _, ok := dlableNames[ad.DownloadName]; !ok {
			sb.WriteString("Always Download's downloadable doesn't exist\n")
		}
	}

	roots := 0
	for _, c := range m.Configurations {
		if c.Name == "" {
			sb.WriteString("Configuration's Name is required\n")
		}
		if c.Description == "" {
			//sb.WriteString(fmt.Sprintf("Configuration's [%s] Description is required\n", c.Name))
		}
		if len(c.Choices) == 0 {
			sb.WriteString(fmt.Sprintf("Configuration's [%s] must have Choices\n", c.Name))
		}
		for _, ch := range c.Choices {
			if ch.Name == "" {
				sb.WriteString(fmt.Sprintf("Configuration's [%s] Choice's Name is required\n", c.Name))
			}
			if ch.NextConfigurationName != nil && *ch.NextConfigurationName == c.Name {
				sb.WriteString(fmt.Sprintf("Configuration's [%s] Choice [%s]'s Next Configuration Name must not be the same as the Configuration's Name\n", c.Name, ch.Name))
			}
			if ch.DownloadFiles != nil && ch.DownloadFiles.DownloadName != "" {
				if _, ok := dlableNames[ch.DownloadFiles.DownloadName]; !ok {
					sb.WriteString(fmt.Sprintf("Configuration's [%s] Choice [%s]'s downloadable doesn't exist\n", c.Name, ch.Name))
				}
			} else if ch.DownloadFiles != nil && ch.DownloadFiles.DownloadName == "" && (len(ch.DownloadFiles.Files) > 0 || len(ch.DownloadFiles.Dirs) > 0) {
				sb.WriteString(fmt.Sprintf("Configuration's [%s] Choice [%s]'s downloadable must be specified\n", c.Name, ch.Name))
			}
		}
		if c.Root {
			roots++
		}
	}
	if len(m.Configurations) > 1 && roots == 0 {
		sb.WriteString("Must have at least one 'Root' Configuration\n")
	}

	return sb.String()
}

func (m *Mod) Supports(game config.Game) error {
	name := config.GameToName(game)
	if len(m.Games) > 0 {
		for _, g := range m.Games {
			if name == g.Name {
				return nil
			}
		}
	}
	return fmt.Errorf("%s does not support %s", m.Name, config.GameNameString(game))
}

func (m *Mod) Merge(from Mod) {
	if m.IsManuallyCreated {
		m.Author = from.Author
		if m.AuthorLink == "" {
			m.AuthorLink = from.AuthorLink
		}
		from.ModCompatibility = m.ModCompatibility
		from.Games = m.Games
		from.Link = m.Link
	} else if from.IsManuallyCreated {
		m.Description = from.Description
		m.ReleaseNotes = from.ReleaseNotes
	}
	return
}

func NewModForVersion(manual *Mod, remote *Mod) *Mod {
	var m Mod
	if manual != nil && manual.IsManuallyCreated {
		m = *manual
		m.Version = remote.Version
		m.Downloadables = remote.Downloadables
	} else {
		m = *remote
	}
	return &m
}

func (m *Mod) DirectoryName() string {
	return util.CreateFileName(string(m.ID))
}

func Sort(mods []*Mod) (sorted []*Mod) {
	var (
		lookup = make(map[string]*Mod)
		sl     = make([]string, len(mods))
		m      *Mod
		i      int
		key    string
	)
	for i, m = range mods {
		key = fmt.Sprintf("%s%s", m.Name, m.ID)
		lookup[key] = m
		sl[i] = key
	}

	sort.Strings(sl)

	sorted = make([]*Mod, len(mods))
	for i, key = range sl {
		sorted[i] = lookup[key]
	}
	return
}

func NewModID(k Kind, modID string) ModID {
	prefix := strings.ToLower(string(k))
	if k == Hosted || strings.HasPrefix(modID, prefix) {
		return ModID(modID)
	}
	return ModID(fmt.Sprintf("%s.%s", prefix, modID))
}
