package premium

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"nusantarare/inti/backend/utils"
)

// coverageUji - satu coverage MBU fixture.
func coverageUji(tsi, rate, mataUang, flagDelete string) CoverageMBU {
	return CoverageMBU{Input: Input{LiniBisnis: LiniMBU, MataUang: mataUang, TSI: tsi, Rate: rate, ProRatePercentCoverage: "100"}, FlagDelete: flagDelete}
}

func teksPerMataUang(t *testing.T, cov []CoverageMBU, urutan []string) map[string]string {
	t.Helper()
	hasil, err := PremiPerMataUang(cov, urutan)
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]string{}
	for _, h := range hasil {
		m[h.MataUang] = utils.FormatDecimal(h.Premi.Amount)
	}
	return m
}

// samaNilai - kesamaan nilai desimal EKSAK, tanpa toleransi. Bukan kesamaan teks:
// CurrencyList.Premium tersimpan tanpa nol di belakang koma ("22536540"), sedangkan
// premi coverage berskala 4 ("6158880.0000"); tipe propertinya `[pertanyaan terbuka]`.
func samaNilai(t *testing.T, a, b string) bool {
	t.Helper()
	da, errA := utils.ParseDecimal(a)
	db, errB := utils.ParseDecimal(b)
	if errA != nil || errB != nil {
		t.Fatalf("angka tak terbaca: %q (%v), %q (%v)", a, errA, b, errB)
	}
	return da.Cmp(db) == 0
}

// TestRekonsiliasiPremiPerMataUangMBU - tiket NB-07, rekonsiliasi eksak pada nilai
// AGREGAT: CurrencyList.Premium lima kasus MBU NB nyata (`testdata/daftarizin/
// mbu_mata_uang.json`, 93 coverage). Empat cocok; satu berbeda - lihat di bawah.
func TestRekonsiliasiPremiPerMataUangMBU(t *testing.T) {
	isi, err := os.ReadFile("testdata/daftarizin/mbu_mata_uang.json")
	if err != nil {
		t.Fatal(err)
	}
	var kasus []struct {
		Kasus                string
		UrutanMasterMataUang []string
		Coverage             []struct{ TSI, Rate, Loading, ProRatePercentCoverage, MataUang, FlagDelete string }
		CurrencyList         []struct{ Name, Premium string }
	}
	if err := json.Unmarshal(isi, &kasus); err != nil {
		t.Fatal(err)
	}
	if len(kasus) != 5 {
		t.Fatalf("%d kasus, mau 5", len(kasus))
	}
	// `[pertanyaan terbuka]` kasus #5: premi keenam coverage cocok L1144 (jumlah
	// 201.301.628), tetapi CurrencyList.Premium tersimpan 309.703.644,8 - tidak
	// terturunkan dari data mana pun di berkas kasus (hipotesis yang gugur: register
	// NB-07). Port langkah 2.6 tetap mengikuti sistem lama (butir 54); dicatat apa
	// adanya, bukan dipaksa cocok.
	selisihTerbuka := map[string][2]string{"kasus MBU NB #5": {"201301628.0000", "309703644.8"}}
	for _, k := range kasus {
		var cov []CoverageMBU
		for _, c := range k.Coverage {
			cov = append(cov, CoverageMBU{Input: Input{LiniBisnis: LiniMBU, MataUang: c.MataUang, TSI: c.TSI, Rate: c.Rate,
				Loading: c.Loading, ProRatePercentCoverage: c.ProRatePercentCoverage}, FlagDelete: c.FlagDelete})
		}
		baru := teksPerMataUang(t, cov, k.UrutanMasterMataUang)
		for _, c := range k.CurrencyList {
			if s, ada := selisihTerbuka[k.Kasus]; ada {
				if baru[c.Name] != s[0] || c.Premium != s[1] {
					t.Errorf("%s %s: baru %s, lama %s - selisih terbuka berubah", k.Kasus, c.Name, baru[c.Name], c.Premium)
				}
				continue
			}
			if !samaNilai(t, baru[c.Name], c.Premium) {
				t.Errorf("%s %s: baru %s, sistem lama %s", k.Kasus, c.Name, baru[c.Name], c.Premium)
			}
		}
	}
}

// TestPremiPerMataUangPembulatanDiDalamLoop - premi tiap coverage dibulatkan
// (L1144, 4 desimal) SEBELUM dijumlahkan, seperti langkah 2.6.2.1.1. Dua coverage
// 0,00006 → 0,0001 + 0,0001 = 0,0002; membulatkan sekali di luar loop memberi 0,0001.
func TestPremiPerMataUangPembulatanDiDalamLoop(t *testing.T) {
	got := teksPerMataUang(t, []CoverageMBU{coverageUji("1", "0.006", "IDR", ""), coverageUji("1", "0.006", "IDR", "")}, []string{"IDR"})
	if got["IDR"] != "0.0002" {
		t.Fatalf("dapat %v, mau IDR 0.0002", got)
	}
}

// TestPremiPerMataUangAturanLangkah26 - langkah 2.6: urut tabel master; FlagDelete 1
// dilewati; mata uang tanpa coverage tidak ditambahkan; coverage bermata uang di luar
// master tidak dijumlahkan.
func TestPremiPerMataUangAturanLangkah26(t *testing.T) {
	cov := []CoverageMBU{
		coverageUji("1000", "1", "USD", ""),
		coverageUji("1000", "2", "IDR", "0"),
		coverageUji("1000", "5", "IDR", "1"), // FlagDelete 1
		coverageUji("1000", "3", "SGD", ""),  // di luar master
	}
	hasil, err := PremiPerMataUang(cov, []string{"IDR", "EUR", "USD"})
	if err != nil {
		t.Fatal(err)
	}
	if len(hasil) != 2 || hasil[0].MataUang != "IDR" || hasil[1].MataUang != "USD" {
		t.Fatalf("urutan/isi %+v, mau IDR lalu USD", hasil)
	}
	if utils.FormatDecimal(hasil[0].Premi.Amount) != "20.0000" || utils.FormatDecimal(hasil[1].Premi.Amount) != "10.0000" {
		t.Fatalf("premi %+v", hasil)
	}
	if _, err := PremiPerMataUang([]CoverageMBU{coverageUji("1", "1", "IDR", "ya")}, []string{"IDR"}); !errors.Is(err, ErrAngkaTakTerbaca) {
		t.Fatalf("FlagDelete bukan angka: galat %v, mau ErrAngkaTakTerbaca", err)
	}
}
