package repository

import (
	"database/sql"
	"regexp"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatycontractout/models"
)

// AC 47: tidak ada identitas mata uang di teks SQL; saringan di-bind.
// Lanjutan 6: tanggal diurai ORACLE dengan ekspresi RDB `GetMasterKursList`
// b85-b86 dan dibandingkan di SQL; satu-satunya literal = kedua topeng.
func TestSQLKursTCO(t *testing.T) {
	q := sqlBerlakuKursTCO("S.K")
	for _, wajib := range []string{
		`TRUNC(TO_TIMESTAMP_TZ(STARTDATE DEFAULT NULL ON CONVERSION ERROR, 'YYYYMMDD"T"HH24MISS.FF3 TZR'))`,
		`TRUNC(TO_TIMESTAMP_TZ(ENDDATE DEFAULT NULL ON CONVERSION ERROR, 'YYYYMMDD"T"HH24MISS.FF3 TZR'))`,
		`TO_DATE(:1, 'YYYYMMDD') BETWEEN MULAI AND AKHIR`,
		`FROM S.K WHERE QUARTER = :2 AND IDCURRENCY = :3`,
	} {
		if !strings.Contains(q, wajib) {
			t.Errorf("tanpa %s:\n%s", wajib, q)
		}
	}
	// Bind berurutan menurut kemunculan - go-ora mengikat menurut posisi.
	if i1, i2, i3 := strings.Index(q, ":1"), strings.Index(q, ":2"), strings.Index(q, ":3"); i1 < 0 || i1 > i2 || i2 > i3 {
		t.Errorf("urutan bind: %d %d %d", i1, i2, i3)
	}
	sisa := strings.NewReplacer("'"+models.FormatTanggalKursTCO+"'", "", "'YYYYMMDD'", "").Replace(q)
	if regexp.MustCompile(`'[^']*'`).MatchString(sisa) {
		t.Errorf("literal di SQL kurs: %s", q)
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Error(err)
	}
	// Tabel: nama DBA, dibaca 12 RDB korpus di enam modul lain;
	// `treatyexchange` hanya di GetMasterKursList (tiket 11, bab lanjutan 6).
	if MasterKursTahunanTCO != "TREATYEXCHANGEYEARLY" {
		t.Errorf("tabel kurs: %s", MasterKursTahunanTCO)
	}
	// Master kurs dan master mata uang terdaftar sebagai dibaca-saja - penjaga
	// TestTCOWarisanHanyaDibaca menolak setiap penulis yang menyebutnya.
	for _, m := range []string{MasterKursTahunanTCO, MasterMataUangTCO} {
		ada := false
		for _, x := range masterDibacaSajaTCO {
			ada = ada || x == m
		}
		if !ada {
			t.Errorf("%s tidak terdaftar dibaca-saja", m)
		}
	}
}

func nsKurs(v ...string) [6]sql.NullString {
	var n [6]sql.NullString
	for i, s := range v {
		n[i] = sql.NullString{String: s, Valid: s != ""}
	}
	return n
}

func tglOracle(t time.Time) sql.NullTime { return sql.NullTime{Time: t, Valid: true} }

// Tanggal dari Oracle (DATE); NULL = Oracle menolak teksnya (atau kosong).
func TestTambahBarisKursTCO(t *testing.T) {
	// go-ora memberi DATE dengan jam dinding Oracle di lokasi zona server.
	wib := time.FixedZone("WIB", 7*3600)
	mulai, akhir := tglOracle(time.Date(2026, 1, 1, 0, 0, 0, 0, wib)), tglOracle(time.Date(2026, 12, 31, 0, 0, 0, 0, wib))
	var h models.HasilMasterKursTCO
	tambahBarisKursTCO(&h, nsKurs("15500,25", "20260101T00000.000 GMT", "20261231T000000.000 GMT", "UJI-USD", "USD", "0"), mulai, akhir, 1)
	// Baris yang tidak berlaku pada tanggal itu tidak dibawa - TOIDR-nya
	// (kosong) tidak pernah diurai.
	tambahBarisKursTCO(&h, nsKurs("", "20120101T000000.000 GMT", "20121231T000000.000 GMT"),
		tglOracle(time.Date(2012, 1, 1, 0, 0, 0, 0, wib)), tglOracle(time.Date(2012, 12, 31, 0, 0, 0, 0, wib)), 0)
	tambahBarisKursTCO(&h, nsKurs("15500", "2019A801T000000.000 GMT", "20191231T000000.000 GMT"), sql.NullTime{}, akhir, 0)
	tambahBarisKursTCO(&h, nsKurs("15500", "20190101T000000.000 GMT", ""), mulai, sql.NullTime{}, 0)

	if len(h.Berlaku) != 1 {
		t.Fatalf("berlaku: %+v", h.Berlaku)
	}
	k := h.Berlaku[0]
	// ⛔ Hari dari medan DATE apa adanya: `.UTC()` menggeser 2026-01-01 00:00
	// WIB ke 2025-12-31.
	if k.ToIDR != nil || k.TeksToIDR != "15500,25" || k.Mulai.Format("2006-01-02") != "2026-01-01" ||
		k.Akhir.Format("2006-01-02") != "2026-12-31" || k.IDCurrency != "UJI-USD" || k.Currency != "USD" || k.Quarter != "0" {
		t.Errorf("baris berlaku: %+v", k)
	}
	mau := []models.BarisKursDitolakTCO{{Kolom: "STARTDATE", Teks: "2019A801T000000.000 GMT"}, {Kolom: "ENDDATE", Teks: ""}}
	if len(h.Ditolak) != 2 || h.Ditolak[0] != mau[0] || h.Ditolak[1] != mau[1] {
		t.Errorf("ditolak: %+v", h.Ditolak)
	}
}
