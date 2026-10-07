package handlers_test

// Uji seam 1 - Utility1 `SaveJsonPolisTreatyIn_Act` (keputusan work owner 06-10-2026: "JSON-nya tidak disimpan,
// tapi tetap insert kolom lainnya"): realisasi SELESAI menulis json_polis tanpa DATA_JSON, ACHIEVEMENT, dan
// TREATYINPRODUCTION di transaksi submit; jalur lain tidak menulis apa pun. Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
)

// sampaiDeptHead - kasus Proporsional berspreading dua baris, sudah di Dept Head bernomor polis.
func (u *uji) sampaiDeptHead() (string, string) {
	u.t.Helper()
	id := u.buat()
	if kode, isi := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusOK {
		u.t.Fatalf("admin: %d %s", kode, isi)
	}
	if kode, isi := u.kirim(id, secHead, putusan("1")); kode != http.StatusOK {
		u.t.Fatalf("Sec Head: %d %s", kode, isi)
	}
	h := u.g.Halaman[id]
	h.Setel("Quotation.ProportionalType", models.ProporsionalPenuh)
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", models.ProporsionalPenuh)
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{
		{"TreatyType": "UJI-JR1", "SharePercentage": "60", "ClaimPercentage": "60", "PremiumSpreaded": "600", "ClaimSpreaded": "0"},
		{"TreatyType": "UJI-JR2", "SharePercentage": "40", "ClaimPercentage": "40", "PremiumSpreaded": "400", "ClaimSpreaded": "0"},
	})
	var n services.NomorPolis
	kode, isi := u.panggil("POST", "/kasus/"+id+"/nomor-polis", deptHead, map[string]any{"halaman": putusan("1")})
	if kode != http.StatusOK || json.Unmarshal([]byte(isi), &n) != nil {
		u.t.Fatalf("nomor polis: %d %s", kode, isi)
	}
	return id, n.PolicyNo
}

func TestSelesaiMenulisJSONPolisCapaianProduksi(t *testing.T) {
	u := baru(t)
	id, nopol := u.sampaiDeptHead()
	if n := len(u.g.JSONPolis) + len(u.g.Capaian) + len(u.g.Produksi); n != 0 {
		t.Fatalf("sebelum selesai nol baris produksi, dapat %d", n)
	}
	if kode, isi := u.kirim(id, deptHead, putusan("1")); kode != http.StatusOK || u.g.Kasus[id].StatusWork != models.StatusSelesai {
		t.Fatalf("Dept Head menyetujui: %d %s", kode, isi)
	}
	kunci := models.KunciInstans(id)
	h := u.g.Halaman[id]

	if len(u.g.JSONPolis) != 1 {
		t.Fatalf("json_polis: %v", u.g.JSONPolis)
	}
	jp := u.g.JSONPolis[0]
	if jp["IDPEGA"] != kunci || jp["NOPOLIS"] != nopol || jp["PRODKE"] != "0" || jp["USERNAME"] != "UJI-DH" ||
		jp["TGL_PROD"] == "" || jp["TGL_PROD"] != h.Ambil("PolicyTreatyIn.ProductionDate") {
		t.Errorf("json_polis = %v", jp)
	}
	if _, ada := jp["DATA_JSON"]; ada {
		t.Error("DATA_JSON tidak boleh ditulis")
	}

	if len(u.g.Produksi) != 2 { // satu per baris spreading
		t.Fatalf("TREATYINPRODUCTION: %d baris", len(u.g.Produksi))
	}
	for i, jr := range []string{"UJI-JR1", "UJI-JR2"} {
		p := u.g.Produksi[i]
		if p["IDPEGA"] != kunci || p["NOPOLIS"] != nopol || p["JN_REAS"] != jr || p["PREMI_OGP"] != "1000" ||
			p["MARKETINGOFFICERCODE"] != "UJI-MO-1" || p["STATEMENT_DATE"] != "2026-10-01" || p["ENDDATE"] != "2026-12-31" ||
			p["LAYERTYPE"] != "0" || p["PROPORTIONALTYPE"] != models.ProporsionalPenuh || p["PROD_DATE"] != jp["TGL_PROD"] {
			t.Errorf("produksi %d = %v", i, p)
		}
	}

	if len(u.g.Capaian) != 1 {
		t.Fatalf("ACHIEVEMENT: %v", u.g.Capaian)
	}
	if c := u.g.Capaian[0]; c["IDPEGA"] != kunci || c["NOPOLIS"] != nopol || c["PREMIUM"] != "1000" || c["PROPORTIONALTYPE"] != models.ProporsionalPenuh {
		t.Errorf("ACHIEVEMENT = %v", c)
	}
	// halaman yang disimpan membawa medan Utility1 (IDNewBisnis = pyID)
	if h.Ambil("PolicyTreatyIn.IDNewBisnis") != id {
		t.Errorf("IDNewBisnis = %q", h.Ambil("PolicyTreatyIn.IDNewBisnis"))
	}
}

func TestBukanSelesaiTidakMenulisProduksi(t *testing.T) {
	u := baru(t)
	// admin menolak = berkas ditutup Ditolak (End3 tanpa Utility1)
	id := u.buat()
	if kode, isi := u.kirim(id, admin, halamanLengkap("0")); kode != http.StatusOK || u.g.Kasus[id].StatusWork != models.StatusDitolak {
		t.Fatalf("admin menolak: %d %s", kode, isi)
	}
	// Dept Head menolak = kembali ke admin
	id2, _ := u.sampaiDeptHead()
	if kode, isi := u.kirim(id2, deptHead, putusan("0")); kode != http.StatusOK || u.g.Kasus[id2].PositionNote != models.PosisiAdmin {
		t.Fatalf("Dept Head menolak: %d %s", kode, isi)
	}
	if n := len(u.g.JSONPolis) + len(u.g.Capaian) + len(u.g.Produksi); n != 0 {
		t.Fatalf("tanpa realisasi selesai nol baris produksi, dapat %d", n)
	}
}

func TestGagalTulisProduksiMembatalkanSubmit(t *testing.T) { // AC 29, 83 - satu transaksi
	u := baru(t)
	id, _ := u.sampaiDeptHead()
	u.g.GagalDi = "SimpanPolisProduksi"
	if kode, _ := u.kirim(id, deptHead, putusan("1")); kode == http.StatusOK {
		t.Fatal("galat tulis produksi harus menggagalkan submit")
	}
	if u.g.Kasus[id].PositionNote != models.PosisiDeptHead || u.g.Kasus[id].StatusWork == models.StatusSelesai {
		t.Fatalf("kasus harus tetap di Dept Head: %+v", u.g.Kasus[id])
	}
	if n := len(u.g.JSONPolis) + len(u.g.Capaian) + len(u.g.Produksi); n != 0 {
		t.Fatalf("baris produksi harus ikut batal, dapat %d", n)
	}
	u.g.GagalDi = ""
	if kode, isi := u.kirim(id, deptHead, putusan("1")); kode != http.StatusOK || len(u.g.JSONPolis) != 1 {
		t.Fatalf("submit ulang: %d %s", kode, isi)
	}
}
