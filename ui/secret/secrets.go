package secret

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config/secrets"
	"github.com/kiamev/moogle-mod-manager/ui/util"
)

const (
	nexusVortexApiAccessUrl = "https://www.nexusmods.com/users/myaccount?tab=api+access"
	cfApiKeyAccessUrl       = "https://console.curseforge.com/#/api-keys"
)

func Show(w fyne.Window) {
	var (
		nwe  = widget.NewPasswordEntry()
		cfwe = widget.NewPasswordEntry()
		n    = secrets.Get(secrets.NexusApiKey)
		cf   = secrets.Get(secrets.CfApiKey)
	)
	nwe.Bind(binding.BindString(&n))
	cfwe.Bind(binding.BindString(&cf))
	d := dialog.NewCustomConfirm("Secrets", "Save", "Cancel", container.NewVBox(
		widget.NewForm(widget.NewFormItem("Nexus Vortex Api Key", nwe)),
		widget.NewLabel("To get a key, follow this link and select [REQUEST AN API KEY] for Vortex. Copy what's generated."),
		util.CreateUrlRow(nexusVortexApiAccessUrl),
		widget.NewForm(widget.NewFormItem("CurseForge Api Key", cfwe)),
		widget.NewLabel("To get a key, follow this link to generate one."),
		util.CreateUrlRow(cfApiKeyAccessUrl)),
		func(ok bool) {
			if ok {
				secrets.Set(secrets.NexusApiKey, n)
				secrets.Set(secrets.CfApiKey, cf)
				_ = secrets.Save()
			}
		}, w)
	d.Resize(fyne.NewSize(800, 400))
	d.Show()
}
