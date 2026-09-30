package handlers

// Uji pintu HTTP dokumen - butir be.

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/claimlife/backend/services"
)

func TestKetigaRuteDokumenTerdaftar(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	for _, rute := range []string{
		`"POST /api/klaim-life/{id}/peserta/{pesertaId}/dokumen"`,
		`"GET /api/dokumen/{dokId}/isi"`,
		`"DELETE /api/klaim-life/{id}/dokumen/{dokId}"`,
	} {
		if !strings.Contains(teks, rute) {
			t.Errorf("rute %s tidak terdaftar di Router", rute)
		}
	}
	// ⛔ Unduh adalah GET dan HANYA GET: berkas yang dapat diminta dengan
	// POST akan mengundang orang menulis ke jalur yang hanya membaca.
	if strings.Contains(teks, `"POST /api/dokumen/{dokId}/isi"`) {
		t.Error("rute unduh terdaftar sebagai POST")
	}
}

func TestPengenalDokumenBukanAngkaDijawab400(t *testing.T) {
	for _, buruk := range []string{"", "abc", "0", "-1", "20260927103000123456789"} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/x", nil)
		r.SetPathValue("dokId", buruk)
		if _, ok := pengenalDokumen(w, r); ok {
			t.Errorf("pengenal %q diterima", buruk)
			continue
		}
		if w.Code != http.StatusBadRequest {
			t.Errorf("pengenal %q: kode = %d, mau 400", buruk, w.Code)
		}
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/x", nil)
	// Cap waktu 17 angka - bentuk `@CurrentDate("yyyyMMddhhmmssSSS")` b648.
	r.SetPathValue("dokId", "20260927103000123")
	n, ok := pengenalDokumen(w, r)
	if !ok || n != 20260927103000123 {
		t.Errorf("pengenal cap waktu -> (%d, %v)", n, ok)
	}
}

func TestGalatDokumenDipetakanKeKodeYangBenar(t *testing.T) {
	for _, u := range []struct {
		nama string
		err  error
		mau  int
	}{
		{"tanpa identitas", inti.ErrTanpaIdentitas, http.StatusUnauthorized},
		{"kasus tertutup", kontrak.ErrKasusSudahTertutup, http.StatusConflict},
		// ⛔ 503, bukan 500 dan bukan 400: folder yang belum disetel adalah
		// keadaan SERVER yang belum siap, bukan permintaan yang salah.
		{"folder belum disetel", unggah.ErrUnggahanDirBelumDisetel,
			http.StatusServiceUnavailable},
		{"kategori belum ada", services.ErrKategoriWajibBelumDiketahui,
			http.StatusServiceUnavailable},
		{"terlalu besar", unggah.ErrBerkasTerlaluBesar,
			http.StatusRequestEntityTooLarge},
		{"kosong", unggah.ErrBerkasKosong, http.StatusBadRequest},
		{"kategori asing", services.ErrKategoriDokumenTidakDikenal,
			http.StatusUnprocessableEntity},
		// ⛔ 409, bukan 404: barisnya ADA, berkasnya belum tertaut. Layar
		// dapat berkata "sedang diproses" alih-alih "tidak ditemukan".
		{"belum terunggah", services.ErrDokumenBelumTerunggah, http.StatusConflict},
		// ⛔ 404, bukan 500 - uji asap DEV (GILIRAN-12): pengenal yang tidak
		// ada dijawab 500 "gagal memproses dokumen", dan itu terbaca sebagai
		// kerusakan server padahal pengenalnya yang salah.
		{"dokumen tidak ada", services.ErrDokumenTidakAda, http.StatusNotFound},
		{"permintaan tidak sah", galat.ErrPermintaanTidakSah, http.StatusBadRequest},
	} {
		w := httptest.NewRecorder()
		if !jawabGalatDokumen(w, u.err) {
			t.Errorf("%s: galat tidak dijawab", u.nama)
			continue
		}
		if w.Code != u.mau {
			t.Errorf("%s: kode = %d, mau %d", u.nama, w.Code, u.mau)
		}
	}
	if jawabGalatDokumen(httptest.NewRecorder(), nil) {
		t.Error("galat nil dijawab")
	}
}

func TestUnggahTanpaDatabaseMenjawab503(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/klaim-life/K-1/peserta/P-1/dokumen", nil)
	r.SetPathValue("id", "K-1")
	r.SetPathValue("pesertaId", "P-1")
	unggahDokumen(services.New(nil), true)(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("kode = %d, mau 503", w.Code)
	}
}
