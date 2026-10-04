package lintasmodul_test

// Tiket R01 (rnwfacin) atas mesin NB (nbfacin) lewat kontrak.TanggaAkseptasiFacIn:
// kasus renewal yang dibuat layanan RNW diterima tangga akseptasi NB TANPA
// penyesuaian - hasilnya identik dengan memanggil tangga NB langsung atas kasus
// yang sama. Masukan: fixture ter-de-identifikasi (tiket NB-15) dan tabel limit
// nyata (butir 42); tanpa basis data.

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/services/acceptance"
	"nusantarare/modul/nbfacin/backend/services/kontrakfacin"
	"nusantarare/modul/rnwfacin/backend/models"
	rnw "nusantarare/modul/rnwfacin/backend/services"
)

const (
	folderLimit = "../../modul/nbfacin/backend/services/acceptance/testdata/limit"
	folderKasus = "../../modul/nbfacin/backend/services/premium/testdata/kasus"
)

func bacaCSV(t *testing.T, nama string) []map[string]string {
	t.Helper()
	f, err := os.Open(filepath.Join(folderLimit, nama+".csv"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = ';'
	baris, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	var hasil []map[string]string
	for _, b := range baris[1:] {
		m := map[string]string{}
		for i, k := range baris[0] {
			m[k] = strings.ReplaceAll(b[i], ".", "") // titik = pemisah ribuan (README fixture)
		}
		hasil = append(hasil, m)
	}
	return hasil
}

func angka(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	d, err := utils.ParseDecimal(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func tanggaNB(t *testing.T) kontrakfacin.Tangga {
	tabel := acceptance.TabelLimit{}
	for _, n := range []acceptance.NamaTabel{acceptance.TabelProperty, acceptance.TabelPropertyNonPreferred,
		acceptance.TabelPropertyPreferredCommercial, acceptance.TabelEngineering, acceptance.TabelNonPropEng} {
		for _, b := range bacaCSV(t, string(n)) {
			tabel[n] = append(tabel[n], acceptance.BarisLimit{Jabatan: acceptance.Jabatan(b["JABATAN"]), TeamGroup: b["TEAM_GROUP"],
				LimitBottom: angka(t, b["LIMIT_BOTTOM"]), LimitBottom2: angka(t, b["LIMIT_BOTTOM2"])})
		}
	}
	var fin []acceptance.BarisFinancial
	for _, b := range bacaCSV(t, "M_LIMIT_FINANCIALINS") {
		fin = append(fin, acceptance.BarisFinancial{Jabatan: acceptance.Jabatan(b["JABATAN"]), LimitBond: angka(t, b["LIMITBOND_BOTTOM"]),
			LimitCreditCL: angka(t, b["LIMITCREDITCL_BOTTOM"]), LimitCreditNCL: angka(t, b["LIMITCREDITNCL_BOTTOM"])})
	}
	return kontrakfacin.Tangga{BentukA: tabel, BentukB: fin}
}

// kasusDariFixture meratakan fixture (akar = halaman OfferFacIn) menjadi jalur
// `pyWorkPage.OfferFacIn.*`; `pyWorkPage.Quotation.*` diisi dari QuotationData
// (butir 49). Daftar (larik) tidak dibaca tangga, jadi dilewati.
func kasusDariFixture(t *testing.T, nama string) models.Kasus {
	t.Helper()
	isi, err := os.ReadFile(filepath.Join(folderKasus, nama+".json"))
	if err != nil {
		t.Fatal(err)
	}
	var akar map[string]any
	if err := json.Unmarshal(isi, &akar); err != nil {
		t.Fatal(err)
	}
	k := models.Kasus{}
	var ratakan func(awalan string, m map[string]any)
	ratakan = func(awalan string, m map[string]any) {
		for kunci, v := range m {
			switch x := v.(type) {
			case string:
				k[awalan+kunci] = x
			case map[string]any:
				ratakan(awalan+kunci+".", x)
			}
		}
	}
	ratakan("pyWorkPage.OfferFacIn.", akar)
	if q, ok := akar["QuotationData"].(map[string]any); ok {
		ratakan("pyWorkPage.Quotation.", q)
	}
	return k
}

type polisFixture map[string]models.Kasus

func (p polisFixture) Baca(_ context.Context, nomor string) (models.Kasus, error) {
	return p[nomor], nil
}

// TestRenewalDiterimaTanggaNBTanpaPenyesuaian - kriteria R01.
func TestRenewalDiterimaTanggaNBTanpaPenyesuaian(t *testing.T) {
	tg := tanggaNB(t)
	naik := 0
	for _, u := range []struct {
		fixture  string
		tsi      string // menimpa TotalTSINusaRe bila terisi, supaya tangga benar-benar naik
		pengguna kontrak.PenggunaFacIn
	}{
		{"rnw-fire-1", "", kontrak.PenggunaFacIn{Jabatan: "SENIORUW"}},
		{"rnw-fire-1", "400000000000", kontrak.PenggunaFacIn{Jabatan: "SENIORUW"}},
		{"nb-fire-1", "400000000000", kontrak.PenggunaFacIn{Jabatan: "UNDERWRITER"}},
		{"nb-fire-1", "", kontrak.PenggunaFacIn{Jabatan: "SENIORUW", AnggotaGrup: true}},
	} {
		lama := kasusDariFixture(t, u.fixture)
		if u.tsi != "" {
			lama["pyWorkPage.OfferFacIn.TotalTSINusaRe"] = u.tsi
		}
		s := rnw.NewService(polisFixture{"POLIS-REKAAN": lama}, tg)
		k, err := s.BuatKasusRenewal(context.Background(), rnw.Masukan{NomorPolis: "POLIS-REKAAN", TanggalRenewal: "2026-12-01"})
		if err != nil {
			t.Fatalf("%s: %v", u.fixture, err)
		}
		lewatRNW, errRNW := s.LangkahAkseptasi(k, u.pengguna)
		langsung, errNB := tg.Langkah(k.Kasus, u.pengguna)
		if errRNW != nil || errNB != nil || lewatRNW != langsung {
			t.Errorf("%s %+v: lewat RNW %+v (%v), NB langsung %+v (%v)", u.fixture, u.pengguna, lewatRNW, errRNW, langsung, errNB)
		}
		if !lewatRNW.Selesai {
			naik++
		}
		t.Logf("%s TSI %q %+v → %+v", u.fixture, u.tsi, u.pengguna, lewatRNW)
	}
	if naik == 0 {
		t.Fatal("tidak satu kasus pun naik tangga - perbandingan tidak menguji apa pun")
	}
}
