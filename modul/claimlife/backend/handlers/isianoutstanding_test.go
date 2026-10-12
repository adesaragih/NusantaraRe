package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/services"
)

func TestRuteIsianOutstandingTerdaftar(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `"PUT /api/klaim-life/{id}/outstanding/isian", simpanIsianOutstanding(`) {
		t.Error("rute Save layar Outstanding tidak terdaftar")
	}
}

func TestGalatIsianOutstandingDipetakan(t *testing.T) {
	for _, u := range []struct {
		nama string
		err  error
		mau  int
	}{
		{"identitas", inti.ErrTanpaIdentitas, http.StatusUnauthorized},
		{"wewenang", inti.ErrTanpaWewenang, http.StatusForbidden},
		{"tertutup", kontrak.ErrKasusSudahTertutup, http.StatusConflict},
		{"tahap", services.ErrIsianBukanOutstanding, http.StatusConflict},
		{"katastrofe terkunci", services.ErrKatastrofeTerkunci, http.StatusConflict},
		{"klaim tidak ada", services.ErrKlaimTidakAda, http.StatusNotFound},
		{"permintaan", galat.ErrPermintaanTidakSah, http.StatusBadRequest},
		{"lain", errors.New("x"), http.StatusInternalServerError},
	} {
		w := httptest.NewRecorder()
		jawabGalatIsianOutstanding(w, u.err)
		if w.Code != u.mau {
			t.Errorf("%s: kode %d, mau %d", u.nama, w.Code, u.mau)
		}
	}
}

// Gerbang tombol bawah layar Outstanding (10-10-2026): Outstanding -> Medical tanpa CLAIM_NO = 409.
func TestGerbangPindahOutstandingDijawab409(t *testing.T) {
	for _, err := range []error{services.ErrPindahMenungguNomorKlaim, services.ErrPindahTypeTanpaMedis,
		services.ErrPindahSudahBernomor} {
		w := httptest.NewRecorder()
		jawabGalatPindah(w, err)
		if w.Code != http.StatusConflict {
			t.Errorf("%v: kode %d, mau 409", err, w.Code)
		}
	}
	w := httptest.NewRecorder()
	jawabGalatPindah(w, nil)
	if w.Code != http.StatusNoContent {
		t.Errorf("berhasil: kode %d, mau 204", w.Code)
	}
}

func TestRutePesertaPolisTerdaftar(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `"GET /api/klaim-life/{id}/peserta-polis", pesertaPolis(svc)`) {
		t.Error("rute View Detail peserta polis tidak terdaftar")
	}
}

// 11-10-2026: Claim RNM + bank dilebur ke tombol Save grid (rute sendiri dicabut); galatnya lewat jawabGalatIsianOutstanding.
func TestGalatClaimRNMLewatSaveGrid(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(isi), "/adjustment/{adjId}/isian") {
		t.Error("rute isian baris masih terdaftar - Claim RNM kini lewat PUT …/outstanding/isian")
	}
	for _, u := range []struct {
		err error
		mau int
	}{
		{models.ErrClaimGrossMelebihiShare, http.StatusUnprocessableEntity},
		{services.ErrAdjustmentTerkunci, http.StatusConflict},
		{services.ErrCatatanAdminTerkunci, http.StatusConflict},
		{galat.ErrPermintaanTidakSah, http.StatusBadRequest},
	} {
		w := httptest.NewRecorder()
		jawabGalatIsianOutstanding(w, u.err)
		if w.Code != u.mau {
			t.Errorf("%v: kode %d, mau %d", u.err, w.Code, u.mau)
		}
	}
}

// 11-10-2026: tanggal peserta sesudah Save to RNM = 409; DOL di luar jendela = 422 (kalimat ValidasiDOL_Act).
func TestGalatTanggalIsianOutstanding(t *testing.T) {
	for _, u := range []struct {
		err error
		mau int
	}{{services.ErrTanggalTerkunciSesudahSaveRNM, http.StatusConflict}, {services.ErrDOLTidakSah, http.StatusUnprocessableEntity}} {
		w := httptest.NewRecorder()
		jawabGalatIsianOutstanding(w, u.err)
		if w.Code != u.mau {
			t.Errorf("%v: kode %d, mau %d", u.err, w.Code, u.mau)
		}
	}
}

func TestRuteSpreadingPolisTerdaftar(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `"GET /api/klaim-life/{id}/peserta/{pesertaId}/spreading-polis", spreadingPolis(svc)`) {
		t.Error("rute spreading polis tidak terdaftar")
	}
}

// 11-10-2026: Save grid menolak DOL di luar jendela dengan kalimat jelas (nomor peserta + rentang), bukan "Invalid DOL".
func TestDOLSaveGridKalimatJelas(t *testing.T) {
	w := httptest.NewRecorder()
	jawabGalatIsianOutstanding(w, services.GalatDOL{PesertaID: "P1", Nomor: 1, Awal: "18/03/2025", Akhir: "18/03/2026"})
	if w.Code != http.StatusUnprocessableEntity ||
		!strings.Contains(w.Body.String(), "Insured number 1: Date of Loss must be within the valuation period (18/03/2025 – 18/03/2026).") {
		t.Errorf("%d %s", w.Code, w.Body.String())
	}
}

// 11-10-2026: peserta Accept tanpa dokumen saat Send to Medical Check - 422 dengan daftar kekurangan.
func TestPindahMedisDokumenKurang422(t *testing.T) {
	w := httptest.NewRecorder()
	jawabGalatPindah(w, &services.PelanggaranPindah{Langkah: "dokumen-medis", Pesan: "a", Kekurangan: []services.KekuranganRNM{
		{NomorPeserta: 1, PesertaID: "P1", Sertifikat: "S1", Nama: "UJI A", Kolom: "Document", Pesan: "a"}}})
	if w.Code != http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), `"kekurangan":[{"nomorPeserta":1`) {
		t.Errorf("%d %s", w.Code, w.Body.String())
	}
}
