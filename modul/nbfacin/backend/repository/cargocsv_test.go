package repository

import (
	"context"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// TestKapalMasterSQL - tiket 74 P6: BrowseShip_RD / BrowseShip2_RD (aktif = REMARK != 'Not Active' OR NULL; nama tak peka huruf,
// ID terbesar), INSERT RDBMASTERSHIP (kolom terisi saja); bind berurutan; ID kapal baru '1' || LPAD(n, 9, '0').
func TestKapalMasterSQL(t *testing.T) {
	for q, mau := range map[string]int{sqlKapalAktifID("S.SHIP"): 2, sqlKapalNama("S.SHIP"): 2, sqlDaftarKapal("S.SHIP"): 5} {
		if n := bindTerbesar(t, q); n != mau {
			t.Errorf("%d bind, mau %d: %s", n, mau, q)
		}
		if err := db.PeriksaSQL(q); err != nil {
			t.Error(err)
		}
	}
	if q := sqlKapalNama("S.SHIP"); !strings.Contains(q, "UPPER(NAME) = UPPER(:1) AND (REMARK != :2 OR REMARK IS NULL) ORDER BY ID DESC FETCH FIRST 1") {
		t.Error(q)
	}
	if q := sqlDaftarKapal("S.SHIP"); !strings.HasSuffix(q, "(ID, NAME, REMARK, TYPE, OPERATOR_ID, INPUT_TIME) VALUES (:1, :2, :3, :4, :5, SYSDATE)") {
		t.Error(q)
	}
	for n, mau := range map[string]string{"7": "1000000007", "123456789": "1123456789", "1234567890": "1123456789"} {
		if got := idKapalBaru(n); got != mau {
			t.Errorf("%s -> %s, mau %s", n, got, mau)
		}
	}
}

// TestSisipCargoBerkelompok - pohon cargo hasil CSV: SEQ_NO cargo lanjut dari `awal`, kapal + AdditionalShip + agen + coverage
// tersisip sekali per cargo lewat array-bind (kolom satu tipe Go, periksaKolomSeragam); SEQ_NO terakhir dibaca per case.
func TestSisipCargoBerkelompok(t *testing.T) {
	ctx := context.Background()
	tc, _ := tabelCargoDari(func(n string) (string, error) { return "S." + n, nil })
	k := models.KapalCargo{ID: "S1", NmShip: "UJI", GRT: "1.5"}
	k.AdditionalShip = []models.KapalTambahan{{ID: "S1", NmShip: "UJI"}}
	cargo := []models.ObjekCargo{
		{TradingID: "T1", TradingNote: "DOMESTIK", IsTopRisk: true, Ship: k, SurveyAgent: models.AgenSurvei{AgentID: "A1"},
			Coverages: []models.CoverageCargo{{TSI: apd.New(15, -1), Currency: "IDR"}}},
		{TradingID: "T2", TradingNote: "EKSPOR", Ship: k, Coverages: []models.CoverageCargo{{Currency: "CNY"}}},
	}
	e := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}, seragam: true}
	if err := sisipCargoBerkelompok(ctx, e, tc, "UJI-NB-1", cargo, uidUrut(), 40); err != nil {
		t.Fatal(err)
	}
	per := map[string]int{}
	var seq []any
	for _, b := range e.baris {
		tabel := strings.Fields(strings.TrimPrefix(b.q, "INSERT INTO "))[0]
		per[tabel]++
		if tabel == tc.cargo {
			seq = append(seq, b.arg[2])
		}
	}
	if per[tc.cargo] != 2 || per[tc.ship] != 2 || per[tc.kapalTambahan] != 2 || per[tc.agen] != 1 || per[tc.cov] != 2 || per[tc.lc] != 0 {
		t.Errorf("baris per tabel: %v", per)
	}
	if len(seq) != 2 || seq[0] != 41 || seq[1] != 42 {
		t.Errorf("SEQ_NO: %v", seq)
	}
	if q := sqlSeqCargoTerakhir("S.T_CARGOLIST"); q != "SELECT NVL(MAX(SEQ_NO), 0) FROM S.T_CARGOLIST WHERE PARENT_ID = :1" {
		t.Error(q)
	}
}
