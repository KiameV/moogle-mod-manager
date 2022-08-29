package mods

type Kind string

const (
	Hosted Kind = "Hosted"
	Nexus  Kind = "Nexus"
)

type ModKind struct {
	Kind Kind `json:"Kind" xml:"Kind"`
}
