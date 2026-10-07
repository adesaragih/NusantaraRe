package handlers_test

// Buka berkas HANYA-LIHAT (laporan work owner 07-10-2026: "saat admin view yang sudah di aksep ada yg ga
// ketarik"): pelaku yang bukan pemegang posisi berkas, dan berkas yang sudah tertutup, tetap melihat halaman
// master `TreatyIn` yang dibaca ulang - Commencement, Termination, subsection NonProp - tanpa satu tulisan pun.
// Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestBukaHanyaLihatMemuatMaster(t *testing.T) {
	u := baru(t)
	kontrakNP(u)
	id := u.buat()
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D-NP"}); kode != http.StatusOK {
		t.Fatalf("pilih bisnis: %d %s", kode, isi)
	}
	periksa := func(apa string, p pelakuUji) {
		t.Helper()
		sebelum, _ := json.Marshal(u.g.Halaman[id])
		kode, isi := u.panggil("GET", "/kasus/"+id, p, nil)
		if kode != http.StatusOK {
			t.Fatalf("%s: %d %s", apa, kode, isi)
		}
		ly := u.layar(isi)
		if ly.BolehKerja {
			t.Fatalf("%s: harus hanya-lihat", apa)
		}
		for j, harap := range map[string]string{
			"TreatyIn.Commencement": "2026-10-01 00:00:00",
			"TreatyIn.Termination":  "2027-09-30 00:00:00",
			"TreatyIn.RNMShare":     "10",
		} {
			if v := ly.Halaman.Ambil(j); v != harap {
				t.Errorf("%s: %s = %q, harap %q", apa, j, v, harap)
			}
		}
		if n := len(ly.Halaman.AmbilDaftar("TreatyIn.LimitShareSummaryList")); n != 1 {
			t.Errorf("%s: grid Share master %d baris, harap 1", apa, n)
		}
		if sesudah, _ := json.Marshal(u.g.Halaman[id]); string(sesudah) != string(sebelum) {
			t.Errorf("%s: buka hanya-lihat mengubah halaman tersimpan", apa)
		}
	}
	// berkas di posisi Admin dibuka Sec Head: bukan pemegang posisinya
	periksa("bukan pemegang", secHead)
	// berkas tertutup dibuka admin pembuatnya
	k := u.g.Kasus[id]
	k.StatusWork = models.StatusSelesai
	u.g.Kasus[id] = k
	periksa("tertutup", admin)
}
