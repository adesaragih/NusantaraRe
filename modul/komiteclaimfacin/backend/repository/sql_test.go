package repository

// Uji teks SQL tanpa Oracle: nol COMMIT / PL-SQL / `%[n]s` (ADR-U-0029, ADR-U-0033), saringan `LINI` ketat + awalan di
// setiap kueri kasus komite (tabel bersama Komite Claim Life / Prop / Non Prop), bind urut kemunculan, daftar kerja
// tingkat berjalan (`KOMITE_URUT = KOMITE_COUNT`, perbaikan prompt §5 butir 1).

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/modul/komiteclaimfacin/backend/models"
)

func semuaSQL() map[string]string {
	return map[string]string{
		"kepala":       sqlKepalaKasus("UJI.G", "UJI.W", true),
		"simpanKepala": sqlSimpanKepala("UJI.G"),
		"tutup":        sqlTutupKasus("UJI.W"),
		"sentuh":       sqlSentuhKasus("UJI.W"),
		"tangga":       sqlTangga("UJI.L"),
		"tambah":       sqlTambahAnggota("UJI.L"),
		"tulis":        sqlTulisAnggota("UJI.L", true),
		"tulisSisa":    sqlTulisAnggota("UJI.L", false),
		"kerja":        sqlDaftarKerja("UJI.G", "UJI.W", "UJI.L", "UJI.C", "UJI.A", 2),
		"os":           sqlSisipOS("UJI.OS", true),
		"osTanpaKonv":  sqlSisipOS("UJI.OS", false),
		"jsonAda":      sqlAdaJSONKlaim("UJI.J"),
		"json":         sqlSisipJSONKlaim("UJI.J"),
		"log":          sqlLogLayanan("UJI.M"),
		"riwayat":      sqlRiwayatAkseptasi("UJI.H"),
		"sub":          sqlUbahSubProgres("UJI.S"),
		"tolak":        sqlSisipKlaimDitolak("UJI.R"),
		"roster":       sqlRosterKomite("UJI.E"),
		"anggotaWB":    sqlEmailAnggotaWorkbasket("UJI.LW", "UJI.WB", "UJI.LG"),
	}
}

func TestSQLTanpaCommitDanIndeksFormat(t *testing.T) {
	larang := regexp.MustCompile(`(?i)\b(COMMIT|BEGIN|EXECUTE|PROC_|PEGA_JSON)\b|%\[`)
	for n, q := range semuaSQL() {
		if larang.MatchString(q) {
			t.Errorf("%s memuat kata terlarang: %s", n, q)
		}
	}
}

// TestBindUrutKemunculan - penampung `:n` naik 1, 2, 3 ... sesuai teks (driver mengikat menurut posisi).
func TestBindUrutKemunculan(t *testing.T) {
	pola := regexp.MustCompile(`:(\d+)`)
	for n, q := range semuaSQL() {
		for i, m := range pola.FindAllStringSubmatch(q, -1) {
			if m[1] != strconv.Itoa(i+1) {
				t.Errorf("%s: penampung ke-%d = :%s", n, i+1, m[1])
				break
			}
		}
	}
}

func TestKueriKasusMenyaringLiniDanAwalan(t *testing.T) {
	for _, q := range []string{sqlKepalaKasus("G", "W", false), sqlDaftarKerja("G", "W", "L", "C", "A", 0)} {
		if !strings.Contains(q, "w.LINI = :") || !strings.Contains(q, "w.ID LIKE :") {
			t.Fatalf("kueri kasus komite tanpa saringan LINI / awalan: %s", q)
		}
	}
	for _, q := range []string{sqlTutupKasus("W"), sqlSentuhKasus("W")} {
		if !strings.Contains(q, "LINI = :") || !strings.Contains(q, "STATUS_WORK IS NULL") {
			t.Fatalf("tulisan work object tanpa saringan: %s", q)
		}
	}
	if q := sqlDaftarKerja("G", "W", "L", "C", "A", 1); !strings.Contains(q, "l.KOMITE_URUT = g.KOMITE_COUNT") ||
		!strings.Contains(q, "LEFT JOIN A a ON a.KOMITE_ID = g.ID") {
		t.Fatalf("daftar kerja: %s", q)
	}
	if models.LiniFacIn != "FACIN" || models.AwalanKomite != "KMT-" {
		t.Fatal("nilai saringan berubah")
	}
}

func TestOSSesuaiProcedure(t *testing.T) {
	q := sqlSisipOS("UJI.OS", true)
	for _, k := range []string{"CASEID", "NOCLAIM", "DATA_JSON", "TANGGAL", "NOPOLIS", "STS_REJECT", "STS_KONVERSI", "STS_DLA"} {
		if !strings.Contains(q, k) {
			t.Fatalf("OS tanpa %s: %s", k, q)
		}
	}
	// jawaban work owner 10-10-2026 (OQ-KCFI-08, preseden Claim Prop): STS_KONVERSI kosong -> kolom tidak disebut
	if q := sqlSisipOS("UJI.OS", false); strings.Contains(q, "STS_KONVERSI") || !strings.Contains(q, "STS_DLA") ||
		strings.Count(q, ":") != 7 {
		t.Fatalf("OS tanpa status konversi: %s", q)
	}
	if strings.Contains(q, "TGL_PROD") || strings.Contains(q, "MASTERID") {
		t.Fatalf("OS menulis kolom yang tidak ditulis procedure (TGL_PROD = trigger): %s", q)
	}
}

func TestKepalaMembacaTeksTutup(t *testing.T) {
	// migrasi 643 (jawaban work owner 10-10-2026 OQ-KCFI-03): teks pop-up TT3 / TT4 di kepala kasus komite
	q := sqlKepalaKasus("UJI.G", "UJI.W", false)
	for _, k := range []string{"g.KOMITE_CIRCUM_CAUSE_OF_LOSS", "g.KOMITE_EXTENT_OF_LOSS", "g.KOMITE_LEGAL_LIABILITY"} {
		if !strings.Contains(q, k) {
			t.Fatalf("kepala kasus tanpa %s: %s", k, q)
		}
	}
}
