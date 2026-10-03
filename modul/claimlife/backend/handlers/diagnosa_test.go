package handlers

// Uji pintu HTTP grid diagnosa - butir bd.

import (
	"encoding/json"
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

func TestKetigaRuteDiagnosaTerdaftar(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	// ⛔ Metodenya disebut satu per satu, dan itu bukan kelebihan: satu rute
	// yang terdaftar dengan metode yang salah tetap membuat Router terkompilasi
	// dan tetap membuat layar diam. Cacat yang persis begitu pernah nyata di
	// modul ini - rute yang ada tanpa pemanggil.
	for _, rute := range []string{
		`"POST /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa"`,
		`"PUT /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa/{diagId}"`,
		`"DELETE /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa/{diagId}"`,
	} {
		if !strings.Contains(teks, rute) {
			t.Errorf("rute %s tidak terdaftar di Router", rute)
		}
	}
	// ⛔ Dan TIDAK ada jalur diagnosa yang berdiri di luar pesertanya.
	// Jalur `/api/diagnosa/{id}` akan membuat pengenal peserta opsional -
	// dan gerbang yang memeriksa peserta menjadi gerbang yang dapat dilewati.
	if strings.Contains(teks, `"/api/diagnosa/`) {
		t.Error("ada rute diagnosa di luar sarang pesertanya")
	}
}

func TestRuteDiagnosaTanpaDatabaseMenjawab503(t *testing.T) {
	svc := services.New(nil)
	for _, u := range []struct {
		nama    string
		metode  string
		jalur   string
		panggil func(*services.Service, bool) http.HandlerFunc
	}{
		{"tambah", http.MethodPost, "/api/klaim-life/K-1/peserta/P-1/diagnosa", tambahDiagnosa},
		{"ubah", http.MethodPut, "/api/klaim-life/K-1/peserta/P-1/diagnosa/7", ubahDiagnosa},
		{"hapus", http.MethodDelete, "/api/klaim-life/K-1/peserta/P-1/diagnosa/7", hapusDiagnosa},
	} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(u.metode, u.jalur, strings.NewReader("{}"))
		r.SetPathValue("id", "K-1")
		r.SetPathValue("pesertaId", "P-1")
		r.SetPathValue("diagId", "7")
		u.panggil(svc, true)(w, r)
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("%s: kode = %d, mau 503", u.nama, w.Code)
		}
		// ⛔ Kunci amplopnya `galat`, tidak pernah `error`. Cacat lintas-lapis
		// PERTAMA di modul ini persis begitu, dan layar membaca `galat`.
		var amplop map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &amplop); err != nil {
			t.Fatalf("%s: badan bukan JSON: %v", u.nama, err)
		}
		if _, ada := amplop["galat"]; !ada {
			t.Errorf("%s: amplop galat tanpa kunci `galat`: %s", u.nama, w.Body)
		}
	}
}

func TestPengenalDiagnosaBukanAngkaDijawab400(t *testing.T) {
	// ⛔ SEBELUM services, dan sebelum database. Pengenal yang tidak terbaca
	// menjadi 0 akan dicari, tidak ketemu, lalu dijawab "bukan milik peserta"
	// - kalimat yang BENAR tentang hal yang SALAH, dan yang membacanya akan
	// mencari baris yang tidak pernah diminta siapa pun.
	for _, buruk := range []string{"", "abc", "0", "-3", "7.5", "9999999999999999999999"} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/x", nil)
		r.SetPathValue("diagId", buruk)
		if _, ok := pengenalDiagnosa(w, r); ok {
			t.Errorf("pengenal %q diterima", buruk)
			continue
		}
		if w.Code != http.StatusBadRequest {
			t.Errorf("pengenal %q: kode = %d, mau 400", buruk, w.Code)
		}
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/x", nil)
	r.SetPathValue("diagId", "42")
	n, ok := pengenalDiagnosa(w, r)
	if !ok || n != 42 {
		t.Errorf("pengenal 42 -> (%d, %v)", n, ok)
	}
}

func TestGalatDiagnosaDipetakanKeKodeYangBenar(t *testing.T) {
	// ⛔ TABEL, dan setiap barisnya pernah menjadi 500 di ronde pertama
	// berkas lain. 500 membuat orang menelepon; 409 membuat orang membaca
	// kalimatnya dan mengerti bahwa pesertanya memang sudah diputus.
	for _, u := range []struct {
		nama string
		err  error
		mau  int
	}{
		{"tanpa identitas", inti.ErrTanpaIdentitas, http.StatusUnauthorized},
		{"tanpa wewenang", inti.ErrTanpaWewenang, http.StatusForbidden},
		{"peserta terkunci", services.ErrDiagnosaTerkunci, http.StatusConflict},
		{"kasus tertutup", kontrak.ErrKasusSudahTertutup, http.StatusConflict},
		{"tahap salah", services.ErrTahapTidakBergridPeserta, http.StatusConflict},
		{"nilai kepanjangan", services.ErrNilaiDiagnosaKepanjangan,
			http.StatusUnprocessableEntity},
		{"permintaan tidak sah", galat.ErrPermintaanTidakSah, http.StatusBadRequest},
	} {
		w := httptest.NewRecorder()
		if !jawabGalatDiagnosa(w, u.err) {
			t.Errorf("%s: galat tidak dijawab", u.nama)
			continue
		}
		if w.Code != u.mau {
			t.Errorf("%s: kode = %d, mau %d", u.nama, w.Code, u.mau)
		}
	}
	// nil berarti "belum dijawab" - pemanggil meneruskan ke jalur berhasil.
	if jawabGalatDiagnosa(httptest.NewRecorder(), nil) {
		t.Error("galat nil dijawab; jalur berhasil tidak akan tercapai")
	}
}
