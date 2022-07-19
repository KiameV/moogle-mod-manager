package mods

import "github.com/kiamev/moogle-mod-manager/config"

type Mod struct {
	ID               string            `json:"ID" xml:"ID"`
	Name             string            `json:"Name" xml:"Name"`
	Author           string            `json:"Author" xml:"Author"`
	Version          string            `json:"Version" xml:"Version"`
	ReleaseDate      string            `json:"ReleaseDate" xml:"ReleaseDate"`
	Category         string            `json:"Category" xml:"Category"`
	Description      string            `json:"Description" xml:"Description"`
	ReleaseNotes     string            `json:"ReleaseNotes" xml:"ReleaseNotes"`
	Link             string            `json:"Link" xml:"Link"`
	Preview          string            `json:"Preview" xml:"Preview"`
	ModCompatibility *ModCompatibility `json:"Compatibility,omitempty" xml:"ModCompatibility,omitempty"`
	Downloadables    []*Download       `json:"Downloadable" xml:"Downloadables"`
	DonationLinks    []*DonationLink   `json:"DonationLink" xml:"DonationLinks"`
	Game             []*Game           `json:"Games" xml:"Games"`

	// Either Download or Configs for a mod
	DownloadFiles  *DownloadFiles   `json:"DownloadFile,omitempty" xml:"DownloadFiles,omitempty"`
	Configurations []*Configuration `json:"Configuration" xml:"Configurations"`
}

type ModCompatibility struct {
	Requires []*ModCompat `json:"Require" xml:"Requires"`
	Forbids  []*ModCompat `json:"Forbid" xml:"Forbids"`
	//OrderConstraints []ModCompat `json:"OrderConstraint"`
}

type ModCompat struct {
	ModID    string          `json:"ModID" xml:"ModID"`
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
	Bundles  InstallType = "Bundles"
	Memoria  InstallType = "Memoria"
	Magicite InstallType = "Magicite"
	BepInEx  InstallType = "BepInEx"
	// DLL Patcher https://discord.com/channels/371784427162042368/518331294858608650/863930606446182420
	DllPatch   InstallType = "DllPatch"
	Compressed InstallType = "Compressed"
)

var InstallTypes = []string{string(Bundles), string(Memoria), string(Magicite), string(BepInEx), string(DllPatch), string(Compressed)}

type Game struct {
	Name     config.GameName `json:"Name" xml:"Name"`
	Versions []string        `json:"Version,omitempty" xml:"GameVersions,omitempty"`
}

type Download struct {
	Name        string      `json:"Name" xml:"Name"`
	Sources     []string    `json:"Source" xml:"Sources"`
	InstallType InstallType `json:"InstallType" xml:"InstallType"`
}

type DownloadFiles struct {
	DownloadName string    `json:"DownloadName" xml:"DownloadName"`
	Files        []ModFile `json:"File" xml:"Files"`
}

type ModFile struct {
	From string `json:"From" xml:"From"`
	To   string `json:"To" xml:"To"`
}

type Configuration struct {
	Name        string   `json:"Name" xml:"Name"`
	Description string   `json:"Description" xml:"Description"`
	Preview     string   `json:"Preview" xml:"Preview"`
	Choices     []Choice `json:"Choice" xml:"Choices"`
}

type Choice struct {
	Description           string        `json:"Description" xml:"Description"`
	Preview               string        `json:"Preview" xml:"Preview"`
	DownloadFiles         DownloadFiles `json:"DownloadFiles" xml:"DownloadFiles"`
	NextConfigurationName *string       `json:"NextConfigurationName,omitempty" xml:"NextConfigurationName"`
}

type DonationLink struct {
	Name string `json:"Name" xml:"Name"`
	Link string `json:"Link" xml:"Link"`
}
