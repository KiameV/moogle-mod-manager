package mods

type (
	//ModCompatOrder string
	ModCompat struct {
		Versions []string `json:"Version,omitempty" xml:"Versions,omitempty"`
		ID       ModID    `json:"ModID,omitempty" xml:"ModID,omitempty"`
		//displayName string           `json:"-" xml:"-"`
	}
)

//const (
//	None   ModCompatOrder = ""
//	Before ModCompatOrder = "Before"
//	After  ModCompatOrder = "After"
//)
//var ModCompatOrders = []string{string(None), string(Before), string(After)}

func (c *ModCompat) ModID() ModID {
	return c.ID
}
