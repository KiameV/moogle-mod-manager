package mods

type Override struct {
	NexusModID       string            `json:"id" xml:"id"`
	Preview          *Preview          `json:"preview" xml:"preview"`
	Description      *string           `json:"description,omitempty" xml:"description,omitempty"`
	ModCompatibility *ModCompatibility `json:"Compatibility,omitempty" xml:"ModCompatibility,omitempty"`
	DonationLinks    []*DonationLink   `json:"DonationLink" xml:"DonationLinks"`
	AlwaysDownload   []*DownloadFiles  `json:"AlwaysDownload,omitempty" xml:"AlwaysDownload,omitempty"`
	Configurations   *[]*Configuration `json:"configurations,omitempty" xml:"configurations,omitempty"`
}

func (o Override) Override(m *Mod) {
	if o.Preview != nil {
		m.Preview = o.Preview
	}
	if o.Description != nil {
		m.Description = *o.Description
	}
	if o.Configurations != nil {
		m.Configurations = *o.Configurations
	}
}
