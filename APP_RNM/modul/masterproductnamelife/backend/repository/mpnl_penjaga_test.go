package repository

// Penjaga MODUL Master Product Name Life (pola penjaga Retro Life ditiru,
// tidak diimpor).
//
//	P1  nol migrasi di rentang 140–179, nol DDL, nol tabel baru
//	    (tiket 01 ditangguhkan; slot menu 960–961 hanya UPDATE DIMIGRASI)
//	-   prosedur PEGA_M_PRODUCT_LIFE / PEGA_M_PRODUCT_INWARD_LIFE tidak dipanggil
//	-   nol COMMIT di teks kode
//	-   nol kata cadangan Oracle sebagai nama kolom
//	-   nol alamat layanan (skema URL, host, IP) dan nol pembacaan env var
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

// berkasModul - isi berkas modul berakhiran tertentu (docs dan node_modules dilewati).
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
	if len(hasil) == 0 {
		t.Fatal("nol berkas modul terbaca; pembacanya yang rusak")
	}
	return hasil
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

// kodeProduksi - berkas .go non-uji dan .sql modul, komentar Go dibuang.
func kodeProduksi(t *testing.T) map[string]string {
	t.Helper()
	hasil := map[string]string{}
	for jalur, isi := range berkasModul(t, ".go", ".sql") {
		if strings.HasSuffix(jalur, "_test.go") {
			continue
		}
		hasil[jalur] = buangKomentarGo(isi)
	}
	return hasil
}

// --- P1: nol migrasi di rentang, nol DDL --------------------------------------

var polaNomorMigrasi = regexp.MustCompile(`^(\d{3})_.*\.sql$`)

func pelanggaranMigrasi(nama string) string {
	m := polaNomorMigrasi.FindStringSubmatch(nama)
	if m == nil {
		return "berkas migrasi tanpa nomor tiga digit"
	}
	n, _ := strconv.Atoi(m[1])
	switch {
	case n >= 140 && n <= 179:
		return "P1: rentang 140-179 tetap kosong (tiket 01 ditangguhkan, nol DDL)"
	case n == 960 || n == 961:
		return ""
	default:
		return "nomor di luar rentang dan slot menu modul ini"
	}
}

func TestMPNLNolMigrasiDiRentang(t *testing.T) {
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

func TestMPNLAturanMigrasiMenggigit(t *testing.T) {
	for _, nama := range []string{"140_m_product_life.sql", "179_x_down.sql", "300_x.sql", "tanpa_nomor.sql"} {
		if pelanggaranMigrasi(nama) == "" {
			t.Errorf("%s seharusnya ditolak", nama)
		}
	}
	for _, nama := range []string{"960_menu_masterproductnamelife.sql", "960_menu_masterproductnamelife_down.sql"} {
		if p := pelanggaranMigrasi(nama); p != "" {
			t.Errorf("%s seharusnya sah: %s", nama, p)
		}
	}
}

var polaDDL = regexp.MustCompile(`(?i)\b(CREATE|ALTER|DROP|TRUNCATE)\s+(TABLE|SEQUENCE|INDEX|VIEW|SYNONYM)\b`)

func adaDDL(teks string) bool { return polaDDL.MatchString(teks) }

func TestMPNLNolDDL(t *testing.T) {
	for jalur, isi := range kodeProduksi(t) {
		if adaDDL(isi) {
			t.Errorf("%s memuat DDL - P1: nol DDL, nol tabel baru", jalur)
		}
	}
}

func TestMPNLAturanDDLMenggigit(t *testing.T) {
	for _, s := range []string{"CREATE TABLE {skema}.PRODUCT_LIFE (ID NUMBER)", "alter  table x add y", "DROP VIEW V"} {
		if !adaDDL(s) {
			t.Errorf("%q seharusnya terbaca DDL", s)
		}
	}
	if adaDDL("UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1'") {
		t.Error("UPDATE slot menu bukan DDL")
	}
}

// --- prosedur lama tidak dipanggil, nol COMMIT --------------------------------

// pelanggaranTulisLama - nama prosedur penulis lama, atau COMMIT, di kode.
func pelanggaranTulisLama(isi string) []string {
	var hasil []string
	for _, kata := range []string{"PEGA_M_PRODUCT_LIFE", "PEGA_M_PRODUCT_INWARD_LIFE"} {
		if strings.Contains(strings.ToUpper(isi), kata) {
			hasil = append(hasil, kata)
		}
	}
	if regexp.MustCompile(`\bCOMMIT\b`).MatchString(isi) {
		hasil = append(hasil, "COMMIT")
	}
	return hasil
}

func TestMPNLProsedurLamaTidakDipanggil(t *testing.T) {
	for jalur, isi := range kodeProduksi(t) {
		if p := pelanggaranTulisLama(isi); len(p) > 0 {
			t.Errorf("%s menyebut %v - prosedur lama tidak dipanggil, transaksi ditutup Go (spec AC 9)", jalur, p)
		}
	}
}

func TestMPNLAturanProsedurMenggigit(t *testing.T) {
	if len(pelanggaranTulisLama("q := `BEGIN POOLDATA.pega_m_product_life(:1); END;`")) == 0 {
		t.Error("pemanggilan prosedur seharusnya tertangkap")
	}
	if len(pelanggaranTulisLama("q := `UPDATE X SET A = 1; COMMIT`")) == 0 {
		t.Error("COMMIT seharusnya tertangkap")
	}
	if len(pelanggaranTulisLama("return tx.Commit()")) != 0 {
		t.Error("metode Commit() Go bukan teks SQL")
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

// semuaKolomFisik - kolom tabel yang disebut SQL modul ini (bertambah tiap paket).
func semuaKolomFisik() [][]string {
	return [][]string{KolomProduk, KolomInward}
}

func TestMPNLNolKataCadanganOracle(t *testing.T) {
	for _, kolom := range semuaKolomFisik() {
		if c := kolomCadangan(kolom); len(c) > 0 {
			t.Errorf("kolom kata cadangan Oracle %v", c)
		}
	}
}

func TestMPNLAturanKataCadanganMenggigit(t *testing.T) {
	if len(kolomCadangan([]string{"ID", "comment"})) != 1 {
		t.Error("COMMENT seharusnya tertangkap")
	}
}

// --- nol alamat layanan, nol env var -------------------------------------------

// polaAlamat - skema URL, `www.`, atau alamat IPv4. Dirakit dari potongan
// supaya berkas ini sendiri tidak memuat teks yang dicarinya.
var polaAlamat = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9+.-]*` + `:` + `//|\bwww` + `\.|\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)

func pelanggaranAlamat(isi string) []string {
	var hasil []string
	hasil = append(hasil, polaAlamat.FindAllString(isi, -1)...)
	if strings.Contains(isi, "os."+"Getenv") || strings.Contains(isi, "os."+"LookupEnv") {
		hasil = append(hasil, "env var")
	}
	return hasil
}

// TestMPNLNolAlamatLayanan - ADR-0013 / brief bab 1: alamat layanan luar
// (`ServiceGoogle`, `LinkService`, penampil `View Office Online`) TIDAK
// dipanggil dan TIDAK ditulis ke berkas apa pun - kode, uji, dan frontend.
func TestMPNLNolAlamatLayanan(t *testing.T) {
	for jalur, isi := range berkasModul(t, ".go", ".sql", ".ts", ".tsx", ".css") {
		if strings.HasSuffix(jalur, "mpnl_penjaga_test.go") {
			continue
		}
		if p := pelanggaranAlamat(isi); len(p) > 0 {
			t.Errorf("%s memuat alamat atau env var %v - alamat layanan di-resolve saat jalan, stub tidak memanggilnya", jalur, p)
		}
	}
}

func TestMPNLAturanAlamatMenggigit(t *testing.T) {
	for _, s := range []string{"const u = \"https:" + "//contoh.invalid/x\"", "host := \"10.0.0." + "1\"", "www" + ".contoh", "os.Get" + "env(\"URL\")"} {
		if len(pelanggaranAlamat(s)) == 0 {
			t.Errorf("%q seharusnya tertangkap", s)
		}
	}
	if len(pelanggaranAlamat("`dd/MM/yyyy` 1.5 versi 2.8.19")) != 0 {
		t.Error("teks biasa bukan alamat")
	}
}

// --- tiruan hanya diimpor uji ---------------------------------------------------

func imporTiruan(jalur, isi string) bool {
	return !strings.HasSuffix(jalur, "_test.go") && !strings.Contains(jalur, "/tiruan/") &&
		strings.Contains(isi, `"nusantarare/modul/masterproductnamelife/backend/tiruan"`)
}

func TestMPNLTiruanHanyaDiUji(t *testing.T) {
	for jalur, isi := range berkasModul(t, ".go") {
		if imporTiruan(jalur, isi) {
			t.Errorf("%s mengimpor paket tiruan - gudang di memori hanya untuk uji", jalur)
		}
	}
}

func TestMPNLAturanTiruanMenggigit(t *testing.T) {
	isi := `import "nusantarare/modul/masterproductnamelife/backend/tiruan"`
	if !imporTiruan("../../backend/services/x.go", isi) {
		t.Error("impor tiruan dari kode produksi seharusnya tertangkap")
	}
	if imporTiruan("../../backend/services/x_test.go", isi) {
		t.Error("impor tiruan dari uji sah")
	}
}
