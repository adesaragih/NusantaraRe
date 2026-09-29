package repository

// Cabut peserta = PENANDA - OQ-M6 ditutup 29-09-2026 (GILIRAN-17).
//
// Tanpa Oracle: bentuk migrasi 022, bentuk penulisnya, dan penjaga bahwa
// SETIAP pembaca/penulis tabel peserta menyaring penandanya.

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestMigrasi022PenandaCabutPeserta(t *testing.T) {
	naik, err := os.ReadFile("migrations/022_kolom_sts_hapus_peserta.sql")
	if err != nil {
		t.Fatal(err)
	}
	turun, err := os.ReadFile("migrations/022_kolom_sts_hapus_peserta_down.sql")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(naik), "ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL ADD (") ||
		!strings.Contains(string(naik), "STS_HAPUS VARCHAR2(1)") {
		t.Errorf("022 tidak menambah STS_HAPUS VARCHAR2(1):\n%s", naik)
	}
	if !strings.Contains(string(turun), "ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL DROP (") ||
		!strings.Contains(string(turun), "STS_HAPUS") {
		t.Errorf("022_down tidak membuang STS_HAPUS:\n%s", turun)
	}
}

// Penulisnya menandai, tidak menghapus - dan hanya sekali, hanya peserta
// milik klaim itu (ADR-U-0031).
func TestSQLCabutPesertaMenandaiBukanMenghapus(t *testing.T) {
	q := sqlCabutPeserta("S.P")
	for _, mau := range []string{"UPDATE S.P SET STS_HAPUS = '1'",
		"WHERE ID = :1 AND CLAIM_ID = :2 AND STS_HAPUS IS NULL"} {
		if !strings.Contains(q, mau) {
			t.Errorf("tanpa %q:\n%s", mau, q)
		}
	}
	if strings.Contains(strings.ToUpper(q), "DELETE") {
		t.Error("cabut peserta menghapus baris")
	}
	if err := PeriksaSQL(q); err != nil {
		t.Error(err)
	}
}

// pengecualianSaringCabut - fungsi yang menyentuh tabel peserta TANPA
// menyaring penandanya, masing-masing dengan alasannya.
var pengecualianSaringCabut = map[string]string{
	"Simpan":               "menyisip peserta baru saat pendaftaran - belum ada yang dapat dicabut",
	"HapusFisik":           "hanya dipakai uji - membersihkan klaim uji utuh",
	"Dampak":               "dampak hapus klaim UTUH (masih 405) - peserta tercabut ikut terhapus",
	"tabel":                "pembantu nama tabel diagnosa; SQL-nya di sqlAmbilDiagnosa (diperiksa di bawah)",
	"PerbaruiTanggalKlaim": "SQL-nya di sqlTanggalKlaim (diperiksa di bawah)",
	"SudahSaveRNM":         "SQL-nya di sqlSudahSaveRNM (diperiksa di bawah)",
}

// wajibSaringCabut - pembangun SQL yang menyentuh tabel peserta lewat
// parameter, sehingga tidak tertangkap pencarian Qualify.
var wajibSaringCabut = []string{"sqlAmbilDiagnosa", "sqlTanggalKlaim", "sqlSudahSaveRNM"}

var polaFungsi = regexp.MustCompile(`(?m)^func (\([^)]*\) )?([A-Za-z0-9_]+)\(`)

// TestSetiapPenyentuhTabelPesertaMenyaringPenandaCabut - "semua pembaca
// menyaring penanda" (OQ-M6). Fungsi BARU yang menyentuh tabel peserta tanpa
// saringan gagal di sini sampai ia menyaring atau dikecualikan dengan alasan.
func TestSetiapPenyentuhTabelPesertaMenyaringPenandaCabut(t *testing.T) {
	masuk, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	badan := map[string]string{}
	for _, m := range masuk {
		if m.IsDir() || !strings.HasSuffix(m.Name(), ".go") || strings.HasSuffix(m.Name(), "_test.go") {
			continue
		}
		isi, err := os.ReadFile(m.Name())
		if err != nil {
			t.Fatal(err)
		}
		s := string(isi)
		lok := polaFungsi.FindAllStringSubmatchIndex(s, -1)
		for i, l := range lok {
			akhir := len(s)
			if i+1 < len(lok) {
				akhir = lok[i+1][0]
			}
			nama := s[l[4]:l[5]]
			badan[m.Name()+":"+nama] = s[l[0]:akhir]
		}
	}
	diperiksa := 0
	for kunci, b := range badan {
		if !strings.Contains(b, `Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")`) {
			continue
		}
		diperiksa++
		nama := kunci[strings.Index(kunci, ":")+1:]
		if _, ada := pengecualianSaringCabut[nama]; ada {
			continue
		}
		if !strings.Contains(b, "STS_HAPUS") {
			t.Errorf("%s menyentuh T_CLAIMLF_PREMIUMLIST_DETAIL tanpa menyaring STS_HAPUS", kunci)
		}
	}
	if diperiksa < 10 {
		t.Errorf("hanya %d fungsi diperiksa; pencariannya rusak", diperiksa)
	}
	for _, nama := range wajibSaringCabut {
		ketemu := false
		for kunci, b := range badan {
			if strings.HasSuffix(kunci, ":"+nama) {
				ketemu = true
				if !strings.Contains(b, "STS_HAPUS IS NULL") {
					t.Errorf("%s tidak menyaring STS_HAPUS IS NULL", kunci)
				}
			}
		}
		if !ketemu {
			t.Errorf("%s tidak ditemukan", nama)
		}
	}
}
