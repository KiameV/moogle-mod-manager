package mods

type Mod struct {
	ID               string           `json:"ID"`
	Name             string           `json:"Name"`
	Author           string           `json:"Author"`
	Version          string           `json:"Version"`
	ReleaseDate      string           `json:"ReleaseDate"`
	Category         string           `json:"Category"`
	Description      string           `json:"Description"`
	ReleaseNotes     string           `json:"ReleaseNotes"`
	Link             string           `json:"Link"`
	Preview          string           `json:"Preview"`
	ModCompatibility ModCompatibility `json:"Compatibility"`
	GameVersions     []string         `json:"GameVersion"`
	Downloadables    []Download       `json:"Downloadable"`
	DonationLinks    []DonationLink   `json:"DonationLink"`

	// Either Download or Configs for a mod
	DownloadFiles  *DownloadFiles  `json:"DownloadFiles,omitempty"`
	Configurations []Configuration `json:"Configuration"`
}

type ModCompatibility struct {
	Requires []ModCompat `json:"Require"`
	Forbids  []ModCompat `json:"Forbid"`
	//OrderConstraints []ModCompat `json:"OrderConstraint"`
}

type ModCompat struct {
	ModID    string   `json:"ModID"`
	Versions []string `json:"Version"`
	Source   string   `json:"Source"`
	//Order   ModCompatOrder `json:"Order"`
}

/*type ModCompatOrder string

const (
	Before ModCompatOrder = "Before"
	After  ModCompatOrder = "After"
)*/

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

type Download struct {
	Name        string   `json:"Name"`
	Sources     []string `json:"Source"`
	InstallType string   `json:"InstallType"`
}

type DownloadFiles struct {
	DownloadName string    `json:"DownloadName"`
	Files        []ModFile `json:"File"`
}

type ModFile struct {
	From string `json:"From"`
	To   string `json:"To"`
}

type Configuration struct {
	Name        string   `json:"Name"`
	Description string   `json:"Description"`
	Preview     string   `json:"Preview"`
	Choices     []Choice `json:"Choice"`
}

type Choice struct {
	Description           string        `json:"Description"`
	Preview               string        `json:"Preview"`
	DownloadFiles         DownloadFiles `json:"DownloadFiles"`
	NextConfigurationName *string       `json:"NextConfigurationName,omitempty"`
}

type DonationLink struct {
	Name string `json:"Name"`
	Link string `json:"Link"`
}
