package mods

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"path/filepath"
)

type (
	ArchiveLocation string
	Download        struct {
		Name    string `json:"Name" xml:"Name"`
		Version string `json:"Version" xml:"Version"`

		Hosted      *HostedDownloadable      `json:"Hosted,omitempty" xml:"Hosted,omitempty"`
		Nexus       *NexusDownloadable       `json:"Nexus,omitempty" xml:"Nexus,omitempty"`
		CurseForge  *CurseForgeDownloadable  `json:"CurseForge,omitempty" xml:"CurseForge,omitempty"`
		GoogleDrive *GoogleDriveDownloadable `json:"GoogleDrive,omitempty" xml:"GoogleDrive,omitempty"`

		DownloadedArchiveLocation *ArchiveLocation `json:"DownloadedLoc,omitempty" xml:"DownloadedLoc,omitempty"`
		//InstallType   InstallType `json:"InstallType" xml:"InstallType"`
	}
	HostedDownloadable struct {
		Sources []string `json:"Source" xml:"Sources"`
	}
	NexusDownloadable struct {
		FileID   int    `json:"FileID"`
		FileName string `json:"FileName"`
	}
	CurseForgeDownloadable struct {
		FileID   int    `json:"FileID"`
		FileName string `json:"FileName"`
		Url      string `json:"Url"`
	}
	GoogleDriveDownloadable struct {
		Name string `json:"Name" xml:"Name"`
		Url  string `json:"Url" xml:"Url"`
	}
)

func (d Download) FileName() (string, error) {
	if d.Nexus != nil {
		return d.Nexus.FileName, nil
	} else if d.CurseForge != nil {
		return d.CurseForge.FileName, nil
	} else if d.GoogleDrive != nil {
		return d.GoogleDrive.Name, nil
	}
	return "", fmt.Errorf("no file name specified for %s", d.Name)
}

func (l *ArchiveLocation) ExtractDir(fileName string) string {
	s := config.PWD
	if l != nil {
		s = filepath.Dir(string(*l))
	}
	return filepath.Join(s, "extracted", fileName)
}
