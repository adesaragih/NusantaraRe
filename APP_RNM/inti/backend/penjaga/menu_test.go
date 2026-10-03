package penjaga

// Penjaga M_NAV_MENU dan ISI menunya - isi awal `inti` 900, bentuk datar 901,
// dan slot menu setiap modul (R3) - TANPA Oracle.
//
// Untuk apa berkas ini: menu aplikasi dirakit dari tabel
// (`PROMPT-MENU-DARI-TABEL-M_NAV_MENU.md`), dan sejak keputusan work owner
// 30-09-2026 (`PROMPT-MENU-DATAR-PER-GROUPMENU.md`) SATU MODUL = SATU BARIS:
// 901 membuang lima butir anak 900 beserta kolom `PARENT_ID`. Isinya bukan data
// uji - ia menu yang dilihat setiap pemakai, dan kelak dasar akses per akun.
//
// ⛔ Yang dibaca: HASIL BERSIH seluruh langkah yang menyentuh M_NAV_MENU -
// 900 (isi awal, milik `inti`), 901 (bentuk datar, milik `inti`), lalu berkas
// slot menu setiap modul - diterapkan menurut urutan pelari ke SKEMA TIRUAN
// (`skemaMenu`): kolom, kunci tamu, indeks, dan baris. Blok berpelindung
// katalog (901) dievaluasi terhadap katalog tiruan itu, jadi langkah yang tidak
// terlindung dan diulang ikut ketahuan (ORA-00904 di skema tiruan). 900 yang
// sudah dijalankan tidak boleh disunting (`T_MIGRASI` mencatat namanya).

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/menu"
	"nusantarare/inti/backend/migrasi"
)

// golonganMenu - golongan yang `GET /api/menu` kenal (`menu.Golongan`, urutan
// tampil sidebar). ⛔ BUKAN salinan: golongan yang ditambah ke CHECK tetapi
// tidak ke `menu.Golongan` akan dibuang `Susun` diam-diam - uji CHECK di bawah
// menagihnya.
var golonganMenu = menu.Golongan

// kelompokMenu adalah satu baris MODUL di isi menu (dulu "kelompok").
type kelompokMenu struct {
	kode, label, golongan, modul string
	urutan                       int
	dimigrasi                    string
}

// butirMenu adalah satu baris butir anak (hanya 900, dibuang 901).
type butirMenu struct {
	kode, label, induk string
	urutan             int
}

var (
	polaIsiKelompok = regexp.MustCompile(`(?s)^INSERT INTO \{skema\}\.M_NAV_MENU \(ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI\)\s+` +
		`SELECT \{skema\}\.SEQ_M_NAV_MENU\.NEXTVAL, NULL, '([^']+)', '([^']+)', '([^']+)', '([^']+)', (\d+), '([01])' FROM DUAL\s+` +
		`WHERE NOT EXISTS \(SELECT 1 FROM \{skema\}\.M_NAV_MENU WHERE KODE = '([^']+)'\)$`)
	// polaIsiKelompokDatar - baris modul DI LUAR korpus (`modulLuarKorpus`): bentuk datar sesudah 901, tanpa
	// PARENT_ID, oleh langkah inti tersendiri (PANDUAN-TIM-PER-MODUL bab 5).
	polaIsiKelompokDatar = regexp.MustCompile(`(?s)^INSERT INTO \{skema\}\.M_NAV_MENU \(ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI\)\s+` +
		`SELECT \{skema\}\.SEQ_M_NAV_MENU\.NEXTVAL, '([^']+)', '([^']+)', '([^']+)', '([^']+)', (\d+), '([01])' FROM DUAL\s+` +
		`WHERE NOT EXISTS \(SELECT 1 FROM \{skema\}\.M_NAV_MENU WHERE KODE = '([^']+)'\)$`)
	polaIsiButir = regexp.MustCompile(`(?s)^INSERT INTO \{skema\}\.M_NAV_MENU \(ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI\)\s+` +
		`SELECT \{skema\}\.SEQ_M_NAV_MENU\.NEXTVAL, k\.ID, '([^']+)', '([^']+)', k\.GROUPMENU, k\.MODUL, (\d+), k\.DIMIGRASI\s+` +
		`FROM \{skema\}\.M_NAV_MENU k\s+WHERE k\.KODE = '([^']+)' AND k\.PARENT_ID IS NULL\s+` +
		`AND NOT EXISTS \(SELECT 1 FROM \{skema\}\.M_NAV_MENU b WHERE b\.KODE = '([^']+)'\)$`)
	// polaUbahDimigrasi - satu-satunya UPDATE yang dikenal, dan satu-satunya
	// isi slot menu modul: modul yang mendapat layar pertamanya menyalakan
	// DIMIGRASI barisnya (panduan bab 6). Diterapkan menurut urutan langkah,
	// jadi DIMIGRASI yang dibaca adalah keadaan SESUDAH migrasi.
	// ⛔ Tanpa `AND PARENT_ID IS NULL` (bentuk sebelum 901): kolom itu dibuang
	// 901, dan bentuk lama mati di ORA-00904 - skema tiruan menagihnya.
	polaUbahDimigrasi = regexp.MustCompile(`(?s)^UPDATE \{skema\}\.M_NAV_MENU SET DIMIGRASI = '([01])', TGL_UBAH = SYSDATE\s+` +
		`WHERE KODE = '([^']+)'$`)
	// polaUbahLabel - nama tampilan baris modul (keputusan work owner 03-10-2026), SAH HANYA untuk baris
	// `labelTampilDisetujui` di slot menu modulnya sendiri (`pelanggaranSlotMenu`).
	polaUbahLabel = regexp.MustCompile(`(?s)^UPDATE \{skema\}\.M_NAV_MENU SET LABEL = '([^']+)', TGL_UBAH = SYSDATE\s+` +
		`WHERE KODE = '([^']+)'$`)
	polaIndeksMenu = regexp.MustCompile(`^CREATE INDEX \{skema\}\.(\w+) ON \{skema\}\.M_NAV_MENU \((\w+)\)$`)
	// polaUbahDimigrasiLama - bentuk slot SEBELUM 901 (menyebut PARENT_ID).
	// Dikenal supaya skema tiruan menolaknya dengan SEBAB yang benar -
	// ORA-00904 sesudah 901 - bukan "bentuk tidak dikenal" (temuan /code-review).
	polaUbahDimigrasiLama = regexp.MustCompile(`(?s)^UPDATE \{skema\}\.M_NAV_MENU SET DIMIGRASI = '([01])', TGL_UBAH = SYSDATE\s+` +
		`WHERE KODE = '([^']+)' AND PARENT_ID IS NULL$`)
	polaKonstrain = regexp.MustCompile(`CONSTRAINT (\w+) `)
)

// labelTampil - nama tampilan baris modul yang BUKAN nama folder korpus.
type labelTampil struct{ folder, tampil string }

// labelTampilDisetujui - SATU-SATUNYA baris modul yang LABEL-nya bukan nama folder korpus: nama tampilan tanpa kata
// "Master" (keputusan work owner 03-10-2026 "ganti nama modul ... hapus kata Master nya", nama tampilan saja - KODE,
// MODUL, folder, rute, dan MODUL_AKTIF tetap). Diubah slot menu modulnya sendiri (961 / 959); mundurnya
// mengembalikan nama folder. Setiap baris lain tetap LABEL = nama folder VERBATIM. Jumlahnya dikunci.
// modulLuarKorpus - baris modul DI LUAR dua puluh folder korpus (PANDUAN-TIM-PER-MODUL bab 5): KODE -> berkas
// migrasi inti yang membuatnya (bentuk datar `polaIsiKelompokDatar`, sesudah 901). Isi awal 900 tidak disunting.
// Langkah inti ini SATU-SATUNYA selain 900/901 yang boleh menyentuh M_NAV_MENU, dan diterapkan skema tiruan
// (`langkahMenu`). Barisnya lahir SESUDAH 903, jadi isi awal hak menu tidak menagihnya - admin memberi hak lewat
// Kelola User.
//
// Keputusan work owner 03-10-2026: modul `marketingofficer`, label "Marketing Officer", kelompok MASTER.
var modulLuarKorpus = map[string]string{
	"marketingofficer": "906_m_nav_menu_marketingofficer.sql",
}

// langkahMenuLuarKorpus menjawab apakah berkas inti `nama` membuat baris modul luar korpus.
func langkahMenuLuarKorpus(nama string) bool {
	for _, b := range modulLuarKorpus {
		if b == nama {
			return true
		}
	}
	return false
}

var labelTampilDisetujui = map[string]labelTampil{
	"masterproductnamelife":   {folder: "Master Product Name Life", tampil: "Product Name Life"},
	"mastercontractretrolife": {folder: "Master Contract Retro Life", tampil: "Contract Retro Life"},
}

// barisMenu - satu baris M_NAV_MENU di skema tiruan. `induk` kosong = baris
// modul; terisi = butir anak (900 saja).
type barisMenu struct {
	kode, label, golongan, modul, dimigrasi, induk string
	urutan                                         int
}

// skemaMenu - M_NAV_MENU tiruan: kolom, kunci tamu ke induk, indeks, baris.
type skemaMenu struct {
	ada   bool
	kolom map[string]bool
	fk    bool
	// konstrain - constraint lain CREATE TABLE (PK, UQ, CK): ada di katalog
	// tiruan, jadi blok yang memeriksanya dievaluasi seperti di Oracle.
	konstrain map[string]bool
	indeks    map[string]string // nama -> kolom
	baris     []barisMenu
}

func skemaMenuKosong() *skemaMenu {
	return &skemaMenu{kolom: map[string]bool{}, konstrain: map[string]bool{}, indeks: map[string]string{}}
}

func (s *skemaMenu) cari(kode string) int {
	for i, b := range s.baris {
		if b.kode == kode {
			return i
		}
	}
	return -1
}

// perluKolom meniru ORA-00904: pernyataan yang menyebut kolom yang tidak ada.
func (s *skemaMenu) perluKolom(k string) error {
	if !s.kolom[k] {
		return fmt.Errorf("ORA-00904: %s tidak ada di M_NAV_MENU", k)
	}
	return nil
}

// terapkan menjalankan satu pernyataan ke skema tiruan. Bentuk yang tidak
// dikenal adalah GALAT: baris yang tidak terbaca penjaga ini tidak terbaca pula
// oleh uji dua arah frontend.
func (s *skemaMenu) terapkan(p string) error {
	if strings.HasPrefix(p, "CREATE SEQUENCE ") {
		return nil
	}
	if nama, kolom := migrasi.KolomCreateTable(p); nama == "M_NAV_MENU" {
		s.ada = true
		for _, k := range kolom {
			s.kolom[k] = true
		}
		s.fk = strings.Contains(p, "CONSTRAINT FK_M_NAV_MENU_INDUK FOREIGN KEY (PARENT_ID)")
		for _, m := range polaKonstrain.FindAllStringSubmatch(p, -1) {
			if m[1] != "FK_M_NAV_MENU_INDUK" {
				s.konstrain[m[1]] = true
			}
		}
		return nil
	}
	if m := polaIndeksMenu.FindStringSubmatch(p); m != nil {
		if err := s.perluKolom(m[2]); err != nil {
			return err
		}
		s.indeks[m[1]] = m[2]
		return nil
	}
	if m := polaIsiKelompok.FindStringSubmatch(p); m != nil {
		if m[7] != m[1] {
			return fmt.Errorf("INSERT baris modul %s memeriksa NOT EXISTS atas %s", m[1], m[7])
		}
		if err := s.perluKolom("PARENT_ID"); err != nil {
			return err
		}
		if s.cari(m[1]) < 0 {
			u, _ := strconv.Atoi(m[5])
			s.baris = append(s.baris, barisMenu{kode: m[1], label: m[2], golongan: m[3], modul: m[4], urutan: u, dimigrasi: m[6]})
		}
		return nil
	}
	if m := polaIsiKelompokDatar.FindStringSubmatch(p); m != nil {
		if m[7] != m[1] {
			return fmt.Errorf("INSERT baris modul %s memeriksa NOT EXISTS atas %s", m[1], m[7])
		}
		if _, luar := modulLuarKorpus[m[1]]; !luar {
			return fmt.Errorf("INSERT datar baris modul %s - hanya modul di luar korpus (modulLuarKorpus)", m[1])
		}
		if s.cari(m[1]) < 0 {
			u, _ := strconv.Atoi(m[5])
			s.baris = append(s.baris, barisMenu{kode: m[1], label: m[2], golongan: m[3], modul: m[4], urutan: u, dimigrasi: m[6]})
		}
		return nil
	}
	if m := polaIsiButir.FindStringSubmatch(p); m != nil {
		if m[5] != m[1] {
			return fmt.Errorf("INSERT butir %s memeriksa NOT EXISTS atas %s", m[1], m[5])
		}
		if err := s.perluKolom("PARENT_ID"); err != nil {
			return err
		}
		i := s.cari(m[4])
		if i < 0 || s.baris[i].induk != "" || s.cari(m[1]) >= 0 {
			return nil // SELECT ... WHERE k.KODE = induk AND NOT EXISTS: nol baris
		}
		k := s.baris[i]
		u, _ := strconv.Atoi(m[3])
		s.baris = append(s.baris, barisMenu{kode: m[1], label: m[2], golongan: k.golongan, modul: k.modul,
			dimigrasi: k.dimigrasi, induk: k.kode, urutan: u})
		return nil
	}
	if m := polaUbahDimigrasiLama.FindStringSubmatch(p); m != nil {
		if err := s.perluKolom("PARENT_ID"); err != nil {
			return err
		}
		i := s.cari(m[2])
		if i < 0 || s.baris[i].induk != "" {
			return fmt.Errorf("UPDATE DIMIGRASI atas %s - baris modul itu tidak (belum) ada", m[2])
		}
		s.baris[i].dimigrasi = m[1]
		return nil
	}
	if m := polaUbahDimigrasi.FindStringSubmatch(p); m != nil {
		i := s.cari(m[2])
		if i < 0 || s.baris[i].induk != "" {
			return fmt.Errorf("UPDATE DIMIGRASI atas %s - baris modul itu tidak (belum) ada", m[2])
		}
		s.baris[i].dimigrasi = m[1]
		return nil
	}
	if m := polaUbahLabel.FindStringSubmatch(p); m != nil {
		i := s.cari(m[2])
		if i < 0 || s.baris[i].induk != "" {
			return fmt.Errorf("UPDATE LABEL atas %s - baris modul itu tidak (belum) ada", m[2])
		}
		s.baris[i].label = m[1]
		return nil
	}
	if pk, ok := migrasi.BacaPerintahKatalog(p); ok {
		if pk.Tabel != "M_NAV_MENU" {
			return fmt.Errorf("blok berpelindung katalog atas %s, bukan M_NAV_MENU", pk.Tabel)
		}
		// ⛔ Perintahnya DIKENALI walau pemeriksaan katalog melewatinya:
		// perintah yang dilewati di skema tiruan dapat berjalan di Oracle
		// (temuan /code-review).
		if !perintahDikenal(pk.Perintah) {
			return fmt.Errorf("perintah blok berpelindung tidak dikenal penjaga: %s", migrasi.RingkasPernyataan(pk.Perintah))
		}
		var ada bool
		switch pk.Katalog {
		case "ALL_TAB_COLUMNS":
			ada = s.kolom[pk.Objek]
		case "ALL_CONSTRAINTS":
			ada = (pk.Objek == "FK_M_NAV_MENU_INDUK" && s.fk) || s.konstrain[pk.Objek]
		case "ALL_INDEXES":
			_, ada = s.indeks[pk.Objek]
		}
		if ada != pk.BilaAda {
			return nil
		}
		return s.terapkanPerintah(pk.Perintah)
	}
	return fmt.Errorf("bentuk tidak dikenal penjaga: %s", migrasi.RingkasPernyataan(p))
}

// perintahDikenal menjawab apakah isi EXECUTE IMMEDIATE salah satu bentuk
// yang `terapkanPerintah` kenal - tanpa menerapkannya.
func perintahDikenal(q string) bool {
	switch q {
	case "DELETE FROM {skema}.M_NAV_MENU WHERE PARENT_ID IS NOT NULL",
		"ALTER TABLE {skema}.M_NAV_MENU DROP CONSTRAINT FK_M_NAV_MENU_INDUK",
		"ALTER TABLE {skema}.M_NAV_MENU DROP COLUMN PARENT_ID",
		"ALTER TABLE {skema}.M_NAV_MENU ADD (PARENT_ID NUMBER(10))",
		"ALTER TABLE {skema}.M_NAV_MENU ADD CONSTRAINT FK_M_NAV_MENU_INDUK FOREIGN KEY (PARENT_ID) REFERENCES {skema}.M_NAV_MENU (ID)":
		return true
	}
	return q == "DROP INDEX {skema}.IX_M_NAV_MENU_PARENT" || polaIndeksMenu.MatchString(q)
}

// terapkanPerintah - isi EXECUTE IMMEDIATE blok berpelindung katalog.
func (s *skemaMenu) terapkanPerintah(q string) error {
	switch {
	case q == "DELETE FROM {skema}.M_NAV_MENU WHERE PARENT_ID IS NOT NULL":
		if err := s.perluKolom("PARENT_ID"); err != nil {
			return err
		}
		var sisa []barisMenu
		for _, b := range s.baris {
			if b.induk == "" {
				sisa = append(sisa, b)
			}
		}
		s.baris = sisa
	case q == "ALTER TABLE {skema}.M_NAV_MENU DROP CONSTRAINT FK_M_NAV_MENU_INDUK":
		if !s.fk {
			return fmt.Errorf("ORA-02443: FK_M_NAV_MENU_INDUK tidak ada")
		}
		s.fk = false
	case strings.HasPrefix(q, "DROP INDEX {skema}."):
		nama := strings.TrimPrefix(q, "DROP INDEX {skema}.")
		if _, ada := s.indeks[nama]; !ada {
			return fmt.Errorf("ORA-01418: indeks %s tidak ada", nama)
		}
		delete(s.indeks, nama)
	case q == "ALTER TABLE {skema}.M_NAV_MENU DROP COLUMN PARENT_ID":
		if err := s.perluKolom("PARENT_ID"); err != nil {
			return err
		}
		// Oracle ikut membuang indeks dan kunci tamu satu-kolom atas kolom itu.
		delete(s.kolom, "PARENT_ID")
		s.fk = false
		for n, k := range s.indeks {
			if k == "PARENT_ID" {
				delete(s.indeks, n)
			}
		}
	case q == "ALTER TABLE {skema}.M_NAV_MENU ADD (PARENT_ID NUMBER(10))":
		if s.kolom["PARENT_ID"] {
			return fmt.Errorf("ORA-01430: PARENT_ID sudah ada")
		}
		s.kolom["PARENT_ID"] = true
	case q == "ALTER TABLE {skema}.M_NAV_MENU ADD CONSTRAINT FK_M_NAV_MENU_INDUK FOREIGN KEY (PARENT_ID) REFERENCES {skema}.M_NAV_MENU (ID)":
		if err := s.perluKolom("PARENT_ID"); err != nil {
			return err
		}
		if s.fk {
			return fmt.Errorf("ORA-02275: FK_M_NAV_MENU_INDUK sudah ada")
		}
		s.fk = true
	default:
		if m := polaIndeksMenu.FindStringSubmatch(q); m != nil {
			if _, ada := s.indeks[m[1]]; ada {
				return fmt.Errorf("ORA-00955: %s sudah ada", m[1])
			}
			return s.terapkan(q)
		}
		return fmt.Errorf("perintah blok berpelindung tidak dikenal penjaga: %s", migrasi.RingkasPernyataan(q))
	}
	return nil
}

// langkahMenu - langkah (maju) yang menyentuh M_NAV_MENU, urutan pelari: 900
// saja bila `hanyaIsiAwal`; selainnya 900, 901, dan slot menu setiap modul.
func langkahMenu(t *testing.T, hanyaIsiAwal bool) []migrasi.Langkah {
	t.Helper()
	langkah, err := migrasi.Daftar(false, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	var jatah map[string]jatahModul
	if !hanyaIsiAwal {
		jatah = jatahSetiapModul(t)
	}
	var out []migrasi.Langkah
	for _, m := range langkah {
		pemilik := berkasMigrasi.modul(m.Nama)
		n, _ := nomorBerkas(m.Nama)
		isiAwal := pemilik == "inti" && strings.HasPrefix(m.Nama, "900_")
		datar := !hanyaIsiAwal && pemilik == "inti" && (strings.HasPrefix(m.Nama, "901_") || langkahMenuLuarKorpus(m.Nama))
		slotModul := !hanyaIsiAwal && pemilik != "inti" && jatah[pemilik].diSlot(n)
		if isiAwal || datar || slotModul {
			out = append(out, m)
		}
	}
	if len(out) == 0 {
		t.Fatal("nol langkah menu; 900_m_nav_menu.sql hilang atau pembacanya rusak")
	}
	return out
}

// pernyataanMenu - pernyataan maju `langkahMenu`, berurutan.
func pernyataanMenu(t *testing.T, hanyaIsiAwal bool) []string {
	t.Helper()
	var out []string
	for _, m := range langkahMenu(t, hanyaIsiAwal) {
		out = append(out, m.Pernyataan...)
	}
	return out
}

// skemaSesudah menerapkan pernyataan ke skema tiruan; galat menjadi t.Errorf.
func skemaSesudah(t *testing.T, s *skemaMenu, nama string, pernyataan []string) *skemaMenu {
	t.Helper()
	for i, p := range pernyataan {
		if err := s.terapkan(p); err != nil {
			t.Errorf("%s pernyataan %d: %v", nama, i, err)
		}
	}
	return s
}

// isiMenu membaca baris modul dan butir HASIL BERSIH langkah menu - lihat
// `langkahMenu` dan `skemaMenu`.
func isiMenu(t *testing.T, hanyaIsiAwal bool) ([]kelompokMenu, []butirMenu) {
	t.Helper()
	s := skemaSesudah(t, skemaMenuKosong(), "langkah menu", pernyataanMenu(t, hanyaIsiAwal))
	var kelompok []kelompokMenu
	var butir []butirMenu
	for _, b := range s.baris {
		if b.induk == "" {
			kelompok = append(kelompok, kelompokMenu{kode: b.kode, label: b.label, golongan: b.golongan,
				modul: b.modul, urutan: b.urutan, dimigrasi: b.dimigrasi})
			continue
		}
		butir = append(butir, butirMenu{kode: b.kode, label: b.label, induk: b.induk, urutan: b.urutan})
	}
	return kelompok, butir
}

// Isi M_NAV_MENU hanya di 900, bentuknya di 901, dan menu modul di slot menu
// pemiliknya: `TestMenuHanyaDi900DanSlotMenuModulnya` dan
// `TestSlotMenuHanyaMenyentuhMenuModulnya` (`rentang_test.go`).

// Rentang 900-949 milik `inti`, dan `inti` hanya bermigrasi di rentang itu.
func TestMigrasiIntiDiRentang900(t *testing.T) {
	milikInti := 0
	for nama := range berkasMigrasi.asal {
		diRentang := nama >= "900_" && nama < "950_"
		if intinya := berkasMigrasi.modul(nama) == "inti"; intinya != diRentang {
			t.Errorf("%s: milik inti=%v, di rentang 900-949=%v - rentang itu milik inti saja", nama, intinya, diRentang)
		}
		if diRentang {
			milikInti++
		}
	}
	if milikInti == 0 {
		t.Fatal("nol berkas migrasi inti; 900_m_nav_menu hilang atau pembacanya rusak")
	}
}

// CHECK GROUPMENU memuat PERSIS empat golongan permintaan work owner, dan
// setiap baris modul memakai salah satunya.
func TestCheckGroupMenu(t *testing.T) {
	polaCheck := regexp.MustCompile(`CONSTRAINT CK_M_NAV_MENU_GROUPMENU CHECK \(GROUPMENU IN \(([^)]*)\)\)`)
	var isi []string
	for _, p := range pernyataanMenu(t, true) {
		if nama, _ := migrasi.KolomCreateTable(p); nama != "M_NAV_MENU" {
			continue
		}
		m := polaCheck.FindStringSubmatch(p)
		if m == nil {
			t.Fatal("CREATE TABLE M_NAV_MENU tanpa CHECK GROUPMENU")
		}
		for _, v := range strings.Split(m[1], ",") {
			isi = append(isi, strings.Trim(strings.TrimSpace(v), "'"))
		}
	}
	if !reflect.DeepEqual(isi, golonganMenu) {
		t.Errorf("CHECK GROUPMENU = %v, mau %v", isi, golonganMenu)
	}
	sah := map[string]bool{}
	for _, g := range golonganMenu {
		sah[g] = true
	}
	kelompok, _ := isiMenu(t, false)
	for _, k := range kelompok {
		if !sah[k.golongan] {
			t.Errorf("baris modul %s: GROUPMENU %q di luar CHECK", k.kode, k.golongan)
		}
	}
}

// Hasil bersih 900 + 901 + slot: DUA PULUH baris, satu per folder modul
// korpus, nol butir anak. LABEL = nama folder VERBATIM - kecuali nama tampilan
// `labelTampilDisetujui` (03-10-2026) -, KODE = MODUL = nama modul backend
// (tabel nama modul: nama FOLDER tanpa spasi, huruf kecil), URUTAN = urutan di
// dalam GROUPMENU.
//
// Menggantikan `TestIsiAwalMenuDuaPuluhKelompok` (900 saja) - brief menu datar
// 30-09-2026 §4: penjaga membaca hasil bersih, bukan isi 900.
func TestMenuBersihDuaPuluhBarisSatuPerModul(t *testing.T) {
	kelompok, butir := isiMenu(t, false)
	if len(butir) != 0 {
		t.Errorf("hasil bersih masih memuat %d butir anak %v - 901 membuangnya (1 modul 1 menu)", len(butir), butir)
	}
	if len(kelompok) != 20+len(modulLuarKorpus) {
		t.Fatalf("hasil bersih memuat %d baris modul, mau %d (satu per folder modul korpus + modulLuarKorpus)",
			len(kelompok), 20+len(modulLuarKorpus))
	}
	kode := map[string]bool{}
	urutanTerakhir := map[string]int{}
	var label []string
	tampil := 0
	for _, k := range kelompok {
		if kode[k.kode] {
			t.Errorf("KODE %s ganda", k.kode)
		}
		kode[k.kode] = true
		// Nama folder baris ini: LABEL-nya, atau - untuk nama tampilan yang disetujui - folder asalnya.
		folder := k.label
		_, luar := modulLuarKorpus[k.kode]
		if lt, ada := labelTampilDisetujui[k.kode]; ada {
			tampil++
			if k.label != lt.tampil {
				t.Errorf("baris %s: LABEL %q, mau nama tampilan %q (labelTampilDisetujui)", k.kode, k.label, lt.tampil)
			}
			folder = lt.folder
		}
		if !luar {
			label = append(label, folder) // modul luar korpus tidak punya folder korpus
		}
		if mau := strings.ToLower(strings.ReplaceAll(folder, " ", "")); k.kode != mau || k.modul != mau {
			t.Errorf("baris %q: KODE %q, MODUL %q, mau keduanya %q (tabel nama modul)", folder, k.kode, k.modul, mau)
		}
		// URUTAN 1, 2, 3, ... di dalam golongannya, menurut urutan berkas.
		if k.urutan != urutanTerakhir[k.golongan]+1 {
			t.Errorf("baris %s: URUTAN %d di %s, mau %d", k.kode, k.urutan, k.golongan, urutanTerakhir[k.golongan]+1)
		}
		urutanTerakhir[k.golongan] = k.urutan
	}
	for _, g := range golonganMenu {
		if urutanTerakhir[g] == 0 {
			t.Errorf("golongan %s tanpa satu pun modul", g)
		}
	}
	if tampil != len(labelTampilDisetujui) {
		t.Errorf("%d nama tampilan terpakai, labelTampilDisetujui memuat %d - pengecualian mati dibuang", tampil, len(labelTampilDisetujui))
	}

	// Folder modul korpus: folder yang memuat Activity DAN Section. Korpus
	// dibaca, tidak pernah ditulis; mesin tanpa korpus MELEWATI bagian ini.
	const akarKorpus = `D:\XML\RNM_BRD`
	entri, err := os.ReadDir(akarKorpus)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v); kesamaan dengan folder korpus tidak terperiksa", err)
	}
	var folder []string
	for _, e := range entri {
		if !e.IsDir() {
			continue
		}
		_, errA := os.Stat(filepath.Join(akarKorpus, e.Name(), "Activity"))
		_, errS := os.Stat(filepath.Join(akarKorpus, e.Name(), "Section"))
		if errA == nil && errS == nil {
			folder = append(folder, e.Name())
		}
	}
	sort.Strings(folder)
	sort.Strings(label)
	if !reflect.DeepEqual(label, folder) {
		t.Errorf("LABEL baris modul tidak sama dengan folder modul korpus:\nlabel  %v\nfolder %v", label, folder)
	}
}

// DIMIGRASI '1' tepat untuk modul yang sudah punya
// `modul/<nama>/backend/modul.go` - bukan daftar tangan yang dapat tertinggal.
// Dibaca dari HASIL BERSIH: modul yang mendapat layar pertamanya menyalakan
// DIMIGRASI barisnya di slot menunya sendiri.
func TestMenuDimigrasiSamaDenganModulBackend(t *testing.T) {
	cocok, err := filepath.Glob(filepath.Join(akarAplikasi, "modul", "*", "backend", "modul.go"))
	if err != nil {
		t.Fatal(err)
	}
	var backend []string
	for _, c := range cocok {
		backend = append(backend, filepath.Base(filepath.Dir(filepath.Dir(c))))
	}
	sort.Strings(backend)
	if len(backend) == 0 {
		t.Fatal("nol modul backend terbaca; pembacanya yang rusak")
	}
	kelompok, _ := isiMenu(t, false)
	var dimigrasi []string
	for _, k := range kelompok {
		if k.dimigrasi == "1" {
			dimigrasi = append(dimigrasi, k.modul)
		}
	}
	sort.Strings(dimigrasi)
	if !reflect.DeepEqual(dimigrasi, backend) {
		t.Errorf("baris DIMIGRASI='1' %v, modul backend %v", dimigrasi, backend)
	}
}

// langkahInti mengambil satu langkah migrasi inti (maju atau mundur) menurut
// nama berkasnya.
func langkahInti(t *testing.T, nama string) []string {
	t.Helper()
	for _, mundur := range []bool{false, true} {
		langkah, err := migrasi.Daftar(mundur, berkasMigrasi)
		if err != nil {
			t.Fatal(err)
		}
		for _, l := range langkah {
			if l.Nama == nama {
				return l.Pernyataan
			}
		}
	}
	t.Fatalf("langkah %s tidak ada", nama)
	return nil
}

// ⛔ SKEMA TIRUAN DARI NOL: 900 lalu 901 menghasilkan 20 baris dan nol kolom
// PARENT_ID; 901 diulang (pelari yang gagal di tengah) tidak berubah; jalur
// mundur mengembalikan kolom, kunci tamu, indeks, dan lima butir PERSIS isi 900;
// mundur diulang pun aman; lalu 901 lagi kembali datar.
//
// Menggantikan `TestIsiAwalMenuLimaButir` (lima butir di bawah induknya): lima
// butir itu kini hanya hidup di jalur mundur 901, dan diuji di sini.
func TestSkemaTiruanDariNol900Lalu901(t *testing.T) {
	isiAwal := langkahInti(t, "900_m_nav_menu.sql")
	datar := langkahInti(t, "901_m_nav_menu_datar.sql")
	mundur := langkahInti(t, "901_m_nav_menu_datar_down.sql")

	s := skemaSesudah(t, skemaMenuKosong(), "900", isiAwal)
	sesudah900 := append([]barisMenu(nil), s.baris...)
	periksa := func(tahap string, baris, butir int, parent bool) {
		t.Helper()
		nButir := 0
		for _, b := range s.baris {
			if b.induk != "" {
				nButir++
			}
		}
		_, ix := s.indeks["IX_M_NAV_MENU_PARENT"]
		if len(s.baris) != baris || nButir != butir || s.kolom["PARENT_ID"] != parent || s.fk != parent || ix != parent {
			t.Errorf("%s: %d baris (%d butir), PARENT_ID %v, FK %v, indeks PARENT %v - mau %d (%d), ketiganya %v",
				tahap, len(s.baris), nButir, s.kolom["PARENT_ID"], s.fk, ix, baris, butir, parent)
		}
		if _, ada := s.indeks["IX_M_NAV_MENU_GROUPMENU"]; !ada {
			t.Errorf("%s: indeks GROUPMENU hilang", tahap)
		}
	}
	periksa("sesudah 900", 25, 5, true)
	skemaSesudah(t, s, "901", datar)
	periksa("sesudah 901", 20, 0, false)
	if len(s.kolom) != 10 {
		t.Errorf("sesudah 901: %d kolom, mau 10 (11 kolom 900 tanpa PARENT_ID)", len(s.kolom))
	}
	skemaSesudah(t, s, "901 diulang", datar)
	periksa("sesudah 901 diulang", 20, 0, false)
	skemaSesudah(t, s, "901 mundur", mundur)
	periksa("sesudah 901 mundur", 25, 5, true)
	if !reflect.DeepEqual(s.baris, sesudah900) {
		t.Errorf("jalur mundur tidak mengembalikan isi 900 persis:\ndapat %v\nmau   %v", s.baris, sesudah900)
	}
	skemaSesudah(t, s, "901 mundur diulang", mundur)
	periksa("sesudah 901 mundur diulang", 25, 5, true)
	skemaSesudah(t, s, "901 lagi", datar)
	periksa("sesudah 901 lagi", 20, 0, false)
}

// Setiap langkah 901 - maju dan mundur - terlindung pemeriksaan katalog atau
// idempoten lewat NOT EXISTS; INSERT jalur mundur adalah salinan HARFIAH INSERT
// butir 900 (label VERBATIM), dan objek yang dikembalikannya berdefinisi sama
// dengan 900.
func TestSetiapLangkah901TerlindungKatalog(t *testing.T) {
	var butir900 []string
	for _, p := range langkahInti(t, "900_m_nav_menu.sql") {
		if polaIsiButir.MatchString(p) {
			butir900 = append(butir900, p)
		}
	}
	if len(butir900) != 5 {
		t.Fatalf("900 memuat %d INSERT butir, mau 5", len(butir900))
	}
	var butirMundur []string
	for _, berkas := range []string{"901_m_nav_menu_datar.sql", "901_m_nav_menu_datar_down.sql"} {
		for i, p := range langkahInti(t, berkas) {
			if pk, ok := migrasi.BacaPerintahKatalog(p); ok {
				if pk.Tabel != "M_NAV_MENU" || !strings.Contains(pk.Perintah, pk.Objek) {
					t.Errorf("%s pernyataan %d: pemeriksaan %s.%s tidak sepadan dengan perintahnya %q",
						berkas, i, pk.Tabel, pk.Objek, pk.Perintah)
				}
				membuang := strings.HasPrefix(pk.Perintah, "DELETE") || strings.Contains(pk.Perintah, " DROP ") ||
					strings.HasPrefix(pk.Perintah, "DROP ")
				if membuang != pk.BilaAda {
					t.Errorf("%s pernyataan %d: %q dijalankan bila objek ada=%v", berkas, i, pk.Perintah, pk.BilaAda)
				}
				continue
			}
			if berkas == "901_m_nav_menu_datar_down.sql" && polaIsiButir.MatchString(p) {
				butirMundur = append(butirMundur, p)
				continue
			}
			t.Errorf("%s pernyataan %d tidak terlindung katalog dan tidak idempoten: %s",
				berkas, i, migrasi.RingkasPernyataan(p))
		}
	}
	if !reflect.DeepEqual(butirMundur, butir900) {
		t.Errorf("INSERT butir jalur mundur 901 bukan salinan harfiah 900:\ndapat %v\nmau   %v", butirMundur, butir900)
	}
	// Definisi objek yang dikembalikan = definisi 900.
	sql900 := strings.Join(langkahInti(t, "900_m_nav_menu.sql"), "\n")
	mundur := strings.Join(langkahInti(t, "901_m_nav_menu_datar_down.sql"), "\n")
	for _, pasangan := range [][2]string{
		{"PARENT_ID    NUMBER(10),", "ADD (PARENT_ID NUMBER(10))"},
		{"CONSTRAINT FK_M_NAV_MENU_INDUK FOREIGN KEY (PARENT_ID) REFERENCES {skema}.M_NAV_MENU (ID)",
			"ADD CONSTRAINT FK_M_NAV_MENU_INDUK FOREIGN KEY (PARENT_ID) REFERENCES {skema}.M_NAV_MENU (ID)"},
		{"CREATE INDEX {skema}.IX_M_NAV_MENU_PARENT ON {skema}.M_NAV_MENU (PARENT_ID)",
			"CREATE INDEX {skema}.IX_M_NAV_MENU_PARENT ON {skema}.M_NAV_MENU (PARENT_ID)"},
	} {
		if !strings.Contains(sql900, pasangan[0]) || !strings.Contains(mundur, pasangan[1]) {
			t.Errorf("definisi 900 %q tidak dikembalikan jalur mundur sebagai %q", pasangan[0], pasangan[1])
		}
	}
}

// Skema tiruan MENGGIGIT: langkah yang tidak terlindung katalog lalu diulang
// mati seperti di Oracle - bukan lulus diam-diam.
func TestSkemaTiruanMenangkapLangkahTanpaPelindung(t *testing.T) {
	for nama, urutan := range map[string][]string{
		"DROP COLUMN diulang": {"ALTER TABLE {skema}.M_NAV_MENU DROP COLUMN PARENT_ID", "ALTER TABLE {skema}.M_NAV_MENU DROP COLUMN PARENT_ID"},
		"DELETE sesudah kolom dibuang": {"ALTER TABLE {skema}.M_NAV_MENU DROP COLUMN PARENT_ID",
			"DELETE FROM {skema}.M_NAV_MENU WHERE PARENT_ID IS NOT NULL"},
		"DROP CONSTRAINT diulang": {"ALTER TABLE {skema}.M_NAV_MENU DROP CONSTRAINT FK_M_NAV_MENU_INDUK",
			"ALTER TABLE {skema}.M_NAV_MENU DROP CONSTRAINT FK_M_NAV_MENU_INDUK"},
	} {
		s := skemaSesudah(t, skemaMenuKosong(), "900", langkahInti(t, "900_m_nav_menu.sql"))
		var galat error
		for _, q := range urutan {
			if err := s.terapkanPerintah(q); err != nil {
				galat = err
			}
		}
		if galat == nil {
			t.Errorf("%s: skema tiruan tidak menolak langkah yang diulang tanpa pelindung", nama)
		}
	}
	// Bentuk UPDATE lama (menyebut PARENT_ID) SAH sebelum 901 dan mati di
	// ORA-00904 sesudahnya - sebabnya kolom yang hilang, bukan bentuknya.
	lama := "UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE\nWHERE KODE = 'nbfacin' AND PARENT_ID IS NULL"
	sebelum := skemaSesudah(t, skemaMenuKosong(), "900", langkahInti(t, "900_m_nav_menu.sql"))
	if err := sebelum.terapkan(lama); err != nil {
		t.Errorf("UPDATE bentuk lama ditolak SEBELUM 901: %v", err)
	}
	s := skemaSesudah(t, skemaMenuKosong(), "900+901",
		append(langkahInti(t, "900_m_nav_menu.sql"), langkahInti(t, "901_m_nav_menu_datar.sql")...))
	if err := s.terapkan(lama); err == nil || !strings.Contains(err.Error(), "ORA-00904") {
		t.Errorf("UPDATE ... AND PARENT_ID IS NULL sesudah 901: %v, mau ORA-00904", err)
	}
	// Blok atas constraint yang ADA (UQ 900) dengan perintah tak dikenal:
	// ditolak walau pemeriksaannya "n = 0" (dilewati) - di Oracle ia dapat
	// berjalan.
	blok := "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS\n" +
		"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_NAV_MENU' AND CONSTRAINT_NAME = 'UQ_M_NAV_MENU_KODE';\n" +
		"  IF n = 0 THEN\n    EXECUTE IMMEDIATE 'DROP TABLE {skema}.M_NAV_MENU';\n  END IF;\nEND;"
	if err := s.terapkan(blok); err == nil {
		t.Error("blok berpelindung dengan perintah tak dikenal diterima karena pemeriksaannya melewatinya")
	}
	if !s.konstrain["UQ_M_NAV_MENU_KODE"] || !s.konstrain["PK_M_NAV_MENU"] || !s.konstrain["CK_M_NAV_MENU_GROUPMENU"] {
		t.Errorf("katalog tiruan tidak memuat constraint 900: %v", s.konstrain)
	}
	if err := s.terapkan("UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE\nWHERE KODE = 'nbfacin'"); err != nil {
		t.Errorf("slot menu bentuk datar ditolak sesudah 901: %v", err)
	}
	if i := s.cari("nbfacin"); i < 0 || s.baris[i].dimigrasi != "1" {
		t.Error("slot menu bentuk datar tidak menyalakan DIMIGRASI nbfacin")
	}
	if err := s.terapkan("UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE\nWHERE KODE = 'tidakada'"); err == nil {
		t.Error("UPDATE atas modul yang tidak ada diterima")
	}
}

// ⛔ Isi awal 903 (M_LOGIN_GO_MENU, Kelola User 01-10-2026): setiap akun yang
// sudah ada mendapat SEMUA menu. "Semua" = KODE setiap baris modul HASIL
// BERSIH menu (900 + 901 + slot) ditambah setiap menu aplikasi
// (`menu.MenuAplikasi`) - persis, tanpa ganda. Daftarnya tertulis harfiah di
// 903 karena langkah inti selain 900/901 tidak boleh menyebut tabel menu
// (`TestMenuHanyaDi900DanSlotMenuModulnya`); penjaga inilah yang menahannya
// tetap sama.
func TestIsiAwalMenuAkunMemuatSemuaMenu(t *testing.T) {
	var isi string
	for _, p := range langkahInti(t, "903_m_login_go_menu.sql") {
		if strings.HasPrefix(p, "INSERT INTO {skema}.M_LOGIN_GO_MENU ") {
			isi = p
		}
	}
	if isi == "" {
		t.Fatal("903 tanpa INSERT isi awal M_LOGIN_GO_MENU")
	}
	var dapat []string
	for _, m := range regexp.MustCompile(`SELECT '([^']+)'(?: AS KODE)? FROM DUAL`).FindAllStringSubmatch(isi, -1) {
		dapat = append(dapat, m[1])
	}
	kelompok, _ := isiMenu(t, false)
	var mau []string
	for _, k := range kelompok {
		if _, luar := modulLuarKorpus[k.kode]; !luar { // lahir sesudah 903 - hak menu lewat Kelola User
			mau = append(mau, k.kode)
		}
	}
	for _, a := range menu.MenuAplikasi {
		mau = append(mau, a.Kode)
	}
	sort.Strings(dapat)
	sort.Strings(mau)
	if len(mau) != 21 || !reflect.DeepEqual(dapat, mau) {
		t.Errorf("isi awal 903 menyebut %v,\nmau %v (dua puluh modul + menu aplikasi)", dapat, mau)
	}
}
