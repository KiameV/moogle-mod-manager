package mods

type Kind string

const (
	Hosted Kind = "Hosted"
	Nexus  Kind = "Nexus"
)

func ToKind(kind Kind) *Kind {
	k := kind
	return &k
}

var Kinds = []string{string(Nexus), string(Hosted)}

type ModKind struct {
	Kind   Kind           `json:"Kind" xml:"Kind"`
	Hosted *HostedModKind `json:"Hosted,omitempty" xml:"Hosted,omitempty"`
	Nexus  *NexusModKind  `json:"Nexus,omitempty" xml:"Nexus,omitempty"`
}

type HostedModKind struct {
	ModFileLinks []string `json:"ModFileLink" xml:"ModFileLink"`
}

type NexusModKind struct {
	ID string `json:"ID" xml:"ID"`
}
