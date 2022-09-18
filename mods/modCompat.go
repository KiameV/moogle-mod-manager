package mods

type ModCompat struct {
	Kind        Kind             `json:"ModKind" xml:"ModKind"`
	Versions    []string         `json:"Version,omitempty" xml:"Versions,omitempty"`
	Hosted      *ModCompatHosted `json:"Hosted,omitempty" xml:"Hosted,omitempty"`
	Nexus       *ModCompatNexus  `json:"Nexus,omitempty" xml:"Nexus,omitempty"`
	Order       *ModCompatOrder  `json:"Order,omitempty" xml:"Order,omitempty"`
	displayName string           `json:"-" xml:"-"`
}

type ModCompatHosted struct {
	ModID ModID `json:"ModID" xml:"ModID"`
}

type ModCompatNexus struct {
	ModID ModID `json:"NexusModID" xml:"NexusModID"`
}

type ModCompatOrder string

const (
	None   ModCompatOrder = ""
	Before ModCompatOrder = "Before"
	After  ModCompatOrder = "After"
)

var ModCompatOrders = []string{string(None), string(Before), string(After)}

func (c *ModCompat) ModID() ModID {
	if c.Kind == Hosted && c.Hosted != nil {
		return c.Hosted.ModID
	} else if c.Kind == Nexus && c.Nexus != nil {
		return c.Nexus.ModID
	}
	return ""
}
