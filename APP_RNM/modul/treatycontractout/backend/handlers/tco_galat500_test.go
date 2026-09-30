package handlers

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ⛔ Sebab uji ini ada, dan itu sungguh terjadi (29-09-2026): grid `Limit MB`
// dua tahun treaty menjawab 500 "gagal memproses permintaan treaty contract
// out", dan sebabnya tidak tercatat DI MANA PUN - tidak di layar, tidak di
// konsol backend. Diagnosanya menuntut membaca DEV langsung (baris yang dibaca:
// `RP` warisan "1.000.000" yang tidak terurai). Galat tak terduga tetap berkalimat umum bagi
// pemakai, tetapi sebab aslinya wajib tercatat di log backend.
func TestGalat500TreatyContractOutMencatatSebabnya(t *testing.T) {
	var catatan bytes.Buffer
	lama := log.Writer()
	log.SetOutput(&catatan)
	defer log.SetOutput(lama)

	w := httptest.NewRecorder()
	jawabGalatTreatyContractOut(w, errors.New(`repository: column RP has value "UJI": cause`))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("kode %d, mau 500", w.Code)
	}
	if !strings.Contains(catatan.String(), `column RP has value "UJI"`) {
		t.Errorf("sebab asli tidak tercatat di log backend; log = %q", catatan.String())
	}
	if strings.Contains(w.Body.String(), "UJI") {
		t.Errorf("rincian internal bocor ke pemakai: %s", w.Body.String())
	}
}
