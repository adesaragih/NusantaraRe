package repository

// Migrasi modul coverlife 085-086 + slot menu 957 (keputusan work owner 08-10-2026 K0, C1, C4) dibaca pengurai PRODUKSI
// pelari (`migrasi.Daftar`, `BacaPerintahKatalog`, `KolomAlterTambah`, `KolomAlterBuang`, `KolomCreateTable`) dan
// berpasangan dengan jalur mundurnya. Nol CREATE TABLE, nol RENAME (C1: nama M_COVER_LIFE TETAP). Urutan pelari 085 ->
// 900 -> 957 dibuktikan atas migrasi inti SUNGGUHAN (K0).

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/migrasi"
)

// berkasModul - seluruh folder `migrations/` modul ini sebagai satu sumber pelari (bentuk yang sama dengan embed).
func berkasModul(t *testing.T) fstest.MapFS {
	t.Helper()
	entri, err := os.ReadDir(filepath.Join("..", "migrations"))
	if err != nil {
		t.Fatal(err)
	}
	sumber := fstest.MapFS{}
	for _, e := range entri {
		isi, err := os.ReadFile(filepath.Join("..", "migrations", e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(isi), "\r") {
			t.Errorf("%s memuat CR", e.Name())
		}
		sumber["migrations/"+e.Name()] = &fstest.MapFile{Data: isi}
	}
	return sumber
}

// langkah - pernyataan migrasi `kunci` (maju atau mundur).
func langkah(t *testing.T, kunci string, mundur bool) []string {
	t.Helper()
	l, err := migrasi.Daftar(mundur, berkasModul(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range l {
		if migrasi.KunciLangkah(x.Nama) == kunci {
			return x.Pernyataan
		}
	}
	t.Fatalf("langkah %s (mundur=%v) tidak ada", kunci, mundur)
	return nil
}

var semuaLangkah = []string{"085_cover_life_kolom", "086_cover_life_satu_tabel", "957_menu_coverlife"}

// viewLama - teks view DEV (fakta WO 08-10-2026, bab 1B prompt), ber-{skema}: dibuat ulang PERSIS oleh 086_down.
const viewLama = "CREATE VIEW {skema}.COVER_LIFE AS SELECT ID, a.JSONDATA.Cover, a.JSONDATA.Note FROM {skema}.M_COVER_LIFE a"

// Objek milik aplikasi lain yang TIDAK boleh disentuh (prompt bab 1B).
var polaJanganDisentuh = regexp.MustCompile(`\b(COVERAGE\w*|COVERNOTE\w*)\b`)

// K0: tepat tiga langkah (085, 086 di rentang 085-089; 957 di slot), masing-masing berpasangan `_down`.
func TestMigrasiTigaLangkahDiJatahK0(t *testing.T) {
	maju, err := migrasi.Daftar(false, berkasModul(t))
	if err != nil {
		t.Fatal(err)
	}
	mundur, err := migrasi.Daftar(true, berkasModul(t))
	if err != nil {
		t.Fatal(err)
	}
	var nama []string
	for _, l := range maju {
		nama = append(nama, migrasi.KunciLangkah(l.Nama))
	}
	if !slices.Equal(nama, semuaLangkah) || len(mundur) != len(maju) {
		t.Errorf("langkah %v (mundur %d), mau %v", nama, len(mundur), semuaLangkah)
	}
}

// C1: nol CREATE TABLE, nol RENAME, nol tabel lain (staging, _OLD, _BARU), nol COMMIT, nol objek "jangan disentuh".
func TestMigrasiTanpaTabelLainDanTanpaRename(t *testing.T) {
	for _, k := range semuaLangkah {
		for _, mundur := range []bool{false, true} {
			for _, p := range langkah(t, k, mundur) {
				if nama, _ := migrasi.KolomCreateTable(p); nama != "" {
					t.Errorf("%s (mundur=%v) membuat tabel %s", k, mundur, nama)
				}
				atas := strings.ToUpper(p)
				if strings.Contains(atas, "RENAME") || strings.Contains(atas, "COMMIT") ||
					regexp.MustCompile(`COVER_LIFE\w*_(OLD|BARU|TMP|TEMP)`).MatchString(p) || polaJanganDisentuh.MatchString(p) {
					t.Errorf("%s menyentuh yang dilarang: %s", k, p)
				}
			}
		}
	}
	if polaJanganDisentuh.MatchString("UPDATE X.M_COVER_LIFE SET COVER = 1") || !polaJanganDisentuh.MatchString("DROP TABLE S.COVERAGE_FACIN") {
		t.Error("pola jangan-disentuh tidak menggigit / terlalu lebar")
	}
}

// 085: SATU ALTER … ADD ( (KolomAlterTambah) - COVER VARCHAR2(200), NOTE VARCHAR2(1000), NULLABLE; mundurnya pelindung
// ORA-00904 lalu DROP (COVER, NOTE).
func TestMigrasi085Kolom(t *testing.T) {
	maju := langkah(t, "085_cover_life_kolom", false)
	nama, kolom := migrasi.KolomAlterTambah(maju[0])
	if len(maju) != 1 || nama != Tabel || !slices.Equal(kolom, KolomTabel[1:]) {
		t.Fatalf("085 = %s %v", nama, kolom)
	}
	s := satuBaris(maju[0])
	if !strings.Contains(s, "COVER VARCHAR2(200), NOTE VARCHAR2(1000) )") || strings.Contains(s, "NOT NULL") {
		t.Errorf("085 lebar / NULLABLE: %s", s)
	}
	turun := langkah(t, "085_cover_life_kolom", true)
	if len(turun) != 2 || satuBaris(turun[0]) != "UPDATE {skema}.M_COVER_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0" ||
		satuBaris(turun[1]) != "ALTER TABLE {skema}.M_COVER_LIFE DROP (COVER, NOTE)" {
		t.Fatalf("085 mundur %q", turun)
	}
	if _, blok := migrasi.BacaPerintahKatalog(turun[0]); blok {
		t.Error("pelindung ORA-00904 tidak boleh blok (blok melewati diri)")
	}
}

// Ekspresi pengisian 086 dan pemeriksaannya - persis seperti yang dikirim pelari.
const (
	perintahIsi = "UPDATE {skema}.M_COVER_LIFE m SET m.COVER = m.JSONDATA.Cover, m.NOTE = m.JSONDATA.Note"
	// Tiga cabang "berbeda" C1.3 per kolom - NULL = NULL bukan berbeda.
	cabangCover     = "(m.COVER IS NULL AND m.JSONDATA.Cover IS NOT NULL) OR (m.COVER IS NOT NULL AND m.JSONDATA.Cover IS NULL) OR m.COVER <> m.JSONDATA.Cover"
	cabangNote      = "(m.NOTE IS NULL AND m.JSONDATA.Note IS NOT NULL) OR (m.NOTE IS NOT NULL AND m.JSONDATA.Note IS NULL) OR m.NOTE <> m.JSONDATA.Note"
	perintahPeriksa = "UPDATE {skema}.M_COVER_LIFE m SET m.ID = NULL WHERE " + cabangCover + " OR " + cabangNote
	perintahBuang   = "ALTER TABLE {skema}.M_COVER_LIFE DROP COLUMN JSONDATA CASCADE CONSTRAINTS"
	perintahView    = "DROP VIEW {skema}.COVER_LIFE"
)

// 086 (C1): isi -> pemeriksaan -> buang JSONDATA -> DROP VIEW TERAKHIR. Tiga blok pertama ALL_TAB_COLUMNS
// M_COVER_LIFE.JSONDATA n > 0, blok keempat ALL_VIEWS COVER_LIFE n > 0 (aman diulang). Ekspresi isi = ekspresi view lama
// dengan alias berbeda SAJA (peka huruf: `Cover`, `Note`).
func TestMigrasi086SatuTabel(t *testing.T) {
	maju := langkah(t, "086_cover_life_satu_tabel", false)
	if len(maju) != 4 {
		t.Fatalf("086 %d pernyataan, mau 4", len(maju))
	}
	for i, mau := range []string{perintahIsi, perintahPeriksa, perintahBuang} {
		pk, ok := migrasi.BacaPerintahKatalog(maju[i])
		if !ok || pk != (migrasi.PerintahKatalog{Katalog: "ALL_TAB_COLUMNS", Tabel: Tabel, Objek: "JSONDATA", BilaAda: true, Perintah: mau}) {
			t.Errorf("086 blok %d: %+v %v", i, pk, ok)
		}
	}
	v, ok := migrasi.BacaPerintahKatalog(maju[3])
	if !ok || v != (migrasi.PerintahKatalog{Katalog: "ALL_VIEWS", Tabel: "COVER_LIFE", Objek: "COVER_LIFE", BilaAda: true, Perintah: perintahView}) {
		t.Errorf("086 blok view %+v %v", v, ok)
	}
	// Ekspresi isi = teks view lama (086_down) dengan alias `m` menggantikan `a` - huruf kunci persis.
	for _, kunci := range []string{"Cover", "Note"} {
		if !strings.Contains(viewLama, "a.JSONDATA."+kunci) || !strings.Contains(perintahIsi, "m.JSONDATA."+kunci) {
			t.Errorf("ekspresi isi %s tidak sama dengan view", kunci)
		}
	}
	if n := strings.Count(perintahPeriksa, "JSONDATA"); n != 6 || n != strings.Count(perintahPeriksa, "m.JSONDATA.Cover")+strings.Count(perintahPeriksa, "m.JSONDATA.Note") {
		t.Errorf("pemeriksaan membaca JSON dengan ekspresi lain: %s", perintahPeriksa)
	}
	for _, salah := range []string{"JSONDATA.COVER", "JSONDATA.cover", "JSONDATA.NOTE", "JSONDATA.note", `JSONDATA."`, "JSON_VALUE", "UPPER("} {
		if strings.Contains(perintahIsi+perintahPeriksa, salah) {
			t.Errorf("086 memuat %q - notasi titik peka huruf, persis seperti view; huruf tidak diubah", salah)
		}
	}
	if nama, kolom := migrasi.KolomAlterBuang(perintahBuang); nama != Tabel || !slices.Equal(kolom, []string{"JSONDATA"}) {
		t.Errorf("KolomAlterBuang = %s %v", nama, kolom)
	}
	gabung := strings.Join(maju, "\n")
	iIsi, iPeriksa, iBuang, iView := strings.Index(gabung, perintahIsi), strings.Index(gabung, "SET m.ID = NULL"),
		strings.Index(gabung, perintahBuang), strings.Index(gabung, perintahView)
	if !(iIsi < iPeriksa && iPeriksa < iBuang && iBuang < iView) {
		t.Error("urutan 086 harus isi -> pemeriksaan -> buang JSONDATA -> DROP VIEW (terakhir)")
	}
	turun := langkah(t, "086_cover_life_satu_tabel", true)
	if len(turun) != 5 || satuBaris(turun[0]) != "UPDATE {skema}.M_COVER_LIFE SET COVER = COVER, NOTE = NOTE WHERE 1 = 0" {
		t.Fatalf("086 mundur %q", turun)
	}
	for i, mau := range map[int]struct {
		objek string
		ada   bool
	}{1: {"JSONDATA", false}, 3: {"ENSURE_M_COVER_LIFE_JSON", false}} {
		pk, ok := migrasi.BacaPerintahKatalog(turun[i])
		if !ok || pk.Tabel != Tabel || pk.Objek != mau.objek || pk.BilaAda != mau.ada {
			t.Errorf("086 mundur %d: %+v %v", i, pk, ok)
		}
	}
	memuat(t, "086 mundur", turun[2], "SET JSONDATA = JSON_OBJECT('Cover' VALUE m.COVER, 'Note' VALUE m.NOTE ABSENT ON NULL RETURNING CLOB)")
	if satuBaris(turun[4]) != viewLama {
		t.Errorf("view dipulihkan\n%s\nmau\n%s", satuBaris(turun[4]), viewLama)
	}
}

// teks - nilai SQL: nil = NULL.
func teks(s string) *string { return &s }

// berbeda meniru SATU kelompok tiga cabang C1.3 dengan logika tiga-nilai Oracle: baris TERPILIH (ID diberi NULL ->
// ORA-01407 -> berhenti) hanya bila salah satu cabang TRUE. `<>` dengan NULL = UNKNOWN (tidak memilih).
func berbeda(kolom, view *string) bool {
	cabang1 := kolom == nil && view != nil
	cabang2 := kolom != nil && view == nil
	cabang3 := kolom != nil && view != nil && *kolom != *view
	return cabang1 || cabang2 || cabang3
}

// C1.3: berhenti HANYA bila COVER atau NOTE berbeda dari view; NULL = NULL sama (NOTE 4/4 DEV lolos).
func TestPemeriksaanC13NullSamaDenganNull(t *testing.T) {
	for _, k := range []struct {
		nama                             string
		cover, coverView, note, noteView *string
		mauBerhenti                      bool
	}{
		{"DEV: cover sama, note NULL di keduanya", teks("UJI JIWA"), teks("UJI JIWA"), nil, nil, false},
		{"cover NULL, view terisi", nil, teks("UJI JIWA"), nil, nil, true},
		{"cover terisi, view NULL", teks("UJI JIWA"), nil, nil, nil, true},
		{"cover beda huruf", teks("UJI JIWA"), teks("UJI Jiwa"), nil, nil, true},
		{"note NULL, view terisi", teks("UJI JIWA"), teks("UJI JIWA"), nil, teks("UJI CATATAN"), true},
		{"note sama", teks("UJI JIWA"), teks("UJI JIWA"), teks("UJI CATATAN"), teks("UJI CATATAN"), false},
	} {
		if dapat := berbeda(k.cover, k.coverView) || berbeda(k.note, k.noteView); dapat != k.mauBerhenti {
			t.Errorf("%s: berhenti=%v, mau %v", k.nama, dapat, k.mauBerhenti)
		}
	}
	if !strings.HasSuffix(perintahPeriksa, "WHERE "+cabangCover+" OR "+cabangNote) {
		t.Errorf("bentuk pemeriksaan %q", perintahPeriksa)
	}
}

// C1.3 (berhenti sebelum DROP JSONDATA): bentuk yang diterima pengurai produksi TIDAK boleh bertanda kutip di teks
// EXECUTE IMMEDIATE; gagal keras lewat ORA-01407 (ID NOT NULL, PK SYS_C009203); pemaksa menulis ID, tidak pernah COVER
// / NOTE.
func TestBentukPemeriksaanC13(t *testing.T) {
	blok := func(perintah string) string {
		return "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS\n" +
			"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_COVER_LIFE' AND COLUMN_NAME = 'JSONDATA';\n" +
			"  IF n > 0 THEN\n    EXECUTE IMMEDIATE '" + perintah + "';\n  END IF;\nEND;"
	}
	if _, ok := migrasi.BacaPerintahKatalog(blok("UPDATE {skema}.M_COVER_LIFE m SET m.ID = NULL WHERE JSON_EXISTS(m.JSONDATA, '$.Cover') AND m.COVER IS NULL")); ok {
		t.Error("blok bertanda kutip diterima pengurai - aturan polaPerintahKatalog berubah?")
	}
	if pk, ok := migrasi.BacaPerintahKatalog(blok(perintahPeriksa)); !ok || pk.Perintah != perintahPeriksa {
		t.Errorf("bentuk pemeriksaan ditolak pengurai: %+v %v", pk, ok)
	}
	if strings.Contains(perintahPeriksa, "SET m.COVER") || strings.Contains(perintahPeriksa, "SET m.NOTE") {
		t.Errorf("pemaksa mengubah nilai yang diperiksa: %q", perintahPeriksa)
	}
}

// 957 (C4): KODE / MODUL coverlife, LABEL 'Cover Life', MASTER TREATY URUTAN 16, DIMIGRASI '1', aman diulang (NOT
// EXISTS); mundurnya membuang hak lalu barisnya.
func TestMigrasi957Menu(t *testing.T) {
	maju := langkah(t, "957_menu_coverlife", false)
	if len(maju) != 1 || !strings.HasPrefix(satuBaris(maju[0]), "INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI) SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, ") ||
		!strings.HasSuffix(satuBaris(maju[0]),
			"'coverlife', 'Cover Life', 'MASTER TREATY', 'coverlife', 16, '1' FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'coverlife')") {
		t.Errorf("957 %q", maju)
	}
	if turun := langkah(t, "957_menu_coverlife", true); len(turun) != 2 ||
		satuBaris(turun[0]) != "DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'coverlife'" ||
		satuBaris(turun[1]) != "DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'coverlife'" {
		t.Errorf("957 mundur %q", turun)
	}
}

// K0: di skema BARU pelari menjalankan 085-086 SEBELUM 900 (tabel Pega, tidak bergantung pada inti), dan 957 SESUDAH 900,
// 901, 909, 949, 951, 955 - urutan nama berkas atas migrasi inti SUNGGUHAN. Jalur mundur kebalikannya.
func TestUrutanPelari085Lalu900Lalu957(t *testing.T) {
	maju, err := migrasi.Daftar(false, inti.SumberMigrasi(), berkasModul(t))
	if err != nil {
		t.Fatal(err)
	}
	posisi := map[string]int{}
	for i, l := range maju {
		posisi[migrasi.KunciLangkah(l.Nama)] = i
	}
	urut := []string{"085_cover_life_kolom", "086_cover_life_satu_tabel", "900_m_nav_menu", "901_m_nav_menu_datar",
		"909_m_nav_menu_master_treaty", "949_m_nav_menu_planlife", "957_menu_coverlife"}
	for i := 1; i < len(urut); i++ {
		a, adaA := posisi[urut[i-1]]
		b, adaB := posisi[urut[i]]
		if !adaA || !adaB || a >= b {
			t.Errorf("%s (%d, %v) harus sebelum %s (%d, %v)", urut[i-1], a, adaA, urut[i], b, adaB)
		}
	}
	mundur, err := migrasi.Daftar(true, inti.SumberMigrasi(), berkasModul(t))
	if err != nil {
		t.Fatal(err)
	}
	var namaMundur []string
	for _, l := range mundur {
		namaMundur = append(namaMundur, migrasi.KunciLangkah(l.Nama))
	}
	if slices.Index(namaMundur, "957_menu_coverlife") > slices.Index(namaMundur, "900_m_nav_menu") ||
		slices.Index(namaMundur, "086_cover_life_satu_tabel") > slices.Index(namaMundur, "085_cover_life_kolom") {
		t.Errorf("urutan mundur %v", namaMundur)
	}
}
