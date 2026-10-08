package handlers_test

// Uji seam 1 - kolom User Create dan Tanggal Create portal (`[keputusan work owner 06-10-2026]`): nama pembuat
// = T_WORK_POLIS.CREATE_OP_NAME (nama tampilan saat Create), tanggal = T_WORK_POLIS.TGL_CREATE. Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestPortalMenampilkanPembuatDanTanggalCreate(t *testing.T) {
	u := baru(t)
	u.g.Nama[admin.akun] = "UJI Pembuat"
	u.buat()
	kode, isi := u.panggil("GET", "/kasus", admin, nil)
	if kode != http.StatusOK {
		t.Fatalf("daftar: %d %s", kode, isi)
	}
	var baris []models.RingkasanKasus
	if err := json.Unmarshal([]byte(isi), &baris); err != nil {
		t.Fatal(err)
	}
	if len(baris) != 1 || baris[0].NamaPembuat != "UJI Pembuat" || baris[0].TglCreate == "" {
		t.Fatalf("User Create / Tanggal Create portal: %+v", baris)
	}
}

// Kolom Type portal (keputusan work owner 06-10-2026): T_POLIS_QUOTATION.PROPORTIONAL_TYPE (diisi Choose Business).
func TestPortalMenampilkanJenisProporsi(t *testing.T) {
	u := baru(t)
	id := u.buat()
	u.g.Halaman[id].Setel(models.HalamanQuotation+".ProportionalType", models.JenisNonProporsional)
	_, isi := u.panggil("GET", "/kasus", admin, nil)
	var baris []models.RingkasanKasus
	if err := json.Unmarshal([]byte(isi), &baris); err != nil {
		t.Fatal(err)
	}
	if len(baris) != 1 || baris[0].JenisProporsi != models.JenisNonProporsional {
		t.Fatalf("kolom Type portal: %+v", baris)
	}
}
