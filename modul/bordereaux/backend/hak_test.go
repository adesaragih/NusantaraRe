package backend

import (
	"testing"

	"nusantarare/modul/bordereaux/backend/handlers"
)

// Akses menu LIHAT: KODE menu yang dibaca handler = nama modul (= KODE M_NAV_MENU), dan pola yang dibebaskan adalah
// rute modul ini (Upload CSV membaca, Submit dijaga layanan).
func TestHakLihatBordereaux(t *testing.T) {
	if handlers.KodeMenu != Nama {
		t.Errorf("KodeMenu %q != Nama %q", handlers.KodeMenu, Nama)
	}
	p := Pendaftaran()
	if p.HakLihat == nil || len(p.HakLihat.Bebas) != 2 ||
		p.HakLihat.Bebas[0] != "POST /api/bordereaux/unggah-csv" || p.HakLihat.Bebas[1] != "POST /api/bordereaux/berkas/{id}/submit" {
		t.Errorf("HakLihat %+v", p.HakLihat)
	}
}
