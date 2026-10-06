package repository

import (
	"errors"
	"strings"
	"testing"
)

func p(s string) *string { return &s }

func TestPeriksaMode(t *testing.T) {
	if err := PeriksaMode(true, true); !errors.Is(err, ErrPindahDiProduksi) {
		t.Errorf("jalankan di produksi: %v", err)
	}
	for _, k := range [][2]bool{{false, true}, {false, false}, {true, false}} {
		if err := PeriksaMode(k[0], k[1]); err != nil {
			t.Errorf("%v: %v", k, err)
		}
	}
}

// K-F2: nilai APA ADANYA seperti `a.JSONDATA.X` - teks tidak dipangkas, angka/boolean literal, kosong/null/objek = NULL.
func TestNilaiJSON(t *testing.T) {
	j := `{"USEDBY":"  UJI A ","TYPE":5,"FLAG":true,"OPERATORID":"","MODIFIEDDATE":null,"X":{"a":1},"Y":[1]}`
	for kunci, mau := range map[string]*string{"USEDBY": p("  UJI A "), "TYPE": p("5"), "FLAG": p("true"), "OPERATORID": nil,
		"MODIFIEDDATE": nil, "X": nil, "Y": nil, "TIDAK": nil} {
		got, ok := NilaiJSON(j, kunci)
		if !ok || !samaTeks(got, mau) {
			t.Errorf("%s: %v %v", kunci, got, ok)
		}
	}
	if _, ok := NilaiJSON("bukan json", "USEDBY"); ok {
		t.Error("JSON rusak harus ok=false")
	}
	if v, ok := NilaiJSON("", "USEDBY"); !ok || v != nil {
		t.Error("JSONDATA kosong = NULL")
	}
}

// Rencana: panjang maksimum per kolom dilaporkan; nilai yang TIDAK MUAT menahan -jalankan (nol pemotongan, tanpa nilai
// di laporan); ID flat sama dilewati, beda menahan; flat saja dibiarkan.
func TestRencanaPindah(t *testing.T) {
	panjang := strings.Repeat("N", 501)
	sumber := []BarisJSON{
		{ID: "101", JSON: `{"USEDBY":"UJI RATE A","TYPE":"L","MODIFIEDDATE":"20240102T030405.000 GMT","OPERATORID":"UJI-A","FLAG":"1"}`},
		{ID: "102", JSON: `{"USEDBY":"UJI RATE B"}`},
		{ID: "103", JSON: `{"USEDBY":"` + panjang + `"}`},
		{ID: "104", JSON: `{"USEDBY":"UJI RATE D"}`},
		{ID: "105", JSON: `{"USEDBY":"UJI RATE E"}`},
		{ID: "106", JSON: `rusak`},
		{ID: "12345678901", JSON: `{"USEDBY":"UJI"}`},
	}
	flat := []RingkasanFlat{
		{ID: "104", UsedBy: p("UJI RATE D")},
		{ID: "105", UsedBy: p("UJI RATE E DIUBAH")},
		{ID: "900", UsedBy: p("UJI APLIKASI")},
	}
	lap, tulis := RencanaPindah(sumber, flat)
	if lap.Sumber != 7 || lap.AkanDitulis != 2 || lap.SudahSama != 1 || lap.FlatSaja != 1 || len(lap.Berbeda) != 1 || lap.Berbeda[0] != "105" {
		t.Errorf("laporan %+v", lap)
	}
	if len(tulis) != 2 || !tulis[0].Sama(RingkasanFlat{ID: "101", UsedBy: p("UJI RATE A"), Type: p("L"),
		ModifiedDate: p("20240102T030405.000 GMT"), OperatorID: p("UJI-A"), Flag: p("1")}) ||
		!tulis[1].Sama(RingkasanFlat{ID: "102", UsedBy: p("UJI RATE B")}) {
		t.Errorf("tulis %+v", tulis)
	}
	if lap.PanjangMaks["USEDBY"] != 501 || lap.PanjangMaks["MODIFIEDDATE"] != 23 || lap.PanjangMaks["ID"] != 11 {
		t.Errorf("panjang maks %v", lap.PanjangMaks)
	}
	if len(lap.TidakMuat) != 2 || lap.TidakMuat[0] != "103 USEDBY: 501 byte > 500" || lap.TidakMuat[1] != "12345678901 ID: 11 byte > 10" {
		t.Errorf("tidak muat %q", lap.TidakMuat)
	}
	if len(lap.Gagal) != 1 || lap.Gagal[0] != "106 JSONDATA: bukan objek JSON" {
		t.Errorf("gagal %q", lap.Gagal)
	}
	if strings.Contains(lap.Teks(), panjang) || strings.Contains(lap.Teks(), "UJI RATE") {
		t.Error("laporan memuat nilai")
	}
	if lap.BolehDitulis() {
		t.Error("tidak muat / berbeda harus menahan")
	}
	bersih, _ := RencanaPindah(sumber[:2], nil)
	if !bersih.BolehDitulis() || bersih.AkanDitulis != 2 {
		t.Errorf("bersih %+v", bersih)
	}
	ulang, tulisUlang := RencanaPindah(sumber[:2], tulis)
	if !ulang.BolehDitulis() || len(tulisUlang) != 0 || ulang.SudahSama != 2 {
		t.Errorf("ulang aman %+v", ulang)
	}
	if !strings.Contains(PeriksaTulisPesan(SqlKunciRingkasan("S.RATE_LIFE_SUMMARY")), "ok") {
		t.Error("kunci tabel flat ditolak penjaga")
	}
}

// PeriksaTulisPesan - pembantu uji: "ok" bila pernyataan boleh atas tabel flat.
func PeriksaTulisPesan(q string) string {
	if err := PeriksaTulis(TabelRingkasan, q); err != nil {
		return err.Error()
	}
	return "ok"
}
