package mods

type ModCompat struct {
	Kind       Kind             `json:"ModKind" xml:"ModKind"`
	Versions   []string         `json:"Version,omitempty" xml:"Versions,omitempty"`
	Hosted     *ModCompatHosted `json:"Hosted,omitempty" xml:"Hosted,omitempty"`
	Nexus      *ModCompatNexus  `json:"Nexus,omitempty" xml:"Nexus,omitempty"`
	CurseForge *ModCompatCF     `json:"CurseForge,omitempty" xml:"CurseForge,omitempty"`
	Order      *ModCompatOrder  `json:"Order,omitempty" xml:"Order,omitempty"`
	//displayName string           `json:"-" xml:"-"`
}

type ModCompatHosted struct {
	ModID ModID `json:"ModID" xml:"ModID"`
}

type ModCompatNexus struct {
	ModID ModID `json:"NexusModID" xml:"NexusModID"`
}

type ModCompatCF struct {
	ModID ModID `json:"CfModID" xml:"CfModID"`
}

type ModCompatOrder string

const (
	None   ModCompatOrder = ""
	Before ModCompatOrder = "Before"
	After  ModCompatOrder = "After"
)

var ModCompatOrders = []string{string(None), string(Before), string(After)}

func (c *ModCompat) ModID() ModID {
	switch c.Kind {
	case Hosted:
		return c.Hosted.ModID
	case Nexus:
		return c.Nexus.ModID
	case CurseForge:
		return c.CurseForge.ModID
	}
	if c.Kind == Hosted && c.Hosted != nil {

	} else if c.Kind == Nexus && c.Nexus != nil {
		return c.Nexus.ModID
	} else if c.Kind == CurseForge && c.CurseForge != nil {
		return c.CurseForge.ModID
	}
	return ""
}
