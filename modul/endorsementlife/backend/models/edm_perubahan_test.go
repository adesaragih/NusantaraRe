package models_test

import (
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
)

func TestAdaPerubahan(t *testing.T) {
	for _, k := range []struct {
		edmType        string
		ditandai, baru int
		mau            bool
	}{
		{models.EdmTypePerubahanData, 0, 0, false},
		{models.EdmTypePerubahanData, 1, 0, true},
		{models.EdmTypePerubahanData, 0, 1, true},
		{models.EdmTypePerubahanData, 2, 3, true},
		{models.EdmTypeBatal, 0, 0, true},
	} {
		if got := models.AdaPerubahan(k.edmType, k.ditandai, k.baru); got != k.mau {
			t.Errorf("AdaPerubahan(%q, %d, %d) = %v", k.edmType, k.ditandai, k.baru, got)
		}
	}
	if (models.PilihanHapus{}).AdaPilihan() || !(models.PilihanHapus{Semua: true}).AdaPilihan() ||
		!(models.PilihanHapus{Pilih: []string{"UJI-D1"}}).AdaPilihan() {
		t.Error("AdaPilihan")
	}
	if models.PesanTanpaPerubahan != "No changes to process. Mark at least one participant for Delete or add participants from CSV." {
		t.Error("pesan berubah")
	}
}
