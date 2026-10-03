package repository

// Uji murni lapisan repository - tanpa Oracle: konversi kolom (ID-14..18),
// bentuk SQL, dan kesepakatan katalog <-> DDL migrasi.

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestPecahAngkaEksak(t *testing.T) {
	for _, tt := range []struct {
		masuk, koef string
		skala       int64
	}{
		{"0", "0", 0},
		{"12.5", "125", 1},
		{"-0.43052837564", "-43052837564", 11},
		{"830.82191780804", "83082191780804", 11},
		{"1000", "1000", 0},
	} {
		d, err := utils.ParseDecimal(tt.masuk)
		if err != nil {
			t.Fatal(err)
		}
		k, s := pecahAngka(d)
		if k != tt.koef || s != tt.skala {
			t.Errorf("%s -> (%s, %d), harap (%s, %d)", tt.masuk, k, s, tt.koef, tt.skala)
		}
	}
}

func TestNilaiTulisKosongJadiNULL(t *testing.T) { // spec-penyimpanan AC 17, 18
	uang := models.Kolom{Properti: "X", Kolom: "X", Golongan: models.GolUang}
	if v, err := nilaiTulis(uang, ""); err != nil || v[0] != nil || v[1] != nil {
		t.Fatalf("uang kosong harus NULL, dapat %v %v", v, err)
	}
	tgl := models.Kolom{Properti: "T", Kolom: "T", Golongan: models.GolTanggal}
	if v, _ := nilaiTulis(tgl, ""); v[0] != nil {
		t.Fatal("tanggal kosong harus NULL")
	}
	kode := models.Kolom{Properti: "K", Kolom: "K", Golongan: models.GolKode, Panjang: 16}
	if v, _ := nilaiTulis(kode, "006"); v[0] != "006" {
		t.Fatal("kode tetap teks: nol di depan bermakna (ID-16)")
	}
	penanda := models.Kolom{Properti: "P", Kolom: "P", Golongan: models.GolPenanda, Panjang: 16}
	if v, _ := nilaiTulis(penanda, "0"); v[0] != "0" {
		t.Fatal("penanda \"0\" berbeda dari kosong (ID-17)")
	}
}

func TestCacahBilanganBulat(t *testing.T) { // ID-14
	c := models.Kolom{Properti: "InstallmentNo", Kolom: "INSTALLMENT_NO", Golongan: models.GolCacah}
	if v, err := nilaiTulis(c, "12"); err != nil || v[0] != "12" {
		t.Fatalf("12: %v %v", v, err)
	}
	if _, err := nilaiTulis(c, "1.5"); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Fatalf("1.5 bukan cacah: %v", err)
	}
	if eks, n := ekspresiTulis(c, 3); eks != "TO_NUMBER(:3)" || n != 1 {
		t.Fatalf("ekspresi %q %d", eks, n)
	}
}

func TestNilaiTulisMenolakMasukanRusak(t *testing.T) {
	uang := models.Kolom{Properti: "PremiOgp", Kolom: "P", Golongan: models.GolUang}
	if _, err := nilaiTulis(uang, "12,5"); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Fatalf("angka berkoma harus ditolak sebagai permintaan tidak sah: %v", err)
	}
	tgl := models.Kolom{Properti: "StartDate", Kolom: "T", Golongan: models.GolTanggal}
	if _, err := nilaiTulis(tgl, "31/12/2026"); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Fatalf("format tanggal kedua ditolak (AC 33): %v", err)
	}
	if v, err := nilaiTulis(tgl, "2026-12-31"); err != nil || v[0] != "2026-12-31 00:00:00" {
		t.Fatalf("tanggal satu format: %v %v", v, err)
	}
	teks := models.Kolom{Properti: "Remark", Kolom: "R", Golongan: models.GolTeks, Panjang: 3}
	if _, err := nilaiTulis(teks, "abcd"); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Fatalf("melebihi panjang kolom: %v", err)
	}
}

func TestNilaiBaca(t *testing.T) {
	uang := models.Kolom{Golongan: models.GolPersen}
	for masuk, harap := range map[string]string{".5": "0.5", "-.25": "-0.25", "12.5": "12.5"} {
		if got := nilaiBaca(uang, sql.NullString{String: masuk, Valid: true}); got != harap {
			t.Errorf("%q -> %q, harap %q", masuk, got, harap)
		}
	}
	tgl := models.Kolom{Golongan: models.GolTanggal}
	if got := nilaiBaca(tgl, sql.NullString{String: "2026-10-03 00:00:00", Valid: true}); got != "2026-10-03" {
		t.Fatalf("tanggal saja: %q", got)
	}
	if got := nilaiBaca(uang, sql.NullString{}); got != "" {
		t.Fatal("NULL dibaca kosong")
	}
}

func TestSQLDaftarKasusPenampungUnik(t *testing.T) {
	q := sqlDaftarKasus("S.W", "S.G", "S.Q", true, true)
	pen := regexp.MustCompile(`:(\d+)`).FindAllStringSubmatch(q, -1)
	lihat := map[string]bool{}
	for _, p := range pen {
		if lihat[p[1]] {
			t.Fatalf("penampung :%s berulang di SQL berpembatas baris:\n%s", p[1], q)
		}
		lihat[p[1]] = true
	}
	if len(lihat) != 7 || !strings.Contains(q, "FETCH FIRST 200 ROWS ONLY") {
		t.Fatalf("SQL daftar kasus:\n%s", q)
	}
}

// Setiap kolom katalog ada di CREATE TABLE tabelnya dengan tipe golongannya,
// dan setiap kolom DDL selain kolom kunci/generasi berasal dari katalog.
func TestKatalogSepakatDenganDDL(t *testing.T) {
	berkas, err := filepath.Glob(filepath.Join("..", "migrations", "3*.sql"))
	if err != nil || len(berkas) == 0 {
		t.Fatalf("migrasi tidak terbaca: %v", err)
	}
	ddl := map[string]map[string]string{}
	pola := regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.(\w+) \((.*?)\n\)`)
	for _, b := range berkas {
		if strings.HasSuffix(b, "_down.sql") {
			continue
		}
		isi, err := os.ReadFile(b)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range pola.FindAllStringSubmatch(string(isi), -1) {
			kol := map[string]string{}
			for _, baris := range strings.Split(m[2], "\n") {
				f := strings.Fields(strings.TrimSuffix(strings.TrimSpace(baris), ","))
				if len(f) >= 2 && f[0] != "CONSTRAINT" {
					kol[f[0]] = f[1]
				}
			}
			ddl[m[1]] = kol
		}
	}
	tipe := func(g models.Golongan) string {
		switch {
		case g.Desimal():
			return "NUMBER(38,8)"
		case g.Tanggal():
			return "DATE"
		case g == models.GolCacah:
			return "NUMBER(10)"
		}
		return "VARCHAR2"
	}
	bukanKatalog := map[string]bool{"ID": true, "POLIS_ID": true, "NOURUT": true, "INSTALMENT_ID": true, "XOL_ID": true,
		"NOPOLIS": true, "PRODKE": true, "NOENDORS": true, "OLD_POLIS_ID": true, "IDPEGA": true, "TGL_INPUT": true,
		"USERNAME": true, "TGL_TUTUP": true}
	for _, tb := range models.SemuaTabel {
		kol, ada := ddl[tb.Nama]
		if !ada {
			t.Errorf("%s tidak dibuat migrasi mana pun", tb.Nama)
			continue
		}
		dariKatalog := map[string]bool{}
		for _, k := range tb.Kolom {
			dariKatalog[k.Kolom] = true
			got, ada := kol[k.Kolom]
			if !ada {
				t.Errorf("%s.%s ada di katalog, tidak di DDL", tb.Nama, k.Kolom)
				continue
			}
			if !strings.HasPrefix(got, tipe(k.Golongan)) {
				t.Errorf("%s.%s bertipe %s, golongan %s menuntut %s", tb.Nama, k.Kolom, got, k.Golongan, tipe(k.Golongan))
			}
		}
		for k := range kol {
			if !dariKatalog[k] && !bukanKatalog[k] {
				t.Errorf("%s.%s ada di DDL, tidak di katalog", tb.Nama, k)
			}
		}
	}
	// AC 64: kolom isApprovedtoDeptHead tidak dibuat; AC 65: nomor surat bukan penanda arah.
	for tb, kol := range ddl {
		for k := range kol {
			if strings.Contains(k, "DEPT_HEAD") || strings.Contains(k, "LETTER") {
				t.Errorf("%s.%s: penanda arah tangga tidak disimpan (AC 64, 65)", tb, k)
			}
		}
	}
}

func TestKolomViewHilangAdalahGalat(t *testing.T) { // AC 89
	tipe := map[string]string{"ID": "VARCHAR2", "LIMITVALUE": "NUMBER", "COMMENCEMENT": "DATE"}
	if _, _, err := pilihKolom("UJI_VIEW", tipe, []string{"ID", "LIMITVALUE", "TREATYID"}); !errors.Is(err, ErrDataKontrakTidakAda) ||
		!strings.Contains(err.Error(), "TREATYID") {
		t.Fatalf("kolom hilang harus galat yang menyebutnya: %v", err)
	}
	eks, ada, err := pilihKolom("UJI_VIEW", tipe, []string{"ID", "LIMITVALUE", "COMMENCEMENT"})
	if err != nil || len(ada) != 3 || !strings.Contains(eks[1], "TM9") || !strings.Contains(eks[2], "TO_CHAR(COMMENCEMENT") || eks[0] != "ID" {
		t.Fatalf("ekspresi menurut tipe: %v %v %v", eks, ada, err)
	}
}
