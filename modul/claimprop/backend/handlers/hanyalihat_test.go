package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/claimprop/backend/handlers"
	"nusantarare/modul/claimprop/backend/models"
)

// Keputusan work owner 09-10-2026: "View more details" layar Komite Claim Prop membuka tampilan klaim Claim Prop READ
// ONLY - juga bagi pemegang assignment klaim itu (`GET /kasus/{id}?lihat=1`: layar terkunci, tanpa pra-proses).
func TestBukaKasusHanyaLihatTerkunciBagiPemegang(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	kode, out := u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka biasa")
	if out["bolehKerja"] != true {
		t.Fatalf("pembuat = pemegang Outstanding Claim: bolehKerja true: %v", out["bolehKerja"])
	}
	kode, out = u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id+"?lihat=1", admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka hanya-lihat")
	if out["bolehKerja"] != false {
		t.Fatalf("?lihat=1: bolehKerja false walau pemegang: %v", out["bolehKerja"])
	}
	b, _ := json.Marshal(out["tata"])
	var ts []models.Tata
	if err := json.Unmarshal(b, &ts); err != nil || len(ts) == 0 {
		t.Fatalf("tata: %v", err)
	}
	terbuka := cari(ts, func(x models.Tata) bool {
		return (x.Jenis == models.JenisTombol && !x.Nonaktif) || (x.Jenis == models.JenisMedan && !x.HanyaBaca)
	})
	if terbuka != nil {
		t.Fatalf("?lihat=1: tombol / medan masih terbuka: %+v", *terbuka)
	}
}
