package handlers

// Uji pintu HTTP form penawaran - tiket 01 bagian 3.

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/services"
)

// Keempat rute terdaftar tanpa tabrakan pola (ServeMux panik bila bertabrakan).
func TestRutePenawaranTerdaftarTanpaTabrakan(t *testing.T) {
	h := Router(services.New(nil), true)
	for _, r := range []struct{ metode, jalur string }{
		{http.MethodGet, "/api/polis-life/NBLF-1/penawaran"},
		{http.MethodPut, "/api/polis-life/NBLF-1/penawaran"},
		{http.MethodGet, "/api/polis-life/cari-ceding?cari=a"},
		{http.MethodGet, "/api/polis-life/cari-pemegang-polis?cari=a"},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(r.metode, r.jalur, strings.NewReader("{}")))
		// Tanpa Oracle setiap rute menjawab 503 - bukan 404/405 rute yang hilang.
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("%s %s = %d, mau 503", r.metode, r.jalur, w.Code)
		}
	}
}

// Isian belum lengkap: kalimat Pega saja, tanpa awalan paket.
func TestIsianBelumLengkapDijawabKalimatPega(t *testing.T) {
	w := httptest.NewRecorder()
	err := fmt.Errorf("%w: %w", galat.ErrPermintaanTidakSah,
		fmt.Errorf("%w: %s", models.ErrIsianPenawaranBelumLengkap,
			models.GabungPesanPenawaran([]string{models.PesanTypeCedingKosong, models.PesanBusinessCodeKosong})))
	if !jawabGalatPenawaran(w, err) || w.Code != http.StatusBadRequest {
		t.Fatalf("kode = %d", w.Code)
	}
	badan := w.Body.String()
	if strings.Contains(badan, "models:") || !strings.Contains(badan, models.PesanTypeCedingKosong) {
		t.Errorf("badan = %s", badan)
	}
}

func TestConfirmPenawaranBelumLengkapDijawab409(t *testing.T) {
	w := httptest.NewRecorder()
	jawabGalatPenawaran(w, fmt.Errorf("%w: %s", services.ErrPenawaranBelumLengkap, models.PesanCommentKosong))
	if w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), models.PesanCommentKosong) {
		t.Errorf("kode = %d, badan = %s", w.Code, w.Body.String())
	}
}

func TestTanggalDiterimaBerbentukTanggal(t *testing.T) {
	if _, err := (isiPenawaran{DateReceived: "30/09/2026"}).keIsian(); err == nil {
		t.Error("tanggal dd/mm/yyyy diterima")
	}
	isi, err := isiPenawaran{DateReceived: "2026-09-30"}.keIsian()
	if err != nil || isi.DateReceived == nil || isi.DateReceived.Day() != 30 {
		t.Errorf("isi = %+v, err = %v", isi, err)
	}
	if isi, _ := (isiPenawaran{}).keIsian(); isi.DateReceived != nil {
		t.Error("tanggal kosong menjadi tanggal nol, bukan nil")
	}
}

func TestBadanPenawaranMengurai059(t *testing.T) {
	isi, err := isiPenawaran{BatasUsiaPeserta: "62", TBC: "25", SumInsured: "6,000,000,000",
		TanggalKonfirmasi: "2026-09-29", PeriodePertanggungan: " 1 YEAR "}.keIsian()
	if err != nil {
		t.Fatal(err)
	}
	if *isi.BatasUsiaPeserta != 62 || *isi.TBC != 25 || isi.SumInsured.Text('f') != "6000000000" ||
		isi.TanggalKonfirmasi.Day() != 29 {
		t.Errorf("isi = %+v", isi)
	}
	for _, salah := range []isiPenawaran{{TBC: "2,5"}, {SumInsured: "enam"}, {TanggalRespon: "31/07/2026"}} {
		if _, err := salah.keIsian(); err == nil {
			t.Errorf("%+v diterima", salah)
		}
	}
}

func TestRuteDataPolisTerdaftar(t *testing.T) {
	h := Router(services.New(nil), true)
	for _, r := range []struct{ metode, jalur string }{
		{http.MethodGet, "/api/polis-life/NBLF-1/data-polis"},
		{http.MethodPut, "/api/polis-life/NBLF-1/data-polis"},
		{http.MethodGet, "/api/polis-life/NBLF-1/cari-produk?cari=a"},
		{http.MethodGet, "/api/polis-life/cari-marketing?cari=a"},
		{http.MethodGet, "/api/polis-life/cari-rislip?cari=a"},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(r.metode, r.jalur, strings.NewReader("{}")))
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("%s %s = %d, mau 503", r.metode, r.jalur, w.Code)
		}
	}
	if _, err := (isiDataPolis{AnnuityInterest: "5,5"}).keIsian(); err == nil {
		t.Error("koma desimal diterima - layar memakai titik")
	}
}
