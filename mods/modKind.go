package mods

type Kind string

const (
	Hosted Kind = "Hosted"
	Nexus  Kind = "Nexus"
)

var Kinds = []string{string(Hosted), string(Nexus)}

type ModKind struct {
	Kind   Kind           `json:"Kind" xml:"Kind"`
	Hosted *HostedModKind `json:"Hosted,omitempty" xml:"Hosted,omitempty"`
	Nexus  *NexusModKind  `json:"Nexus,omitempty" xml:"Nexus,omitempty"`
}

type HostedModKind struct {
	Version      string   `json:"Version" xml:"Version"`
	ModFileLinks []string `json:"ModFileLink" xml:"ModFileLink"`
}

type NexusModKind struct {
	ID string `json:"ID" xml:"ID"`
}
