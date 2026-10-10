package tiruan

import "nusantarare/modul/claimnonprop/backend/services"

var (
	_ services.Gudang = (*Gudang)(nil)
	_ services.Acuan  = (*Acuan)(nil)
)
