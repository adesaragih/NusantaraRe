package repository

// SQL Aggregate TANPA Oracle: kolom, urutan bind, saringan Pega, dan nol COMMIT.

import (
	"regexp"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/aggregate/backend/models"
)

func satuBaris(q string) string { return strings.Join(strings.Fields(q), " ") }

// Sisip: ID, TANGGAL_INPUT, USER_INPUT, lalu 40 kolom grid; angka = dua bind (koefisien, skala); nol COMMENCEMENT.
func TestSisipKolomDanBind(t *testing.T) {
	q := satuBaris(sqlSisip("S.AGGREGATE"))
	if strings.Contains(q, "COMMENCEMENT") {
		t.Errorf("sisip mengisi COMMENCEMENT: %s", q)
	}
	if !strings.HasPrefix(q, "INSERT INTO S.AGGREGATE (ID, TANGGAL_INPUT, USER_INPUT, ASSESMENT_ZONE, M_TREATY_ID,") {
		t.Errorf("sisip: %s", q[:120])
	}
	bind := regexp.MustCompile(`:(\d+)`).FindAllStringSubmatch(q, -1)
	if len(bind) != 72 || bind[71][1] != "72" {
		t.Errorf("bind %d, mau 72 berurutan", len(bind))
	}
	b := models.Baris{"ASSESMENT_ZONE": "1.1 UJI", "AS_AT": "30-09-2024", "TO_USD": "0.00006060606060606061",
		"TOTAL_IN_AMOUNT": "357477219", "RNM_VALUE": "35747721.9", "REMARK": ""}
	n, err := NilaiSisip("AGG-1", "2026-10-04 09:00:00", "UJI-ADMIN", b)
	if err != nil {
		t.Fatal(err)
	}
	if len(n) != 72 || n[0] != "AGG-1" || n[2] != "UJI-ADMIN" || n[3] != "1.1 UJI" || n[7] != "30-09-2024" {
		t.Errorf("nilai awal %v", n[:8])
	}
	// TO_USD (kolom ke-11 grid): koefisien dan skala 20.
	if n[13] != "6060606060606061" || n[14] != int64(20) {
		t.Errorf("TO_USD %v %v", n[13], n[14])
	}
	if n[len(n)-1] != nil {
		t.Errorf("REMARK kosong = NULL, dapat %v", n[len(n)-1])
	}
	if _, err := NilaiSisip("AGG-1", "x", "x", models.Baris{"TO_USD": "abc"}); err == nil {
		t.Error("angka tak terurai diterima")
	}
}

func TestPecahAngka(t *testing.T) {
	for masuk, mau := range map[string][2]any{"1500000000.10": {"150000000010", int64(2)}, "0.5": {"5", int64(1)},
		"0": {"0", int64(0)}, "-12.5": {"-125", int64(1)}} {
		d, _, _ := apd.NewFromString(masuk)
		koef, skala := PecahAngka(d)
		if koef != mau[0] || skala != mau[1] {
			t.Errorf("%s -> %v %v, mau %v", masuk, koef, skala, mau)
		}
	}
	if koef, _ := PecahAngka(nil); koef != nil {
		t.Error("nil bukan NULL")
	}
}

// Master ID: CEDING LIKE, (Proportional + PROPERTY) atau NonProportional, tanpa ORDER BY (Obj-Browse Pega).
func TestCariTreatySaringanPega(t *testing.T) {
	q := satuBaris(sqlCariTreaty("S.V"))
	if !strings.Contains(q, "WHERE CEDING LIKE :1 ESCAPE '\\' AND ((PROPORTIONTYPE = :2 AND TREATYGROUP = :3) OR PROPORTIONTYPE = :4)") ||
		strings.Contains(q, "ORDER BY") {
		t.Errorf("cari treaty: %s", q)
	}
	if PolaCari(" a%b_ ") != `%A\%B\_%` {
		t.Errorf("pola %q", PolaCari(" a%b_ "))
	}
}

// Treaty Year dan kurs: baris pertama tanpa urutan (keputusan work owner: ikuti XML).
func TestTahunDanKursBarisPertama(t *testing.T) {
	if q := satuBaris(sqlTahunTreaty("S.T")); q != "SELECT TREATYYEAR FROM S.T WHERE STARTDATE <= :1 AND ENDDATE >= :2 FETCH FIRST 1 ROWS ONLY" {
		t.Errorf("tahun: %s", q)
	}
	if q := satuBaris(sqlKurs("S.K")); !strings.HasSuffix(q, "WHERE TREATYYEAR = :1 AND CURRENCY = :2 FETCH FIRST 1 ROWS ONLY") {
		t.Errorf("kurs: %s", q)
	}
}

// Kunci daftar: DECODE menyamakan NULL; urutan bind kunci.
func TestKunciDanHapus(t *testing.T) {
	q := satuBaris(sqlHapus("S.AGGREGATE"))
	for _, k := range []string{"DECODE(TRUNC(TANGGAL_INPUT), TO_DATE(:1, 'DD-MM-YYYY'), 1, 0) = 1", "DECODE(CEDING_NAME, :3, 1, 0) = 1",
		"DECODE(UW_YEAR, :6, 1, 0) = 1"} {
		if !strings.Contains(q, k) {
			t.Errorf("hapus tanpa %q: %s", k, q)
		}
	}
	n := NilaiKunci(models.Kunci{TanggalInput: "01-03-2025", CedingName: "UJI"})
	if n[0] != "01-03-2025" || n[1] != nil || n[2] != "UJI" || len(n) != 6 {
		t.Errorf("nilai kunci %v", n)
	}
}

func TestSeluruhSQLLolosPeriksa(t *testing.T) {
	for _, q := range []string{sqlSisip("S.T"), sqlDaftar("S.T"), sqlHitungDaftar("S.T"), sqlRincian("S.T"), sqlHapus("S.T"),
		sqlRingkasan("S.T"), sqlNomor("S.Q"), sqlAdaID("S.T"), sqlSekarang, sqlCariTreaty("S.V"),
		sqlAmbilTreaty("S.V"), sqlDaftarZona("S.Z"), sqlTahunTreaty("S.T"), sqlKurs("S.K")} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Error(err)
		}
	}
}
