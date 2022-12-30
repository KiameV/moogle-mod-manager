package mods

import "strings"

type (
	Kind       string
	NexusModID int
	CfModID    int
)

const (
	Nexus        Kind = "Nexus"
	CurseForge   Kind = "CurseForge"
	HostedAt     Kind = "HostedAt"
	HostedGitHub Kind = "GitHub"
)

type (
	GitHub struct {
		Owner   string `json:"Owner"`
		Repo    string `json:"Repo"`
		Version string `json:"Version"`
	}
	Kinds   []Kind
	ModKind struct {
		Kinds        Kinds       `json:"Kinds" xml:"Kinds"`
		NexusID      *NexusModID `json:"NexusID,omitempty" xml:"NexusID,omitempty"`
		CurseForgeID *CfModID    `json:"CurseForgeID,omitempty" xml:"CurseForgeID,omitempty"`
		GitHub       *GitHub     `json:"Github,omitempty" xml:"Github,omitempty"`
	}
)

var SubKinds = []string{
	string(HostedAt),
	string(HostedGitHub),
}

func (k Kind) Is(kind Kind) bool {
	return k == kind
}

func (k *Kinds) Is(kind Kind) bool {
	for _, i := range *k {
		if i == kind {
			return true
		}
	}
	return false
}

func (k *Kinds) Add(kind Kind) {
	if !k.Is(kind) {
		*k = append(*k, kind)
	}
}

func (k *Kinds) Remove(kind Kind) {
	for i, v := range *k {
		if v == kind {
			*k = append((*k)[:i], (*k)[i+1:]...)
			break
		}
	}
}

func (k *Kinds) IsHosted() bool {
	return k.Is(HostedAt) || k.Is(HostedGitHub)
}

func (k *Kinds) String() string {
	sb := strings.Builder{}
	for i, v := range *k {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(string(v))
	}
	return sb.String()
}
