package repository_test

import (
	"testing"

	"nusantarare/modul/treatyin/backend/repository"
)

// ⭐ Setiap kolom akar Share yang layar baca WAJIB ada di peta pendaratan
// `T_TREATY_REVISION` — pola `TestKolomRevisiAdaDiPeta`: kolom yang dibaca
// tetapi tak dipetakan tidak pernah terisi pemuat, dan nol galat menyebutnya.
func TestKolomShareRevisiAdaDiPeta(t *testing.T) {
	var kolom, kunci map[string]bool
	for _, p := range repository.PetaPendaratan {
		if p.Tabel != "T_TREATY_REVISION" {
			continue
		}
		kolom, kunci = map[string]bool{}, map[string]bool{}
		for i, k := range p.Kolom {
			kolom[k] = true
			kunci[k] = p.Kunci[i] != ""
		}
		for i, k := range p.Kunci {
			kunci[k] = i < len(p.Kolom)
		}
	}
	if kolom == nil {
		t.Fatal("T_TREATY_REVISION tidak ada di peta")
	}
	for _, k := range repository.KolomShareRevisi {
		if !kolom[k[0]] {
			t.Errorf("kolom %s dibaca layar tetapi tidak ada di peta", k[0])
		}
		if !kunci[k[1]] {
			t.Errorf("kunci %s dibaca layar tetapi tidak ada di peta", k[1])
		}
	}
}
