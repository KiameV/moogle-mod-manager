package confirm

import (
	"github.com/kiamev/moogle-mod-manager/mods"
)

type cfConfirmer struct{}

func newCfConfirmer(_ Params) Confirmer { return nil }

func (_ *cfConfirmer) ConfirmDownload(done func(result mods.Result)) error {
	done(mods.Ok)
	return nil
}
