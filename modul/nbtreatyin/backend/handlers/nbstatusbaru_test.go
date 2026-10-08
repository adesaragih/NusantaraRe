package handlers_test

// Uji seam 1 - berkas baru ber-NBStatus "New" (`[keputusan work owner 06-10-2026]` "begitu create buat aja
// New"; XML: connector Start1 tidak menulis NBStatus, kolom Status portal kosong sampai Submit).

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestBerkasBaruBerstatusNew(t *testing.T) {
	u := baru(t)
	id := u.buat()
	if got := u.g.Halaman[id].Ambil("NBStatus"); got != models.NBStatusBaru {
		t.Errorf("NBStatus tersimpan = %q, harap %q", got, models.NBStatusBaru)
	}
	kode, isi := u.panggil("GET", "/kasus", admin, nil)
	if kode != http.StatusOK {
		t.Fatalf("daftar: %d %s", kode, isi)
	}
	var baris []models.RingkasanKasus
	if err := json.Unmarshal([]byte(isi), &baris); err != nil {
		t.Fatal(err)
	}
	if len(baris) != 1 || baris[0].NBStatus != "New" {
		t.Fatalf("kolom Status portal: %+v", baris)
	}
}
