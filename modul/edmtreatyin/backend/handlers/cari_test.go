package handlers_test

// Kotak saring portal (perintah work owner 07-10-2026 "pencarian ... buat bisa mencari nomor nb/edm, insured name dll,
// intinya buat searchnya itu sangat berguna"): nomor kasus, nomor polis, insured name, ... - tanpa beda huruf, setiap
// kata cocok dengan salah satu kolom. Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"nusantarare/modul/edmtreatyin/backend/models"
)

func (u *uji) cari(kata string) []models.RingkasanKasus {
	u.t.Helper()
	kode, isi := u.panggil("GET", "/kasus?cari="+url.QueryEscape(kata), admin, nil)
	u.wajib(kode, isi, http.StatusOK, "cari "+kata)
	var out []models.RingkasanKasus
	if err := json.Unmarshal([]byte(isi), &out); err != nil {
		u.t.Fatal(err)
	}
	return out
}

func TestCariPortalNomorKasusPolisInsured(t *testing.T) {
	u := baru(t)
	k := u.buat("1")
	u.pilihBisnis(k.ID) // CopyGeneralDataEDM: InsuredName dari generasi lama
	for _, kata := range []string{k.ID, "edmt-", "uji-pol-0001", "satu tertanggung", "  UJI   sob ", "admin"} {
		if d := u.cari(kata); len(d) != 1 || d[0].ID != k.ID {
			t.Errorf("cari %q: %+v", kata, d)
		}
	}
	for _, kata := range []string{"tertanggung dua", "UJI-POL-9999", "100%"} {
		if d := u.cari(kata); len(d) != 0 {
			t.Errorf("cari %q tidak boleh menemukan: %+v", kata, d)
		}
	}
}
