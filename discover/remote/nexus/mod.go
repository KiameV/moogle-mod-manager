package nexus

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"time"
)

type nexusMod struct {
	ModID           int              `json:"mod_id"`
	Name            string           `json:"name"`
	Summary         string           `json:"summary"`
	Description     string           `json:"description"`
	PictureUrl      string           `json:"picture_url"`
	CreatedTime     time.Time        `json:"created_time"`
	UpdatedTime     time.Time        `json:"updated_time"`
	Version         string           `json:"version"`
	GamePath        config.NexusPath `json:"domain_name"`
	CategoryID      int              `json:"category_id"`
	Author          string           `json:"author"`
	AuthorLink      string           `json:"author_link"`
	HasAdultContent bool             `json:"contains_adult_content"`
	Available       bool             `json:"available"`
	Link            string           `json:"-"`
}

type NexusFile struct {
	FileID      int    `json:"file_id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	IsPrimary   bool   `json:"is_primary"`
	FileName    string `json:"file_name"`
	ModVersion  string `json:"mod_version"`
	Description string `json:"description"`
}

type fileParent struct {
	Files []NexusFile `json:"files"`
}

/*
{
   "files":[
      {
         "id":[
            232,
            4335
         ],
         "uid":18618683228392,
         "file_id":232,
         "name":"FFVIPR_SNES_Battle_Backgrounds_metalliguy",
         "version":"1.1",
         "category_id":1,
         "category_name":"MAIN",
         "is_primary":false,
         "size":22333,
         "file_name":"FFVIPR_SNES_Battle_Backgrounds_metalliguy-48-1-1-1657022362.rar",
         "uploaded_timestamp":1657022362,
         "uploaded_time":"2022-07-05T11:59:22.000+00:00",
         "mod_version":"1.1",
         "external_virus_scan_url":"https://www.virustotal.com/gui/file/b396a748611f317853f0b89da7ad4398216d943f73ff19739206e72430d1a02f/detection/f-b396a748611f317853f0b89da7ad4398216d943f73ff19739206e72430d1a02f-1657022438",
         "description":"Updated to support patch 1.0.6",
         "size_kb":22333,
         "size_in_bytes":22868776,
         "changelog_html":"Updated to support patch 1.0.6",
         "content_preview_link":"https://file-metadata.nexusmods.com/file/nexus-files-s3-meta/4335/48/FFVIPR_SNES_Battle_Backgrounds_metalliguy-48-1-1-1657022362.rar.json"
      }
   ],
   "file_updates":[
      {
         "old_file_id":168,
         "new_file_id":232,
         "old_file_name":"FFVIPR_SNES_Battle_Backgrounds_metalliguy-48-1-0-1651437383.rar",
         "new_file_name":"FFVIPR_SNES_Battle_Backgrounds_metalliguy-48-1-1-1657022362.rar",
         "uploaded_timestamp":1657022362,
         "uploaded_time":"2022-07-05T11:59:22.000+00:00"
      }
   ]
}  ],
}
*/
