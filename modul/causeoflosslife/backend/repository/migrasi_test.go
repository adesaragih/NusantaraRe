package repository

// Migrasi modul causeoflosslife 090-092 + slot menu 955 (keputusan work owner 08-10-2026 K0-K5) dibaca pengurai
// PRODUKSI pelari (`migrasi.Daftar`, `BacaPerintahKatalog`, `KolomAlterTambah`, `KolomAlterBuang`, `KolomCreateTable`)
// dan berpasangan dengan jalur mundurnya. Nol CREATE TABLE CAUSEOFLOSS_LIFE (K1: RENAME, bukan tabel baru). Urutan
// pelari 090 -> 900 -> 955 dibuktikan atas migrasi inti SUNGGUHAN (K0).

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

var semuaLangkah = []string{"090_causeofloss_life_ganti_nama", "091_causeofloss_life_kolom", "092_causeofloss_life_satu_tabel",
	"955_menu_causeoflosslife"}

// viewLama - teks view DEV (fakta WO 08-10-2026, bab 1 prompt), ber-{skema}: dibuat ulang PERSIS oleh 090_down.
const viewLama = "CREATE VIEW {skema}.CAUSEOFLOSS_LIFE AS SELECT a.ID, a.JSONDATA.CauseofLoss FROM {skema}.M_CAUSEOFLOSS_LIFE a"

// Objek milik aplikasi lain yang TIDAK boleh disentuh (prompt bab 1).
var polaJanganDisentuh = regexp.MustCompile(`\b(M_CAUSE_OF_LOSS\w*|D_CAUSE_OF_LOSS\w*|V_M_CAUSE_OF_LOSS\w*|V_D_CAUSE_OF_LOSS\w*|` +
	`T_LISTCAUSEOFLOSS|PEGA_M_CAUSE_OF_LOSS|PEGA_D_CAUSE_OF_LOSS|PEGA_M_CAUSEOFLOSS_LIFE)\b`)

// K0: tepat empat langkah (090, 091, 092 di rentang 090-099; 955 di slot), masing-masing berpasangan `_down`.
func TestMigrasiEmpatLangkahDiJatahK0(t *testing.T) {
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

// K1: nol CREATE TABLE, nol tabel lain (staging, _OLD, _BARU), nol objek "jangan disentuh".
func TestMigrasiTanpaTabelLain(t *testing.T) {
	for _, k := range semuaLangkah {
		for _, mundur := range []bool{false, true} {
			for _, p := range langkah(t, k, mundur) {
				if nama, _ := migrasi.KolomCreateTable(p); nama != "" {
					t.Errorf("%s (mundur=%v) membuat tabel %s", k, mundur, nama)
				}
				if regexp.MustCompile(`CAUSEOFLOSS_LIFE\w*_(OLD|BARU|TMP|TEMP)`).MatchString(p) || polaJanganDisentuh.MatchString(p) {
					t.Errorf("%s menyentuh objek lain: %s", k, p)
				}
				if strings.Contains(strings.ToUpper(p), "COMMIT") {
					t.Errorf("%s memuat COMMIT", k)
				}
			}
		}
	}
	if polaJanganDisentuh.MatchString("UPDATE X.CAUSEOFLOSS_LIFE") || !polaJanganDisentuh.MatchString("DROP VIEW S.V_D_CAUSE_OF_LOSS_X") {
		t.Error("pola jangan-disentuh tidak menggigit")
	}
}

// 090: blok ALL_VIEWS (DROP VIEW hanya bila VIEW) lalu blok RENAME berpelindung tabel sumber; mundurnya: pelindung
// ORA-00904, RENAME balik berpelindung JSONDATA, CREATE VIEW PERSIS teks DEV.
func TestMigrasi090GantiNama(t *testing.T) {
	maju := langkah(t, "090_causeofloss_life_ganti_nama", false)
	if len(maju) != 2 {
		t.Fatalf("090 %d pernyataan, mau 2", len(maju))
	}
	v, ok := migrasi.BacaPerintahKatalog(maju[0])
	if !ok || v != (migrasi.PerintahKatalog{Katalog: "ALL_VIEWS", Tabel: Tabel, Objek: Tabel, BilaAda: true, Perintah: "DROP VIEW {skema}.CAUSEOFLOSS_LIFE"}) {
		t.Errorf("blok view %+v %v", v, ok)
	}
	r, ok := migrasi.BacaPerintahKatalog(maju[1])
	if !ok || r != (migrasi.PerintahKatalog{Katalog: "ALL_TAB_COLUMNS", Tabel: "M_CAUSEOFLOSS_LIFE", Objek: "ID", BilaAda: true,
		Perintah: "ALTER TABLE {skema}.M_CAUSEOFLOSS_LIFE RENAME TO CAUSEOFLOSS_LIFE"}) {
		t.Errorf("blok rename %+v %v", r, ok)
	}
	turun := langkah(t, "090_causeofloss_life_ganti_nama", true)
	if len(turun) != 3 || satuBaris(turun[0]) != "UPDATE {skema}.CAUSEOFLOSS_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0" {
		t.Fatalf("090 mundur %q", turun)
	}
	if _, blok := migrasi.BacaPerintahKatalog(turun[0]); blok {
		t.Error("pelindung ORA-00904 tidak boleh blok (blok melewati diri)")
	}
	b, ok := migrasi.BacaPerintahKatalog(turun[1])
	if !ok || b != (migrasi.PerintahKatalog{Katalog: "ALL_TAB_COLUMNS", Tabel: Tabel, Objek: "JSONDATA", BilaAda: true,
		Perintah: "ALTER TABLE {skema}.CAUSEOFLOSS_LIFE RENAME TO M_CAUSEOFLOSS_LIFE"}) {
		t.Errorf("mundur rename %+v %v", b, ok)
	}
	if satuBaris(turun[2]) != viewLama {
		t.Errorf("view dipulihkan\n%s\nmau\n%s", satuBaris(turun[2]), viewLama)
	}
}

// 091: SATU ALTER … ADD ( (KolomAlterTambah) - CAUSEOFLOSS VARCHAR2(200) NULLABLE (pola Benefit 943); mundurnya
// pelindung ORA-00904 lalu DROP (CAUSEOFLOSS).
func TestMigrasi091Kolom(t *testing.T) {
	maju := langkah(t, "091_causeofloss_life_kolom", false)
	nama, kolom := migrasi.KolomAlterTambah(maju[0])
	if len(maju) != 1 || nama != Tabel || !slices.Equal(kolom, KolomTabel[1:]) {
		t.Fatalf("091 = %s %v", nama, kolom)
	}
	if !regexp.MustCompile(`\bCAUSEOFLOSS\s+VARCHAR2\(200\)\s*\)`).MatchString(satuBaris(maju[0])) || strings.Contains(maju[0], "NOT NULL") {
		t.Errorf("091 lebar / NULLABLE: %s", satuBaris(maju[0]))
	}
	turun := langkah(t, "091_causeofloss_life_kolom", true)
	if len(turun) != 2 || satuBaris(turun[0]) != "UPDATE {skema}.CAUSEOFLOSS_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0" ||
		satuBaris(turun[1]) != "ALTER TABLE {skema}.CAUSEOFLOSS_LIFE DROP (CAUSEOFLOSS)" {
		t.Fatalf("091 mundur %q", turun)
	}
}

// Ekspresi pengisian 092 dan pemeriksaannya - persis seperti yang dikirim pelari.
const (
	ekspresiViewLama = "a.JSONDATA.CauseofLoss"
	perintahIsi      = "UPDATE {skema}.CAUSEOFLOSS_LIFE m SET m.CAUSEOFLOSS = m.JSONDATA.CauseofLoss"
	// Tiga cabang "berbeda" K1.4 - NULL = NULL bukan berbeda.
	cabangKolomNull = "(m.CAUSEOFLOSS IS NULL AND m.JSONDATA.CauseofLoss IS NOT NULL)"
	cabangViewNull  = "(m.CAUSEOFLOSS IS NOT NULL AND m.JSONDATA.CauseofLoss IS NULL)"
	cabangBeda      = "m.CAUSEOFLOSS <> m.JSONDATA.CauseofLoss"
	perintahPeriksa = "UPDATE {skema}.CAUSEOFLOSS_LIFE m SET m.ID = NULL WHERE " + cabangKolomNull + " OR " + cabangViewNull + " OR " + cabangBeda
	perintahBuang   = "ALTER TABLE {skema}.CAUSEOFLOSS_LIFE DROP COLUMN JSONDATA CASCADE CONSTRAINTS"
)

// 092 (K1/K2): isi -> pemeriksaan -> buang, ketiganya blok ALL_TAB_COLUMNS CAUSEOFLOSS_LIFE.JSONDATA n > 0 (aman
// diulang). Ekspresi isi = ekspresi view lama dengan alias berbeda SAJA (peka huruf: `CauseofLoss`), jadi nilainya
// identik dengan yang view keluarkan - termasuk NULL untuk `"CauseofLoss":""` (100001).
func TestMigrasi092SatuTabel(t *testing.T) {
	maju := langkah(t, "092_causeofloss_life_satu_tabel", false)
	if len(maju) != 3 {
		t.Fatalf("092 %d pernyataan, mau 3", len(maju))
	}
	for i, mau := range []string{perintahIsi, perintahPeriksa, perintahBuang} {
		pk, ok := migrasi.BacaPerintahKatalog(maju[i])
		if !ok || pk != (migrasi.PerintahKatalog{Katalog: "ALL_TAB_COLUMNS", Tabel: Tabel, Objek: "JSONDATA", BilaAda: true, Perintah: mau}) {
			t.Errorf("092 blok %d: %+v %v", i, pk, ok)
		}
	}
	// Ekspresi isi = teks view lama (090_down) dengan alias `m` menggantikan `a` - huruf kunci persis.
	if !strings.Contains(viewLama, ekspresiViewLama) || !strings.HasSuffix(perintahIsi, "= m."+strings.TrimPrefix(ekspresiViewLama, "a.")) {
		t.Errorf("ekspresi isi %q tidak sama dengan view %q", perintahIsi, ekspresiViewLama)
	}
	// Setiap pembacaan JSON di pemeriksaan memakai ekspresi yang SAMA (bukan JSON_VALUE / huruf lain).
	if n := strings.Count(perintahPeriksa, "JSONDATA"); n != strings.Count(perintahPeriksa, "m.JSONDATA.CauseofLoss") || n != 3 {
		t.Errorf("pemeriksaan membaca JSON dengan ekspresi lain: %s", perintahPeriksa)
	}
	for _, salah := range []string{"JSONDATA.CAUSEOFLOSS", "JSONDATA.causeofloss", "JSONDATA.CauseOfLoss", `JSONDATA."`, "JSON_VALUE", "UPPER("} {
		if strings.Contains(perintahIsi+perintahPeriksa, salah) {
			t.Errorf("092 memuat %q - notasi titik peka huruf, persis seperti view; huruf tidak diubah", salah)
		}
	}
	if nama, kolom := migrasi.KolomAlterBuang(perintahBuang); nama != Tabel || !slices.Equal(kolom, []string{"JSONDATA"}) {
		t.Errorf("KolomAlterBuang = %s %v", nama, kolom)
	}
	gabung := strings.Join(maju, "\n")
	if !(strings.Index(gabung, perintahIsi) < strings.Index(gabung, "SET m.ID = NULL") &&
		strings.Index(gabung, "SET m.ID = NULL") < strings.Index(gabung, perintahBuang)) {
		t.Error("urutan 092 harus isi -> pemeriksaan -> buang")
	}
	turun := langkah(t, "092_causeofloss_life_satu_tabel", true)
	if len(turun) != 3 {
		t.Fatalf("092 mundur %d pernyataan, mau 3", len(turun))
	}
	for i, mau := range map[int]struct {
		objek string
		ada   bool
	}{0: {"JSONDATA", false}, 2: {"ENSURE_M_CAUSEOFLOSS_LIFE_JSON", false}} {
		pk, ok := migrasi.BacaPerintahKatalog(turun[i])
		if !ok || pk.Tabel != Tabel || pk.Objek != mau.objek || pk.BilaAda != mau.ada {
			t.Errorf("092 mundur %d: %+v %v", i, pk, ok)
		}
	}
	memuat(t, "092 mundur", turun[1], "SET JSONDATA = JSON_OBJECT('CauseofLoss' VALUE m.CAUSEOFLOSS ABSENT ON NULL RETURNING CLOB)")
}

// teks - nilai SQL: nil = NULL.
func teks(s string) *string { return &s }

// berbeda meniru WHERE pemeriksaan 092 dengan logika tiga-nilai Oracle (TRUE / FALSE / UNKNOWN): baris TERPILIH (ID
// diberi NULL -> ORA-01407 -> berhenti) hanya bila salah satu cabang TRUE. `<>` dengan NULL = UNKNOWN (tidak memilih).
func berbeda(kolom, view *string) bool {
	cabang1 := kolom == nil && view != nil
	cabang2 := kolom != nil && view == nil
	cabang3 := kolom != nil && view != nil && *kolom != *view // UNKNOWN bila salah satu NULL
	return cabang1 || cabang2 || cabang3
}

// K1.4: berhenti HANYA bila nilai kolom berbeda dari view; NULL = NULL sama (100001 `"CauseofLoss":""` lolos).
// Dibuktikan atas tiruan logika tiga-nilai; perilaku Oracle sungguhan: uji -tags=db TestDBPemeriksaanK14 (menunggu
// skema uji).
func TestPemeriksaanK14NullSamaDenganNull(t *testing.T) {
	for _, k := range []struct {
		nama        string
		kolom, view *string
		mauBerhenti bool
	}{
		{"100001: NULL di keduanya", nil, nil, false},
		{"isi sama", teks("UJI SAKIT"), teks("UJI SAKIT"), false},
		{"kolom NULL, view terisi", nil, teks("UJI SAKIT"), true},
		{"kolom terisi, view NULL", teks("UJI SAKIT"), nil, true},
		{"isi berbeda", teks("UJI SAKIT"), teks("UJI Sakit"), true},
	} {
		if dapat := berbeda(k.kolom, k.view); dapat != k.mauBerhenti {
			t.Errorf("%s: berhenti=%v, mau %v", k.nama, dapat, k.mauBerhenti)
		}
	}
	// Tiruan di atas memang tiruan teks perintahnya: tiga cabang, digabung OR, bentuk persis.
	if !strings.HasSuffix(perintahPeriksa, "WHERE "+cabangKolomNull+" OR "+cabangViewNull+" OR "+cabangBeda) {
		t.Errorf("bentuk pemeriksaan %q", perintahPeriksa)
	}
}

// K1.4 (berhenti sebelum DROP JSONDATA): bentuk yang diterima pengurai produksi TIDAK boleh bertanda kutip di teks
// EXECUTE IMMEDIATE - versi JSON_EXISTS(…, '$.CauseofLoss') / DECODE dengan literal ditolak pengurai, versi notasi titik
// diterima. Gagal keras lewat ORA-01407 (ID NOT NULL, PK SYS_C008825) - bukan galat yang ditoleransi pelari (hanya
// ORA-00955 pada CREATE maju, ORA-00942 / ORA-02289 mundur). Pemaksa menulis ID, tidak pernah CAUSEOFLOSS.
func TestBentukPemeriksaanK14(t *testing.T) {
	blok := func(perintah string) string {
		return "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS\n" +
			"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'CAUSEOFLOSS_LIFE' AND COLUMN_NAME = 'JSONDATA';\n" +
			"  IF n > 0 THEN\n    EXECUTE IMMEDIATE '" + perintah + "';\n  END IF;\nEND;"
	}
	if _, ok := migrasi.BacaPerintahKatalog(blok("UPDATE {skema}.CAUSEOFLOSS_LIFE m SET m.ID = NULL WHERE JSON_EXISTS(m.JSONDATA, '$.CauseofLoss') AND m.CAUSEOFLOSS IS NULL")); ok {
		t.Error("blok bertanda kutip diterima pengurai - aturan polaPerintahKatalog berubah?")
	}
	if pk, ok := migrasi.BacaPerintahKatalog(blok(perintahPeriksa)); !ok || pk.Perintah != perintahPeriksa {
		t.Errorf("bentuk pemeriksaan ditolak pengurai: %+v %v", pk, ok)
	}
	if !strings.Contains(perintahPeriksa, "SET m.ID = NULL WHERE") || strings.Contains(perintahPeriksa, "SET m.CAUSEOFLOSS") {
		t.Errorf("bentuk pemeriksaan %q", perintahPeriksa)
	}
}

// 955 (K5): KODE / MODUL causeoflosslife, LABEL 'Cause Of Loss Life', MASTER TREATY URUTAN 14, DIMIGRASI '1', aman
// diulang (NOT EXISTS); mundurnya membuang hak lalu barisnya.
func TestMigrasi955Menu(t *testing.T) {
	maju := langkah(t, "955_menu_causeoflosslife", false)
	if len(maju) != 1 || !strings.HasPrefix(satuBaris(maju[0]), "INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI) SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, ") ||
		!strings.HasSuffix(satuBaris(maju[0]),
			"'causeoflosslife', 'Cause Of Loss Life', 'MASTER TREATY', 'causeoflosslife', 14, '1' FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'causeoflosslife')") {
		t.Errorf("955 %q", maju)
	}
	if turun := langkah(t, "955_menu_causeoflosslife", true); len(turun) != 2 ||
		satuBaris(turun[0]) != "DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'causeoflosslife'" ||
		satuBaris(turun[1]) != "DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'causeoflosslife'" {
		t.Errorf("955 mundur %q", turun)
	}
}

// K0: di skema BARU pelari menjalankan 090-092 SEBELUM 900 (tabel tidak bergantung pada inti), dan 955 SESUDAH 900,
// 901 (menu datar), 909 (CHECK MASTER TREATY), dan 949 (Plan URUTAN 13) - urutan nama berkas atas migrasi inti
// SUNGGUHAN. Jalur mundur kebalikannya: 955 dibongkar sebelum 900_down.
func TestUrutanPelari090Lalu900Lalu955(t *testing.T) {
	maju, err := migrasi.Daftar(false, inti.SumberMigrasi(), berkasModul(t))
	if err != nil {
		t.Fatal(err)
	}
	posisi := map[string]int{}
	for i, l := range maju {
		posisi[migrasi.KunciLangkah(l.Nama)] = i
	}
	urut := []string{"090_causeofloss_life_ganti_nama", "091_causeofloss_life_kolom", "092_causeofloss_life_satu_tabel",
		"900_m_nav_menu", "901_m_nav_menu_datar", "909_m_nav_menu_master_treaty", "949_m_nav_menu_planlife", "955_menu_causeoflosslife"}
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
	if slices.Index(namaMundur, "955_menu_causeoflosslife") > slices.Index(namaMundur, "900_m_nav_menu") ||
		slices.Index(namaMundur, "090_causeofloss_life_ganti_nama") < slices.Index(namaMundur, "900_m_nav_menu") {
		t.Errorf("urutan mundur %v", namaMundur)
	}
}
