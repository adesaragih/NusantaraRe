package tiruan_test

import (
	"nusantarare/modul/komiteclaimfacin/backend/services"
	"nusantarare/modul/komiteclaimfacin/backend/tiruan"
)

var (
	_ services.Gudang = (*tiruan.Gudang)(nil)
	_ services.Acuan  = (*tiruan.Acuan)(nil)
)
