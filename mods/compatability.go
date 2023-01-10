package mods

type ModCompatibility struct {
	Requires []*ModCompat `json:"Require" xml:"Requires"`
	Forbids  []*ModCompat `json:"Forbid" xml:"Forbids"`
	//OrderConstraints []ModCompat `json:"OrderConstraint"`
}

func (c *ModCompatibility) HasItems() bool {
	return c != nil && (len(c.Requires) > 0 || len(c.Forbids) > 0)
}
