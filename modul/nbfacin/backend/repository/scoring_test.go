package repository

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

func tabelScoringUji() tabelScoring {
	return tabelScoring{work: "UJI.W", general: "UJI.G", sr: "UJI.SR", hasil: "UJI.H", data: "UJI.D", faktor: "UJI.F", opsi: "UJI.O",
		seqOpsi: "UJI.SEQ_O"}
}

// TestSQLSisipScoring - jumlah bind = jumlah kolom tiap INSERT; skor / persen lewat angkaMasuk; ID opsi dari sequence
// berskema; 16 bendera hasil berurutan kolomHasilScoring.
func TestSQLSisipScoring(t *testing.T) {
	tb := tabelScoringUji()
	for _, k := range []struct {
		q     string
		kolom int
	}{{sqlSisipScoringRisk(tb.sr), 5}, {sqlSisipHasilScoring(tb.hasil), 18}, {sqlSisipDataScoring(tb.data), 9},
		{sqlSisipFaktorScoring(tb.faktor), 9}, {sqlSisipOpsiScoring(tb.opsi, tb.seqOpsi), 7}} {
		m := regexp.MustCompile(`\(([^)]*)\) VALUES \((.*)\)$`).FindStringSubmatch(k.q)
		if m == nil {
			t.Fatalf("bentuk: %s", k.q)
		}
		kolom := strings.Split(m[1], ", ")
		nilai := regexp.MustCompile(`TO_NUMBER\(:\d+[^)]*\)|:\d+|UJI\.SEQ_O\.NEXTVAL`).FindAllString(m[2], -1)
		if len(kolom) != k.kolom || len(nilai) != len(kolom) {
			t.Errorf("%d kolom / %d nilai: %s", len(kolom), len(nilai), k.q)
		}
	}
	if !strings.Contains(sqlSisipOpsiScoring(tb.opsi, tb.seqOpsi), "VALUES (UJI.SEQ_O.NEXTVAL, :1,") {
		t.Error("ID opsi harus dari sequence berskema")
	}
	for _, q := range []string{sqlSisipScoringRisk(tb.sr), sqlSisipDataScoring(tb.data), sqlSisipFaktorScoring(tb.faktor)} {
		if !strings.Contains(q, "TO_NUMBER(") {
			t.Errorf("skor tanpa angkaMasuk: %s", q)
		}
	}
	arg := argHasilScoring(models.HasilScoring{Good: true, PresidentDirector: true})
	if len(arg) != len(kolomHasilScoring) || arg[0] != "true" || arg[15] != "true" || arg[1] != "false" ||
		kolomHasilScoring[0] != "GOOD" || kolomHasilScoring[15] != "PRESIDENT_DIRECTOR" {
		t.Errorf("hasil: %v", arg)
	}
}

// TestSQLHapusScoring - anak -> induk (FK), semua satu bind ID kasus; baca hanya opsi tercentang.
func TestSQLHapusScoring(t *testing.T) {
	tb := tabelScoringUji()
	q := sqlHapusScoring(tb)
	urut := []string{"DELETE FROM UJI.O ", "DELETE FROM UJI.F ", "DELETE FROM UJI.D ", "DELETE FROM UJI.H ", "DELETE FROM UJI.SR "}
	if len(q) != len(urut) {
		t.Fatalf("%d pernyataan", len(q))
	}
	for i, s := range q {
		if !strings.HasPrefix(s, urut[i]) || strings.Count(s, ":1") < 1 || strings.Contains(s, ":2") {
			t.Errorf("hapus[%d] %q", i, s)
		}
	}
	if !strings.Contains(sqlBacaOpsiScoring(tb), "o.CHECK_BOX = 'true'") || !strings.Contains(sqlKepalaScoring(tb.work), "LINI = :2") {
		t.Error("baca opsi / kepala")
	}
}

// TestFaktorMenurutHalaman - 24 pasangan (FACTOR_GROUP, FACTOR_NAME) unik; faktor langsung bergrup kosong (loader V-30).
func TestFaktorMenurutHalaman(t *testing.T) {
	if len(faktorMenurutHalaman) != 24 || faktorMenurutHalaman[[2]string{"", "Occupation"}] != "occupation" ||
		faktorMenurutHalaman[[2]string{"ObjectConditions", "Stockinthepolicycover"}] != "stock" ||
		faktorMenurutHalaman[[2]string{"LossRatio", "LossRatio"}] != "lossRatio" {
		t.Errorf("%v", faktorMenurutHalaman)
	}
	opsi := 0
	for _, f := range models.DaftarFaktorScoring {
		opsi += f.JumlahOpsi
	}
	if opsi != 101 { // 202 baris T_SCORING_OPTION untuk dua kolom (kasus_test loader)
		t.Errorf("%d opsi per kolom, mau 101", opsi)
	}
}

// TestSQLScoringMemakaiKolomMigrasi - setiap kolom yang ditulis / dibaca ada di CREATE TABLE migrasi 205.
func TestSQLScoringMemakaiKolomMigrasi(t *testing.T) {
	kolom := map[string]map[string]bool{}
	for _, m := range regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.(\w+) \((.*?)\n\)`).FindAllStringSubmatch(bacaMigrasi(t, "205_scoring_risk.sql"), -1) {
		kolom[m[1]] = map[string]bool{}
		for _, b := range strings.Split(m[2], "\n") {
			if f := strings.Fields(b); len(f) > 0 && f[0] != "CONSTRAINT" {
				kolom[m[1]][f[0]] = true
			}
		}
	}
	tb := tabelScoring{sr: TabelScoringRisk, hasil: TabelHasilScoring, data: TabelDataScoring, faktor: TabelFaktorScoring,
		opsi: TabelOpsiScoring, seqOpsi: "SEQ"}
	if len(kolom) != 5 {
		t.Fatalf("%d tabel terbaca", len(kolom))
	}
	for _, q := range []string{sqlSisipScoringRisk(tb.sr), sqlSisipHasilScoring(tb.hasil), sqlSisipDataScoring(tb.data),
		sqlSisipFaktorScoring(tb.faktor), sqlSisipOpsiScoring(tb.opsi, tb.seqOpsi)} {
		m := regexp.MustCompile(`INSERT INTO (\w+) \(([^)]*)\)`).FindStringSubmatch(q)
		for _, c := range strings.Split(m[2], ", ") {
			if !kolom[m[1]][c] {
				t.Errorf("%s.%s ditulis, tidak ada di 205", m[1], c)
			}
		}
	}
	for tabel, q := range map[string]string{TabelDataScoring: sqlBacaDataScoring(tb), TabelFaktorScoring: sqlBacaFaktorScoring(tb),
		TabelOpsiScoring: sqlBacaOpsiScoring(tb)} {
		for _, m := range regexp.MustCompile(`\b[dfo]\.([A-Z_]+)`).FindAllStringSubmatch(q, -1) {
			if !kolom[tabel][m[1]] && !kolom[TabelFaktorScoring][m[1]] && !kolom[TabelScoringRisk][m[1]] {
				t.Errorf("%s dibaca, tidak ada di 205", m[1])
			}
		}
	}
}
