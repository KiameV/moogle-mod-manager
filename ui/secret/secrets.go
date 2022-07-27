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

const nexusVortexApiAccessUrl = "https://www.nexusmods.com/users/myaccount?tab=api+access"

func Show(w fyne.Window) {
	var (
		pwe  = widget.NewPasswordEntry()
		u, _ = url.Parse(nexusVortexApiAccessUrl)
		sct  = config.GetSecrets()
	)
	pwe.Bind(binding.BindString(&sct.NexusApiKey))
	d := dialog.NewCustomConfirm("Secrets", "Save", "Cancel", container.NewVBox(
		widget.NewForm(widget.NewFormItem("Nexus Vortex Api Key", pwe)),
		widget.NewLabel("To get a key, follow this link and select [REQUEST AN API KEY] for Vortex. Copy what's generated."),
		widget.NewHyperlink(nexusVortexApiAccessUrl, u)),
		func(ok bool) {
			if ok {
				sct.NexusApiKey = pwe.Text
				_ = sct.Save()
			}
		}, w)
	d.Resize(fyne.NewSize(800, 400))
	d.Show()
}
