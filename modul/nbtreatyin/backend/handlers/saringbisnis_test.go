package handlers_test

// Keputusan work owner 06-10-2026: saringan kolom popup Choose Business dicari di SERVER - dipakai di query
// SEBELUM batas 500 baris (RD `BrowseTreatyJoinEDM` pyMaxRecords 500, urut TREATYID ASC). Tanpa itu kontrak di
// luar 500 baris pertama (mis. kontrak terbaru, TREATYID terbesar) tidak pernah dapat dipilih. Saringan = teks
// "memuat", tidak membedakan huruf besar/kecil; kolom di luar daftar kolom popup diabaikan.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

// kontrakBanyak - 600 baris Proportional UJI-T0000..UJI-T0599, lalu satu kontrak terbaru UJI-T9999.
func kontrakBanyak(u *uji) {
	for i := 0; i < 600; i++ {
		id := fmt.Sprintf("UJI-D%04d", i)
		u.g.Kontrak[id] = models.BarisKontrak{"ID": id, "TREATYID": fmt.Sprintf("UJI-T%04d", i), "PROPORTIONTYPE": "Proportional",
			"TREATYCONTRACTNAME": "UJI LAMA"}
	}
	u.g.Kontrak["UJI-D9999"] = models.BarisKontrak{"ID": "UJI-D9999", "TREATYID": "UJI-T9999", "PROPORTIONTYPE": "Proportional",
		"TREATYCONTRACTNAME": "UJI Kontrak Terbaru"}
}

func (u *uji) daftarBisnisSaring(id string, h *models.Halaman, saringan map[string]string) []models.BarisKontrak {
	u.t.Helper()
	kode, isi := u.panggil("POST", "/kasus/"+id+"/bisnis", admin, map[string]any{"halaman": h, "saringan": saringan})
	if kode != http.StatusOK {
		u.t.Fatalf("daftar bisnis: %d %s", kode, isi)
	}
	var out []models.BarisKontrak
	if err := json.Unmarshal([]byte(isi), &out); err != nil {
		u.t.Fatal(err)
	}
	return out
}

func halamanProporsional() *models.Halaman {
	h := models.HalamanBaru()
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", "Proportional")
	return h
}

func TestDaftarBisnisSaringanDiServerSebelumBatas500(t *testing.T) {
	u := baru(t)
	id := u.buat()
	kontrakBanyak(u)
	h := halamanProporsional()

	semua := u.daftarBisnisSaring(id, h, nil)
	if len(semua) != models.BatasDaftarBisnis {
		t.Fatalf("tanpa saringan %d baris, harap %d", len(semua), models.BatasDaftarBisnis)
	}
	for _, b := range semua {
		if b["ID"] == "UJI-D9999" {
			t.Fatal("prasyarat: kontrak terbaru seharusnya di luar 500 baris pertama")
		}
	}
	for nama, saring := range map[string]map[string]string{
		"Treaty Offer ID":           {"TREATYID": "9999"},
		"Contract Name huruf kecil": {"TREATYCONTRACTNAME": "kontrak terbaru"},
		"dua kolom":                 {"TREATYID": "UJI-T9", "TREATYCONTRACTNAME": "terbaru"},
	} {
		out := u.daftarBisnisSaring(id, h, saring)
		if len(out) != 1 || out[0]["ID"] != "UJI-D9999" {
			t.Errorf("%s: %d baris, harap tepat UJI-D9999", nama, len(out))
		}
	}
	if out := u.daftarBisnisSaring(id, h, map[string]string{"BUKAN_KOLOM": "x"}); len(out) != models.BatasDaftarBisnis {
		t.Errorf("kolom di luar popup harus diabaikan: %d baris", len(out))
	}
}

func TestPilihBisnisKontrakDiLuar500BarisPertamaDiterima(t *testing.T) {
	u := baru(t)
	id := u.buat()
	kontrakBanyak(u)
	h := halamanProporsional()
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D9999", "halaman": h}); kode != http.StatusOK {
		t.Fatalf("pilih kontrak terbaru: %d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.NoOffer"); got != "UJI-T9999" {
		t.Errorf("NoOffer = %q, harap UJI-T9999", got)
	}
}
