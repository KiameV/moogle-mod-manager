package mods

type ModCompatibility struct {
	Requires []*ModCompat `json:"Require,omitempty" xml:"Requires,omitempty"`
	Forbids  []*ModCompat `json:"Forbid,omitempty" xml:"Forbids,omitempty"`
	// OrderConstraints []ModCompat `json:"OrderConstraint"`
}

func (c *ModCompatibility) HasItems() bool {
	return c != nil && (len(c.Requires) > 0 || len(c.Forbids) > 0)
}
