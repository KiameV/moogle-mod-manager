package confirm

import (
	"github.com/kiamev/moogle-mod-manager/mods"
)

type bypassConfirmer struct {
	Params
}

func newBypassConfirmer(params Params) Confirmer {
	return &bypassConfirmer{Params: params}
}

func (_ *bypassConfirmer) Downloads(done func(mods.Result)) (err error) {
	done(mods.Ok)
	return nil
}
