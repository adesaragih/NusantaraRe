package handlers

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/services"
)

// TestCsvObjekHandler - POST /api/nbfacin/kasus/{caseId}/objek/csv (tiket 66): kunci jawaban dua arah dengan api.ts
// HasilCsvObjekServer (instrumen diuji atas jawaban yang sudah diketahui: 2 kunci baris / pesan); 401 tanpa identitas; 400
// badan bukan JSON / baris kosong ("Please Upload CSV"); 503 tanpa basis data.
func TestCsvObjekHandler(t *testing.T) {
	ts, gob := kunciInterfaceTS(t, "HasilCsvObjekServer"), kunciTagJSON(hasilCsvObjek{})
	if !reflect.DeepEqual(ts, []string{"baris", "pesan"}) || !reflect.DeepEqual(gob, ts) {
		t.Fatalf("ts %v, go %v", ts, gob)
	}
	jalur := "/api/nbfacin/kasus/UJI-NB-1/objek/csv"
	for nama, c := range map[string]struct {
		badan, pelaku string
		kode          int
		isi           string
	}{
		"401":        {`{"baris":[{"NO":"1"}]}`, "", 401, ""},
		"400 json":   {`{"baris":[1]}`, "UJI-USER", 400, "baris CSV"},
		"503":        {`{"baris":[{"NO":"1"}]}`, "UJI-USER", 503, ""},
		"bukan json": {`bukan`, "UJI-USER", 400, "baris CSV"},
	} {
		kode, isi := minta(t, services.Baru(nil), "POST", jalur, c.badan, c.pelaku)
		if kode != c.kode || !strings.Contains(isi, c.isi) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	// Ukuran badan (RALAT D7): 3000 baris × 61 kunci ber-isi 30 byte (± 8,5 MiB) diterima pembaca badan (lanjut ke layanan:
	// 503 tanpa basis data); di atas 16 MiB -> 400.
	badanBesar := func(baris, lebar int) string {
		sel := `"` + strings.Repeat("X", lebar) + `"`
		var kunci []string
		for _, k := range services.KolomCsvObjek {
			kunci = append(kunci, `"`+k+`":`+sel)
		}
		satu := "{" + strings.Join(kunci, ",") + "}"
		return `{"baris":[` + strings.TrimSuffix(strings.Repeat(satu+",", baris), ",") + `]}`
	}
	if b := badanBesar(services.BatasBarisCsvObjek, 30); len(b) < 8<<20 || len(b) > batasBadanCsvObjek {
		t.Fatalf("instrumen: %d byte", len(b))
	}
	if kode, isi := minta(t, services.Baru(nil), "POST", jalur, badanBesar(services.BatasBarisCsvObjek, 30), "UJI-USER"); kode != 503 {
		t.Errorf("3000 baris: %d %s", kode, isi)
	}
	if kode, _ := minta(t, services.Baru(nil), "POST", jalur, badanBesar(services.BatasBarisCsvObjek, 100), "UJI-USER"); kode != 400 {
		t.Errorf("> 16 MiB: %d", kode)
	}
}
