package handlers_test

// Uji jalur HTTP tiket 16, 17, 18, 19.

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// Tiket 17 — pencarian lewat nomor warisan.
func TestCariNomorWarisanMenjawab200(t *testing.T) {
	g := gudangTiruan{warisan: []models.Kontrak{{ID: 7, NomorKontrakWarisan: "TRI-2019-0007"}}}
	h := handlers.RouterDengan(services.LayananDengan(g), true, true)

	w := kirim(t, h, http.MethodGet, handlers.Prefix+"/kontrak/cari?nomorWarisan=TRI-2019-0007", "AKUN-UJI", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("mau 200, dapat %d (badan %s)", w.Code, w.Body.String())
	}
	var k []models.Kontrak
	if err := json.Unmarshal(w.Body.Bytes(), &k); err != nil {
		t.Fatalf("badan bukan daftar JSON: %v", err)
	}
	if len(k) != 1 || k[0].ID != 7 {
		t.Errorf("hasil tidak sesuai: %+v", k)
	}
}

// Hasil kosong terkirim sebagai `[]`, bukan `null` — pembaca JavaScript yang
// menerima null jatuh pada `.map`, dan jatuhnya di layar.
func TestCariNomorWarisanKosongTerkirimSebagaiLarik(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{warisan: []models.Kontrak{}}), true, true)

	w := kirim(t, h, http.MethodGet, handlers.Prefix+"/kontrak/cari?nomorWarisan=TIDAK-ADA", "AKUN-UJI", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("mau 200, dapat %d", w.Code)
	}
	if got := w.Body.String(); got != "[]\n" && got != "[]" {
		t.Fatalf("mau badan `[]`, dapat %q", got)
	}
}

// Nomor kosong → 422: permintaan yang belum lengkap, bukan "cari semua".
func TestCariNomorWarisanTanpaKataKunciMenjawab422(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	w := kirim(t, h, http.MethodGet, handlers.Prefix+"/kontrak/cari", "AKUN-UJI", nil)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("mau 422, dapat %d (badan %s)", w.Code, w.Body.String())
	}
}

// Tiket 18 — ruas beku ditolak 422, ruas lain diterima 204.
func TestUbahKontrakRuasBekuMenjawab422(t *testing.T) {
	g := gudangTiruan{kontrak: models.KontrakDenganVersi{Kontrak: models.Kontrak{
		ID: 7, IDCedant: 1, IDAsalBisnis: 2, SifatProporsi: models.SifatNonProporsional,
		TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31"}}}
	h := handlers.RouterDengan(services.LayananDengan(g), true, true)

	badan, _ := json.Marshal(services.MasukanUbahKontrak{
		IDCedant: 99, IDAsalBisnis: 2, SifatProporsi: models.SifatNonProporsional,
		TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31"})
	w := kirim(t, h, http.MethodPut, handlers.Prefix+"/kontrak/7", "AKUN-UJI", badan)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("mau 422, dapat %d (badan %s)", w.Code, w.Body.String())
	}
}

func TestUbahKontrakRuasBukanBekuMenjawab204(t *testing.T) {
	g := gudangTiruan{kontrak: models.KontrakDenganVersi{Kontrak: models.Kontrak{
		ID: 7, IDCedant: 1, IDAsalBisnis: 2, SifatProporsi: models.SifatNonProporsional,
		TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31"}}}
	h := handlers.RouterDengan(services.LayananDengan(g), true, true)

	badan, _ := json.Marshal(services.MasukanUbahKontrak{
		NomorKontrakWarisan: "TRI-BARU", IDCedant: 1, IDAsalBisnis: 2,
		SifatProporsi: models.SifatNonProporsional,
		TanggalMulai:  "2026-01-01", TanggalBerakhir: "2026-12-31"})
	w := kirim(t, h, http.MethodPut, handlers.Prefix+"/kontrak/7", "AKUN-UJI", badan)

	if w.Code != http.StatusNoContent {
		t.Fatalf("mau 204, dapat %d (badan %s)", w.Code, w.Body.String())
	}
}

// Tiket 19 — versi menyimpang 422, versi selaras 201.
func TestTambahVersiMenyimpangMenjawab422(t *testing.T) {
	g := gudangTiruan{kontrak: models.KontrakDenganVersi{Kontrak: models.Kontrak{
		ID: 7, IDCedant: 1, IDAsalBisnis: 2, SifatProporsi: models.SifatNonProporsional,
		TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31"}}}
	h := handlers.RouterDengan(services.LayananDengan(g), true, true)

	badan, _ := json.Marshal(services.MasukanVersiTambahan{
		NomorUrutVersi: 2, KeadaanSiklusHidup: "DRAFT", NamaKontrak: "V2",
		IDCedant: 99, IDAsalBisnis: 2, SifatProporsi: models.SifatNonProporsional,
		TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31"})
	w := kirim(t, h, http.MethodPost, handlers.Prefix+"/kontrak/7/versi", "AKUN-UJI", badan)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("mau 422, dapat %d (badan %s)", w.Code, w.Body.String())
	}
}

func TestTambahVersiSelarasMenjawab201(t *testing.T) {
	g := gudangTiruan{kontrak: models.KontrakDenganVersi{Kontrak: models.Kontrak{
		ID: 7, IDCedant: 1, IDAsalBisnis: 2, SifatProporsi: models.SifatNonProporsional,
		TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31"}}}
	h := handlers.RouterDengan(services.LayananDengan(g), true, true)

	badan, _ := json.Marshal(services.MasukanVersiTambahan{
		NomorUrutVersi: 2, KeadaanSiklusHidup: "DRAFT", NamaKontrak: "V2",
		IDCedant: 1, IDAsalBisnis: 2, SifatProporsi: models.SifatNonProporsional,
		TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31"})
	w := kirim(t, h, http.MethodPost, handlers.Prefix+"/kontrak/7/versi", "AKUN-UJI", badan)

	if w.Code != http.StatusCreated {
		t.Fatalf("mau 201, dapat %d (badan %s)", w.Code, w.Body.String())
	}
}

// Tiket 16 — peringatan menumpang jawaban 201, bukan menggagalkannya.
func TestPeringatanKunciAlamiTidakMenggagalkanPenyimpanan(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	w := kirim(t, h, http.MethodPost, handlers.Prefix+"/kontrak", "AKUN-UJI", badanSah())

	if w.Code != http.StatusCreated {
		t.Fatalf("mau 201, dapat %d", w.Code)
	}
	var h2 services.HasilBuatKontrak
	if err := json.Unmarshal(w.Body.Bytes(), &h2); err != nil {
		t.Fatalf("badan bukan JSON: %v", err)
	}
	if h2.IDKontrak == 0 {
		t.Error("pengenal tidak dikembalikan")
	}
}
