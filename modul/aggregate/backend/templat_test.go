package backend

import (
	"testing"

	"nusantarare/inti/backend/templat"
)

// Slot templat Aggregate lolos katalog: berkas bawaan 38 kolom = kolom yang dibaca Upload CSV.
func TestSlotTemplatAggregate(t *testing.T) {
	p := Pendaftaran()
	if _, err := templat.NewKatalog(p.Templat); err != nil {
		t.Fatal(err)
	}
	if len(p.Templat) != 1 || p.Templat[0].JumlahKolom != 38 || p.Templat[0].Kode != KodeTemplat || p.Templat[0].Menu != Nama {
		t.Errorf("slot %+v", p.Templat)
	}
}
