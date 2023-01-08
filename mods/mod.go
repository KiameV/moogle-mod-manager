package mods

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/util"
	"sort"
	"strings"
)

type (
	SelectType string
	Category   string
	ModID      string
	ModName    string
)

func (n ModName) Contains(text string) bool {
	return strings.Contains(strings.ToLower(string(n)), strings.ToLower(text))
}

const (
	Auto   SelectType = "Auto"
	Select SelectType = "Select"
	Radio  SelectType = "Radio"
)

var SelectTypes = []string{string(Auto), string(Select), string(Radio)}

const (
	BattleScene        Category = "Battle Scene"
	EnemySprite        Category = "Enemy Sprite"
	GameOverhauls      Category = "GameDef Overhauls"
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
	string(Fonts),
	string(PlayerNpcSprites),
	string(ScriptText),
	string(Soundtrack),
	string(TileSet),
	string(TitleScreen),
	string(UIGeneral),
	string(UiMenuPortraits),
	string(UiWindowFrames),
	string(UiTextBoxPortraits),
	string(Utility),
}

const (
	BugKindCosmetic Kind = "Cosmetic"
	BugKindAnnoying Kind = "Annoying"
	BugKindCrash    Kind = "Critical"
)

type (
	BugKind   string
	BugReport struct {
		Kind        BugKind `json:"Kind" xml:"Kind"`
		Description string  `json:"Description" xml:"Description"`
	}
	ModDef struct {
		ModID               ModID               `json:"ID" xml:"ID"`
		Name                ModName             `json:"Name" xml:"Name"`
		Author              string              `json:"Author" xml:"Author"`
		AuthorLink          string              `json:"AuthorLink" xml:"AuthorLink"`
		ReleaseDate         string              `json:"ReleaseDate" xml:"ReleaseDate"`
		Category            Category            `json:"Category" xml:"Category"`
		Description         string              `json:"Description" xml:"Description"`
		ReleaseNotes        string              `json:"ReleaseNotes" xml:"ReleaseNotes"`
		Link                string              `json:"Link" xml:"Link"`
		Version             string              `json:"Version" xml:"Version"`
		InstallType_        *config.InstallType `json:"InstallType,omitempty" xml:"InstallType,omitempty"`
		Preview             *Preview            `json:"Preview,omitempty" xml:"Preview,omitempty"`
		Previews            []*Preview          `json:"Previews,omitempty" xml:"Previews,omitempty"`
		ModKind             ModKind             `json:"ModKind" xml:"ModKind"`
		ModCompatibility    *ModCompatibility   `json:"Compatibility,omitempty" xml:"ModCompatibility,omitempty"`
		Downloadables       []*Download         `json:"Downloadable" xml:"Downloadables"`
		DonationLinks       []*DonationLink     `json:"DonationLink" xml:"DonationLinks"`
		Games               []*Game             `json:"Games" xml:"Games"`
		AlwaysDownload      []*DownloadFiles    `json:"AlwaysDownload,omitempty" xml:"AlwaysDownload,omitempty"`
		Configurations      []*Configuration    `json:"Configuration,omitempty" xml:"Configurations,omitempty"`
		ConfigSelectionType SelectType          `json:"ConfigSelectionType" xml:"ConfigSelectionType"`
		Hide                bool                `json:"Hide" xml:"Hide"`
		VerifiedAsWorking   bool                `json:"VerifiedAsWorking" xml:"VerifiedAsWorking"`
		Bugs                []BugReport         `json:"Bugs,omitempty" xml:"Bugs,omitempty"`
		IsManuallyCreated   bool                `json:"IsManuallyCreated" xml:"IsManuallyCreated"`
	}
)

func NewMod(def *ModDef) *Mod {
	return &Mod{ModDef: def}
}

type Mod struct {
	*ModDef
}

func (m *Mod) ID() ModID {
	return m.ModID
}

func (m *Mod) Kinds() Kinds {
	return m.ModKind.Kinds
}

func (m *Mod) InstallType(game config.GameDef) config.InstallType {
	i := game.DefaultInstallType()
	if m.InstallType_ != nil {
		i = *m.InstallType_
	}
	return i
}

func (m *Mod) BranchName() string {
	return fmt.Sprintf("%s_%s", m.ModID, m.Version)
}

func (m *Mod) Save(to string) error {
	return util.SaveToFile(to, m.ModDef, '\n')
}

type Size struct {
	X int `json:"X" xml:"X"`
	Y int `json:"Y" xml:"Y"`
}

type ModCompatibility struct {
	Requires []*ModCompat `json:"Require" xml:"Requires"`
	Forbids  []*ModCompat `json:"Forbid" xml:"Forbids"`
	//OrderConstraints []ModCompat `json:"OrderConstraint"`
}

func (c *ModCompatibility) HasItems() bool {
	return c != nil && (len(c.Requires) > 0 || len(c.Forbids) > 0)
}

var InstallTypes = []string{string(config.Move), string(config.MoveToArchive)}

type Game struct {
	ID       config.GameID    `json:"Name" xml:"Name"`
	Versions []config.Version `json:"Version,omitempty" xml:"GameVersions,omitempty"`
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

func (f *DownloadFiles) HasArchive() []string {
	var s []string
	for _, file := range f.Files {
		if file.ToArchive == nil {
			s = append(s, file.From)
		}
	}
	for _, dir := range f.Dirs {
		if dir.ToArchive == nil {
			s = append(s, dir.From)
		}
	}
	return s
}

type ModFile struct {
	From      string  `json:"From" xml:"From"`
	To        string  `json:"To" xml:"To"`
	ToArchive *string `json:"ToArchive,omitempty" xml:"ToArchive,omitempty"`
}

type ModDir struct {
	From      string  `json:"From" xml:"From"`
	To        string  `json:"To" xml:"To"`
	Recursive bool    `json:"Recursive" xml:"Recursive"`
	ToArchive *string `json:"ToArchive,omitempty" xml:"ToArchive,omitempty"`
}

type Configuration struct {
	Name        string    `json:"Name" xml:"Name"`
	Description string    `json:"Description" xml:"Description"`
	Preview     *Preview  `json:"Preview,omitempty" xml:"Preview,omitempty"`
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
	if m.ModID == "" {
		sb.WriteString("ModID is required\n")
	}
	if m.Name == "" {
		sb.WriteString("Name is required\n")
	}

	if m.Hide {
		if m.ModKind.Kinds.IsHosted() {
			sb.WriteString("Cannot Hide hosted mods\n")
		}
		return sb.String()
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
		sb.WriteString("Category is required\n")
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

	//kinds := m.ModKind.Kinds
	dlableNames := make(map[string]bool)
	for _, d := range m.Downloadables {
		if d.Name == "" {
			sb.WriteString("Downloadables' name is required\n")
		}
		// TODO add more validations
		dlableNames[d.Name] = true
	}
	/*if len(m.Downloadables) == 0 {
		sb.WriteString("Must have at least one Downloadables\n")
	}
	for _, d := range m.Downloadables {
		if d.Name == "" {
			sb.WriteString("Downloadables' name is required\n")
		}
		if kinds.Is(HostedAt) {
			sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Hosted is required\n", d.Name))
		} else { // GitHub
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
		if kinds.Is(Nexus) {
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
		}
		if kinds.Is(CurseForge) {
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
	}*/

	if len(m.AlwaysDownload) == 0 && len(m.Configurations) == 0 {
		sb.WriteString("One \"Always Download\", at least one \"Configuration\" or both are required\n")
	}

	for _, ad := range m.AlwaysDownload {
		if ad.IsEmpty() {
			sb.WriteString(fmt.Sprintf("AlwaysDownload [%s] Must have at least one File or Dir specified\n", ad.DownloadName))
		}
		if m.InstallType_.Is(config.MoveToArchive) {
			if f := ad.HasArchive(); len(f) > 0 {
				sb.WriteString("AlwaysDownload missing archives for " + strings.Join(f, ", "))
			}
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
		//if c.Description == "" {
		//	sb.WriteString(fmt.Sprintf("Configuration's [%s] Description is required\n", c.Name))
		//}
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
				if ch.DownloadFiles.IsEmpty() {
					sb.WriteString(fmt.Sprintf("Configuration's [%s] Choice [%s]'s must have at least one File or Dir specified\n", c.Name, ch.Name))
				}
				if m.InstallType_.Is(config.MoveToArchive) {
					if f := ch.DownloadFiles.HasArchive(); len(f) > 0 {
						sb.WriteString(fmt.Sprintf("Configuration's [%s] Choice [%s]'s missing archives for "+strings.Join(f, ", "), c.Name, ch.Name))
					}
				}
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

func (m *Mod) Supports(game config.GameDef) error {
	if len(m.Games) > 0 {
		for _, g := range m.Games {
			if game.ID() == g.ID {
				return nil
			}
		}
	}
	return fmt.Errorf("%s does not support %s", m.Name, game.Name())
}

//func (m *Mod) Merge(from Mod) {
//	if m.IsManuallyCreated {
//		m.Author = from.Author
//		if m.AuthorLink == "" {
//			m.AuthorLink = from.AuthorLink
//		}
//		from.ModCompatibility = m.ModCompatibility
//		from.Games = m.Games
//		from.Link = m.Link
//	} else if from.IsManuallyCreated {
//		m.Description = from.Description
//		m.ReleaseNotes = from.ReleaseNotes
//	}
//}

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

func (m *Mod) Mod() *Mod {
	return m
}

func (m *Mod) LoadFromFile(file string) (err error) {
	if err = util.LoadFromFile(file, m); err == nil {
		if m.Preview != nil && len(m.Previews) == 0 {
			m.Previews = []*Preview{m.Preview}
		}
	}
	return
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
		key = fmt.Sprintf("%s%s", m.Name, m.ModID)
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
	if k == HostedAt || k == HostedGitHub || strings.HasPrefix(modID, prefix) {
		return ModID(modID)
	}
	return ModID(fmt.Sprintf("%s.%s", prefix, modID))
}
