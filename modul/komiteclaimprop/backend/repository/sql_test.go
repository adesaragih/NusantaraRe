package repository

// Uji perakit SQL (tanpa basis data): setiap kueri lolos `db.PeriksaSQL`, menyaring LINI PROP ketat dan awalan TKMT-,
// dan nol `COMMIT` / pemanggilan procedure.

import (
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/komiteclaimprop/backend/models"
)

func semuaSQL(t *testing.T) map[string]string {
	t.Helper()
	os, _, err := sqlSisipOS("S.OS_AKSEPTASI_KLAIM", models.BarisOSAkseptasi{CaseID: "CLMP-UJI001", Value: "25.5",
		GrossValue: "100", EstimationDate: time.Now()}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return map[string]string{
		"kepala":        sqlKepalaKasus("S.T_GENERAL_KOMITE", "S.T_WORK_CLAIM", false),
		"kepalaKunci":   sqlKepalaKasus("S.T_GENERAL_KOMITE", "S.T_WORK_CLAIM", true),
		"simpanKepala":  sqlSimpanKepala("S.T_GENERAL_KOMITE"),
		"tutup":         sqlTutupKasus("S.T_WORK_CLAIM"),
		"sentuh":        sqlSentuhKasus("S.T_WORK_CLAIM"),
		"tangga":        sqlTangga("S." + tabelTangga),
		"tulisAnggota":  sqlTulisAnggota("S."+tabelTangga, true),
		"tulisAnggota2": sqlTulisAnggota("S."+tabelTangga, false),
		"daftarKerja":   sqlDaftarKerja("S.G", "S.W", "S.L", "S.C", "S.A", 2),
		"os":            os,
		"adaJSON":       sqlAdaJSONKlaim("S.JSON_KLAIM"),
		"sisipJSON":     sqlSisipJSONKlaim("S.JSON_KLAIM"),
		"log":           sqlLogLayanan("S.MONITORING_KLAIM_LOG"),
		"riwayat":       sqlRiwayatAkseptasi("S.HISTORYAKSEPTASIPEGA"),
		"emailPelaku":   sqlEmailPelaku("S.M_LOGIN_GO"),
		"dokumenKlaim":  sqlSisipDokumenKlaim("S." + TabelDokumenKlaim),
	}
}

func TestSQLKomiteLolosPemeriksaan(t *testing.T) {
	for nama, q := range semuaSQL(t) {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
		u := strings.ToUpper(q)
		if strings.Contains(u, "COMMIT") || strings.Contains(u, "BEGIN") || strings.Contains(q, "%[") {
			t.Errorf("%s: COMMIT / blok PL/SQL / %%[n]s dilarang: %s", nama, q)
		}
	}
}

func TestSQLKomiteMenyaringLiniPropKetat(t *testing.T) {
	s := semuaSQL(t)
	for _, nama := range []string{"kepala", "daftarKerja"} {
		q := s[nama]
		if !strings.Contains(q, "w.LINI = :") || !strings.Contains(q, "w.ID LIKE :") {
			t.Errorf("%s: saringan LINI / awalan TKMT- hilang", nama)
		}
		if strings.Contains(strings.ToUpper(q), "LINI IS NULL") {
			t.Errorf("%s: LINI harus KETAT (tanpa OR LINI IS NULL)", nama)
		}
	}
	if !strings.Contains(s["daftarKerja"], "MIN(l2.KOMITE_URUT)") {
		t.Error("daftar kerja: baris tangga PERTAMA yang menunggu (KomiteRouter S6.1)")
	}
	if !strings.Contains(s["tulisAnggota"], "KOMITE_APPROVAL = :7") || !strings.Contains(s["simpanKepala"],
		"KOMITE_COUNT = :8") {
		t.Error("tulisan keputusan bersyarat keadaan yang dibaca (dua klik tidak sama-sama menang)")
	}
	if !strings.Contains(s["kepalaKunci"], "FOR UPDATE") {
		t.Error("Submit mengunci kepala kasus")
	}
}

func TestPecahAngkaEksak(t *testing.T) {
	_, args, err := sqlSisipOS("S.OS", models.BarisOSAkseptasi{CaseID: "X", Value: "1234.5600"}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	var koef string
	var skala int64
	for i, a := range args {
		if s, ok := a.(string); ok && s == "12345600" {
			koef, skala = s, args[i+1].(int64)
		}
	}
	if koef != "12345600" || skala != 4 {
		t.Fatalf("angka eksak tanpa float: %v", args)
	}
}
