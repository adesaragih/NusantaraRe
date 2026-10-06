package handlers_test

// Uji rute tiket 32, 40, 41 — tanpa Oracle.

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// Memakai `kirim` dari `rute_kontrak_test.go` — satu pembangun permintaan
// untuk seluruh paket uji ini.
func badan(s string) []byte {
	if s == "" {
		return nil
	}
	return []byte(s)
}

// ------------------------------------------------------------------ tiket 32

func TestRutePemulihanMenolakPersenDiLuarRentang(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	w := kirim(t, h, http.MethodPut, "/api/treaty-in/layer/7/pemulihan", "AKUN-UJI", badan(`[{"nomorUrut":1,"persenPemulihan":"101"}]`))
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("kode %d, mau 422. Badan: %s", w.Code, w.Body.String())
	}
	// Penolakan yang tidak menyebut angkanya membuat pengirimnya menebak.
	if !strings.Contains(w.Body.String(), "101") {
		t.Errorf("badan tidak menyebut nilai yang ditolak: %s", w.Body.String())
	}
}

func TestRutePemulihanMenerimaTarifBerbeda(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	w := kirim(t, h, http.MethodPut, "/api/treaty-in/layer/7/pemulihan", "AKUN-UJI", badan(`[{"nomorUrut":1,"persenPemulihan":"100","persenTambahan":"0"},
		  {"nomorUrut":2,"persenPemulihan":"100","persenTambahan":"100"}]`))
	if w.Code != http.StatusNoContent {
		t.Fatalf("kode %d, mau 204. Badan: %s", w.Code, w.Body.String())
	}
}

func TestRutePemulihanMenolakPengenalLayerBukanAngka(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	w := kirim(t, h, http.MethodPut, "/api/treaty-in/layer/abc/pemulihan", "AKUN-UJI", badan(`[]`))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("kode %d, mau 400", w.Code)
	}
	// Kalimatnya menyebut layer, bukan kontrak.
	if !strings.Contains(w.Body.String(), "layer id") {
		t.Errorf("kalimat galat menyebut hal yang salah: %s", w.Body.String())
	}
}

// ------------------------------------------------------------------ tiket 40

// Versi pertama menjawab 200 dengan `"ada": false` — BUKAN 404.
//
// 404 berarti "sumber daya ini tidak ada"; yang benar di sini "pertanyaannya
// terjawab, dan jawabannya tidak ada versi sebelumnya".
func TestRuteVersiPertamaMenjawab200DenganAdaFalse(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{versiDasar: nil}), true, true)

	w := kirim(t, h, http.MethodGet, "/api/treaty-in/versi/5/sebelumnya", "AKUN-UJI", badan(""))
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d, mau 200. Badan: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"ada":false`) {
		t.Errorf("badan tidak menyatakan ketiadaan: %s", w.Body.String())
	}
}

func TestRuteVersiSebelumnyaMenjawabVersiDasar(t *testing.T) {
	g := gudangTiruan{versiDasar: &models.VersiKontrak{ID: 11, IDKontrak: 3, NamaKontrak: "Treaty 2025"}}
	h := handlers.RouterDengan(services.LayananDengan(g), true, true)

	w := kirim(t, h, http.MethodGet, "/api/treaty-in/versi/12/sebelumnya", "AKUN-UJI", badan(""))
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d, mau 200", w.Code)
	}
	for _, petik := range []string{`"ada":true`, "Treaty 2025"} {
		if !strings.Contains(w.Body.String(), petik) {
			t.Errorf("badan tidak memuat %q: %s", petik, w.Body.String())
		}
	}
}

// ------------------------------------------------------------------ tiket 41

func TestRuteBentukLamaMenyatakanNomorLamaTidakAda(t *testing.T) {
	g := gudangTiruan{kontrak: models.KontrakDenganVersi{Kontrak: models.Kontrak{
		ID: 7, IDCedant: 2, IDAsalBisnis: 3, SifatProporsi: models.SifatProporsional,
		TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31",
	}}}
	h := handlers.RouterDengan(services.LayananDengan(g), true, true)

	w := kirim(t, h, http.MethodGet, "/api/treaty-in/kontrak/7/bentuk-lama", "AKUN-UJI", badan(""))
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d, mau 200. Badan: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"nomorLamaAda":false`) {
		t.Errorf("badan tidak menyatakan ketiadaan nomor lama: %s", w.Body.String())
	}
	// Dan ia TIDAK mengarang nomornya.
	if strings.Contains(w.Body.String(), `"nomorLama"`) {
		t.Errorf("nomor lama dikarang untuk kontrak sistem baru: %s", w.Body.String())
	}
}

func TestRuteBentukLamaMenjawabNomorWarisan(t *testing.T) {
	g := gudangTiruan{kontrak: models.KontrakDenganVersi{Kontrak: models.Kontrak{
		ID: 7, NomorKontrakWarisan: "TR-2019-0042", IDCedant: 2, IDAsalBisnis: 3,
		SifatProporsi: models.SifatNonProporsional, TanggalMulai: "2019-01-01", TanggalBerakhir: "2019-12-31",
	}}}
	h := handlers.RouterDengan(services.LayananDengan(g), true, true)

	w := kirim(t, h, http.MethodGet, "/api/treaty-in/kontrak/7/bentuk-lama", "AKUN-UJI", badan(""))
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d, mau 200", w.Code)
	}
	for _, petik := range []string{`"nomorLamaAda":true`, "TR-2019-0042", "2019-01-01"} {
		if !strings.Contains(w.Body.String(), petik) {
			t.Errorf("badan tidak memuat %q: %s", petik, w.Body.String())
		}
	}
}

// Ketiga rute menolak permintaan tanpa identitas — 401, dan sebelum apa pun
// dibaca.
func TestKetigaRuteBaruMenolakTanpaIdentitas(t *testing.T) {
	// adaDB true: 503 dijawab SEBELUM identitas, jadi tanpa basis data uji
	// ini mengukur hal yang salah.
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	kasus := []struct{ cara, jalur, badan string }{
		{http.MethodPut, "/api/treaty-in/layer/7/pemulihan", `[{"nomorUrut":1,"persenPemulihan":"50"}]`},
		{http.MethodGet, "/api/treaty-in/versi/5/sebelumnya", ""},
		{http.MethodGet, "/api/treaty-in/kontrak/7/bentuk-lama", ""},
	}
	for _, k := range kasus {
		w := kirim(t, h, k.cara, k.jalur, "", badan(k.badan))
		if w.Code != http.StatusUnauthorized {
			t.Errorf("%s %s: kode %d, mau 401. Badan: %s", k.cara, k.jalur, w.Code, w.Body.String())
		}
	}
}

// ------------------------------------------------------------------ tiket 42

func TestRuteArsipMenyimpanDanMenunjukkanAda(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{
		bukti: models.BuktiArsip{Cacah: 2, TerakhirDikirim: "2026-10-02T09:15:00Z",
			Tujuan: []string{"PEGA_TREATY_IN"}},
	}), true, true)

	w := kirim(t, h, http.MethodPost, "/api/treaty-in/kontrak/7/arsip", "AKUN-UJI",
		badan(`{"tujuan":"PEGA_TREATY_IN","dikirimPada":"2026-10-02T09:15:00Z","berhasil":true,"muatan":"{}"}`))
	if w.Code != http.StatusNoContent {
		t.Fatalf("simpan: kode %d, mau 204. Badan: %s", w.Code, w.Body.String())
	}

	w = kirim(t, h, http.MethodGet, "/api/treaty-in/kontrak/7/arsip", "AKUN-UJI", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("bukti: kode %d, mau 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"cacah":2`) {
		t.Errorf("badan tidak menyebut cacahnya: %s", w.Body.String())
	}
	// ⛔ Dan ia TIDAK membawa muatan — INV-61 terlihat di permukaan HTTP juga.
	for _, terlarang := range []string{"muatan", "payload"} {
		if strings.Contains(strings.ToLower(w.Body.String()), terlarang) {
			t.Errorf("jawaban memuat %q; arsip tidak punya jalur baca: %s", terlarang, w.Body.String())
		}
	}
}

func TestRuteArsipMenolakTujuanAsing(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	w := kirim(t, h, http.MethodPost, "/api/treaty-in/kontrak/7/arsip", "AKUN-UJI",
		badan(`{"tujuan":"GUDANG_LAIN","dikirimPada":"2026-10-02T09:15:00Z","muatan":"{}"}`))
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("kode %d, mau 422. Badan: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "GUDANG_LAIN") {
		t.Errorf("penolakan tidak menyebut tujuan yang ditolak: %s", w.Body.String())
	}
}

// ------------------------------------------------------------- opsi kepala

func TestRuteOpsiKepala(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)
	if w := kirim(t, h, http.MethodGet, "/api/treaty-in/warisan/opsi-kepala", "", nil); w.Code != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: kode %d, mau 401", w.Code)
	}
	w := kirim(t, h, http.MethodGet, "/api/treaty-in/warisan/opsi-kepala", "AKUN-UJI", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"value":"underwriting","label":"Underwriting Year"`) {
		t.Errorf("kode %d badan %s", w.Code, w.Body.String())
	}
}

// ------------------------------------------------- Apply Reporting Period

func TestRuteHitungPeriodePelaporan(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)
	jalur := "/api/treaty-in/hitung/periode-pelaporan"
	isi := badan(`{"mulai":"06-10-2026","akhir":"06-10-2026","periode":"quarter","penyerahan":"12","konfirmasi":"12","pelunasan":"12"}`)
	if w := kirim(t, h, http.MethodPost, jalur, "", isi); w.Code != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d, mau 401", w.Code)
	}
	if w := kirim(t, h, http.MethodPost, jalur, "AKUN-UJI", badan(`{`)); w.Code != http.StatusBadRequest {
		t.Errorf("badan rusak: %d, mau 400", w.Code)
	}
	w := kirim(t, h, http.MethodPost, jalur, "AKUN-UJI", isi)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"jatuhTempoKirimAsli":"20270117"`) {
		t.Errorf("kode %d badan %s", w.Code, w.Body.String())
	}
}
