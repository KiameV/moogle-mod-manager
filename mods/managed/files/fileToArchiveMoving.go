package files

type ArchiveMover interface {
}

func NewArchiveMover() ArchiveMover {
	return &archiveMover{}
}

type archiveMover struct{}

