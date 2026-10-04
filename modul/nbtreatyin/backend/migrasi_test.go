package backend

// K18 (PROMPT putaran 3 bab 2; koreksi WO 04-10-2026) - TANPA Oracle.
//
// `POOLDATA.T_GENERAL_POLIS` sudah ada, dibuat migrasi `182_t_general_polis`
// milik `nbfacin` (di luar repo ini) dengan TUJUH kolom: ID, IDPEGA, COB_GROUP,
// START_DATE_TIME, OFFERING_DATE, END_DATE_TIME, FOLLOWING. Diagram grilling
// memberi nama itu ke Treaty In (PERMINTAAN-TIM-INTI C10).
//
// ⛔ Koreksi WO: migrasi 320 TIDAK diberi blok PL/SQL penjaga. Penjaganya sudah
// ada di inti: `praTerbangBentuk` (`inti/backend/migrasi/migrasi.go`) memeriksa
// SELURUH langkah yang belum tercatat SEBELUM satu pernyataan pun dikirim -
// setiap CREATE TABLE diurai `KolomCreateTable`, dan bila tabelnya sudah ada,
// kolom DDL dibandingkan kolom katalog Oracle lewat `SelisihKolom`; selisih
// apa pun menghentikan migrasi dengan galat:
//
//	repository: migrasi <langkah>: tabel <T> sudah ada di skema <S> tetapi
//	BENTUKNYA BERBEDA - kolom yang diminta migrasi tetapi tidak ada: [...];
//	kolom yang ada tetapi tidak diminta: [...]. Migrasi dihentikan sebelum
//	satu pernyataan pun dikirim; tidak ada yang diubah
//
// `praTerbangBentuk` sendiri butuh Oracle (katalog); uji di sini membuktikan
// dua bagian murninya atas migrasi 320 yang sebenarnya (berkas tertanam yang
// sama dengan `-migrate`): (1) CREATE TABLE 320 terurai PERSIS menjadi kolom
// T_GENERAL_POLIS modul ini, dan (2) bentuk tujuh kolom FacIn menghasilkan
// selisih tidak kosong di kedua arah - isi pesan galat di atas.

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul/nbtreatyin/backend/models"
)

const langkah320 = "320_t_general_polis.sql"

// kolomFacIn182 - bentuk `POOLDATA.T_GENERAL_POLIS` yang dibuat migrasi 182
// `nbfacin` (PROMPT putaran 3 bab 2 K18, katalog dicek 04-10-2026).
var kolomFacIn182 = []string{"ID", "IDPEGA", "COB_GROUP", "START_DATE_TIME", "OFFERING_DATE", "END_DATE_TIME", "FOLLOWING"}

// kolomKunciGeneralPolis - kolom T_GENERAL_POLIS di luar katalog medan
// (`models.TabelGeneralPolis`): kunci bersama T_WORK_POLIS (ID, ID-7), kunci
// generasi (NOPOLIS, PRODKE, NOENDORS, OLD_POLIS_ID; ID-8, ID-9), dan kolom
// json_polis (IDPEGA, TGL_INPUT, USERNAME; ID-21) - urutan DDL 320.
var kolomKunciGeneralPolis = []string{"ID", "NOPOLIS", "PRODKE", "NOENDORS", "OLD_POLIS_ID", "IDPEGA", "TGL_INPUT", "USERNAME"}

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

// (1) Pengurai pra-terbang membaca 320 persis: nama T_GENERAL_POLIS, kolom =
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

// (2) Bentuk FacIn tujuh kolom diberikan ke pembanding inti `SelisihKolom`
// (yang dipanggil `praTerbangBentuk`): kedua daftar dalam pesan galat terisi -
// pra-terbang menghentikan migrasi 320 sebelum apa pun berubah. Nilai harapan
// ditulis tangan dari tujuh kolom K18, bukan dihitung ulang.
func TestPraTerbangMenolakTGeneralPolisBentukFacIn(t *testing.T) {
	_, kolom := kolomDDL320(t)
	kurang, lebih := migrasi.SelisihKolom(kolom, kolomFacIn182)

	// "kolom yang ada tetapi tidak diminta" - lima kolom FacIn, terurut.
	if harap := []string{"COB_GROUP", "END_DATE_TIME", "FOLLOWING", "OFFERING_DATE", "START_DATE_TIME"}; !reflect.DeepEqual(lebih, harap) {
		t.Errorf("kolom yang ada tetapi tidak diminta = %v, harap %v", lebih, harap)
	}
	// "kolom yang diminta migrasi tetapi tidak ada" - seluruh kolom 320 kecuali
	// dua yang kebetulan bernama sama (ID, IDPEGA).
	if len(kurang) != len(kolom)-2 {
		t.Errorf("kolom yang diminta tetapi tidak ada: %d, harap %d (320 tanpa ID, IDPEGA)", len(kurang), len(kolom)-2)
	}
	for _, k := range []string{"NOPOLIS", "PRODKE", "OLD_POLIS_ID", "POSITION_NOTE", "PREMI_OGP", "OVERIDDING_COMM_ONP"} {
		if !mengandung(kurang, k) {
			t.Errorf("%s tidak disebut sebagai kolom yang diminta tetapi tidak ada: %v", k, kurang)
		}
	}
	for _, k := range []string{"ID", "IDPEGA"} {
		if mengandung(kurang, k) {
			t.Errorf("%s ada di kedua bentuk, tidak boleh disebut kurang", k)
		}
	}
	// Syarat penolakan `praTerbangBentuk`: len(kurang) > 0 || len(lebih) > 0.
	if len(kurang) == 0 && len(lebih) == 0 {
		t.Fatal("bentuk FacIn lolos pra-terbang - migrasi 320 akan melewati tabel milik nbfacin")
	}
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
