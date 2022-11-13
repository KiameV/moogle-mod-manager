package mods

type Download struct {
	Name    string `json:"Name" xml:"Name"`
	Version string `json:"Version" xml:"Version"`

	Hosted     *HostedDownloadable     `json:"Hosted,omitempty" xml:"Hosted,omitempty"`
	Nexus      *RemoteDownloadable     `json:"Nexus,omitempty" xml:"Nexus,omitempty"`
	CurseForge *CurseForgeDownloadable `json:"CurseForge,omitempty" xml:"CurseForge,omitempty"`

	DownloadedArchiveLocation *string `json:"DownloadedLoc,omitempty" xml:"DownloadedLoc,omitempty"`
	//InstallType   InstallType `json:"InstallType" xml:"InstallType"`
}

type HostedDownloadable struct {
	Sources []string `json:"Source" xml:"Sources"`
}

type RemoteDownloadable struct {
	FileID   int    `json:"FileID"`
	FileName string `json:"FileName"`
}

type CurseForgeDownloadable struct {
	RemoteDownloadable
	Url string `json:"Url"`
}
