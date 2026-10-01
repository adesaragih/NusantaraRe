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
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/models"
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
	return append([][]string{KolomProduk, KolomInward, KolomLampiran, KolomObjek, KolomOutbox,
		KolomKontrakTreaty, KolomTahunTreaty, KolomRate}, KolomMaster()...)
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

// --- master dibaca saja; sumber rate tidak dibaca ------------------------------

var polaTulis = regexp.MustCompile(`(?i)\b(INSERT\s+INTO|UPDATE|DELETE\s+FROM|DELETE|MERGE\s+INTO)\b`)

// tulisMaster - fungsi yang menulis DAN menyebut objek master.
func tulisMaster(isi string) []string {
	var hasil []string
	for _, blok := range strings.Split(isi, "\nfunc ") {
		if !polaTulis.MatchString(blok) {
			continue
		}
		for _, m := range DaftarMasterDibacaSaja {
			if regexp.MustCompile(`\b` + m + `\b`).MatchString(blok) {
				hasil = append(hasil, m)
			}
		}
	}
	return hasil
}

func TestMPNLMasterDibacaSaja(t *testing.T) {
	for jalur, isi := range kodeProduksi(t) {
		if !strings.Contains(jalur, "/repository/") {
			continue
		}
		if m := tulisMaster(isi); len(m) > 0 {
			t.Errorf("%s menulis master %v - master rujukan dibaca saja", jalur, m)
		}
	}
}

func TestMPNLAturanMasterMenggigit(t *testing.T) {
	isi := "\nfunc sqlX(t string) string {\n\treturn fmt.Sprintf(`UPDATE %s SET NAME = :1`, t) // CLIENT\n}\n"
	if len(tulisMaster(isi)) == 0 {
		t.Error("UPDATE yang menyebut CLIENT seharusnya tertangkap")
	}
}

// TestMPNLRateDibacaKolomRDSaja - K1 keputusan work owner 01-10-2026 (OQ-MPNL-03): kedua view rate dibaca
// SAJA, kolom RD saja (`BrowseRateLifeSummary` b692–b717, `BrowseRateLife_RD` b747–b791), nol `SELECT *`,
// nol `JSONDATA`; `View Rate` berkunci `IDUSEDBY`, urut dan batas RD.
func TestMPNLRateDibacaKolomRDSaja(t *testing.T) {
	kolomRD := map[string][]string{
		MasterRIRate: {"ID", "USEDBY", "OPERATORID", "MODIFIEDDATE", "TYPE"},
		MasterRate:   {"ID", "IDUSEDBY", "USEDBY", "GENDER", "CONTRACT", "AGE", "RATE", "TYPE"},
	}
	s := sumberMaster[models.MasterRIRate]
	kasus := []struct{ objek, q string }{
		{MasterRIRate, s.sqlCari("S.V")},
		{MasterRIRate, s.sqlAmbil("S.V")},
		{MasterRate, sqlDaftarRate("S.V")},
	}
	for _, k := range kasus {
		q := strings.Join(strings.Fields(k.q), " ")
		if strings.Contains(q, "*") || strings.Contains(strings.ToUpper(q), "JSONDATA") {
			t.Errorf("%s: SELECT * / JSONDATA dilarang: %s", k.objek, q)
		}
		for _, kol := range strings.Split(q[len("SELECT "):strings.Index(q, " FROM ")], ",") {
			ada := false
			for _, r := range kolomRD[k.objek] {
				ada = ada || r == strings.TrimSpace(kol)
			}
			if !ada {
				t.Errorf("%s: kolom %s bukan kolom RD %v", k.objek, kol, kolomRD[k.objek])
			}
		}
	}
	q := strings.Join(strings.Fields(sqlDaftarRate("S.V")), " ")
	for _, w := range []string{"WHERE IDUSEDBY = :1", "ORDER BY ID DESC, RATE ASC", "FETCH FIRST 501 ROWS ONLY"} {
		if !strings.Contains(q, w) {
			t.Errorf("View Rate tanpa %q: %s", w, q)
		}
	}
	if BatasRate != 500 || strings.Join(KolomRate, ",") != "ID,USEDBY,GENDER,CONTRACT,AGE,RATE" {
		t.Errorf("BatasRate %d / KolomRate %v", BatasRate, KolomRate)
	}
}

// TestMPNLSetiapSQLMasterAdalahSelect - setiap SQL yang disusun untuk objek di DaftarMasterDibacaSaja (ketujuh
// pemilih, PLAN LIST, kedua view rate) berawal SELECT dan lolos penjaga runtime periksaBacaSaja.
func TestMPNLSetiapSQLMasterAdalahSelect(t *testing.T) {
	kasus := map[string]string{"plan cari": sqlCariPlan("S.V"), "plan ambil": sqlAmbilPlan("S.V"), "rate": sqlDaftarRate("S.V")}
	objek := map[string]string{"plan cari": MasterJenisPlan, "plan ambil": MasterJenisPlan, "rate": MasterRate}
	for jenis, s := range sumberMaster {
		kasus[string(jenis)+" cari"], objek[string(jenis)+" cari"] = s.sqlCari("S.V"), s.objek
		kasus[string(jenis)+" ambil"], objek[string(jenis)+" ambil"] = s.sqlAmbil("S.V"), s.objek
	}
	if len(kasus) < 17 {
		t.Fatalf("hanya %d SQL master terbaca - pembacanya yang rusak", len(kasus))
	}
	for nama, q := range kasus {
		if err := periksaBacaSaja(objek[nama], q); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
	}
}

// Uji gigit K1: tulisan ke view rate (dan master lain) ditolak SEBELUM sampai ke Oracle; tabel produk
// milik modul ini tetap boleh ditulis.
func TestMPNLPeriksaBacaSajaMenolakTulisanKeView(t *testing.T) {
	for _, objek := range []string{MasterRate, MasterRIRate, MasterAgent} {
		for _, q := range []string{
			"UPDATE S.V SET RATE = :1 WHERE ID = :2",
			"insert into S.V (ID) values (:1)",
			"  DELETE FROM S.V WHERE ID = :1",
			"MERGE INTO S.V USING DUAL ON (1 = 1) WHEN MATCHED THEN UPDATE SET RATE = :1",
		} {
			if err := periksaBacaSaja(objek, q); !errors.Is(err, ErrMasterBacaSaja) {
				t.Errorf("%s: %q harus ditolak: %v", objek, q, err)
			}
		}
	}
	if err := periksaBacaSaja(TabelProduk, "UPDATE S.T SET JSONDATA = :1 WHERE ID = :2"); err != nil {
		t.Errorf("tabel produk modul ini boleh ditulis: %v", err)
	}
}
