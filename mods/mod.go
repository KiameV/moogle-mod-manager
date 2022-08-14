package mods

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods/managed/cache"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type SelectType string

const (
	Auto   SelectType = "Auto"
	Select SelectType = "Select"
	Radio  SelectType = "Radio"
)

var SelectTypes = []string{string(Auto), string(Select), string(Radio)}

type Mod struct {
	ID                  string            `json:"ID" xml:"ID"`
	Name                string            `json:"Name" xml:"Name"`
	Author              string            `json:"Author" xml:"Author"`
	AuthorLink          string            `json:"AuthorLink" xml:"AuthorLink"`
	ReleaseDate         string            `json:"ReleaseDate" xml:"ReleaseDate"`
	Category            string            `json:"Category" xml:"Category"`
	Description         string            `json:"Description" xml:"Description"`
	ReleaseNotes        string            `json:"ReleaseNotes" xml:"ReleaseNotes"`
	Link                string            `json:"Link" xml:"Link"`
	Version             string            `json:"Version"`
	Preview             *Preview          `json:"Preview,omitempty" xml:"Preview,omitempty"`
	ModKind             *ModKind          `json:"ModKind" xml:"ModKind"`
	ModCompatibility    *ModCompatibility `json:"Compatibility,omitempty" xml:"ModCompatibility,omitempty"`
	Downloadables       []*Download       `json:"Downloadable" xml:"Downloadables"`
	DonationLinks       []*DonationLink   `json:"DonationLink" xml:"DonationLinks"`
	Games               []*Game           `json:"Games" xml:"Games"`
	AlwaysDownload      []*DownloadFiles  `json:"AlwaysDownload,omitempty" xml:"AlwaysDownload,omitempty"`
	Configurations      []*Configuration  `json:"Configuration,omitempty" xml:"Configurations,omitempty"`
	ConfigSelectionType SelectType        `json:"ConfigSelectionType" xml:"ConfigSelectionType"`
}

type Preview struct {
	Url   *string       `json:"Url,omitempty" xml:"Url,omitempty"`
	Local *string       `json:"Local,omitempty" xml:"Local,omitempty"`
	Size  Size          `json:"Size,omitempty" xml:"Size,omitempty"`
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
		size := fyne.Size{Width: float32(p.Size.X), Height: float32(p.Size.Y)}
		p.img.SetMinSize(size)
		p.img.Resize(size)
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
	return len(c.Requires) > 0 || len(c.Forbids) > 0
}

type ModCompat struct {
	ModID    string          `json:"ModID" xml:"ModID"`
	Name     string          `json:"Name" xml:"Name"`
	Versions []string        `json:"Version,omitempty" xml:"Versions,omitempty"`
	Source   string          `json:"Source" xml:"Source"`
	Order    *ModCompatOrder `json:"Order,omitempty" xml:"Order,omitempty"`
}

type ModCompatOrder string

const (
	None   ModCompatOrder = ""
	Before ModCompatOrder = "Before"
	After  ModCompatOrder = "After"
)

var ModCompatOrders = []string{string(None), string(Before), string(After)}

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
	From string `json:"From" xml:"From"`
	To   string `json:"To" xml:"To"`
}

type ModDir struct {
	From      string `json:"From" xml:"From"`
	To        string `json:"To" xml:"To"`
	Recursive bool   `json:"Recursive" xml:"Recursive"`
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
		sb.WriteString("Category is required\n")
	}
	if m.Description == "" {
		sb.WriteString("Description is required\n")
	}
	if m.Link == "" {
		sb.WriteString("Link is required\n")
	}

	if m.Preview != nil {
		if m.Preview.Size.X <= 50 || m.Preview.Size.Y <= 50 {
			sb.WriteString("Preview size must be greater than 50\n")
		}
	}

	kind := m.ModKind.Kind
	if kind == Hosted {
		h := m.ModKind.Hosted
		if h == nil {
			sb.WriteString("Hosted is required\n")
		} else {
			if len(h.ModFileLinks) == 0 {
				sb.WriteString("Hosted 'Mod File' Links is required\n")
			}
			for _, mfl := range h.ModFileLinks {
				if strings.HasSuffix(mfl, ".json") == false && strings.HasSuffix(mfl, ".xml") == false {
					sb.WriteString(fmt.Sprintf("Hosted 'Mod File' Link [%s] must be json or xml\n", mfl))
				}
			}
		}
	} else { // nexus
		n := m.ModKind.Nexus
		if n == nil {
			sb.WriteString("Nexus is required\n")
		} else {
			if n.ID == "" {
				sb.WriteString("Nexus Mod ID is required\n")
			}
		}
	}

	dlableNames := make(map[string]bool)
	if len(m.Downloadables) == 0 {
		sb.WriteString("Must have at least one Downloadables\n")
	}
	for _, d := range m.Downloadables {
		if d.Name == "" {
			sb.WriteString("Downloadables' name is required\n")
		}
		if strings.Index(d.Name, " ") != -1 {
			sb.WriteString(fmt.Sprintf("Downloadables [%s]'s name cannot contain spaces\n"))
		}
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
		} else {
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
			sb.WriteString(fmt.Sprintf("Configuration's [%s] Description is required\n", c.Name))
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
			}
		}
		if c.Root {
			roots++
		}
	}
	if len(m.Configurations) > 1 && roots == 0 {
		sb.WriteString("Must have at least one 'Root' Configuration\n")
	} else if roots > 1 {
		sb.WriteString("Only one 'Root' Configuration is allowed\n")
	}

	return sb.String()
}

func (m *Mod) Supports(game config.Game) error {
	gs := " " + config.String(game)
	for _, g := range m.Games {
		if strings.HasSuffix(string(g.Name), gs) {
			return nil
		}
	}
	return fmt.Errorf("%s does not support %s", m.Name, config.GameNameString(game))
}
