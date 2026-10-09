package tiruan

import (
	"testing"

	"nusantarare/modul/claimfacin/backend/repository"
	"nusantarare/modul/claimfacin/backend/services"
)

// Gudang / acuan tiruan dan Oracle memenuhi antarmuka services yang sama.
func TestTiruanMemenuhiAntarmuka(t *testing.T) {
	var _ services.Gudang = Baru()
	var _ services.Acuan = AcuanBaru()
	var _ services.Acuan = (*repository.Acuan)(nil)
}
