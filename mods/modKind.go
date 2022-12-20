package mods

type (
	Kind    string
	SubKind string
)

const (
	Hosted     Kind = "Hosted"
	Nexus      Kind = "Nexus"
	CurseForge Kind = "CurseForge"

	HostedBlank  SubKind = ""
	HostedAt     SubKind = "HostedAt"
	HostedGitHub SubKind = "GitHub"
)

type (
	GitHub struct {
		Owner string `json:"Owner"`
		Repo  string `json:"Repo"`
	}
	ModKind struct {
		Kind    Kind     `json:"Kind" xml:"Kind"`
		SubKind *SubKind `json:"SubKind,omitempty" xml:"SubKind,omitempty"`
		GitHub  *GitHub  `json:"Github,omitempty" xml:"Github,omitempty"`
	}
)

var SubKinds = []string{
	string(HostedAt),
	string(HostedGitHub),
}

func NewModKind(kind Kind, subKind SubKind) ModKind {
	var sk *SubKind
	if kind == Hosted {
		sk = &subKind
	}
	return ModKind{
		Kind:    kind,
		SubKind: sk,
	}
}

func (k Kind) Is(kind Kind) bool {
	return k == kind
}

func (sk *SubKind) Get() SubKind {
	if sk != nil {
		return *sk
	}
	return HostedBlank
}

func (sk *SubKind) Is(value SubKind) bool {
	if sk != nil {
		return *sk == value
	}
	return false
}
