package repository

// Migrasi modul diseaselife 080-081 + slot menu 951 (keputusan work owner 08-10-2026 K0, D1, D4) dibaca pengurai
// PRODUKSI pelari (`migrasi.Daftar`, `BacaSequenceDariKueri`, `BacaPerintahKatalog`, `KolomCreateTable`,
// `KolomAlterTambah`) dan berpasangan dengan jalur mundurnya. Nol CREATE TABLE, nol RENAME, nol kolom baru, nol tulis
// data (D1: tabel sudah flat; baris uji 102051 dihapus WO lewat berkas docs/sql, bukan migrasi - D2). Urutan pelari
// 080 -> 900 -> 951 -> 955 dibuktikan atas migrasi inti SUNGGUHAN (K0).

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

var semuaLangkah = []string{"080_seq_disease_life", "081_disease_life_pk", "951_menu_diseaselife"}

// Objek yang TIDAK boleh disentuh (D2: M_DISEASE_LIFE JSON lama + prosedurnya; D1.1: sequence Pega; tabel klaim).
var polaJanganDisentuh = regexp.MustCompile(`\b(M_DISEASE_LIFE\w*|PEGA_M_DISEASE_LIFE|PEGA_DISEASE_LIFE|T_CLAIMLF_\w+)\b`)

// K0: tepat tiga langkah (080, 081 di rentang 080-084; 951 di slot), masing-masing berpasangan `_down`.
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

// D1: nol CREATE TABLE, nol RENAME, nol kolom baru, nol INSERT / UPDATE / DELETE atas DISEASE_LIFE (data tidak
// disentuh migrasi), nol COMMIT, nol objek "jangan disentuh".
func TestMigrasiTanpaTabelLainDanTanpaTulisData(t *testing.T) {
	for _, k := range semuaLangkah {
		for _, mundur := range []bool{false, true} {
			for _, p := range langkah(t, k, mundur) {
				atas := strings.ToUpper(p)
				if nama, _ := migrasi.KolomCreateTable(p); nama != "" {
					t.Errorf("%s (mundur=%v) membuat tabel %s", k, mundur, nama)
				}
				if nama, _ := migrasi.KolomAlterTambah(p); nama != "" {
					t.Errorf("%s menambah kolom %s", k, nama)
				}
				if strings.Contains(atas, "RENAME") || strings.Contains(atas, "COMMIT") || polaJanganDisentuh.MatchString(p) {
					t.Errorf("%s menyentuh yang dilarang: %s", k, p)
				}
				if regexp.MustCompile(`(?i)(INSERT INTO|UPDATE|DELETE FROM)\s+\{skema\}\.DISEASE_LIFE\b`).MatchString(p) {
					t.Errorf("%s menulis data DISEASE_LIFE: %s", k, p)
				}
			}
		}
	}
	if !polaJanganDisentuh.MatchString("UPDATE S.M_DISEASE_LIFE_SEQ") || polaJanganDisentuh.MatchString("ALTER TABLE {skema}.DISEASE_LIFE") {
		t.Error("pola jangan-disentuh tidak menggigit / terlalu lebar")
	}
}

// 080 (D1.1): bentuk PERSIS 923 - blok sequence-dari-kueri (ALL_SEQUENCES dulu), awal = ID angka tertinggi + 1 dihitung
// di basis data, INCREMENT BY 1 NOCACHE NOCYCLE; mundurnya DROP SEQUENCE biasa.
func TestMigrasi080Sequence(t *testing.T) {
	maju := langkah(t, "080_seq_disease_life", false)
	if len(maju) != 1 {
		t.Fatalf("080 %d pernyataan, mau 1", len(maju))
	}
	s, ok := migrasi.BacaSequenceDariKueri(maju[0])
	if !ok || s != (migrasi.SequenceDariKueri{Nama: "SEQ_DISEASE_LIFE",
		Kueri: "SELECT NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0) + 1 FROM {skema}.DISEASE_LIFE",
		Opsi:  "INCREMENT BY 1 NOCACHE NOCYCLE"}) {
		t.Errorf("080 %+v %v", s, ok)
	}
	// Bentuk sama dengan 923 (kueri awal untuk SEQ_M_RATE_LIFE) - hanya nama tabel dan sequence berbeda.
	isi923, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "inti", "backend", "migrations", "923_seq_rate_life.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi923), "SELECT NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0) + 1 INTO awal FROM {skema}.M_RATE_LIFE;") ||
		!strings.Contains(maju[0], "SELECT NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0) + 1 INTO awal FROM {skema}.DISEASE_LIFE;") {
		t.Error("080 menyimpang dari bentuk 923")
	}
	if strings.Contains(maju[0], "M_DISEASE_LIFE_SEQ") || strings.Contains(maju[0], "LPAD") {
		t.Error("080 tidak boleh memakai rumus / sequence Pega")
	}
	if turun := langkah(t, "080_seq_disease_life", true); len(turun) != 1 || satuBaris(turun[0]) != "DROP SEQUENCE {skema}.SEQ_DISEASE_LIFE" {
		t.Errorf("080 mundur %q", turun)
	}
}

// 081 (D1.2): SATU blok berpelindung ALL_CONSTRAINTS (n = 0) berisi ALTER ... ADD CONSTRAINT PK_DISEASE_LIFE PRIMARY
// KEY (ID) - tanpa tanda kutip di teks EXECUTE IMMEDIATE; galatnya sendiri (ORA-02437 kembar / ORA-01449 NULL) adalah
// penghenti, nol baris dihapus. Mundurnya blok berpelindung yang sama (n > 0) DROP CONSTRAINT.
func TestMigrasi081PK(t *testing.T) {
	maju := langkah(t, "081_disease_life_pk", false)
	if len(maju) != 1 {
		t.Fatalf("081 %d pernyataan, mau 1", len(maju))
	}
	pk, ok := migrasi.BacaPerintahKatalog(maju[0])
	if !ok || pk != (migrasi.PerintahKatalog{Katalog: "ALL_CONSTRAINTS", Tabel: Tabel, Objek: "PK_DISEASE_LIFE", BilaAda: false,
		Perintah: "ALTER TABLE {skema}.DISEASE_LIFE ADD CONSTRAINT PK_DISEASE_LIFE PRIMARY KEY (ID)"}) {
		t.Errorf("081 %+v %v", pk, ok)
	}
	if strings.Contains(pk.Perintah, "'") {
		t.Error("teks EXECUTE IMMEDIATE bertanda kutip")
	}
	turun := langkah(t, "081_disease_life_pk", true)
	b, ok := migrasi.BacaPerintahKatalog(turun[0])
	if len(turun) != 1 || !ok || b != (migrasi.PerintahKatalog{Katalog: "ALL_CONSTRAINTS", Tabel: Tabel, Objek: "PK_DISEASE_LIFE", BilaAda: true,
		Perintah: "ALTER TABLE {skema}.DISEASE_LIFE DROP CONSTRAINT PK_DISEASE_LIFE"}) {
		t.Errorf("081 mundur %+v %v", b, ok)
	}
	// Nama objek <= 30 byte (Oracle < 12.2 ORA-00972; penjaga inti TestPengenalTidakLebihDari30Byte).
	for _, n := range []string{"SEQ_DISEASE_LIFE", "PK_DISEASE_LIFE"} {
		if len(n) > 30 {
			t.Errorf("%s > 30 byte", n)
		}
	}
}

// 951 (D4): KODE / MODUL diseaselife, LABEL 'Disease Life', MASTER TREATY URUTAN 15, DIMIGRASI '1', aman diulang (NOT
// EXISTS); mundurnya membuang hak lalu barisnya.
func TestMigrasi951Menu(t *testing.T) {
	maju := langkah(t, "951_menu_diseaselife", false)
	if len(maju) != 1 || !strings.HasPrefix(satuBaris(maju[0]), "INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI) SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, ") ||
		!strings.HasSuffix(satuBaris(maju[0]),
			"'diseaselife', 'Disease Life', 'MASTER TREATY', 'diseaselife', 15, '1' FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'diseaselife')") {
		t.Errorf("951 %q", maju)
	}
	if turun := langkah(t, "951_menu_diseaselife", true); len(turun) != 2 ||
		satuBaris(turun[0]) != "DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'diseaselife'" ||
		satuBaris(turun[1]) != "DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'diseaselife'" {
		t.Errorf("951 mundur %q", turun)
	}
}

// K0: di skema BARU pelari menjalankan 080-081 SEBELUM 900 (hanya membaca / mengubah objek Pega), dan 951 SESUDAH 900,
// 901 (menu datar), 909 (CHECK MASTER TREATY), 949 (Plan) - urutan nama berkas atas migrasi inti SUNGGUHAN. Jalur
// mundur kebalikannya: 951 dibongkar sebelum 900_down, 081_down sebelum 080_down.
func TestUrutanPelari080Lalu900Lalu951(t *testing.T) {
	maju, err := migrasi.Daftar(false, inti.SumberMigrasi(), berkasModul(t))
	if err != nil {
		t.Fatal(err)
	}
	posisi := map[string]int{}
	for i, l := range maju {
		posisi[migrasi.KunciLangkah(l.Nama)] = i
	}
	urut := []string{"080_seq_disease_life", "081_disease_life_pk", "900_m_nav_menu", "901_m_nav_menu_datar",
		"909_m_nav_menu_master_treaty", "949_m_nav_menu_planlife", "951_menu_diseaselife"}
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
	if slices.Index(namaMundur, "951_menu_diseaselife") > slices.Index(namaMundur, "900_m_nav_menu") ||
		slices.Index(namaMundur, "081_disease_life_pk") > slices.Index(namaMundur, "080_seq_disease_life") {
		t.Errorf("urutan mundur %v", namaMundur)
	}
}
