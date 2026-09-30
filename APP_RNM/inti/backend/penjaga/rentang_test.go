package penjaga

// Penjaga rentang migrasi (R2) dan slot menu (R3) - struktur tim satu folder
// per modul (30-09-2026). TANPA Oracle, TANPA satu pun nama modul.
//
// Untuk apa berkas ini: pelari migrasi mengurutkan nama berkas SEBAGAI TEKS di
// semua sumber sekaligus (`inti/backend/migrasi` `Daftar`). Dua pengembang yang
// sama-sama mengambil "nomor berikutnya" akan bertabrakan, dan nomor empat digit
// salah urut (`1000_` sebelum `101_`). Karena itu setiap modul MENYATAKAN
// rentang migrasi dan slot menunya di `MODUL.md` (`Rentang migrasi`,
// `Slot menu`), dan penjaga ini menegakkannya:
//
//   - nomor SELALU tiga digit;
//   - setiap berkas migrasi modul berada di rentang atau slot menu modulnya;
//     berkas `inti` di 900-949;
//   - dua `MODUL.md` tidak berbagi nomor, dan tidak ada yang memakai 900-949;
//   - baris `M_NAV_MENU` hanya di 900 (isi awal milik `inti`) dan di slot menu
//     modul pemiliknya - dan slot itu hanya menyentuh menu modul itu sendiri.
//     `inti` boleh mengubah BENTUK tabel itu (901, menu datar 30-09-2026),
//     tetapi tidak menambah baris di luar 900.

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/migrasi"
)

// Rentang tetap milik aplikasi, bukan milik modul mana pun.
const (
	awalInti, akhirInti = 900, 949 // tabel lintas modul milik `inti`
	awalSlot, akhirSlot = 950, 999 // slot menu - dibagi per modul lewat MODUL.md
)

var polaNomorBerkas = regexp.MustCompile(`^(\d{3})_[a-z0-9_]+\.sql$`)

// nomorBerkas membaca nomor tiga digit di depan nama berkas migrasi.
func nomorBerkas(nama string) (int, bool) {
	m := polaNomorBerkas.FindStringSubmatch(nama)
	if m == nil {
		return 0, false
	}
	n, _ := strconv.Atoi(m[1])
	return n, true
}

// jatahModul - rentang migrasi dan slot menu setiap modul, dari MODUL.md.
type jatahModul struct {
	migrasi, slot [2]int
}

func (j jatahModul) diMigrasi(n int) bool { return n >= j.migrasi[0] && n <= j.migrasi[1] }
func (j jatahModul) diSlot(n int) bool    { return n >= j.slot[0] && n <= j.slot[1] }

// jatahSetiapModul membaca jatah setiap modul dan menolak bentuk yang salah.
func jatahSetiapModul(t *testing.T) map[string]jatahModul {
	t.Helper()
	hasil := map[string]jatahModul{}
	for _, m := range bacaModulMD(t) {
		if nama := strings.Trim(m.kunci["Nama modul"], "`"); nama != m.folder {
			t.Errorf("%s: `Nama modul` %q, mau nama foldernya %q", m.jalur, nama, m.folder)
		}
		mig, slot := rentangModul(t, m)
		hasil[m.folder] = jatahModul{migrasi: mig, slot: slot}
	}
	return hasil
}

// ⛔ R2: jatah setiap modul berbentuk benar dan tidak berbagi nomor.
func TestRentangMigrasiDanSlotMenuTidakBerbagiNomor(t *testing.T) {
	jatah := jatahSetiapModul(t)
	type rentang struct {
		pemilik string
		awal    int
		akhir   int
	}
	semua := []rentang{{"inti (900-949)", awalInti, akhirInti}}
	for m, j := range jatah {
		if j.migrasi[0] < 1 || j.migrasi[1] >= awalInti {
			t.Errorf("modul %s: rentang migrasi %03d-%03d di luar 001-899 (900-999 milik inti dan slot menu)",
				m, j.migrasi[0], j.migrasi[1])
		}
		if j.slot[0] < awalSlot || j.slot[1] > akhirSlot {
			t.Errorf("modul %s: slot menu %03d-%03d di luar %d-%d", m, j.slot[0], j.slot[1], awalSlot, akhirSlot)
		}
		semua = append(semua, rentang{"modul " + m + " (rentang migrasi)", j.migrasi[0], j.migrasi[1]},
			rentang{"modul " + m + " (slot menu)", j.slot[0], j.slot[1]})
	}
	sort.Slice(semua, func(i, k int) bool { return semua[i].awal < semua[k].awal })
	for i := 1; i < len(semua); i++ {
		if a, b := semua[i-1], semua[i]; b.awal <= a.akhir {
			t.Errorf("%s %03d-%03d dan %s %03d-%03d berbagi nomor", a.pemilik, a.awal, a.akhir, b.pemilik, b.awal, b.akhir)
		}
	}
	t.Logf("jatah %d modul terbaca", len(jatah))
}

// ⛔ R2: setiap berkas migrasi berada di rentang atau slot menu modulnya, dan
// nomornya tiga digit.
func TestSetiapMigrasiDiRentangAtauSlotModulnya(t *testing.T) {
	jatah := jatahSetiapModul(t)
	diperiksa := 0
	for nama, jalur := range berkasMigrasi.asal {
		diperiksa++
		n, ok := nomorBerkas(nama)
		if !ok {
			t.Errorf("%s: nama berkas migrasi harus NNN_nama.sql - nomor tiga digit, huruf kecil", filepath.ToSlash(jalur))
			continue
		}
		pemilik := pemilikJalur(jalur)
		if pemilik == "inti" {
			if n < awalInti || n > akhirInti {
				t.Errorf("%s: migrasi inti di luar %d-%d", nama, awalInti, akhirInti)
			}
			continue
		}
		j, ada := jatah[pemilik]
		switch {
		case !ada:
			t.Errorf("%s: modul %s tanpa rentang migrasi di MODUL.md", nama, pemilik)
		case !j.diMigrasi(n) && !j.diSlot(n):
			t.Errorf("%s: nomor %03d di luar rentang migrasi %03d-%03d dan slot menu %03d-%03d modul %s",
				nama, n, j.migrasi[0], j.migrasi[1], j.slot[0], j.slot[1], pemilik)
		}
	}
	if diperiksa < 20 {
		t.Fatalf("hanya %d berkas migrasi terbaca; pembacanya yang rusak", diperiksa)
	}
}

// ⛔ R3: baris M_NAV_MENU hanya di 900 (isi awal, milik inti) dan di slot menu
// modul pemiliknya. 901-949 tetap milik inti - tetapi bukan untuk menu.
//
// Menggantikan `TestMenuHanyaDiMigrasiInti` (menu dulu hanya boleh di
// `inti/migrations` 900-949): folder bersama itu membuat dua pengembang yang
// sama-sama menambah `901_…` bertabrakan, dan membuat pengembang modul
// menyentuh folder inti.
func TestMenuHanyaDi900DanSlotMenuModulnya(t *testing.T) {
	jatah := jatahSetiapModul(t)
	menyentuh := 0
	for nama, isi := range seluruhSQL(t, false) {
		if !strings.Contains(strings.ToUpper(isi), "M_NAV_MENU") {
			continue
		}
		menyentuh++
		pemilik := berkasMigrasi.modul(nama)
		n, _ := nomorBerkas(nama)
		switch {
		case pemilik == "inti":
			// 901 (menu datar) mengubah BENTUK tabel itu - membuang butir anak
			// dan PARENT_ID - tanpa menambah baris. INSERT hanya di 900.
			if !strings.HasPrefix(nama, "900_") && strings.Contains(strings.ToUpper(isi), "INSERT INTO {SKEMA}.M_NAV_MENU") {
				t.Errorf("%s (inti) menambah baris M_NAV_MENU - isi menu inti hanya di 900; menu modul di slot menunya", nama)
			}
		case !jatah[pemilik].diSlot(n):
			t.Errorf("%s (modul %s) menyentuh M_NAV_MENU di luar slot menunya %03d-%03d",
				nama, pemilik, jatah[pemilik].slot[0], jatah[pemilik].slot[1])
		}
	}
	if menyentuh == 0 {
		t.Fatal("nol migrasi menyentuh M_NAV_MENU; 900_m_nav_menu hilang atau pembacanya rusak")
	}
}

// ⛔ R3: berkas di slot menu sebuah modul HANYA menu modul itu - butir di bawah
// kelompoknya sendiri dan penanda DIMIGRASI-nya - dalam bentuk yang dibaca
// penjaga isi menu (bentuk bab 6 PANDUAN-DEPLOY). Tidak ada tabel, tidak ada
// kelompok baru, tidak ada menu modul lain.
func TestSlotMenuHanyaMenyentuhMenuModulnya(t *testing.T) {
	jatah := jatahSetiapModul(t)
	maju, mundur := pernyataanPerLangkah(t, false), pernyataanPerLangkah(t, true)
	slot := 0
	for nama := range maju {
		pemilik := berkasMigrasi.modul(nama)
		n, _ := nomorBerkas(nama)
		if pemilik == "inti" || !jatah[pemilik].diSlot(n) {
			continue
		}
		slot++
		turun := strings.TrimSuffix(nama, ".sql") + "_down.sql"
		for _, alasan := range pelanggaranSlotMenu(pemilik, maju[nama], mundur[turun]) {
			t.Errorf("%s: %s", nama, alasan)
		}
	}
	t.Logf("%d berkas slot menu diperiksa", slot)
}

// pernyataanPerLangkah - pernyataan setiap langkah migrasi (maju atau
// mundur) semua modul, per nama berkas, apa adanya dari pelari.
func pernyataanPerLangkah(t *testing.T, mundur bool) map[string][]string {
	t.Helper()
	langkah, err := migrasi.Daftar(mundur, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string][]string{}
	for _, l := range langkah {
		out[l.Nama] = l.Pernyataan
	}
	return out
}

var polaHapusButir = regexp.MustCompile(`(?s)^DELETE FROM \{skema\}\.M_NAV_MENU WHERE KODE = '([^']+)'$`)

// pelanggaranSlotMenu menjawab mengapa pernyataan maju/mundur sebuah berkas
// slot menu modul `modul` tidak sah; kosong = sah.
func pelanggaranSlotMenu(modul string, maju, mundur []string) []string {
	var alasan []string
	butir := map[string]bool{}
	if len(maju) == 0 {
		alasan = append(alasan, "berkas slot menu tanpa satu pun pernyataan")
	}
	for _, p := range maju {
		switch m := polaUbahDimigrasi.FindStringSubmatch(p); {
		case m != nil:
			if m[2] != modul {
				alasan = append(alasan, fmt.Sprintf("mengubah DIMIGRASI kelompok %s - slot ini milik %s", m[2], modul))
			}
			continue
		}
		if m := polaIsiButir.FindStringSubmatch(p); m != nil {
			if m[4] != modul {
				alasan = append(alasan, fmt.Sprintf("menyisipkan butir %s di bawah kelompok %s - slot ini milik %s", m[1], m[4], modul))
			}
			butir[m[1]] = true
			continue
		}
		alasan = append(alasan, "pernyataan di luar bentuk menu modul (UPDATE DIMIGRASI kelompoknya, "+
			"INSERT butir di bawahnya): "+ringkas(p))
	}
	for _, p := range mundur {
		if m := polaUbahDimigrasi.FindStringSubmatch(p); m != nil {
			if m[2] != modul {
				alasan = append(alasan, fmt.Sprintf("jalur mundur mengubah DIMIGRASI kelompok %s - slot ini milik %s", m[2], modul))
			}
			continue
		}
		if m := polaHapusButir.FindStringSubmatch(p); m != nil {
			if !butir[m[1]] {
				alasan = append(alasan, fmt.Sprintf("jalur mundur menghapus %s yang tidak disisipkan langkah majunya", m[1]))
			}
			continue
		}
		alasan = append(alasan, "pernyataan mundur di luar bentuk menu modul: "+ringkas(p))
	}
	return alasan
}

func ringkas(p string) string {
	if len(p) > 90 {
		return p[:90] + "…"
	}
	return p
}

// Aturan slot menu diuji dua arah, supaya penjaga di atas tidak lulus karena
// aturannya longgar - hari ini belum ada satu pun berkas slot menu.
func TestAturanSlotMenuMenggigit(t *testing.T) {
	ubah := func(kode, dim string) string {
		return "UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '" + dim + "', TGL_UBAH = SYSDATE\nWHERE KODE = '" + kode + "' AND PARENT_ID IS NULL"
	}
	sisip := func(butir, induk string) string {
		return "INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)\n" +
			"SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, k.ID, '" + butir + "', 'Label', k.GROUPMENU, k.MODUL, 1, k.DIMIGRASI\n" +
			"FROM {skema}.M_NAV_MENU k\nWHERE k.KODE = '" + induk + "' AND k.PARENT_ID IS NULL\n" +
			"AND NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU b WHERE b.KODE = '" + butir + "')"
	}
	hapus := func(butir string) string { return "DELETE FROM {skema}.M_NAV_MENU WHERE KODE = '" + butir + "'" }
	for _, k := range []struct {
		nama         string
		maju, mundur []string
		sah          bool
	}{
		{"bentuk bab 6, modulnya sendiri", []string{ubah("alfa", "1"), sisip("alfa-inbox", "alfa")},
			[]string{hapus("alfa-inbox"), ubah("alfa", "0")}, true},
		{"butir di bawah modul lain", []string{sisip("beta-inbox", "beta")}, []string{hapus("beta-inbox")}, false},
		{"DIMIGRASI modul lain", []string{ubah("beta", "1")}, []string{ubah("beta", "0")}, false},
		{"membuat tabel di slot", []string{"CREATE TABLE {skema}.T_ALFA (ID NUMBER(10))"}, nil, false},
		{"mundur menghapus butir lain", []string{sisip("alfa-inbox", "alfa")}, []string{hapus("alfa-lain")}, false},
		{"slot kosong", nil, nil, false},
	} {
		if dapat := len(pelanggaranSlotMenu("alfa", k.maju, k.mundur)) == 0; dapat != k.sah {
			t.Errorf("%s: sah=%v, mau %v (%v)", k.nama, dapat, k.sah, pelanggaranSlotMenu("alfa", k.maju, k.mundur))
		}
	}
}

// Aturan rentang, dua arah: bentuk nomor dan bentuk rentang di MODUL.md.
func TestAturanRentangMenggigit(t *testing.T) {
	for nama, sah := range map[string]bool{
		"001_t_work_claim.sql": true, "900_m_nav_menu.sql": true, "962_menu_nbfacin_down.sql": true,
		"1000_x.sql": false, "01_x.sql": false, "001-x.sql": false, "001_X.sql": false,
	} {
		if _, ok := nomorBerkas(nama); ok != sah {
			t.Errorf("nomor berkas %s: sah=%v, mau %v", nama, ok, sah)
		}
	}
	for nilai, sah := range map[string]bool{
		"`001-029`": true, "050-099": true, "`30-49`": false, "`0300-0319`": false, "`319-300`": false, "001–029": false,
	} {
		if _, _, err := rentangNomor(nilai); (err == nil) != sah {
			t.Errorf("rentang %s: sah=%v, mau %v (%v)", nilai, err == nil, sah, err)
		}
	}
}

// ⛔ R3: di skema uji dari NOL, 900 (CREATE TABLE M_NAV_MENU + isi awal)
// berjalan SEBELUM slot menu 95x mana pun - pelari mengurutkan nama berkas
// sebagai teks, dan tiga digit menjamin urutan itu. Diperiksa lewat pelari
// yang sama dengan `-migrate`, atas migrasi inti sungguhan dan satu modul
// tiruan yang punya berkas slot menu, bukan dengan membandingkan teks sendiri.
func TestSlotMenuBerjalanSesudah900(t *testing.T) {
	for m, j := range jatahSetiapModul(t) {
		if j.slot[0] <= 900 {
			t.Errorf("modul %s: slot menu %03d tidak sesudah 900", m, j.slot[0])
		}
	}
	tiruan := fstest.MapFS{
		"migrations/952_menu_tiruan.sql":      {Data: []byte("UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE\nWHERE KODE = 'tiruan' AND PARENT_ID IS NULL\n/\n")},
		"migrations/952_menu_tiruan_down.sql": {Data: []byte("UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '0', TGL_UBAH = SYSDATE\nWHERE KODE = 'tiruan' AND PARENT_ID IS NULL\n/\n")},
		"migrations/030_tiruan.sql":           {Data: []byte("CREATE TABLE {skema}.T_TIRUAN (ID NUMBER(10))\n/\n")},
		"migrations/030_tiruan_down.sql":      {Data: []byte("DROP TABLE {skema}.T_TIRUAN\n/\n")},
	}
	langkah, err := migrasi.Daftar(false, inti.SumberMigrasi(), tiruan)
	if err != nil {
		t.Fatal(err)
	}
	var urut []string
	for _, l := range langkah {
		urut = append(urut, l.Nama)
	}
	// 901 (menu datar, milik inti) berjalan sesudah 900 dan SEBELUM slot mana
	// pun - slot menu karena itu melihat tabel yang sudah datar.
	if mau := []string{"030_tiruan.sql", "900_m_nav_menu.sql", "901_m_nav_menu_datar.sql", "952_menu_tiruan.sql"}; strings.Join(urut, ",") != strings.Join(mau, ",") {
		t.Errorf("urutan pelari %v, mau %v", urut, mau)
	}
}
