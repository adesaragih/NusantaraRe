package repository

// Uji skema warisan dan pengurai nilainya - TANPA Oracle (tiket 01, tco4).
//
// Fixture SINTETIS: nol nama orang, nol nomor polis nyata, nol potongan data
// produksi; seluruh nilai berawalan UJI.

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/pkg/utils"
)

func TestUraiTanggalWarisanTCOMengenalBentukYangDikenal(t *testing.T) {
	kasus := map[string]string{
		"20260131":                "2026-01-31 00:00:00",
		"20260131T170000.000 GMT": "2026-01-31 00:00:00",
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

// ⛔ Oracle membulatkan NUMBER(38,8) DIAM-DIAM. Yang melampaui delapan angka
// di belakang koma harus ditolak di sini, sebelum satu digit pun hilang.
func TestUraiDesimalWarisanTCOMenolakPresisiMelampaui(t *testing.T) {
	_, catatan, ok := UraiDesimalWarisanTCO("0.123456789")
	if ok || !strings.Contains(catatan, "belakang koma") {
		t.Errorf("sembilan desimal harus ditolak dengan sebab; ok=%v catatan=%q", ok, catatan)
	}
	// Nol di ekor BUKAN presisi: 0.500000000 tetap 0.5.
	if _, _, ok := UraiDesimalWarisanTCO("0.500000000"); !ok {
		t.Error("nol di ekor tidak melampaui presisi")
	}
	_, catatan, ok = UraiDesimalWarisanTCO("1234567890123456789012345678901")
	if ok || !strings.Contains(catatan, "depan koma") {
		t.Errorf("tiga puluh satu digit bulat harus ditolak; ok=%v catatan=%q", ok, catatan)
	}
}

func TestEkorIdentitasTCO(t *testing.T) {
	if n, ok := EkorIdentitasTCO("1000042", LebarIdentitasTCO); !ok || n != 42 {
		t.Errorf("1000042 -> %d %v, mau 42 true", n, ok)
	}
	if n, ok := EkorIdentitasTCO("10000042", LebarIdentitasKlausulTCO); !ok || n != 42 {
		t.Errorf("10000042 lebar 7 -> %d %v, mau 42 true", n, ok)
	}
	for _, buruk := range []string{"10000042", "2000001", "UJI-1", "", "100000"} {
		if _, ok := EkorIdentitasTCO(buruk, LebarIdentitasTCO); ok {
			t.Errorf("%q seharusnya bukan identitas lebar 6", buruk)
		}
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
	nama := regexp.MustCompile(`(?i)^(T_)?(M?TREATY|PROPORTIONAL|SEQ_T_TREATY|SEQ_T_MTREATY|SEQ_T_PROPORTIONAL)`)
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
