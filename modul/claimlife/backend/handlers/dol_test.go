package handlers

// Uji pintu HTTP dialog Edit Date - DOL (tiket 06) dan tiga tanggal klaim
// (sensus §3.1, 28-09-2026).

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/services"
)

func TestRuteTanggalKlaimTerdaftarDiSarangPeserta(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi),
		`"PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-klaim"`) {
		t.Error("rute tanggal klaim tidak terdaftar di Router")
	}
}

// TestGalatTanggalDipetakanKeKodeYangBenar - SATU terjemahan untuk kedua rute
// dialog. ⛔ Sebelum 28-09-2026 rute DOL menjawab 500 atas kasus yang sudah
// ditutup: `ErrKasusSudahTertutup` tidak ada di daftarnya.
func TestGalatTanggalDipetakanKeKodeYangBenar(t *testing.T) {
	for _, u := range []struct {
		nama string
		err  error
		mau  int
	}{
		{"tanpa identitas", inti.ErrTanpaIdentitas, http.StatusUnauthorized},
		{"tanpa wewenang", inti.ErrTanpaWewenang, http.StatusForbidden},
		{"kasus tertutup", kontrak.ErrKasusSudahTertutup, http.StatusConflict},
		{"tahap salah", services.ErrTahapTidakBolehUbahTanggal, http.StatusConflict},
		{"terkunci sesudah Save to RNM", services.ErrTanggalTerkunciSesudahSaveRNM, http.StatusConflict},
		{"tahap tak dikenal", services.ErrTahapTidakDikenal, http.StatusUnprocessableEntity},
		{"permintaan tidak sah", galat.ErrPermintaanTidakSah, http.StatusBadRequest},
	} {
		w := httptest.NewRecorder()
		jawabGalatTanggal(w, u.err, "gagal")
		if w.Code != u.mau {
			t.Errorf("%s: kode = %d, mau %d", u.nama, w.Code, u.mau)
		}
	}
}

// TestTanggalKlaimMenolakTanggalAsing - bentuk yang tidak dikenal ditolak 400
// di pintu, sebelum layanan dipanggil; kosong SAH (kolomnya menjadi NULL).
func TestTanggalKlaimMenolakTanggalAsing(t *testing.T) {
	if _, err := uraiTanggalOpsional(""); err != nil {
		t.Errorf("kosong ditolak: %v", err)
	}
	if w, err := uraiTanggalOpsional("  "); err != nil || w != nil {
		t.Errorf("spasi saja harus kosong: %v %v", w, err)
	}
	if _, err := uraiTanggalOpsional("2026-13-45"); err == nil {
		t.Error("tanggal mustahil diterima")
	}
	if w, err := uraiTanggalOpsional("2026-03-02"); err != nil || w == nil {
		t.Errorf("tanggal sah ditolak: %v", err)
	}
}

func TestRuteTanggalKlaimTanpaDatabaseMenjawab503(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut,
		"/api/klaim-life/K-1/peserta/P-1/tanggal-klaim", strings.NewReader("{}"))
	r.SetPathValue("id", "K-1")
	r.SetPathValue("pesertaId", "P-1")
	setTanggalKlaim(services.New(nil), true)(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("kode = %d, mau 503", w.Code)
	}
}

// OQ-M5 (GILIRAN-17): rute tolak membaca `Remarks` dari badan JSON; badan
// yang bukan JSON ditolak di pintu.
func TestUraiAlasanTolak(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(`{"komentar":" UJI alasan "}`))
	if a, err := uraiAlasanTolak(r); err != nil || a != " UJI alasan " {
		t.Errorf("badan sah: %q %v", a, err)
	}
	r = httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(`bukan json`))
	if _, err := uraiAlasanTolak(r); err == nil {
		t.Error("badan bukan JSON diterima")
	}
}

// OQ-M6 (GILIRAN-17): rute cabut peserta terdaftar di sarang peserta, dan
// galatnya diterjemahkan sekali.
func TestRuteCabutPesertaDanGalatnya(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `"POST /api/klaim-life/{id}/peserta/{pesertaId}/cabut"`) {
		t.Error("rute cabut peserta tidak terdaftar di Router")
	}
	for _, u := range []struct {
		err error
		mau int
	}{
		{inti.ErrTanpaIdentitas, http.StatusUnauthorized},
		{inti.ErrTanpaWewenang, http.StatusForbidden},
		{kontrak.ErrKasusSudahTertutup, http.StatusConflict},
		{services.ErrPesertaTidakDapatDicabut, http.StatusConflict},
		{galat.ErrPermintaanTidakSah, http.StatusBadRequest},
		{services.ErrTahapTidakDikenal, http.StatusUnprocessableEntity},
	} {
		w := httptest.NewRecorder()
		jawabGalatCabut(w, u.err)
		if w.Code != u.mau {
			t.Errorf("%v: kode = %d, mau %d", u.err, w.Code, u.mau)
		}
	}
}
