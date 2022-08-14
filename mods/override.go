package mods

type Override struct {
	NexusModID     string            `json:"id" xml:"id"`
	Description    *string           `json:"description,omitempty" xml:"description,omitempty"`
	ReleaseNotes   *string           `json:"releaseNotes,omitempty" xml:"releaseNotes,omitempty"`
	Configurations *[]*Configuration `json:"configurations,omitempty" xml:"configurations,omitempty"`
}

func (o Override) Override(m *Mod) {
	if o.Description != nil {
		m.Description = *o.Description
	}
	if o.ReleaseNotes != nil {
		m.ReleaseNotes = *o.ReleaseNotes
	}
	if o.Configurations != nil {
		m.Configurations = *o.Configurations
	}
}
