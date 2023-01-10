package archive

import "github.com/kiamev/moogle-mod-manager/config"

type (
	Archive           string
	ArchiveFileBackup struct {
		Relative string
	}
	ArchiveFileBackups struct {
		Archives map[Archive][]ArchiveFileBackup
	}
	BackupFinder struct {
		Backups map[config.GameID]ArchiveFileBackups
	}
)

func NewBackupFinder() *BackupFinder {
	return &BackupFinder{
		Backups: make(map[config.GameID]ArchiveFileBackups),
	}
}

func (b *BackupFinder) Find(gameID config.GameID, relative string) (afb ArchiveFileBackup, archive Archive, found bool) {
	var (
		a ArchiveFileBackups
		f []ArchiveFileBackup
	)
	if a, found = b.Backups[gameID]; found {
		if f, found = a.Archives[archive]; found {
			found = false
			for _, i := range f {
				if i.Relative == relative {
					afb = i
					found = true
					break
				}
			}
		}
	}
	return
}
