package repository

// Penjaga NAMA JUJUR modul Treaty Contract Out (penyimpangan sadar 8; AC 13,
// AC 40, AC 50 spec).
//
// Tiga nama korpus yang BERBOHONG - rule `_Old` yang justru terbaru, activity
// `testingKurs` yang justru jalur produksi, teks galat `"JSON_KLAIM"` sisa
// salin-tempel - tidak boleh merambat ke pengenal sistem baru. Nama rule
// sumbernya tetap boleh disebut sebagai BUKTI di komentar.
//
// ⚠️ Batasnya dinyatakan: yang dipindai adalah kode Go modul ini (berkas
// `tco_*.go` dan `rute_treaty_contract_out.go`, bukan test) sesudah komentar
// sebaris dan string literal dibuang. Sisi React dijaga
// `frontend/src/modul/treaty-contract-out/components/namaJujur.test.ts`.

import (
	"regexp"
	"strings"
	"testing"
)

var (
	polaStringGo   = regexp.MustCompile("`[^`]*`|\"(?:[^\"\\\\]|\\\\.)*\"")
	polaNamaBohong = regexp.MustCompile(`(\b|[a-z_])Old([^a-z]|$)|_OLD_|\bOLD\b|[Tt]esting|JSON_KLAIM`)
)

func berkasModulTCO(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	for nama, isi := range berkasGoSelainTest(t) {
		dasar := nama[strings.LastIndex(nama, "/")+1:]
		if strings.HasPrefix(dasar, "tco_") || dasar == "rute_treaty_contract_out.go" {
			out[nama] = isi
		}
	}
	if len(out) < 5 {
		t.Fatalf("hanya %d berkas modul terbaca; pembacanya yang rusak", len(out))
	}
	return out
}

func TestTCONolNamaBohongDiPengenal(t *testing.T) {
	for nama, isi := range berkasModulTCO(t) {
		kode := polaStringGo.ReplaceAllString(buangKomentarGo(isi), `""`)
		for _, baris := range strings.Split(kode, "\n") {
			if m := polaNamaBohong.FindString(baris); m != "" {
				t.Errorf("%s memuat pengenal %q di baris: %s\n"+
					"Nama yang berbohong di korpus tidak dibawa (penyimpangan sadar 8)",
					nama, m, strings.TrimSpace(baris))
			}
		}
	}
}

// Instrumen: pola benar-benar menangkap bentuk yang dilarang, dan tidak
// menuduh kata yang kebetulan memuatnya.
func TestPolaNamaBohongMasihMenggigit(t *testing.T) {
	for _, buruk := range []string{"BrowseReinsuranceTypeOld()", "kursTesting", "var Old int", "ERR_JSON_KLAIM", "x_OLD_y"} {
		if !polaNamaBohong.MatchString(buruk) {
			t.Errorf("%q lolos", buruk)
		}
	}
	for _, aman := range []string{"Golden", "hold", "Threshold", "bold := 1", "holder"} {
		if polaNamaBohong.MatchString(aman) {
			t.Errorf("%q dituduh", aman)
		}
	}
	// String literal dibuang sebelum pemindaian - nama rule sumber boleh
	// disebut sebagai bukti.
	if polaNamaBohong.MatchString(polaStringGo.ReplaceAllString(`x := "BrowseReinsuranceType_RD_Old"`, `""`)) {
		t.Error("string literal ikut dituduh")
	}
}
