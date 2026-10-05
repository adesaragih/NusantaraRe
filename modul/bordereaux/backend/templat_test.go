package backend

import (
	"testing"

	"nusantarare/inti/backend/templat"
)

// 29 slot templat lolos katalog Template Manager: berkas bawaan (header saja) = jumlah kolom yang dibaca Upload CSV.
func TestSlotTemplatBordereaux(t *testing.T) {
	slot := SlotTemplat()
	if len(slot) != 29 {
		t.Fatalf("slot %d", len(slot))
	}
	if _, err := templat.NewKatalog(slot); err != nil {
		t.Fatal(err)
	}
	for _, s := range slot {
		if s.Menu != Nama || s.Kode[:len(AwalanTemplat)] != AwalanTemplat {
			t.Errorf("slot %s menu %s", s.Kode, s.Menu)
		}
	}
}
