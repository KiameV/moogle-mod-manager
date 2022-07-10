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
	GameVersions     []string         `json:"GameVersions"`

	// Either Files or Configs for a mod
	Files          []ModFile       `json:"Files"`
	Configurations []Configuration `json:"Configurations"`
}

type ModCompatibility struct {
	Require          []ModCompat `json:"Require"`
	Forbid           []ModCompat `json:"Forbid"`
	OrderConstraints []ModCompat `json:"OrderConstraints"`
}

type ModCompat struct {
	ModID    string         `json:"ModID"`
	Versions []string       `json:"Versions"`
	Order    ModCompatOrder `json:"Order"`
}

type ModCompatOrder string

const (
	Before ModCompatOrder = "Before"
	After  ModCompatOrder = "After"
)

type ModFile struct {
	IsDir bool   `json:"IsDir"`
	From  string `json:"From"`
	To    string `json:"To"`
}

type Configuration struct {
	Name        string   `json:"Name"`
	Description string   `json:"Description"`
	Choices     []Choice `json:"Choices"`
}

type Choice struct {
	Description           string    `json:"Description"`
	Preview               string    `json:"Preview"`
	Files                 []ModFile `json:"Files"`
	NextConfigurationName *string   `json:"NextConfigurationName"`
}
