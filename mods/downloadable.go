package mods

type Download struct {
	Name    string `json:"Name" xml:"Name"`
	Version string `json:"Version" xml:"Version"`
	//InstallType   InstallType `json:"InstallType" xml:"InstallType"`
	Hosted        *HostedDownloadable `json:"Hosted,omitempty" xml:"Hosted,omitempty"`
	Nexus         *NexusDownloadable  `json:"Nexus,omitempty" xml:"Nexus,omitempty"`
	DownloadedLoc string              `json:"DownloadedLoc,omitempty" xml:"DownloadedLoc,omitempty"`
}

type HostedDownloadable struct {
	Sources []string `json:"Source" xml:"Sources"`
}

type NexusDownloadable struct {
	FileID   int    `json:"FileID"`
	FileName string `json:"FileName"`
}
