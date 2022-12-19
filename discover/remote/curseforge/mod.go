package curseforge

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"time"
)

type cfModResponse struct {
	Data cfMod `json:"data"`
}

type cfMod struct {
	ModID        int             `json:"id"`
	Name         string          `json:"name"`
	Summary      string          `json:"summary"`
	Category     []category      `json:"categories"`
	Links        links           `json:"links"`
	Author       []author        `json:"authors"`
	Logo         logo            `json:"logo"`
	Screenshots  []screenshot    `json:"screenshots"`
	CreatedTime  time.Time       `json:"dateCreated"`
	UpdatedTime  time.Time       `json:"dateModified"`
	GameVersions []string        `json:"gameVersions"`
	Game         config.CfGameID `json:"gameId"`
}

func (m *cfMod) Version() string {
	return m.UpdatedTime.Format("2006.01.02.15.04.05")
}

type CfFile struct {
	FileID      int       `json:"id"`
	Name        string    `json:"displayName"`
	DownloadUrl string    `json:"downloadUrl"`
	FileDate    time.Time `json:"fileDate"`
}

func (f CfFile) Version() string {
	return f.FileDate.Format("2006.01.02.15.04.05")
}

type category struct {
	Name string `json:"name"`
}

func (c *category) toCategory() (r mods.Category, err error) {
	switch c.Name {
	case "Script/Text":
		r = mods.ScriptText
	case "Enemy Sprite":
		r = mods.EnemySprite
	case "Window Frames":
		r = mods.UiWindowFrames
	case "General":
		r = mods.General
	case "Tile Set":
		r = mods.TileSet
	case "Textbox Portraits":
		r = mods.UiTextBoxPortraits
	case "Game Overhauls":
		r = mods.GameOverhauls
	case "Player/NPC Sprite":
		r = mods.PlayerNpcSprites
	case "Title Screen":
		r = mods.TitleScreen
	case "Menu Portraits":
		r = mods.UiMenuPortraits
	case "Soundtrack":
		r = mods.Soundtrack
	case "Utility":
		r = mods.Utility
	case "Gameplay":
		r = mods.Gameplay
	case "Battle Scene":
		r = mods.BattleScene
	case "Fonts":
		r = mods.Fonts
	case "UI":
		r = mods.UIGeneral
	default:
		err = fmt.Errorf("unknown category: " + c.Name)
	}
	return
}

type links struct {
	WebsiteUrl string `json:"websiteUrl"`
}

type screenshot struct {
	ThumbnailUrl string `json:"thumbnailUrl"`
	Url          string `json:"url"`
}

type logo struct {
	ThumbnailUrl string `json:"thumbnailUrl"`
	Url          string `json:"url"`
}

type author struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type fileParent struct {
	Files []CfFile `json:"data"`
}

type description struct {
	Data string `json:"data"`
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
