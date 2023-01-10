package mods

type DownloadFiles struct {
	DownloadName string `json:"DownloadName" xml:"DownloadName"`
	// IsInstallAll is used by nexus mods when a mod.xml is not used
	Files []*ModFile `json:"File,omitempty" xml:"Files,omitempty"`
	Dirs  []*ModDir  `json:"Dir,omitempty" xml:"Dirs,omitempty"`
}

func (f *DownloadFiles) IsEmpty() bool {
	return len(f.Files) == 0 && len(f.Dirs) == 0
}

func (f *DownloadFiles) HasArchive() []string {
	var s []string
	for _, file := range f.Files {
		if file.ToArchive == nil {
			s = append(s, file.From)
		}
	}
	for _, dir := range f.Dirs {
		if dir.ToArchive == nil {
			s = append(s, dir.From)
		}
	}
	return s
}
