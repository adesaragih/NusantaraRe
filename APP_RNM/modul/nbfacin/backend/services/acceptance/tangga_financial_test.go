package acceptance

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

// muatFinancial - M_LIMIT_FINANCIALINS dari fixture `testdata/limit`:
// SENIOR UNDERWRITER 0/0/0, KADIV KEUANGAN 0/0/0, DIREKTUR MARKETING
// 5.000.000.001/0/0, DIREKTUR TEKNIK 10.000.000.001/5.000.000.001/10.000.000.001
// (LIMITBOND_BOTTOM / LIMITCREDITCL_BOTTOM / LIMITCREDITNCL_BOTTOM).
func muatFinancial(t *testing.T) []BarisFinancial {
	t.Helper()
	baris := bacaCSVLimit(t, filepath.Join(folderLimit, "M_LIMIT_FINANCIALINS.csv"))
	kol := map[string]int{}
	for i, k := range baris[0] {
		kol[k] = i
	}
	tanpaTitik := func(s string) string {
		if !angkaLimit.MatchString(s) {
			t.Fatalf("angka %q bukan bentuk ekspor yang dikenal", s)
		}
		return strings.ReplaceAll(s, ".", "")
	}
	var hasil []BarisFinancial
	for _, b := range baris[1:] {
		hasil = append(hasil, BarisFinancial{Jabatan: Jabatan(b[kol["JABATAN"]]),
			LimitBond:      desimal(t, tanpaTitik(b[kol["LIMITBOND_BOTTOM"]])),
			LimitCreditCL:  desimal(t, tanpaTitik(b[kol["LIMITCREDITCL_BOTTOM"]])),
			LimitCreditNCL: desimal(t, tanpaTitik(b[kol["LIMITCREDITNCL_BOTTOM"]]))})
	}
	return hasil
}

// bond - kasus Bond (BusinessOldId 36 → IsLimitSBondKBG), TSI 1 M.
var bond = kasusTangga{
	jTSI: "1000000000", jTopRisk: "0", jAdaTop: "false", jStatus: "1", jBisnis: "Bonding",
	jBisnisLama: "36", jBanding: "false", jReject: "false", jAntrean: "ReasFacInUnderwritingFinancial",
}

// TestNextFinancialSatuLangkah - langkah 22: satu langkah naik menurut antrean
// saat ini, atas daftar yang sudah disaring WHERE SQL (butir 46).
func TestNextFinancialSatuLangkah(t *testing.T) {
	tabel := muatFinancial(t)
	for _, u := range []struct {
		nama  string
		kasus kasusTangga
		mau   Transisi
	}{
		{"UW Financial → Kadiv Keuangan", bond,
			Transisi{JabatanTujuan: "KADIVFINANCIAL", Antrean: "ReasFacInFinDivHead", PositionNoteDitulis: true}},
		{"Kadiv Keuangan → Direktur Marketing", bond.dengan(jAntrean, "ReasFacInFinDivHead", jTSI, "6000000000"),
			Transisi{JabatanTujuan: "DIREKTURMARKETING", Antrean: "ReasFacInMarketingDirector", PositionNoteDitulis: true}},
		{"Direktur Marketing → Direktur Teknik", bond.dengan(jAntrean, "ReasFacInMarketingDirector", jTSI, "20000000000"),
			Transisi{JabatanTujuan: "DIREKTURTEKNIK", Antrean: "ReasFacInTechnicalDirector", PositionNoteDitulis: true}},
		// WHERE LIMITBOND_BOTTOM <= CARID2: 5 M < 5.000.000.001 → DIREKTUR MARKETING tidak di daftar.
		{"Kadiv Keuangan, TSI di bawah limit Direktur Marketing", bond.dengan(jAntrean, "ReasFacInFinDivHead", jTSI, "5000000000"),
			Transisi{Selesai: true}},
		// Batas inklusif: LIMITBOND_BOTTOM 5.000.000.001 <= CARID2 5.000.000.001.
		{"Kadiv Keuangan, TSI tepat di limit Direktur Marketing", bond.dengan(jAntrean, "ReasFacInFinDivHead", jTSI, "5000000001"),
			Transisi{JabatanTujuan: "DIREKTURMARKETING", Antrean: "ReasFacInMarketingDirector", PositionNoteDitulis: true}},
		{"Direktur Teknik, tidak ada gerbang", bond.dengan(jAntrean, "ReasFacInTechnicalDirector", jTSI, "20000000000"),
			Transisi{Selesai: true}},
		// Kredit Cash Loan: kolom LIMITCREDITCL_BOTTOM (DIREKTUR TEKNIK 5.000.000.001).
		{"Kredit CL", bond.dengan(jBisnisLama, "C2", jAntrean, "ReasFacInMarketingDirector", jTSI, "6000000000"),
			Transisi{JabatanTujuan: "DIREKTURTEKNIK", Antrean: "ReasFacInTechnicalDirector", PositionNoteDitulis: true}},
		// Trade Credit memakai SQL Kredit Cash Loan (langkah 13).
		{"Trade Credit", bond.dengan(jBisnisLama, "49", jAntrean, "ReasFacInMarketingDirector", jTSI, "6000000000"),
			Transisi{JabatanTujuan: "DIREKTURTEKNIK", Antrean: "ReasFacInTechnicalDirector", PositionNoteDitulis: true}},
		// Kredit Non Cash Loan: LIMITCREDITNCL_BOTTOM 10.000.000.001 > 6 M.
		{"Kredit NCL", bond.dengan(jBisnisLama, "46", jAntrean, "ReasFacInMarketingDirector", jTSI, "6000000000"),
			Transisi{Selesai: true}},
	} {
		got, err := NextFinancial(u.kasus, tabel, Pengguna{})
		if err != nil || got != u.mau {
			t.Errorf("%s: dapat %+v (%v), mau %+v", u.nama, got, err, u.mau)
		}
	}
	if got, err := NextFinancial(bond, tabel, Pengguna{AnggotaGrup: true}); err != nil || got != (Transisi{Selesai: true}) {
		t.Errorf("anggota grup: dapat %+v (%v)", got, err)
	}
}

// TestDuaJalurTerpisah - bentuk standar tidak dipakai untuk lini financial, dan
// sebaliknya (tiket 12).
func TestDuaJalurTerpisah(t *testing.T) {
	if _, err := Next(bond, muatLimit(t), penggunaBiasa("SENIORUW")); !errors.Is(err, ErrBentukB) {
		t.Errorf("Next atas kasus Bond: galat %v, mau ErrBentukB", err)
	}
	if _, err := NextFinancial(fireTG1, muatFinancial(t), Pengguna{}); !errors.Is(err, ErrBukanBentukB) {
		t.Errorf("NextFinancial atas kasus FIRE: galat %v, mau ErrBukanBentukB", err)
	}
}

// TestNextFinancialEjaanBerspasi - pencocokan jabatan memakai ejaan berspasi
// tabel FINANCIALINS; tidak ada normalisasi.
func TestNextFinancialEjaanBerspasi(t *testing.T) {
	tanpaSpasi := []BarisFinancial{{Jabatan: "KADIVKEUANGAN", LimitBond: desimal(t, "0"),
		LimitCreditCL: desimal(t, "0"), LimitCreditNCL: desimal(t, "0")}}
	if got, err := NextFinancial(bond, tanpaSpasi, Pengguna{}); err != nil || got != (Transisi{Selesai: true}) {
		t.Errorf("KADIVKEUANGAN tanpa spasi: dapat %+v (%v), mau selesai", got, err)
	}
}

// TestNextFinancialTafsirLimit - langkah 22.2.1 `Local.Limit = .CARID2` membaca
// kolom yang tidak dipilih SQL (limit dipetakan ke CARI2). Hanya berpengaruh bila
// Local.TotalTSI ≠ CARID2 - endorsement dengan selisih top risk - dan kedua tafsir
// berbeda hasil → galat (A26).
func TestNextFinancialTafsirLimit(t *testing.T) {
	// CARID2 = |20 M − 0| (top risk), TotalTSI = |2 M − 1 M| = 1 M. DIREKTUR TEKNIK
	// (10.000.000.001 ≤ 20 M) masuk daftar; 1 M ≥ Limit benar bila Limit kosong
	// = 0, salah bila Limit = 10.000.000.001.
	endorse := bond.dengan(jAntrean, "ReasFacInMarketingDirector", jStatus, "3", jTSI, "2000000000",
		jTSILama, "1000000000", jTopRisk, "20000000000", jTopLama, "0")
	if _, err := NextFinancial(endorse, muatFinancial(t), Pengguna{}); !errors.Is(err, ErrTafsirLimit) {
		t.Errorf("galat %v, mau ErrTafsirLimit", err)
	}
	// TotalTSI = CARID2 → kedua tafsir sama.
	if got, err := NextFinancial(endorse.dengan(jTopRisk, "0"), muatFinancial(t), Pengguna{}); err != nil || got != (Transisi{Selesai: true}) {
		t.Errorf("tanpa selisih top risk: dapat %+v (%v)", got, err)
	}
}

// TestLangkah3SebelumLangkah4Dan5 - urutan korpus: anggota grup keluar di langkah 3,
// sebelum langkah 4-5 membaca top risk dan OldData; data itu tidak boleh membuat
// galat (temuan tinjauan NB-12).
func TestLangkah3SebelumLangkah4Dan5(t *testing.T) {
	rusak := kasusTangga{jTSI: "1000000000", jTopRisk: "", jStatus: "3", jTSILama: "", jBisnisLama: "36"}
	grup := Pengguna{Jabatan: "SENIORUW", AnggotaGrup: true}
	if got, err := Next(rusak, muatLimit(t), grup); err != nil || got != (Transisi{Selesai: true}) {
		t.Errorf("Next: dapat %+v (%v)", got, err)
	}
	if got, err := NextFinancial(rusak, muatFinancial(t), grup); err != nil || got != (Transisi{Selesai: true}) {
		t.Errorf("NextFinancial: dapat %+v (%v)", got, err)
	}
	// Langkah 2 tetap mendahului langkah 3.
	if _, err := Next(rusak.dengan(jTSI, ""), muatLimit(t), grup); err == nil {
		t.Error("TSI kosong: mau galat langkah 2")
	}
}

// TestNextFinancialDaftarKosong - tidak ada baris lolos WHERE → 22.1 LetterNo = ""
// → Decision23 Else → selesai.
func TestNextFinancialDaftarKosong(t *testing.T) {
	if got, err := NextFinancial(bond, nil, Pengguna{}); err != nil || got != (Transisi{Selesai: true}) {
		t.Errorf("tabel kosong: dapat %+v (%v)", got, err)
	}
}

// TestLimitKosongDitolak - kolom limit kosong tidak ditebak sebagai NULL SQL.
func TestLimitKosongDitolak(t *testing.T) {
	fin := []BarisFinancial{{Jabatan: "KADIV KEUANGAN", LimitCreditCL: desimal(t, "0"), LimitCreditNCL: desimal(t, "0")}}
	if _, err := NextFinancial(bond, fin, Pengguna{}); !errors.Is(err, ErrLimitKosong) {
		t.Errorf("bentuk B: galat %v, mau ErrLimitKosong", err)
	}
	// Kontrol: LIMIT_BOTTOM2 kosong tidak dibaca jalur non-banding.
	a := TabelLimit{TabelProperty: {{Jabatan: "SENIORUW", TeamGroup: "1", LimitBottom: desimal(t, "1")}}}
	if _, err := Next(fireTG1, a, penggunaBiasa("SENIORUW")); err != nil {
		t.Fatalf("kontrol: %v", err)
	}
	a[TabelProperty] = append(a[TabelProperty], BarisLimit{Jabatan: "KADIVTEKNIK", TeamGroup: "1", LimitBottom2: desimal(t, "5")})
	if _, err := Next(fireTG1, a, penggunaBiasa("SENIORUW")); !errors.Is(err, ErrLimitKosong) {
		t.Errorf("bentuk A: galat %v, mau ErrLimitKosong", err)
	}
}
