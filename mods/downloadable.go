package mods

type Download struct {
	Name    string `json:"Name" xml:"Name"`
	Version string `json:"Version" xml:"Version"`

	Hosted *HostedDownloadable `json:"Hosted,omitempty" xml:"Hosted,omitempty"`
	Nexus  *NexusDownloadable  `json:"Nexus,omitempty" xml:"Nexus,omitempty"`

	DownloadedArchiveLocation string `json:"DownloadedLoc,omitempty" xml:"DownloadedLoc,omitempty"`
	//InstallType   InstallType `json:"InstallType" xml:"InstallType"`
}

type HostedDownloadable struct {
	Sources []string `json:"Source" xml:"Sources"`
}

type NexusDownloadable struct {
	FileID   int    `json:"FileID"`
	FileName string `json:"FileName"`
}
