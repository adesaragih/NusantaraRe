package backend

// K18 (PROMPT putaran 3 bab 2; koreksi WO 04-10-2026) - TANPA Oracle.
//
// `POOLDATA.T_GENERAL_POLIS` MILIK `nbfacin` (migrasi `182_t_general_polis`, TUJUH kolom: ID, IDPEGA,
// COB_GROUP, START_DATE_TIME, OFFERING_DATE, END_DATE_TIME, FOLLOWING). Perintah work owner 05-10-2026:
// NB Treaty In TIDAK memakai T_GENERAL_POLIS - tabel induknya `T_GENERAL_POLIS_TREATY` (migrasi
// `320_t_general_polis_treaty`). DEV yang sudah menjalankan 320 lama dipindah lewat skrip transisi
// (SCRIPT-TABEL-KOLOM-BARU.xlsx sheet NB TREATY, dijalankan manusia).
//
// ⛔ Koreksi WO: migrasi 320 TIDAK diberi blok PL/SQL penjaga. Penjaga bentuk ada di inti:
// `praTerbangBentuk` (`inti/backend/migrasi/migrasi.go`) memeriksa SELURUH langkah yang belum tercatat
// SEBELUM satu pernyataan pun dikirim - setiap CREATE TABLE diurai `KolomCreateTable`, dan bila tabelnya
// sudah ada, kolom DDL dibandingkan kolom katalog Oracle lewat `SelisihKolom`.
//
// Uji di sini membuktikan, atas berkas migrasi tertanam yang sama dengan `-migrate`: (1) CREATE TABLE 320
// terurai PERSIS menjadi kolom T_GENERAL_POLIS_TREATY modul ini, dan (2) tidak satu pun pernyataan migrasi
// modul ini menyebut T_GENERAL_POLIS milik nbfacin; setiap FK induk anak menunjuk T_GENERAL_POLIS_TREATY.

import (
	"io/fs"
	"path"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul/nbtreatyin/backend/models"
)

const langkah320 = "320_t_general_polis_treaty.sql"

// kolomKunciGeneralPolis - kolom T_GENERAL_POLIS_TREATY di luar katalog medan
// (`models.TabelGeneralPolis`): kunci bersama T_WORK_POLIS (ID, ID-7), kunci
// generasi (NOPOLIS, PRODKE, NOENDORS, OLD_POLIS_ID; ID-8, ID-9), dan kolom
// json_polis (IDPEGA, TGL_INPUT, USERNAME; ID-21) - urutan DDL 320.
// IDPEGA DIBUANG (keputusan work owner 06-10-2026: kasus baru menyimpan IDPEGA = ID).
var kolomKunciGeneralPolis = []string{"ID", "NOPOLIS", "PRODKE", "NOENDORS", "OLD_POLIS_ID", "TGL_INPUT", "USERNAME"}

// kolomDDL320 menguraikan CREATE TABLE langkah 320 dengan pengurai
// pra-terbang inti; gagal bila jumlahnya bukan tepat satu.
func kolomDDL320(t *testing.T) (string, []string) {
	t.Helper()
	p, err := migrasi.PernyataanLangkah(berkasMigrasi, langkah320)
	if err != nil {
		t.Fatal(err)
	}
	var nama string
	var kolom []string
	n := 0
	for _, s := range p {
		if tb, k := migrasi.KolomCreateTable(s); tb != "" {
			n++
			nama, kolom = tb, k
		} else if strings.Contains(strings.ToUpper(s), "CREATE TABLE") {
			// pra-terbang menolak CREATE TABLE yang tidak terurai dengan galat
			// lain ("tidak dapat diurai") - bukan pesan bentuk yang dimaksud K18.
			t.Fatalf("%s: CREATE TABLE tidak terurai KolomCreateTable:\n%s", langkah320, s)
		}
	}
	if n != 1 {
		t.Fatalf("%s: %d CREATE TABLE terurai, harap tepat 1", langkah320, n)
	}
	return nama, kolom
}

// kolomDeklarasi320 membaca kolom CREATE TABLE 320 dengan cara LAIN dari
// pengurai inti: badan dipecah pada koma tingkat atas (bukan per baris), lalu
// setiap deklarasi selain CONSTRAINT diambil nama pertamanya. Pengurai inti
// membaca SATU kolom per baris; kolom kedua di baris yang sama tidak terlihat
// olehnya - dan karena itu juga tidak dibandingkan pra-terbang.
func kolomDeklarasi320(t *testing.T) []string {
	t.Helper()
	p, err := migrasi.PernyataanLangkah(berkasMigrasi, langkah320)
	if err != nil {
		t.Fatal(err)
	}
	var kolom []string
	for _, s := range p {
		m := migrasi.PolaCreateTabel.FindStringSubmatch(s)
		if m == nil {
			continue
		}
		dalam, awal := 0, 0
		var potong []string
		for i, c := range m[2] {
			switch c {
			case '(':
				dalam++
			case ')':
				dalam--
			case ',':
				if dalam == 0 {
					potong = append(potong, m[2][awal:i])
					awal = i + 1
				}
			}
		}
		potong = append(potong, m[2][awal:])
		for _, d := range potong {
			f := strings.Fields(strings.ToUpper(d))
			if len(f) == 0 || f[0] == "CONSTRAINT" {
				continue
			}
			kolom = append(kolom, f[0])
		}
	}
	return kolom
}

// (1) Pengurai pra-terbang membaca 320 persis: nama T_GENERAL_POLIS_TREATY, kolom =
// kunci + katalog medan, dan = setiap deklarasi kolom di badan DDL (nol kolom
// tersembunyi dari pra-terbang; baris CONSTRAINT bukan kolom). Sesudah 320
// berjalan, katalog Oracle = kolom DDL ini, sehingga menjalankan ulang tidak
// ditolak (SelisihKolom kosong -> `bentukCocok`).
func TestMigrasi320TeruraiPraTerbangPersisKolomGeneralPolis(t *testing.T) {
	nama, kolom := kolomDDL320(t)
	if nama != models.TabelGeneralPolis.Nama {
		t.Fatalf("tabel terurai %q, harap %q", nama, models.TabelGeneralPolis.Nama)
	}
	harap := append([]string{}, kolomKunciGeneralPolis...)
	for _, k := range models.TabelGeneralPolis.Kolom {
		harap = append(harap, k.Kolom)
	}
	if !reflect.DeepEqual(kolom, harap) {
		t.Fatalf("kolom terurai pra-terbang (%d):\n %v\nharap kunci + katalog (%d):\n %v", len(kolom), kolom, len(harap), harap)
	}
	if dekl := kolomDeklarasi320(t); !reflect.DeepEqual(kolom, dekl) {
		t.Fatalf("pra-terbang membaca %d kolom, badan DDL mendeklarasikan %d - kolom yang tidak terbaca "+
			"pengurai inti tidak pernah dibandingkan bentuknya:\n terurai %v\n deklarasi %v", len(kolom), len(dekl), kolom, dekl)
	}
	if kurang, lebih := migrasi.SelisihKolom(kolom, kolom); len(kurang) != 0 || len(lebih) != 0 {
		t.Fatalf("bentuk 320 lawan dirinya sendiri harus cocok: kurang %v, lebih %v", kurang, lebih)
	}
}

// (2) Perintah WO 05-10-2026: tabel induk NB Treaty In = T_GENERAL_POLIS_TREATY. T_GENERAL_POLIS milik
// nbfacin (182) tidak disebut satu pernyataan pun migrasi modul ini (komentar `--` tidak dihitung), dan FK
// induk setiap anak langsung (321, 323, 325, 326) menunjuk T_GENERAL_POLIS_TREATY. Nama harapan ditulis
// tangan, bukan dibaca dari katalog.
func TestMigrasiMemakaiTGeneralPolisTreatyBukanTabelFacIn(t *testing.T) {
	if nama, _ := kolomDDL320(t); nama != "T_GENERAL_POLIS_TREATY" {
		t.Fatalf("320 membuat %q, harap T_GENERAL_POLIS_TREATY", nama)
	}
	milikFacIn := regexp.MustCompile(`(?i)\bT_GENERAL_POLIS\b`)
	berkas, err := fs.Glob(berkasMigrasi, "migrations/*.sql")
	if err != nil || len(berkas) == 0 {
		t.Fatalf("berkas migrasi tidak terbaca: %v", err)
	}
	for _, b := range berkas {
		p, err := migrasi.PernyataanLangkah(berkasMigrasi, path.Base(b))
		if err != nil {
			t.Fatal(err)
		}
		for _, s := range p {
			if milikFacIn.MatchString(tanpaKomentar(s)) {
				t.Errorf("%s menyebut T_GENERAL_POLIS milik nbfacin:\n%s", path.Base(b), s)
			}
		}
	}
	for _, l := range []string{"321_t_polis_quotation.sql", "323_t_polis_instalment.sql", "325_t_polis_spreading.sql", "326_t_polis_xol.sql"} {
		p, err := migrasi.PernyataanLangkah(berkasMigrasi, l)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(strings.Join(p, "\n"), "FOREIGN KEY (POLIS_ID) REFERENCES {skema}.T_GENERAL_POLIS_TREATY (ID)") {
			t.Errorf("%s: FK POLIS_ID tidak menunjuk T_GENERAL_POLIS_TREATY", l)
		}
	}
}

// tanpaKomentar membuang komentar baris `--` dari satu pernyataan.
func tanpaKomentar(s string) string {
	baris := strings.Split(s, "\n")
	for i, x := range baris {
		if j := strings.Index(x, "--"); j >= 0 {
			baris[i] = x[:j]
		}
	}
	return strings.Join(baris, "\n")
}

// Koreksi WO 04-10-2026: 320 tetap DDL biasa - nol blok PL/SQL (penjaga ada
// di pra-terbang inti, bukan di berkas migrasi).
func TestMigrasi320TanpaBlokPLSQL(t *testing.T) {
	p, err := migrasi.PernyataanLangkah(berkasMigrasi, langkah320)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range p {
		atas := strings.ToUpper(s)
		if !strings.HasPrefix(atas, "CREATE TABLE ") && !strings.HasPrefix(atas, "CREATE UNIQUE INDEX ") {
			t.Errorf("%s memuat pernyataan selain CREATE TABLE / CREATE UNIQUE INDEX:\n%s", langkah320, s)
		}
		for _, kata := range []string{"BEGIN", "DECLARE", "EXECUTE IMMEDIATE"} {
			if strings.Contains(atas, kata) {
				t.Errorf("%s memuat %s - koreksi WO 04-10-2026: tanpa blok PL/SQL", langkah320, kata)
			}
		}
	}
}

func mengandung(daftar []string, s string) bool {
	for _, x := range daftar {
		if x == s {
			return true
		}
	}
	return false
}
