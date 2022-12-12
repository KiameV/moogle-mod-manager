package mods

import (
	"fmt"
)

type Download struct {
	Name    string `json:"Name" xml:"Name"`
	Version string `json:"Version" xml:"Version"`

	Hosted     *HostedDownloadable     `json:"Hosted,omitempty" xml:"Hosted,omitempty"`
	Nexus      *RemoteDownloadable     `json:"Nexus,omitempty" xml:"Nexus,omitempty"`
	CurseForge *CurseForgeDownloadable `json:"CurseForge,omitempty" xml:"CurseForge,omitempty"`

	DownloadedArchiveLocation *string `json:"DownloadedLoc,omitempty" xml:"DownloadedLoc,omitempty"`
	//InstallType   InstallType `json:"InstallType" xml:"InstallType"`
}

func (d Download) FileName() (string, error) {
	if d.Nexus != nil {
		return d.Nexus.FileName, nil
	} else if d.CurseForge != nil {
		return d.CurseForge.FileName, nil
	}
	return "", fmt.Errorf("no file name specified for %s", d.Name)
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
