package repository

// Penjaga MODUL Master Contract Retro Life (pola `TestTCONolTabelBaru` Treaty
// Contract Out ditiru, tidak diimpor).
//
//	K1  nol tabel baru; DDL hanya constraint dan indeks di migrasi 100–139
//	    (keputusan work owner 01-10-2026; slot menu 958–959 hanya UPDATE)
//	-   master rujukan dibaca saja
//	-   nol kata cadangan Oracle sebagai nama kolom telanjang
//	-   paket `tiruan` hanya diimpor uji
//
// Setiap penjaga punya uji gigit: aturan yang tidak pernah terbukti merah
// tidak menjaga apa pun.

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// akarModul - folder modul ini dari folder paket repository.
const akarModul = "../.."

func berkasModul(t *testing.T, akhiran ...string) map[string]string {
	t.Helper()
	hasil := map[string]string{}
	err := filepath.Walk(akarModul, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && (info.Name() == "node_modules" || info.Name() == "docs") {
			return filepath.SkipDir
		}
		for _, a := range akhiran {
			if !info.IsDir() && strings.HasSuffix(p, a) {
				isi, err := os.ReadFile(p)
				if err != nil {
					return err
				}
				hasil[filepath.ToSlash(p)] = string(isi)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return hasil
}

// --- K1: nol migrasi di rentang, nol DDL -----------------------------------

var polaNomorMigrasi = regexp.MustCompile(`^(\d{3})_.*\.sql$`)

func pelanggaranMigrasi(nama string) string {
	m := polaNomorMigrasi.FindStringSubmatch(nama)
	if m == nil {
		return "berkas migrasi tanpa nomor tiga digit"
	}
	n, _ := strconv.Atoi(m[1])
	switch {
	case n >= 100 && n <= 139:
		// Keputusan work owner 01-10-2026: ALTER TABLE diizinkan untuk relasi
		// (constraint dan indeks) - isinya dijaga TestMCRLDDLHanyaRelasi.
		return ""
	case n == 958 || n == 959:
		return ""
	default:
		return "nomor di luar rentang dan slot menu modul ini"
	}
}

func TestMCRLNolMigrasiDiRentang(t *testing.T) {
	berkas, err := filepath.Glob(filepath.Join(akarModul, "backend", "migrations", "*.sql"))
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range berkas {
		if p := pelanggaranMigrasi(filepath.Base(b)); p != "" {
			t.Errorf("%s: %s", filepath.Base(b), p)
		}
	}
}

func TestMCRLAturanMigrasiMenggigit(t *testing.T) {
	for _, nama := range []string{"140_x.sql", "099_x_down.sql", "300_x.sql", "tanpa_nomor.sql"} {
		if pelanggaranMigrasi(nama) == "" {
			t.Errorf("%s seharusnya ditolak", nama)
		}
	}
	for _, nama := range []string{"100_relasi_lima_tabel_life.sql", "139_x_down.sql",
		"958_menu_mastercontractretrolife.sql", "958_menu_mastercontractretrolife_down.sql"} {
		if p := pelanggaranMigrasi(nama); p != "" {
			t.Errorf("%s seharusnya sah: %s", nama, p)
		}
	}
}

var polaDDL = regexp.MustCompile(`(?i)\b(CREATE|ALTER|DROP|TRUNCATE)\s+(TABLE|SEQUENCE|INDEX|VIEW|SYNONYM)\b`)

func adaDDL(teks string) bool { return polaDDL.MatchString(teks) }

func TestMCRLNolDDL(t *testing.T) {
	for jalur, isi := range berkasModul(t, ".go", ".sql") {
		if strings.HasSuffix(jalur, "_test.go") {
			continue // tiruan skema uji `db` membuat tabel TIRUAN di skema uji
		}
		if strings.Contains(jalur, "/migrations/1") {
			continue // migrasi relasi 100-139 dijaga TestMCRLDDLHanyaRelasi
		}
		if adaDDL(buangKomentarGo(isi)) {
			t.Errorf("%s memuat DDL - kode modul tidak pernah menjalankan DDL", jalur)
		}
	}
}

// --- Relasi (keputusan work owner 01-10-2026): DDL hanya constraint dan indeks

// polaDDLTerlarang - DDL yang tetap dilarang di migrasi relasi: tabel baru atau
// dibuang, kolom ditambah, dibuang, atau diganti nama, dan TRUNCATE. Prosedur
// Pega masih membaca dan menulis setiap kolom lima tabel warisan.
var polaDDLTerlarang = regexp.MustCompile(`(?i)\b(CREATE\s+TABLE|DROP\s+TABLE|TRUNCATE|DROP\s+COLUMN|RENAME|ADD\s*\(|CREATE\s+(SEQUENCE|VIEW|SYNONYM)|DROP\s+(SEQUENCE|VIEW|SYNONYM))\b`)

func pelanggaranDDLRelasi(teks string) string {
	if m := polaDDLTerlarang.FindString(teks); m != "" {
		return "DDL di luar relasi: " + m
	}
	return ""
}

func TestMCRLDDLHanyaRelasi(t *testing.T) {
	n := 0
	for jalur, isi := range berkasModul(t, ".sql") {
		if !strings.Contains(jalur, "/migrations/1") {
			continue
		}
		n++
		if p := pelanggaranDDLRelasi(isi); p != "" {
			t.Errorf("%s: %s", jalur, p)
		}
	}
	if n == 0 {
		t.Fatal("nol berkas migrasi relasi 100-139 terbaca - penjaga lulus hampa")
	}
}

func TestMCRLAturanDDLRelasiMenggigit(t *testing.T) {
	for _, s := range []string{"CREATE TABLE {skema}.X (ID NUMBER)", "ALTER TABLE {skema}.X ADD (Y NUMBER)",
		"ALTER TABLE {skema}.X DROP COLUMN Y", "DROP TABLE {skema}.X", "ALTER TABLE {skema}.X RENAME COLUMN A TO B"} {
		if pelanggaranDDLRelasi(s) == "" {
			t.Errorf("%q seharusnya ditolak", s)
		}
	}
	for _, s := range []string{"ALTER TABLE {skema}.X ADD CONSTRAINT PK_X PRIMARY KEY (ID)",
		"ALTER TABLE {skema}.X MODIFY (A NOT NULL)", "CREATE INDEX {skema}.IX_X ON {skema}.X (A)",
		"ALTER TABLE {skema}.X DROP CONSTRAINT PK_X", "DROP INDEX {skema}.IX_X"} {
		if p := pelanggaranDDLRelasi(s); p != "" {
			t.Errorf("%q seharusnya sah: %s", s, p)
		}
	}
}

func TestMCRLAturanDDLMenggigit(t *testing.T) {
	for _, s := range []string{"CREATE TABLE {skema}.X (ID NUMBER)", "alter  table x add y", "DROP SEQUENCE S"} {
		if !adaDDL(s) {
			t.Errorf("%q seharusnya terbaca DDL", s)
		}
	}
	if adaDDL("UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1'") {
		t.Error("UPDATE slot menu bukan DDL")
	}
}

func buangKomentarGo(isi string) string {
	var b strings.Builder
	for _, baris := range strings.Split(isi, "\n") {
		if strings.HasPrefix(strings.TrimSpace(baris), "//") {
			continue
		}
		b.WriteString(baris)
		b.WriteByte('\n')
	}
	return b.String()
}

// --- master dibaca saja ------------------------------------------------------

var polaTulis = regexp.MustCompile(`(?i)\b(INSERT\s+INTO|UPDATE|DELETE\s+FROM|MERGE\s+INTO)\s+%s`)

// tulisMaster - fungsi SQL yang menulis ke master. Setiap fungsi `sql…` di
// berkas repository non-uji dibaca; yang menulis dan menyebut objek master
// adalah pelanggaran.
func tulisMaster(isi string) []string {
	var hasil []string
	for _, blok := range strings.Split(isi, "\nfunc ") {
		if !polaTulis.MatchString(blok) {
			continue
		}
		for _, m := range DaftarMasterDibacaSaja {
			if strings.Contains(blok, m) {
				hasil = append(hasil, m)
			}
		}
	}
	return hasil
}

func TestMCRLMasterDibacaSaja(t *testing.T) {
	for jalur, isi := range berkasModul(t, ".go") {
		if !strings.Contains(jalur, "/repository/") || strings.HasSuffix(jalur, "_test.go") {
			continue
		}
		if m := tulisMaster(buangKomentarGo(isi)); len(m) > 0 {
			t.Errorf("%s menulis master %v - master rujukan dibaca saja", jalur, m)
		}
	}
}

func TestMCRLAturanMasterMenggigit(t *testing.T) {
	isi := "\nfunc sqlX(t string) string {\n\treturn fmt.Sprintf(`UPDATE %s SET NOTE = :1`, t) // AGENT\n}\n"
	if len(tulisMaster(isi)) == 0 {
		t.Error("UPDATE yang menyebut AGENT seharusnya tertangkap")
	}
}

// --- kata cadangan Oracle ----------------------------------------------------

// kataCadanganOracle - disalin dari daftar penjaga inti (pl6, 28-09-2026);
// sengaja tidak lengkap, sama seperti sumbernya.
var kataCadanganOracle = map[string]bool{
	"INITIAL": true, "LEVEL": true, "SIZE": true, "DATE": true, "NUMBER": true, "COMMENT": true, "ORDER": true,
	"GROUP": true, "CHECK": true, "DEFAULT": true, "ACCESS": true, "AUDIT": true, "CLUSTER": true, "COLUMN": true,
	"OPTION": true, "ROW": true, "ROWID": true, "SESSION": true, "SHARE": true, "START": true, "SUCCESSFUL": true,
	"SYNONYM": true, "TABLE": true, "UID": true, "USER": true, "VALIDATE": true, "VALUES": true, "VIEW": true,
	"MODE": true, "RESOURCE": true, "ONLINE": true, "OFFLINE": true, "INCREMENT": true, "MINUS": true,
	"PRIOR": true, "PUBLIC": true, "FILE": true, "RAW": true, "LONG": true, "UNIQUE": true, "INDEX": true,
}

func kolomCadangan(kolom []string) []string {
	var hasil []string
	for _, k := range kolom {
		if kataCadanganOracle[strings.ToUpper(k)] {
			hasil = append(hasil, k)
		}
	}
	return hasil
}

func TestMCRLNolKataCadanganOracle(t *testing.T) {
	semua := [][]string{KolomTahun, KolomKontrak, KolomReinsurer, KolomSecurity, KolomBusiness,
		{"ID", "NOTE", "FLAG", "CLIENTNAME", "STATUSACTIVE", "OLDID"}}
	for _, kolom := range semua {
		if c := kolomCadangan(kolom); len(c) > 0 {
			t.Errorf("kolom kata cadangan Oracle %v", c)
		}
	}
}

func TestMCRLAturanKataCadanganMenggigit(t *testing.T) {
	if len(kolomCadangan([]string{"ID", "level"})) != 1 {
		t.Error("LEVEL seharusnya tertangkap")
	}
}

// --- tiruan hanya diimpor uji --------------------------------------------------

func imporTiruan(jalur, isi string) bool {
	return !strings.HasSuffix(jalur, "_test.go") && !strings.Contains(jalur, "/tiruan/") &&
		strings.Contains(isi, `"nusantarare/modul/mastercontractretrolife/backend/tiruan"`)
}

func TestMCRLTiruanHanyaDiUji(t *testing.T) {
	for jalur, isi := range berkasModul(t, ".go") {
		if imporTiruan(jalur, isi) {
			t.Errorf("%s mengimpor paket tiruan - gudang di memori hanya untuk uji", jalur)
		}
	}
}

func TestMCRLAturanTiruanMenggigit(t *testing.T) {
	isi := `import "nusantarare/modul/mastercontractretrolife/backend/tiruan"`
	if !imporTiruan("../../backend/services/x.go", isi) {
		t.Error("impor tiruan dari kode produksi seharusnya tertangkap")
	}
	if imporTiruan("../../backend/services/x_test.go", isi) {
		t.Error("impor tiruan dari uji sah")
	}
}
