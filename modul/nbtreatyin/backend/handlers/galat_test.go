package handlers

// Uji seam 1 - pemetaan galat ke HTTP: 422 data kontrak/master (AC 37) hanya
// membawa teks galat dasarnya; nama objek berskema dan rincian katalog yang
// dibungkus repository (`repository/acuan.go` tipeKolomObjek/pilihKolom) tidak
// sampai ke klien (temuan tinjauan P9). Fixture UJI-.

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/services"
)

func TestGalat422TanpaNamaSkema(t *testing.T) {
	for _, dasar := range galatDataDasar {
		w := httptest.NewRecorder()
		tulisGalat(w, fmt.Errorf("services: UJI-konteks: %w",
			fmt.Errorf("%w: objek UJI_SKEMA.TREATYINDETAILJOINEDM tidak terbaca di katalog", dasar)))
		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("%v: kode %d", dasar, w.Code)
		}
		isi := w.Body.String()
		if strings.Contains(isi, "UJI_SKEMA") || strings.Contains(isi, "UJI-konteks") {
			t.Errorf("%v: rincian bungkusan bocor ke klien: %s", dasar, isi)
		}
		if !strings.Contains(isi, dasar.Error()) {
			t.Errorf("%v: teks galat dasar wajib tampil (AC 37): %s", dasar, isi)
		}
	}
	if galatData(services.ErrKasusTertutup) != nil {
		t.Fatal("galat lain bukan 422 data")
	}
}
