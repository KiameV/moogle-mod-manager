package secret

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"net/url"
)

const (
	nexusVortexApiAccessUrl = "https://www.nexusmods.com/users/myaccount?tab=api+access"
	cfApiKeyAccessUrl       = "https://console.curseforge.com/#/api-keys"
)

func Show(w fyne.Window) {
	var (
		nwe    = widget.NewPasswordEntry()
		cfwe   = widget.NewPasswordEntry()
		nu, _  = url.Parse(nexusVortexApiAccessUrl)
		cfu, _ = url.Parse(cfApiKeyAccessUrl)
		sct    = config.GetSecrets()
	)
	nwe.Bind(binding.BindString(&sct.NexusApiKey))
	cfwe.Bind(binding.BindString(&sct.CfApiKey))
	d := dialog.NewCustomConfirm("Secrets", "Save", "Cancel", container.NewVBox(
		widget.NewForm(widget.NewFormItem("Nexus Vortex Api Key", nwe)),
		widget.NewLabel("To get a key, follow this link and select [REQUEST AN API KEY] for Vortex. Copy what's generated."),
		widget.NewHyperlink(nexusVortexApiAccessUrl, nu),
		widget.NewForm(widget.NewFormItem("CurseForge Api Key", cfwe)),
		widget.NewLabel("To get a key, follow this link to generate one."),
		widget.NewHyperlink(cfApiKeyAccessUrl, cfu)),
		func(ok bool) {
			if ok {
				sct.NexusApiKey = nwe.Text
				sct.CfApiKey = cfwe.Text
				_ = sct.Save()
			}
		}, w)
	d.Resize(fyne.NewSize(800, 400))
	d.Show()
}
