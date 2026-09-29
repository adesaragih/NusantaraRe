package repository

// Uji skema warisan dan pengurai nilainya - TANPA Oracle (tiket 01, tco4).
//
// Fixture SINTETIS: nol nama orang, nol nomor polis nyata, nol potongan data
// produksi; seluruh nilai berawalan UJI.

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/utils"
)

func TestUraiTanggalWarisanTCOMengenalBentukYangDikenal(t *testing.T) {
	kasus := map[string]string{
		"20260131": "2026-01-31 00:00:00",
		// RALAT tco4: stempel GMT dibaca di zona Jakarta - 17:00 GMT = 00:00 WIB esok.
		"20260131T170000.000 GMT": "2026-02-01 00:00:00",
		"20260131T000000.000 GMT": "2026-01-31 00:00:00",
		"31/01/2026":              "2026-01-31 00:00:00",
		"2026-01-31":              "2026-01-31 00:00:00",
		"2026-01-31 10:20:30":     "2026-01-31 10:20:30",
		"31-JAN-26":               "2026-01-31 00:00:00",
		"31-Jan-2026":             "2026-01-31 00:00:00",
		" 20260131 ":              "2026-01-31 00:00:00",
	}
	for masuk, mau := range kasus {
		hasil, ok := UraiTanggalWarisanTCO(masuk)
		if !ok {
			t.Errorf("%q tidak terurai", masuk)
			continue
		}
		if dapat := utils.FormatTanggalWaktu(hasil); dapat != mau {
			t.Errorf("%q -> %q, mau %q", masuk, dapat, mau)
		}
	}
}

func TestUraiTanggalWarisanTCOKosongDanTakDikenal(t *testing.T) {
	if hasil, ok := UraiTanggalWarisanTCO("   "); !ok || !hasil.IsZero() {
		t.Errorf("kosong harus KOSONG dan ok, dapat %v %v", hasil, ok)
	}
	for _, buruk := range []string{"2026", "31-13-2026", "20261340", "besok", "1/2/2026"} {
		if _, ok := UraiTanggalWarisanTCO(buruk); ok {
			t.Errorf("%q seharusnya TIDAK terurai - ia harus dilaporkan, bukan ditebak", buruk)
		}
	}
}

func TestUraiDesimalWarisanTCO(t *testing.T) {
	baik := map[string]string{
		"12.5": "12.5", "12,5": "12.5", ".5": "0.5", "1000000000.12345678": "1000000000.12345678",
		" 7 ": "7", "0": "0", "-3.25": "-3.25",
	}
	for masuk, mau := range baik {
		d, catatan, ok := UraiDesimalWarisanTCO(masuk)
		if !ok {
			t.Errorf("%q ditolak: %s", masuk, catatan)
			continue
		}
		if dapat := utils.FormatDecimal(d); dapat != mau {
			t.Errorf("%q -> %q, mau %q", masuk, dapat, mau)
		}
	}
	if d, _, ok := UraiDesimalWarisanTCO(""); !ok || d != nil {
		t.Error("kosong harus KOSONG (nil) dan ok")
	}
	buruk := []string{"1.000,5", "abc", "NaN", "Infinity", "1e400"}
	for _, b := range buruk {
		if _, _, ok := UraiDesimalWarisanTCO(b); ok {
			t.Errorf("%q seharusnya ditolak", b)
		}
	}
}

// Titik pemisah RIBUAN - data DEV `PROPORTIONALARRG` (satu baris Limit MB yang
// dibaca): `RP` dan `USD` = "1.000.000", sedangkan `PCTME`/`MORERP` baris yang
// sama 1000000. Dua titik atau lebih dengan kelompok tepat tiga angka
// TIDAK mungkin desimal (desimal hanya punya satu titik) - bukan tebakan.
// Tanpa ini satu baris mematikan seluruh grid jenis itu (500).
func TestUraiDesimalWarisanTCORibuanTitik(t *testing.T) {
	for masuk, mau := range map[string]string{
		"1.000.000": "1000000", "12.345.678": "12345678", "-1.000.000": "-1000000", " 250.000.000 ": "250000000",
		"1.000": "1.000", // SATU titik tetap desimal
	} {
		d, catatan, ok := UraiDesimalWarisanTCO(masuk)
		if !ok || utils.FormatDecimal(d) != mau {
			t.Errorf("%q -> %v (%s), mau %q", masuk, d, catatan, mau)
		}
	}
	// Kelompok yang bukan tiga angka, atau kelompok pertama lebih dari tiga,
	// bukan bentuk ribuan - tetap ditolak, tidak ditebak.
	for _, b := range []string{"1.00.000", "1.000.00", "1000.000.000", "1.000.000.", ".000.000", "1.000.000,5",
		"000.000.000", "01.000.000"} {
		if _, _, ok := UraiDesimalWarisanTCO(b); ok {
			t.Errorf("%q seharusnya ditolak", b)
		}
	}
}

// ⛔ Oracle membulatkan NUMBER(38,8) DIAM-DIAM. Yang melampaui delapan angka
// di belakang koma harus ditolak di sini, sebelum satu digit pun hilang.
// tco4 (temuan /code-review): teks warisan VARCHAR2 tanpa batas skala - nilai
// berdesimal panjang yang diketik di Pega tetap terbaca utuh.
func TestUraiDesimalWarisanTCOTanpaBatasSkala(t *testing.T) {
	for _, teks := range []string{"33.3333333333", "33,3333333333", "1234567890123456789012345678901.5"} {
		d, catatan, ok := UraiDesimalWarisanTCO(teks)
		if !ok || d == nil {
			t.Errorf("%q ditolak: %s", teks, catatan)
		}
	}
	if d, _, _ := UraiDesimalWarisanTCO("33.3333333333"); d.Text('f') != "33.3333333333" {
		t.Errorf("digit hilang: %s", d.Text('f'))
	}
}

// Penjaga instrumen: pembaca mengenal SELURUH kolom yang procedure tulis.
func TestKolomWarisanTCOSesuaiProcedure(t *testing.T) {
	cacah := map[string]int{
		warisanTahunTCO: 10, warisanKontrakTCO: 8, warisanReinsurerTCO: 19,
		warisanSecurityTCO: 7, warisanBusinessTCO: 12, warisanKlausulTCO: 35,
	}
	for tabel, n := range cacah {
		if len(KolomWarisanTCO(tabel)) != n {
			t.Errorf("%s: %d kolom, mau %d (parameter procedure)", tabel, len(KolomWarisanTCO(tabel)), n)
		}
	}
	// Sembilan kolom khusus induk adalah sembilan kolom TERAKHIR klausul.
	kol := KolomWarisanTCO(warisanKlausulTCO)
	induk := []string{"ID_OCCUPATION", "OCCUPATION", "ID_CLAUSE", "CLAUSE", "TREATYLIMIT",
		"COINS_MIN", "COINS_MAX", "MORERP", "MOREUSD"}
	for i, k := range induk {
		if kol[26+i] != k {
			t.Errorf("kolom ke-%d = %s, mau %s", 27+i, kol[26+i], k)
		}
	}
}

// tco4 (keputusan work owner 29-09-2026): NOL tabel baru untuk modul ini -
// rentang migrasinya 300-319 kosong, dan tidak satu pun berkas migrasi mana
// pun membuat tabel atau sequence bernama Treaty Contract Out.
func TestTCONolTabelBaru(t *testing.T) {
	pola := regexp.MustCompile(`(?i)CREATE\s+(TABLE|SEQUENCE)\s+\{skema\}\.(\w+)`)
	nama := regexp.MustCompile(`(?i)^(SEQ_)?(T_)?(M?TREATY|PROPORTIONAL)`)
	berkas := 0
	for n, teks := range seluruhSQL(t, false) {
		berkas++
		if n >= "300_" && n < "320_" {
			t.Errorf("%s: berkas migrasi di rentang Treaty Contract Out 300-319 - tco4 menolak tabel baru", n)
		}
		for _, m := range pola.FindAllStringSubmatch(teks, -1) {
			if nama.MatchString(m[2]) {
				t.Errorf("%s membuat %s %s - tco4: modul ini memakai tabel warisan", n, m[1], m[2])
			}
		}
	}
	if berkas == 0 {
		t.Fatal("nol berkas migrasi terbaca; pembacanya yang rusak")
	}
}

// Tepi tulis: stempel Pega, tanggal YYYYMMDD, dan desimal - dua arah.
func TestTepiTulisWarisanTCODuaArah(t *testing.T) {
	w := time.Date(2026, 9, 29, 5, 6, 7, 891_000_000, time.UTC)
	if s := StempelPegaTCO(w); s != "20260929T050607.891 GMT" {
		t.Errorf("stempel %q", s)
	}
	if b, ok := UraiWaktuWarisanTCO(StempelPegaTCO(w)); !ok || !b.Equal(w) {
		t.Errorf("stempel pulang %v %v", b, ok)
	}
	if StempelPegaTCO(time.Time{}) != "" {
		t.Error("waktu nol harus teks kosong (NULL)")
	}
	// OQ-TCO-01 (lanjutan 4, data DEV 182/182): tanggal tahun = delapan angka.
	tgl := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := TanggalYYYYMMDDTCO(tgl)
	if s != "20260101" {
		t.Errorf("tanggal tahun %q, mau 20260101", s)
	}
	if b, ok := uraiYYYYMMDDTCO(s); !ok || !b.Equal(tgl) {
		t.Errorf("tanggal pulang %v %v, mau %v", b, ok, tgl)
	}
	if TanggalYYYYMMDDTCO(time.Time{}) != "" {
		t.Error("tanggal nol harus teks kosong (NULL)")
	}
	for masuk, mau := range map[string]string{"12.50": "12.5", "1000000000.12345678": "1000000000.12345678",
		"100": "100", "0.0": "0", "-3.10": "-3.1"} {
		d, _ := utils.ParseDecimal(masuk)
		if dapat := TulisDesimalWarisanTCO(d); dapat != mau {
			t.Errorf("desimal %s -> %v, mau %s", masuk, dapat, mau)
		}
		pulang, _, ok := UraiDesimalWarisanTCO(mau)
		if !ok || pulang.Cmp(d) != 0 {
			t.Errorf("desimal %s tidak pulang utuh", masuk)
		}
	}
	if TulisDesimalWarisanTCO(nil) != nil {
		t.Error("nil harus NULL")
	}
}

// tco4: nol nama tabel baru modul (T_ + TREATY…/MTREATY…/PROPORTIONAL…) di KODE
// Go dan frontend. Komentar - catatan sejarah dan ralat - dibuang lebih dulu.
// Polanya dirakit dari potongan supaya berkas ini sendiri tidak memuatnya.
func TestTCONolNamaTabelBaruDiKode(t *testing.T) {
	pola := regexp.MustCompile("T" + "_" + `(TREATY|MTREATY|PROPORTIONAL)\w*`) // tanpa \b: SEQ_ + nama ikut
	if !pola.MatchString("SELECT ID FROM S."+"T"+"_TREATYYEAR") || !pola.MatchString("S.SEQ_"+"T"+"_TREATYYEAR") ||
		pola.MatchString("SELECT ID FROM S.TREATYYEAR") {
		t.Fatal("pola penjaga tidak menggigit atau menuduh nama warisan")
	}
	ekor := regexp.MustCompile(`(^|\s)//.*$`)
	blok := regexp.MustCompile(`(?s)/\*.*?\*/`)
	berkas := 0
	for _, akar := range []string{akarModul + "/inti", akarModul + "/modul", akarModul + "/cmd", akarModul + "/frontend/src"} {
		err := filepath.Walk(filepath.FromSlash(akar), func(jalur string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if info.Name() == "node_modules" || info.Name() == "dist" {
					return filepath.SkipDir
				}
				return nil
			}
			if ext := filepath.Ext(jalur); ext != ".go" && ext != ".ts" && ext != ".tsx" {
				return nil
			}
			isi, err := os.ReadFile(jalur)
			if err != nil {
				return err
			}
			berkas++
			kode := blok.ReplaceAllString(string(isi), "")
			for i, baris := range strings.Split(kode, "\n") {
				baris = ekor.ReplaceAllString(baris, "")
				if filepath.Ext(jalur) != ".go" && strings.HasPrefix(strings.TrimSpace(baris), "*") {
					continue // badan JSDoc
				}
				if m := pola.FindString(baris); m != "" {
					t.Errorf("%s:%d menyebut %s - tco4: modul ini memakai tabel warisan", filepath.ToSlash(jalur), i+1, m)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if berkas < 100 {
		t.Fatalf("hanya %d berkas terbaca; pembacanya yang rusak", berkas)
	}
}

// OQ-TCO-01: kolom tanggal TREATYYEAR HANYA menerima YYYYMMDD - stempel Pega
// dan bentuk lain ditolak dengan galat yang menyebut kolom dan bentuknya.
func TestTanggalTahunHanyaYYYYMMDD(t *testing.T) {
	for _, baik := range []string{"20260101", " 20261231 ", ""} {
		if _, err := tanggalTahunWarisanTeks(sqlNull(baik), "STARTDATE"); err != nil {
			t.Errorf("%q ditolak: %v", baik, err)
		}
	}
	for _, buruk := range []string{"20251231T170000.000 GMT", "01/01/2026", "2026-01-01", "2026011", "20261340"} {
		_, err := tanggalTahunWarisanTeks(sqlNull(buruk), "STARTDATE")
		if err == nil || !strings.Contains(err.Error(), "bukan YYYYMMDD") || !strings.Contains(err.Error(), "STARTDATE") {
			t.Errorf("%q: %v", buruk, err)
		}
	}
}

func sqlNull(s string) sql.NullString { return sql.NullString{String: s, Valid: s != ""} }
