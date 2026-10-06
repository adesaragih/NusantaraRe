package repository

import (
	"errors"
	"slices"
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

func jsonRing(usedby, modified, flag string) string {
	s := `{"USEDBY":"` + usedby + `"`
	if modified != "" {
		s += `,"MODIFIEDDATE":"` + modified + `"`
	}
	if flag != "" {
		s += `,"FLAG":"` + flag + `"`
	}
	return s + `}`
}

// Pemindahan penuh (tanpa -sejak): FLAG dan TYPE disalin APA ADANYA (AP/PM/PY dari data DEV); panjang maksimum per kolom
// dilaporkan; nilai yang TIDAK MUAT menahan -jalankan (nol pemotongan, tanpa nilai di laporan); cacah = saat itu.
func TestRencanaPindahPenuh(t *testing.T) {
	panjang := strings.Repeat("N", 501)
	sumber := []BarisJSON{
		{ID: "101", JSON: `{"USEDBY":"UJI RATE A","TYPE":"L","MODIFIEDDATE":"20240102T030405.000 GMT","OPERATORID":"UJI-A","FLAG":"AP"}`},
		{ID: "102", JSON: jsonRing("UJI RATE B", "20240102T030405.000 GMT", "PM")},
		{ID: "103", JSON: jsonRing(panjang, "", "")},
		{ID: "104", JSON: jsonRing("UJI RATE D", "", "PY")},
		{ID: "106", JSON: `rusak`},
		{ID: "12345678901", JSON: jsonRing("UJI", "", "")},
	}
	lap, sisip, ubah, err := RencanaPindah(sumber, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if lap.CacahSumber != 6 || lap.CacahFlat != 0 || len(sisip) != 3 || len(ubah) != 0 || !slices.Equal(lap.Baru, []string{"101", "102", "104"}) {
		t.Errorf("laporan %+v sisip %d", lap, len(sisip))
	}
	if !sisip[0].Sama(RingkasanFlat{ID: "101", UsedBy: p("UJI RATE A"), Type: p("L"), ModifiedDate: p("20240102T030405.000 GMT"),
		OperatorID: p("UJI-A"), Flag: p("AP")}) || *sisip[1].Flag != "PM" || *sisip[2].Flag != "PY" || sisip[2].Type != nil {
		t.Errorf("FLAG/TYPE apa adanya %+v", sisip)
	}
	if lap.PanjangMaks["USEDBY"] != 501 || lap.PanjangMaks["MODIFIEDDATE"] != 23 || lap.PanjangMaks["ID"] != 11 || lap.PanjangMaks["FLAG"] != 2 {
		t.Errorf("panjang maks %v", lap.PanjangMaks)
	}
	if len(lap.TidakMuat) != 2 || lap.TidakMuat[0] != "103 USEDBY: 501 byte > 500" || lap.TidakMuat[1] != "12345678901 ID: 11 byte > 10" {
		t.Errorf("tidak muat %q", lap.TidakMuat)
	}
	if len(lap.Gagal) != 1 || lap.Gagal[0] != "106 JSONDATA: bukan objek JSON" || lap.BolehDitulis() {
		t.Errorf("gagal %q", lap.Gagal)
	}
	if strings.Contains(lap.Teks(), panjang) || strings.Contains(lap.Teks(), "UJI RATE") {
		t.Error("laporan memuat nilai")
	}
	if lap.BatasBerikut != "20240102T030405.000 GMT" {
		t.Errorf("batas berikut %q", lap.BatasBerikut)
	}
}

// Delta: baru / berubah (sumber lebih baru) / sama / konflik (flat lebih baru, sama waktu, atau stempel tak terbaca -
// TIDAK ditimpa) / dilewati (lebih lama dari -sejak dan tidak ada di flat = dihapus aplikasi) / flat saja.
func TestRencanaPindahDelta(t *testing.T) {
	const lama, tengah, baru = "20261001T000000.000 GMT", "20261005T000000.000 GMT", "20261006T040628.169 GMT"
	sumber := []BarisJSON{
		{ID: "201", JSON: jsonRing("UJI SAMA", lama, "AP")},
		{ID: "202", JSON: jsonRing("UJI PEGA UBAH", baru, "PM")},      // flat tengah -> berubah
		{ID: "203", JSON: jsonRing("UJI APLIKASI UBAH", lama, "AP")},  // flat baru (aplikasi) -> konflik
		{ID: "204", JSON: jsonRing("UJI WAKTU SAMA", tengah, "AP")},   // isi beda, waktu sama -> konflik
		{ID: "205", JSON: jsonRing("UJI TANPA WAKTU", "", "AP")},      // flat beda, stempel kosong -> konflik
		{ID: "206", JSON: jsonRing("UJI PEGA BARU", baru, "PY")},      // tidak di flat, >= sejak -> baru
		{ID: "207", JSON: jsonRing("UJI DIHAPUS APLIKASI", lama, "")}, // tidak di flat, < sejak -> dilewati
		{ID: "208", JSON: jsonRing("UJI LAMA TANPA WAKTU", "", "")},   // tidak di flat, tanpa stempel -> dilewati
	}
	flat := []RingkasanFlat{
		{ID: "201", UsedBy: p("UJI SAMA"), ModifiedDate: p(lama), Flag: p("AP")},
		{ID: "202", UsedBy: p("UJI PEGA"), ModifiedDate: p(tengah), Flag: p("AP")},
		{ID: "203", UsedBy: p("UJI APLIKASI SUDAH UBAH"), ModifiedDate: p(baru), Flag: p("AP")},
		{ID: "204", UsedBy: p("UJI WAKTU SAMA X"), ModifiedDate: p(tengah), Flag: p("AP")},
		{ID: "205", UsedBy: p("UJI LAIN"), Flag: p("AP")},
		{ID: "900", UsedBy: p("UJI RINGKASAN APLIKASI"), ModifiedDate: p(baru)},
	}
	lap, sisip, ubah, err := RencanaPindah(sumber, flat, tengah)
	if err != nil {
		t.Fatal(err)
	}
	for nama, k := range map[string][2][]string{
		"baru":     {lap.Baru, {"206"}},
		"berubah":  {lap.Berubah, {"202"}},
		"sama":     {lap.Sama, {"201"}},
		"konflik":  {lap.Konflik, {"203", "204", "205"}},
		"dilewati": {lap.Dilewati, {"207", "208"}},
	} {
		if !slices.Equal(k[0], k[1]) {
			t.Errorf("%s %v, mau %v", nama, k[0], k[1])
		}
	}
	if lap.FlatSaja != 1 || lap.CacahSumber != 8 || lap.CacahFlat != 6 || !lap.BolehDitulis() || lap.BatasBerikut != baru {
		t.Errorf("laporan %+v", lap)
	}
	if len(sisip) != 1 || sisip[0].ID != "206" || len(ubah) != 1 || ubah[0].ID != "202" || *ubah[0].Flag != "PM" || *ubah[0].UsedBy != "UJI PEGA UBAH" {
		t.Errorf("sisip %+v ubah %+v", sisip, ubah)
	}
	// Tanpa -sejak: ringkasan yang tidak ada di flat disisip (pemindahan penuh pertama).
	if lap2, s2, _, _ := RencanaPindah(sumber, flat, ""); len(lap2.Dilewati) != 0 || len(s2) != 3 {
		t.Errorf("tanpa -sejak %+v", lap2)
	}
	// Ulang sesudah menulis: semua sama, nol tulisan.
	semua := append(append([]RingkasanFlat{}, flat...), sisip...)
	for i := range semua {
		if semua[i].ID == "202" {
			semua[i] = ubah[0]
		}
	}
	if lap3, s3, u3, _ := RencanaPindah(sumber[:2], semua, ""); len(s3)+len(u3) != 0 || len(lap3.Sama) != 2 {
		t.Errorf("ulang %+v", lap3)
	}
	if _, _, _, err := RencanaPindah(sumber, flat, "kemarin"); !errors.Is(err, ErrSejakTidakTerbaca) {
		t.Errorf("-sejak salah %v", err)
	}
	if !strings.Contains(lap.Teks(), "konflik (TIDAK ditimpa, periksa): 3 203 204 205") || !strings.Contains(lap.Teks(), "batas delta berikutnya (-sejak): "+baru) {
		t.Errorf("teks:\n%s", lap.Teks())
	}
	if err := PeriksaTulis(TabelRingkasan, SqlPerbaruiRingkasanFlat("S.RATE_LIFE_SUMMARY")); err != nil {
		t.Error(err)
	}
	if err := PeriksaTulis(TabelRingkasan, SqlKunciRingkasan("S.RATE_LIFE_SUMMARY")); err != nil {
		t.Error(err)
	}
}
