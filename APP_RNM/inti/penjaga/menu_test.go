package penjaga

// Penjaga migrasi 900 - M_NAV_MENU dan isi awalnya - TANPA Oracle.
//
// Untuk apa berkas ini: menu aplikasi kini dirakit dari tabel
// (`PROMPT-MENU-DARI-TABEL-M_NAV_MENU.md`). Isi awalnya bukan data uji - ia
// menu yang akan dilihat setiap pemakai, dan kelak dasar akses per akun. Satu
// kelompok yang hilang atau satu butir yang salah induk tidak menggagalkan
// satu pun uji lain; berkas inilah yang menangkapnya.

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/migrasi"
)

// golonganMenu - isi CHECK GROUPMENU, urutan tampil di sidebar (permintaan
// work owner 30-09-2026).
var golonganMenu = []string{"TREATY", "FACULTATIVE", "KLAIM", "MASTER"}

// kelompokMenu adalah satu baris kelompok modul di isi awal 900.
type kelompokMenu struct {
	kode, label, golongan, modul string
	urutan                       int
	dimigrasi                    string
}

// butirMenu adalah satu baris butir menu di isi awal 900.
type butirMenu struct {
	kode, label, induk string
	urutan             int
}

var (
	polaIsiKelompok = regexp.MustCompile(`(?s)^INSERT INTO \{skema\}\.M_NAV_MENU \(ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI\)\s+` +
		`SELECT \{skema\}\.SEQ_M_NAV_MENU\.NEXTVAL, NULL, '([^']+)', '([^']+)', '([^']+)', '([^']+)', (\d+), '([01])' FROM DUAL\s+` +
		`WHERE NOT EXISTS \(SELECT 1 FROM \{skema\}\.M_NAV_MENU WHERE KODE = '([^']+)'\)$`)
	polaIsiButir = regexp.MustCompile(`(?s)^INSERT INTO \{skema\}\.M_NAV_MENU \(ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI\)\s+` +
		`SELECT \{skema\}\.SEQ_M_NAV_MENU\.NEXTVAL, k\.ID, '([^']+)', '([^']+)', k\.GROUPMENU, k\.MODUL, (\d+), k\.DIMIGRASI\s+` +
		`FROM \{skema\}\.M_NAV_MENU k\s+WHERE k\.KODE = '([^']+)' AND k\.PARENT_ID IS NULL\s+` +
		`AND NOT EXISTS \(SELECT 1 FROM \{skema\}\.M_NAV_MENU b WHERE b\.KODE = '([^']+)'\)$`)
)

// pernyataanMenu - pernyataan maju langkah 900, apa adanya.
func pernyataanMenu(t *testing.T) []string {
	t.Helper()
	langkah, err := migrasi.Daftar(false, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range langkah {
		if m.Nama == "900_m_nav_menu.sql" {
			return m.Pernyataan
		}
	}
	t.Fatal("langkah 900_m_nav_menu.sql tidak ada di daftar migrasi")
	return nil
}

// isiAwalMenu membaca kelompok dan butir dari INSERT langkah 900.
//
// ⛔ Setiap INSERT WAJIB cocok dengan salah satu bentuk: INSERT berbentuk lain
// adalah baris yang tidak terbaca penjaga ini - dan tidak terbaca pula oleh
// uji dua arah frontend yang membaca berkas yang sama.
func isiAwalMenu(t *testing.T) ([]kelompokMenu, []butirMenu) {
	t.Helper()
	var kelompok []kelompokMenu
	var butir []butirMenu
	for i, p := range pernyataanMenu(t) {
		if !strings.HasPrefix(p, "INSERT") {
			continue
		}
		if m := polaIsiKelompok.FindStringSubmatch(p); m != nil {
			// Idempoten: yang diperiksa NOT EXISTS adalah KODE yang dimasukkan.
			if m[7] != m[1] {
				t.Errorf("INSERT kelompok %s memeriksa NOT EXISTS atas %s", m[1], m[7])
			}
			u, _ := strconv.Atoi(m[5])
			kelompok = append(kelompok, kelompokMenu{kode: m[1], label: m[2], golongan: m[3], modul: m[4], urutan: u, dimigrasi: m[6]})
			continue
		}
		if m := polaIsiButir.FindStringSubmatch(p); m != nil {
			if m[5] != m[1] {
				t.Errorf("INSERT butir %s memeriksa NOT EXISTS atas %s", m[1], m[5])
			}
			u, _ := strconv.Atoi(m[3])
			butir = append(butir, butirMenu{kode: m[1], label: m[2], induk: m[4], urutan: u})
			continue
		}
		t.Errorf("pernyataan %d langkah 900 adalah INSERT yang bentuknya tidak dikenal penjaga: %s",
			i, migrasi.RingkasPernyataan(p))
	}
	return kelompok, butir
}

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
// setiap kelompok isi awal memakai salah satunya.
func TestCheckGroupMenu(t *testing.T) {
	polaCheck := regexp.MustCompile(`CONSTRAINT CK_M_NAV_MENU_GROUPMENU CHECK \(GROUPMENU IN \(([^)]*)\)\)`)
	var isi []string
	for _, p := range pernyataanMenu(t) {
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
	kelompok, _ := isiAwalMenu(t)
	for _, k := range kelompok {
		if !sah[k.golongan] {
			t.Errorf("kelompok %s: GROUPMENU %q di luar CHECK", k.kode, k.golongan)
		}
	}
}

// Dua puluh kelompok - satu per folder modul korpus, LABEL = nama folder
// VERBATIM, KODE = MODUL = nama modul backend (tabel nama modul: nama folder
// tanpa spasi, huruf kecil).
func TestIsiAwalMenuDuaPuluhKelompok(t *testing.T) {
	kelompok, _ := isiAwalMenu(t)
	if len(kelompok) != 20 {
		t.Fatalf("isi awal memuat %d kelompok, mau 20 (satu per folder modul korpus)", len(kelompok))
	}
	kode := map[string]bool{}
	urutanTerakhir := map[string]int{}
	var label []string
	for _, k := range kelompok {
		if kode[k.kode] {
			t.Errorf("KODE kelompok %s ganda", k.kode)
		}
		kode[k.kode] = true
		label = append(label, k.label)
		if mau := strings.ToLower(strings.ReplaceAll(k.label, " ", "")); k.kode != mau || k.modul != mau {
			t.Errorf("kelompok %q: KODE %q, MODUL %q, mau keduanya %q (tabel nama modul)", k.label, k.kode, k.modul, mau)
		}
		// URUTAN 1, 2, 3, ... di dalam golongannya, menurut urutan berkas.
		if k.urutan != urutanTerakhir[k.golongan]+1 {
			t.Errorf("kelompok %s: URUTAN %d di %s, mau %d", k.kode, k.urutan, k.golongan, urutanTerakhir[k.golongan]+1)
		}
		urutanTerakhir[k.golongan] = k.urutan
	}
	for _, g := range golonganMenu {
		if urutanTerakhir[g] == 0 {
			t.Errorf("golongan %s tanpa satu pun kelompok", g)
		}
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
		t.Errorf("LABEL kelompok tidak sama dengan folder modul korpus:\nlabel  %v\nfolder %v", label, folder)
	}
}

// DIMIGRASI '1' tepat untuk modul yang sudah punya folder backend
// `modul/<nama>/modul.go` - bukan daftar tangan yang dapat tertinggal.
func TestIsiAwalMenuDimigrasiSamaDenganModulBackend(t *testing.T) {
	cocok, err := filepath.Glob(filepath.Join(akarAplikasi, "modul", "*", "modul.go"))
	if err != nil {
		t.Fatal(err)
	}
	var backend []string
	for _, c := range cocok {
		backend = append(backend, filepath.Base(filepath.Dir(c)))
	}
	sort.Strings(backend)
	if len(backend) == 0 {
		t.Fatal("nol modul backend terbaca; pembacanya yang rusak")
	}
	kelompok, _ := isiAwalMenu(t)
	var dimigrasi []string
	for _, k := range kelompok {
		if k.dimigrasi == "1" {
			dimigrasi = append(dimigrasi, k.modul)
		}
	}
	sort.Strings(dimigrasi)
	if !reflect.DeepEqual(dimigrasi, backend) {
		t.Errorf("kelompok DIMIGRASI='1' %v, modul backend %v", dimigrasi, backend)
	}
}

// Lima butir, masing-masing di bawah kelompok yang ADA dan sudah dimigrasi.
// GROUPMENU, MODUL, dan DIMIGRASI butir dibaca dari induknya (bentuk
// `polaIsiButir`), jadi keduanya tidak dapat berbeda.
func TestIsiAwalMenuLimaButir(t *testing.T) {
	kelompok, butir := isiAwalMenu(t)
	induk := map[string]kelompokMenu{}
	for _, k := range kelompok {
		induk[k.kode] = k
	}
	mau := map[string]string{
		"inbox": "claimlife", "register": "claimlife", "premiumlist": "premiumlistlife",
		"komite": "komiteclaimlife", "tco-tahun": "treatycontractout",
	}
	dapat := map[string]string{}
	for _, b := range butir {
		dapat[b.kode] = b.induk
		k, ada := induk[b.induk]
		switch {
		case !ada:
			t.Errorf("butir %s: induk %s bukan kelompok isi awal", b.kode, b.induk)
		case k.dimigrasi != "1":
			t.Errorf("butir %s di bawah %s yang belum dimigrasi", b.kode, b.induk)
		}
		if _, bentrok := induk[b.kode]; bentrok {
			t.Errorf("KODE butir %s sama dengan KODE kelompok (UNIQUE)", b.kode)
		}
		if strings.TrimSpace(b.label) == "" || b.urutan < 1 {
			t.Errorf("butir %s: LABEL %q, URUTAN %d", b.kode, b.label, b.urutan)
		}
	}
	if !reflect.DeepEqual(dapat, mau) {
		t.Errorf("butir isi awal %v, mau %v", dapat, mau)
	}
}
