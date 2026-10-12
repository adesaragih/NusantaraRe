package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/services"
)

// 12-10-2026 (opsi A) - pop-up Send Claim to Committee: rute pratinjau + Remarks; kalimat galat Inggris.

func TestRuteKomiteTerdaftar(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, rute := range []string{
		"\"POST /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite\",\n\t\tserahkanKomite(svc, stubPelaku))",
		"\"GET /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite\",\n\t\tpratinjauKomite(svc, stubPelaku))",
	} {
		if !strings.Contains(string(isi), rute) {
			t.Errorf("rute tidak terdaftar: %s", rute)
		}
	}
}

func TestGalatKomiteDipetakan(t *testing.T) {
	for _, u := range []struct {
		err     error
		kode    int
		kalimat string
	}{
		{inti.ErrTanpaIdentitas, http.StatusUnauthorized, "The request has no user identity."},
		{inti.ErrTanpaWewenang, http.StatusForbidden, "Your role is not allowed to send this claim type to the committee."},
		{services.ErrRemarksKomiteKosong, http.StatusUnprocessableEntity, "Remarks is required."},
		{services.ErrRemarksKomiteKepanjangan, http.StatusUnprocessableEntity, "Remarks cannot exceed 4000 characters."},
		{services.ErrRekeningBelumLengkap, http.StatusUnprocessableEntity, "Name of bank cannot be empty"},
		{services.ErrSerahkanBukanClaimAnalis, http.StatusConflict, "Send Claim to Committee is only available at the Claim Analyst stage."},
		{services.ErrBarisSudahDiserahkan, http.StatusConflict, "This row has already been sent to the committee."},
		{services.ErrBarisBukanOutstanding, http.StatusConflict, "Only Outstanding rows can be sent to the committee."},
		{services.ErrMataUangKlaimCampur, http.StatusUnprocessableEntity, "The rows of this claim use more than one currency; sending to the committee requires a single currency."},
		{fmt.Errorf("%w: %w", galat.ErrPermintaanTidakSah, services.ErrNilaiKlaimKomiteKosong), http.StatusUnprocessableEntity, "Claim RNM of this row is empty; the committee level cannot be determined."},
		{services.ErrRosterKomiteKosong, http.StatusUnprocessableEntity, "No committee level covers this claim amount."},
		{kontrak.ErrKasusSudahTertutup, http.StatusConflict, "The case is already closed and cannot be changed."},
		{services.ErrTypeTidakDikenal, http.StatusUnprocessableEntity, "The claim type is not recognized; sending to the committee requires a valid type."},
		{errors.New("x"), http.StatusInternalServerError, "Failed to send the row to the committee."},
	} {
		w := httptest.NewRecorder()
		if !jawabGalatKomite(w, u.err) {
			t.Errorf("%v: tidak dijawab", u.err)
			continue
		}
		if w.Code != u.kode || !strings.Contains(w.Body.String(), u.kalimat) {
			t.Errorf("%v: %d %s, mau %d %q", u.err, w.Code, w.Body.String(), u.kode, u.kalimat)
		}
	}
	if jawabGalatKomite(httptest.NewRecorder(), nil) {
		t.Error("galat nil dijawab")
	}
}

// 12-10-2026 - pesan galat handler Claim Life yang tampil di layar berbahasa Inggris (tanpa campuran Indonesia).
func TestPesanGalatHandlerBerbahasaInggris(t *testing.T) {
	berkas, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	pola := regexp.MustCompile(`galat\.Tulis\(w, http\.Status[A-Za-z]+,\s*"([^"]*)"`)
	indonesia := regexp.MustCompile(`\b(tidak|belum|hanya|sudah|gagal|wajib|bukan|dapat|harus|kasus|klaim|peserta|baris|tahap)\b`)
	diperiksa := 0
	for _, b := range berkas {
		if b.IsDir() || !strings.HasSuffix(b.Name(), ".go") || strings.HasSuffix(b.Name(), "_test.go") {
			continue
		}
		isi, err := os.ReadFile(b.Name())
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range pola.FindAllStringSubmatch(string(isi), -1) {
			diperiksa++
			if indonesia.MatchString(m[1]) {
				t.Errorf("%s: pesan berbahasa Indonesia: %q", b.Name(), m[1])
			}
		}
	}
	if diperiksa < 50 {
		t.Fatalf("hanya %d pesan terbaca; pembacanya yang rusak", diperiksa)
	}
}
