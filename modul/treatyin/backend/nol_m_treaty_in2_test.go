package backend_test

// ⛔⛔ PENJAGA PENCABUTAN `M_TREATY_IN2`, 5 Oktober 2026.
//
// Keputusan pemilik proses: `M_TREATY_IN.JSONDATA` adalah SATU-SATUNYA sumber
// isi kontrak, dan `M_TREATY_IN2` punya **nol jalur baca**.
//
// ⚠️ Larangan yang hanya ditulis di komentar dilanggar dalam enam bulan.
// Berkas ini mewujudkannya sebagai uji: ia menyapu SELURUH berkas Go
// repositori ini dan merah bila ada yang memasang kuerinya kembali.
//
// ⛔ Tabelnya TETAP ADA di Oracle dan TETAP terlarang disentuh — nol `DROP`,
// nol `INSERT`, nol `UPDATE`, nol migrasi. Yang dicabut jalur bacanya di
// aplikasi, bukan tabelnya; `TestWarisanHanyaDibaca` tetap menamainya, dan
// uji ini TIDAK menggantikan penjaga itu.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// polaKueri mencari penyebutan `M_TREATY_IN2` di dalam konteks SQL.
//
// ⚠️ Sengaja LEBAR: apa pun yang menyebut namanya di luar daftar-izin di
// bawah dianggap pelanggaran. Pola sempit yang hanya mencari `FROM
// M_TREATY_IN2` akan melewatkan `JOIN`, nama tabel yang dirakit dari
// potongan, dan konstanta bernama.
var polaKueri = regexp.MustCompile(`M_TREATY_IN2`)

// diizinkan — berkas yang BOLEH menyebut namanya, masing-masing dengan
// sebabnya. Daftar ini sengaja pendek dan sengaja eksplisit.
var diizinkan = map[string]string{
	// Berkas ini sendiri — ia harus menyebut nama yang dilarangnya.
	filepath.Join("modul", "treatyin", "backend", "nol_m_treaty_in2_test.go"): "penjaga ini sendiri",

	// Penjaga yang menamainya sebagai tabel WARISAN yang tidak boleh
	// ditulis. Perannya berlawanan dengan jalur baca: ia melarang.
	filepath.Join("modul", "treatyin", "backend", "migrasi_invarian_test.go"):   "daftar tabel warisan baca-saja",
	filepath.Join("modul", "treatyin", "backend", "migrasi_pendaratan_test.go"): "daftar tabel warisan tanpa kunci asing",

	// Komentar yang menyatakan tabel ini TIDAK disentuh. Briefing pencabutan
	// menyatakan komentar itu tetap benar dan boleh tinggal.
	filepath.Join("modul", "treatyin", "backend", "repository", "pendaratan_muat.go"):               "komentar: yang TIDAK dimuat",
	filepath.Join("modul", "treatyin", "backend", "repository", "warisan_layer_dokumen.go"):         "komentar: riwayat pencabutan",
	filepath.Join("modul", "treatyin", "backend", "models", "warisan_layer.go"):                     "komentar: riwayat pencabutan",
	filepath.Join("modul", "treatyin", "backend", "models", "warisan_kontrak.go"):                   "komentar: riwayat pencabutan",
	filepath.Join("modul", "treatyin", "backend", "repository", "warisan_layer_dokumen_test.go"):    "komentar: uji lama yang dicabut, beserta sebabnya",
	filepath.Join("modul", "treatyin", "backend", "repository", "warisan_layer_dokumen_db_test.go"): "uji yang MEMBUKTIKAN pencabutannya — kueri di uji, bukan di jalur baca",
}

// ⭐ NOL berkas Go menyebut `M_TREATY_IN2` di luar daftar-izin.
//
// ⛔ Daftar-izinnya TIDAK boleh tumbuh tanpa alasan. Uji kedua di bawah
// menolak entri yang tidak terpakai, sehingga daftar yang ditambah
// "sekadar untuk lewat" akan terlihat begitu berkasnya berhenti menyebutnya.
func TestNolKueriMTreatyIn2(t *testing.T) {
	akar := akarRepo(t)
	var pelanggar []string

	err := filepath.Walk(akar, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			switch info.Name() {
			case "node_modules", ".git", "unggahan":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(p, ".go") {
			return nil
		}
		rel, _ := filepath.Rel(akar, p)
		if _, ok := diizinkan[rel]; ok {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		if polaKueri.Match(b) {
			pelanggar = append(pelanggar, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range pelanggar {
		t.Errorf("%s menyebut M_TREATY_IN2.\n"+
			"⛔ Tabel itu DICABUT sebagai sumber 5 Oktober 2026 — satu-satunya sumber\n"+
			"   isi kontrak adalah `M_TREATY_IN.JSONDATA`. Baris layer terurai di\n"+
			"   `repository/warisan_layer_dokumen.go`.\n"+
			"   Bila penyebutan ini memang sah (komentar, daftar larangan), tambahkan\n"+
			"   berkasnya ke `diizinkan` DI SINI beserta sebabnya.", p)
	}
}

// ⛔ Daftar-izin yang TIDAK TERPAKAI juga ditolak.
//
// Entri basi membuat daftar itu tumbuh diam-diam sampai larangannya tidak
// menjaga apa pun lagi.
func TestDaftarIzinMTreatyIn2TidakBasi(t *testing.T) {
	akar := akarRepo(t)
	for rel, sebab := range diizinkan {
		b, err := os.ReadFile(filepath.Join(akar, rel))
		if err != nil {
			t.Errorf("daftar-izin memuat %s (%s) yang tidak ada lagi — cabut entrinya", rel, sebab)
			continue
		}
		if !polaKueri.Match(b) {
			t.Errorf("daftar-izin memuat %s (%s) yang TIDAK lagi menyebut M_TREATY_IN2 — "+
				"cabut entrinya supaya daftarnya tetap menjaga sesuatu", rel, sebab)
		}
	}
}

// akarRepo menaiki pohon sampai menemukan `go.mod`.
func akarRepo(t *testing.T) string {
	t.Helper()
	d, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
			return d
		}
		d = filepath.Dir(d)
	}
	t.Fatal("go.mod tidak ditemukan dari direktori uji")
	return ""
}
