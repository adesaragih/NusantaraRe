package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// kolomSemuaMigrasi - kolom tiap tabel dari SELURUH migrasi maju nbfacin (CREATE TABLE + ALTER TABLE … ADD (…)); cara kedua
// terhadap tabel kolom Go (bukan daftar yang ditulis ulang).
func kolomSemuaMigrasi(t *testing.T) map[string]map[string]bool {
	t.Helper()
	berkas, err := filepath.Glob("../migrations/*.sql")
	if err != nil || len(berkas) == 0 {
		t.Fatalf("migrasi: %v", err)
	}
	buat := regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.([A-Z_0-9]+) \((.*?)\n\)`)
	tambah := regexp.MustCompile(`(?s)ALTER TABLE \{skema\}\.([A-Z_0-9]+) ADD \((.*?)\n\)`)
	hasil := map[string]map[string]bool{}
	for _, f := range berkas {
		if strings.HasSuffix(f, "_down.sql") {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		s := strings.ReplaceAll(string(b), "\r", "")
		for _, pola := range []*regexp.Regexp{buat, tambah} {
			for _, m := range pola.FindAllStringSubmatch(s, -1) {
				if hasil[m[1]] == nil {
					hasil[m[1]] = map[string]bool{}
				}
				for _, baris := range strings.Split(m[2], "\n") {
					if f := strings.Fields(baris); len(f) > 0 && f[0] != "CONSTRAINT" {
						hasil[m[1]][f[0]] = true
					}
				}
			}
		}
	}
	return hasil
}

// TestKolomCargoAdaDiMigrasi - setiap kolom yang ditulis / dibaca jalur cargo ada di migrasinya (256 + 193 / 194 / 198 / 214),
// dan setiap kolom tabel BARU 256 dipakai (A\B dan B\A).
func TestKolomCargoAdaDiMigrasi(t *testing.T) {
	mig := kolomSemuaMigrasi(t)
	pakai := map[string][]string{
		TabelCargoList:         append([]string{"ID", "PARENT_ID", "SEQ_NO", "ROW_UID", "IS_TOP_RISK"}, namaKolom(kolomCargoData)...),
		TabelShip256:           append([]string{"ID", "PARENT_ID"}, namaKolom(kolomShip)...),
		TabelLC:                append([]string{"ID", "PARENT_ID", "LC"}, namaKolom(kolomLC)...),
		TabelSurveyAgent:       append([]string{"ID", "PARENT_ID"}, namaKolom(kolomAgen)...),
		TabelAdditionalShip:    append([]string{"ID", "PARENT_ID", "SEQ_NO", "ROW_UID"}, namaKolom(kolomKapalTambahan)...),
		TabelAdditionalCoverag: append([]string{"ID", "PARENT_ID", "SEQ_NO", "ROW_UID"}, namaKolom(kolomWarranty)...),
		TabelCoverageList:      append([]string{"CURRENCY_CODE"}, namaKolom(kolomCoverageCargoData)...),
		TabelDeductibleList:    append([]string{"CURRENCY"}, namaKolom(kolomDeductibleCargoData)...),
	}
	for tabel, kolom := range pakai {
		for _, k := range kolom {
			if !mig[tabel][k] {
				t.Errorf("%s.%s dipakai kode, tidak ada di migrasi", tabel, k)
			}
		}
	}
	kolomSistem := map[string]bool{"IDPEGA": true, "COB_GROUP": true}
	for _, tabel := range []string{TabelCargoList, TabelShip256, TabelLC, TabelSurveyAgent, TabelAdditionalShip, TabelAdditionalCoverag} {
		dipakai := map[string]bool{}
		for _, k := range pakai[tabel] {
			dipakai[k] = true
		}
		for k := range mig[tabel] {
			if !dipakai[k] && !kolomSistem[k] {
				t.Errorf("%s.%s di migrasi 256 tidak dipakai kode", tabel, k)
			}
		}
	}
}

// TestSQLCargo - bind berurutan tanpa celah; baca / hapus satu bind ID case; hapus anak sebelum induk.
func TestSQLCargo(t *testing.T) {
	tc, err := tabelCargoDari(func(n string) (string, error) { return "S." + n, nil })
	if err != nil {
		t.Fatal(err)
	}
	for nama, q := range map[string]string{"cargo": sqlBacaCargo(tc), "tambahan": sqlBacaKapalTambahan(tc), "cov": sqlBacaCoverageCargo(tc),
		"ded": sqlBacaDeductibleCargo(tc), "war": sqlBacaWarranty(tc), "spread": sqlBacaSpreadingCargo(tc)} {
		if n := bindTerbesar(t, q); n != 1 {
			t.Errorf("%s: %d bind: %s", nama, n, q)
		}
	}
	hapus := sqlHapusCargo(tc)
	urut := []string{"S.T_SPREADINGLIST", "S.T_DEDUCTIBLELIST", "S.T_ADDITIONALCOVERAGE", "S.T_COVERAGELIST", "S.T_ADDITIONALSHIP",
		"S.T_SHIP", "S.T_LC", "S.T_SURVEYAGENT", "S.T_CARGOLIST"}
	if len(hapus) != len(urut) {
		t.Fatalf("%d hapus", len(hapus))
	}
	for i, q := range hapus {
		if !strings.HasPrefix(q, "DELETE FROM "+urut[i]+" ") || bindTerbesar(t, q) != 1 {
			t.Errorf("hapus %d: %s", i, q)
		}
	}
	if !strings.Contains(hapus[3], "PARENT_TABLE = 'T_CARGOLIST'") || !strings.Contains(sqlBacaCoverageCargo(tc), "v.PARENT_TABLE = 'T_CARGOLIST'") {
		t.Error("coverage cargo tanpa saringan PARENT_TABLE")
	}
}

// TestTulisCargo - traversal sisip: setiap baris lolos periksaBindBaris (bentuk bind Oracle), bind = jumlah nilai, induk
// tersambung; halaman kosong tidak disisipkan; coverage ber-PARENT_TABLE T_CARGOLIST; mata uang opsional; DESCRIPTIONS terikat.
func TestTulisCargo(t *testing.T) {
	tc, _ := tabelCargoDari(func(n string) (string, error) { return n, nil })
	d := desimalUji
	cargo := []models.ObjekCargo{
		{IsTopRisk: true, TradingID: "T1", GoodID: "G1", SailDate: "20261001T050000.000 GMT",
			Ship: models.KapalCargo{ID: "S9", NmShip: "UJI", GRT: "1.5", AdditionalShip: []models.KapalTambahan{{Register: "R1"}, {Register: "R2"}}},
			LC:   models.LCCargo{LC: true, Bank: "B"}, SurveyAgent: models.AgenSurvei{AgentID: "A1"},
			Coverages: []models.CoverageCargo{{Coverage: "C1", Currency: "IDR", TSI: d("500000000000000"), Rate: d("0.1"), Premium: d("500000000000"),
				WarrantyExtension: []models.WarrantyExtension{{Coverage: "W1", CoverageNote: "N"}},
				Deductibles: []models.DeductibleCargo{{Deductible: models.Deductible{TypeDeductible: "1", PctDeductible: d("10")}, Descriptions: "DESK"},
					{Deductible: models.Deductible{Currency: "USD"}}},
				Spreading: []models.BarisSpreading{{TreatyType: "QS", SharePercentage: d("100")}}}}},
		{Coverages: []models.CoverageCargo{{Coverage: "C2"}}},
		{Ship: models.KapalCargo{AdditionalShip: []models.KapalTambahan{{Register: "R3"}}}}, // kapal hanya Additional Ship
	}
	type sisip struct {
		q, tabel string
		arg      []any
	}
	var semua []sisip
	n := 0
	w := &tulisCargo{t: tc, nomor: func(seq string) (string, error) { n++; return fmt.Sprint(n), nil }, uid: uidUrut(),
		sisip: func(q, tabel string, arg []any) error {
			if err := periksaBindBaris(q, tabel, [][]any{arg}, 0); err != nil {
				return err
			}
			if b := bindTerbesar(t, q); b != len(arg) {
				t.Errorf("%s: %d bind, %d nilai", tabel, b, len(arg))
			}
			semua = append(semua, sisip{q, tabel, arg})
			return nil
		}}
	if err := w.sisipCargo("K1", cargo); err != nil {
		t.Fatal(err)
	}
	per := map[string][]sisip{}
	for _, s := range semua {
		per[s.tabel] = append(per[s.tabel], s)
	}
	for tabel, mau := range map[string]int{TabelCargoList: 3, TabelShip256: 2, TabelAdditionalShip: 3, TabelLC: 1, TabelSurveyAgent: 1,
		TabelCoverageList: 2, TabelDeductibleList: 2, TabelAdditionalCoverag: 1, TabelSpreadingList: 1} {
		if len(per[tabel]) != mau {
			t.Errorf("%s: %d baris, mau %d", tabel, len(per[tabel]), mau)
		}
	}
	c1, c2 := per[TabelCargoList][0].arg, per[TabelCargoList][1].arg
	if c1[1] != "K1" || c1[2] != 1 || c1[4] != "true" || c2[2] != 2 || c2[4] != "false" {
		t.Errorf("cargo: %v / %v", c1[:5], c2[:5])
	}
	ship := per[TabelShip256][0].arg
	if ship[1] != c1[0] || per[TabelAdditionalShip][1].arg[1] != ship[0] || per[TabelAdditionalShip][1].arg[2] != 2 {
		t.Errorf("kapal: %v / %v", ship[:2], per[TabelAdditionalShip][1].arg[:3])
	}
	if lc := per[TabelLC][0].arg; lc[1] != c1[0] || lc[2] != "true" {
		t.Errorf("LC: %v", lc)
	}
	v1, v2 := per[TabelCoverageList][0], per[TabelCoverageList][1]
	if v1.arg[1] != c1[0] || v1.arg[2] != "T_CARGOLIST" || v1.arg[3] != "CargoList/CoverageList" || v1.arg[6] != "IDR" ||
		!strings.Contains(v1.q, "CURRENCY_CODE") || v2.arg[1] != c2[0] || strings.Contains(v2.q, "CURRENCY_CODE") {
		t.Errorf("coverage: %v / %v", v1.arg[:7], v2.arg[:6])
	}
	d1, d2 := per[TabelDeductibleList][0], per[TabelDeductibleList][1]
	if d1.arg[1] != v1.arg[0] || d1.arg[3] != "CargoList/CoverageList/DeductibleList" || d1.arg[len(d1.arg)-1] != "DESK" ||
		strings.Contains(d1.q, "CURRENCY,") || !strings.Contains(d2.q, "CURRENCY") || d2.arg[6] != "USD" || d2.arg[len(d2.arg)-1] != nil {
		t.Errorf("deductible: %v / %v", d1.arg, d2.arg)
	}
	if s := per[TabelSpreadingList][0].arg; s[1] != v1.arg[0] || s[3] != "CargoList/CoverageList/SpreadingList" || s[13] != "IDR" {
		t.Errorf("spreading: %v", s)
	}
	if wx := per[TabelAdditionalCoverag][0].arg; wx[1] != v1.arg[0] || wx[4] != "W1" || wx[5] != "N" {
		t.Errorf("warranty: %v", wx)
	}
	// Instrumen periksaBindBaris benar-benar memeriksa: TSI 9 desimal ditolak sebelum Oracle.
	salah := []models.ObjekCargo{{Coverages: []models.CoverageCargo{{TSI: d("1.123456789")}}}}
	if err := w.sisipCargo("K1", salah); err == nil || !strings.Contains(err.Error(), "TSI") {
		t.Errorf("TSI 9 desimal lolos: %v", err)
	}
}
