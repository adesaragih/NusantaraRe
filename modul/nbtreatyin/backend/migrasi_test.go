package backend

// T_GENERAL_POLIS BERSAMA FacIn + Treaty In - TANPA Oracle.
//
// ⛔ KEPUTUSAN WORK OWNER 04-10-2026 (mengalahkan bab 0 butir 11 dan K18
// PROMPT putaran 3): `T_GENERAL_POLIS` adalah SATU tabel bersama FacIn dan
// Treaty In. Tabel DASARNYA dibuat migrasi `182_t_general_polis` milik
// `nbfacin` (sudah dijalankan di POOLDATA; berkasnya belum di repo) dengan
// tujuh kolom: ID (VARCHAR2(32), PK), IDPEGA (VARCHAR2(50)), COB_GROUP,
// START_DATE_TIME, OFFERING_DATE, END_DATE_TIME, FOLLOWING. Migrasi 320 modul
// ini TIDAK membuat tabel itu: ia hanya MENAMBAH kolom Treaty lewat
// `ALTER TABLE {skema}.T_GENERAL_POLIS ADD (` biasa, lalu constraint dan
// indeks Treaty. Jalur mundurnya hanya membuang milik Treaty.
//
// Yang dibuktikan di sini atas berkas tertanam yang SAMA dengan `-migrate`:
//
//  1. 320 nol CREATE TABLE - pra-terbang inti (`praTerbangBentuk`, hanya
//     membandingkan CREATE TABLE) tidak lagi menolak tabel dasar FacIn;
//  2. kolom yang DITAMBAHKAN 320 terbaca pengurai inti `KolomAlterTambah`
//     (pembanding STRUKTUR inti memakainya) PERSIS sama dengan pembacaan
//     cara lain, dan ditambah kolom dasar yang dipakai (ID, IDPEGA) sama
//     PERSIS dengan kolom Treaty (kunci + katalog medan);
//  3. tidak satu pun kolom dasar FacIn ditambahkan, dibuang, atau disebut
//     jalur mundur.

import (
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul/nbtreatyin/backend/models"
)

const (
	langkah320      = "320_t_general_polis.sql"
	langkah320Turun = "320_t_general_polis_down.sql"
)

// kolomDasarFacIn182 - tujuh kolom tabel dasar yang dibuat migrasi 182
// `nbfacin` (keputusan WO 04-10-2026; katalog POOLDATA dicek 04-10-2026).
var kolomDasarFacIn182 = []string{"ID", "IDPEGA", "COB_GROUP", "START_DATE_TIME", "OFFERING_DATE", "END_DATE_TIME", "FOLLOWING"}

// kolomDasarDipakaiTreaty - kolom tabel dasar yang DIBACA/DITULIS Treaty In:
// ID (kunci bersama T_WORK_POLIS, ID-7) dan IDPEGA (json_polis, ID-21).
var kolomDasarDipakaiTreaty = []string{"ID", "IDPEGA"}

// kolomKunciGeneralPolis - kolom T_GENERAL_POLIS Treaty di luar katalog medan
// (`models.TabelGeneralPolis`): kunci bersama (ID), kunci generasi (NOPOLIS,
// PRODKE, NOENDORS, OLD_POLIS_ID; ID-8, ID-9), dan kolom json_polis (IDPEGA,
// TGL_INPUT, USERNAME; ID-21).
var kolomKunciGeneralPolis = []string{"ID", "NOPOLIS", "PRODKE", "NOENDORS", "OLD_POLIS_ID", "IDPEGA", "TGL_INPUT", "USERNAME"}

func pernyataan(t *testing.T, nama string) []string {
	t.Helper()
	p, err := migrasi.PernyataanLangkah(berkasMigrasi, nama)
	if err != nil {
		t.Fatal(err)
	}
	if len(p) == 0 {
		t.Fatalf("%s: nol pernyataan", nama)
	}
	return p
}

// kolomAlter320 - kolom yang ditambahkan 320 menurut pengurai INTI.
func kolomAlter320(t *testing.T) []string {
	t.Helper()
	var kolom []string
	n := 0
	for _, s := range pernyataan(t, langkah320) {
		if tb, k := migrasi.KolomAlterTambah(s); tb != "" {
			n++
			if tb != models.TabelGeneralPolis.Nama {
				t.Fatalf("%s menambah kolom ke %s, harap %s", langkah320, tb, models.TabelGeneralPolis.Nama)
			}
			kolom = append(kolom, k...)
		}
	}
	if n != 1 {
		t.Fatalf("%s: %d ALTER TABLE ... ADD ( terurai, harap tepat 1", langkah320, n)
	}
	return kolom
}

// kolomAlterCaraLain membaca badan `ADD (...)` dengan cara LAIN dari pengurai
// inti: dipecah pada koma TINGKAT ATAS (bukan setiap koma), lalu nama pertama
// tiap deklarasi. Pengurai inti memecah pada setiap koma dan hanya mengenali
// potongan yang diawali nama kolom - kolom yang luput darinya tidak terlihat
// pembanding STRUKTUR mana pun.
func kolomAlterCaraLain(t *testing.T) []string {
	t.Helper()
	pola := regexp.MustCompile(`(?is)^ALTER\s+TABLE\s+\{skema\}\.T_GENERAL_POLIS\s+ADD\s*\((.*)\)\s*$`)
	var kolom []string
	for _, s := range pernyataan(t, langkah320) {
		m := pola.FindStringSubmatch(strings.TrimSpace(s))
		if m == nil {
			continue
		}
		dalam, awal := 0, 0
		var potong []string
		for i, c := range m[1] {
			switch c {
			case '(':
				dalam++
			case ')':
				dalam--
			case ',':
				if dalam == 0 {
					potong = append(potong, m[1][awal:i])
					awal = i + 1
				}
			}
		}
		potong = append(potong, m[1][awal:])
		for _, d := range potong {
			if f := strings.Fields(strings.ToUpper(d)); len(f) > 0 {
				kolom = append(kolom, f[0])
			}
		}
	}
	return kolom
}

func terurut(s []string) []string {
	out := append([]string{}, s...)
	sort.Strings(out)
	return out
}

func mengandung(daftar []string, s string) bool {
	for _, x := range daftar {
		if x == s {
			return true
		}
	}
	return false
}

// (1) Nol CREATE TABLE di 320: tabelnya milik dasar bersama (nbfacin 182).
// Setiap pernyataan adalah ALTER atas T_GENERAL_POLIS atau CREATE UNIQUE INDEX
// atasnya; nol PL/SQL (penjaga inti `pelanggaranBlokPLSQL` menolak kolom baru
// lewat blok di jalur maju).
func TestMigrasi320HanyaMengubahTabelBersama(t *testing.T) {
	for _, s := range pernyataan(t, langkah320) {
		atas := strings.ToUpper(s)
		if tb, _ := migrasi.KolomCreateTable(s); tb != "" || strings.Contains(atas, "CREATE TABLE") {
			t.Errorf("%s membuat tabel - T_GENERAL_POLIS adalah tabel dasar bersama nbfacin 182 (WO 04-10-2026):\n%s", langkah320, s)
		}
		if !strings.HasPrefix(atas, "ALTER TABLE {SKEMA}.T_GENERAL_POLIS ") &&
			!strings.HasPrefix(atas, "CREATE UNIQUE INDEX {SKEMA}.UQ_GENERAL_POLIS_NOPOLIS ON {SKEMA}.T_GENERAL_POLIS ") {
			t.Errorf("%s memuat pernyataan selain ALTER TABLE / CREATE UNIQUE INDEX atas T_GENERAL_POLIS:\n%s", langkah320, s)
		}
		for _, kata := range []string{"BEGIN", "DECLARE", "EXECUTE IMMEDIATE", " MODIFY"} {
			if strings.Contains(atas, kata) {
				t.Errorf("%s memuat %q - kolom Treaty ditambah ALTER ... ADD ( biasa, kolom FacIn tidak diubah", langkah320, kata)
			}
		}
		// ALTER bukan "pernyataan buat": galatnya (mis. ORA-01430 kolom sudah
		// ada) TIDAK ditelan pelari migrasi inti.
		if strings.HasPrefix(atas, "ALTER") && migrasi.PernyataanBuat(s) {
			t.Errorf("ALTER dianggap pernyataan CREATE oleh pelari inti: %s", migrasi.RingkasPernyataan(s))
		}
	}
}

// (2) Kolom yang ditambahkan 320 + kolom dasar yang dipakai = kolom Treaty
// (kunci + katalog medan) PERSIS; pengurai inti membaca setiap deklarasi.
func TestMigrasi320MenambahPersisKolomTreaty(t *testing.T) {
	alter := kolomAlter320(t)
	if lain := kolomAlterCaraLain(t); !reflect.DeepEqual(alter, lain) {
		t.Fatalf("pengurai inti membaca %d kolom, badan ADD mendeklarasikan %d - kolom yang luput tidak "+
			"terlihat pembanding STRUKTUR:\n inti %v\n deklarasi %v", len(alter), len(lain), alter, lain)
	}
	for _, k := range alter {
		if mengandung(kolomDasarFacIn182, k) {
			t.Errorf("320 menambah kolom dasar FacIn %s - kolom itu milik nbfacin 182", k)
		}
	}
	harap := append([]string{}, kolomKunciGeneralPolis...)
	for _, k := range models.TabelGeneralPolis.Kolom {
		harap = append(harap, k.Kolom)
	}
	dapat := append(append([]string{}, kolomDasarDipakaiTreaty...), alter...)
	if !reflect.DeepEqual(terurut(dapat), terurut(harap)) {
		t.Fatalf("ALTER 320 + kolom dasar dipakai (%d):\n %v\nharap kunci + katalog Treaty (%d):\n %v",
			len(dapat), terurut(dapat), len(harap), terurut(harap))
	}
	// Nilai tangan: 6 kolom kunci/generasi/json_polis Treaty + katalog.
	if len(alter) != 6+len(models.TabelGeneralPolis.Kolom) {
		t.Fatalf("320 menambah %d kolom, harap 6 + %d katalog", len(alter), len(models.TabelGeneralPolis.Kolom))
	}
	for _, k := range models.TabelGeneralPolis.Kolom {
		if mengandung(kolomDasarFacIn182, k.Kolom) {
			t.Errorf("katalog Treaty memuat kolom dasar FacIn %s", k.Kolom)
		}
	}
}

// Constraint dan indeks Treaty (diagram F12-F16; ID-7..ID-9): FK OLD_POLIS_ID,
// UQ OLD_POLIS_ID, FK ID -> T_WORK_POLIS (shared PK; ASUMSI: 182 belum
// memilikinya - PERMINTAAN C10), indeks unik NOPOLIS/PRODKE bernomor.
func TestMigrasi320ConstraintTreaty(t *testing.T) {
	gabung := strings.Join(pernyataan(t, langkah320), "\n/\n")
	for _, w := range []string{
		"ALTER TABLE {skema}.T_GENERAL_POLIS ADD CONSTRAINT FK_GENERAL_POLIS_OLD FOREIGN KEY (OLD_POLIS_ID) REFERENCES {skema}.T_WORK_POLIS (ID)",
		"ALTER TABLE {skema}.T_GENERAL_POLIS ADD CONSTRAINT UQ_GENERAL_POLIS_OLD UNIQUE (OLD_POLIS_ID)",
		"ALTER TABLE {skema}.T_GENERAL_POLIS ADD CONSTRAINT FK_GENERAL_POLIS_WORK FOREIGN KEY (ID) REFERENCES {skema}.T_WORK_POLIS (ID)",
		"CREATE UNIQUE INDEX {skema}.UQ_GENERAL_POLIS_NOPOLIS ON {skema}.T_GENERAL_POLIS (CASE WHEN NOPOLIS IS NOT NULL THEN NOPOLIS END, CASE WHEN NOPOLIS IS NOT NULL THEN PRODKE END)",
	} {
		if !strings.Contains(gabung, w) {
			t.Errorf("%s tidak memuat:\n %s", langkah320, w)
		}
	}
	// Kunci utama milik tabel dasar (182) - 320 tidak membuatnya ulang.
	if strings.Contains(strings.ToUpper(gabung), "PRIMARY KEY") {
		t.Errorf("%s membuat PRIMARY KEY - PK(ID) milik tabel dasar nbfacin 182", langkah320)
	}
}

// polaPengenal - kata pengenal SQL utuh.
var polaPengenal = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_$#]*`)

// (3) Jalur mundur 320 hanya membuang milik Treaty: nol DROP TABLE, nol kolom
// dasar FacIn disebut (kata utuh), dan kolom yang dibuang = kolom yang
// ditambahkan 320 persis.
func TestMigrasi320TurunHanyaMembuangMilikTreaty(t *testing.T) {
	turun := pernyataan(t, langkah320Turun)
	polaBuang := regexp.MustCompile(`(?is)^ALTER\s+TABLE\s+\{skema\}\.T_GENERAL_POLIS\s+DROP\s*\((.*)\)\s*$`)
	var dibuang []string
	for _, s := range turun {
		atas := strings.ToUpper(s)
		if strings.Contains(atas, "DROP TABLE") || strings.Contains(atas, "TRUNCATE") || strings.Contains(atas, "DELETE") {
			t.Errorf("%s membuang tabel/baris - tabel dasar dan baris FacIn tidak disentuh:\n%s", langkah320Turun, s)
		}
		for _, kata := range polaPengenal.FindAllString(atas, -1) {
			if mengandung(kolomDasarFacIn182, kata) {
				t.Errorf("%s menyebut kolom dasar FacIn %s:\n%s", langkah320Turun, kata, s)
			}
		}
		if m := polaBuang.FindStringSubmatch(strings.TrimSpace(s)); m != nil {
			for _, k := range strings.Split(m[1], ",") {
				if k = strings.ToUpper(strings.TrimSpace(k)); k != "" {
					dibuang = append(dibuang, k)
				}
			}
		}
	}
	if alter := kolomAlter320(t); !reflect.DeepEqual(terurut(dibuang), terurut(alter)) {
		t.Fatalf("kolom dibuang jalur mundur (%d) != kolom ditambah 320 (%d):\n %v\n %v",
			len(dibuang), len(alter), terurut(dibuang), terurut(alter))
	}
	gabung := strings.Join(turun, "\n")
	for _, w := range []string{"DROP INDEX {skema}.UQ_GENERAL_POLIS_NOPOLIS", "DROP CONSTRAINT FK_GENERAL_POLIS_WORK",
		"DROP CONSTRAINT FK_GENERAL_POLIS_OLD", "DROP CONSTRAINT UQ_GENERAL_POLIS_OLD"} {
		if !strings.Contains(gabung, w) {
			t.Errorf("%s tidak memuat %q", langkah320Turun, w)
		}
	}
}
