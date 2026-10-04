package repository

// Uji murni bagian repository pemuat dokumen lama (tiket 22) - tanpa Oracle:
// bentuk SQL, dan penjaga "satu antarmuka penyimpanan" (spec-penyimpanan
// ID-3, AC 56): jalur pemuat TIDAK punya penulisan sendiri ke tabel katalog.

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestSQLPemuatLamaBerskemaTanpaCommit(t *testing.T) {
	const tabel = "UJI_SKEMA.JSON_POLIS"
	for nama, q := range map[string]string{
		"kunci":  sqlKunciJSONPolis(tabel),
		"lain":   sqlHitungJSONPolisLain(tabel),
		"baca":   sqlBacaJSONPolis(tabel),
		"idpega": sqlIDPegaKasus("UJI_SKEMA.T_GENERAL_POLIS"),
		"datar":  sqlSetelKolomDatarLama("UJI_SKEMA.T_GENERAL_POLIS"),
		// F3: penjaga dobel salinan SuggestList menurut IDPEGA.
		"usulan": sqlAdaUsulanIDPega("UJI_SKEMA.HISTORYAKSEPTASIPRODUCTION"),
	} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
		if !strings.Contains(q, "UJI_SKEMA.") {
			t.Errorf("%s tanpa skema eksplisit (AC 47): %s", nama, q)
		}
	}
	// Generasi endorsemen bukan lingkup NB (edmtreatyin tiket 10).
	if !strings.Contains(sqlKunciJSONPolis(tabel), "'0'") {
		t.Error("kunci JSON_POLIS harus menyaring PRODKE 0")
	}
	// Penjaga dobel F3: dicari per IDPEGA (satu penampung), baca saja.
	q := strings.ToUpper(sqlAdaUsulanIDPega("UJI_SKEMA.HISTORYAKSEPTASIPRODUCTION"))
	if !strings.Contains(q, "WHERE IDPEGA = :1") || strings.Contains(q, ":2") || !strings.HasPrefix(strings.TrimSpace(q), "SELECT") {
		t.Errorf("penjaga dobel salinan usulan: %s", q)
	}
}

// TestPemuatTanpaJalurTulisTerpisah - AC 56: isi 8 tabel diagram ditulis
// pemuat HANYA lewat SimpanHalaman / SisipKasus / SetelNomorPolis /
// TutupKasus (antarmuka yang sama dengan jalur biasa). Berkas pemuat hanya
// boleh MEMBACA JSON_POLIS dan memperbarui kolom datar json_polis (ID-21)
// yang tidak pernah ditulis jalur biasa.
func TestPemuatTanpaJalurTulisTerpisah(t *testing.T) {
	isi, err := os.ReadFile("lama.go")
	if err != nil {
		t.Fatal(err)
	}
	s := strings.ToUpper(string(isi))
	for _, kata := range []string{"INSERT ", "DELETE ", "MERGE ", "JSON_POLIS SET", "CREATE "} {
		if strings.Contains(s, kata) {
			t.Errorf("lama.go memuat %q - jalur tulis terpisah dilarang (AC 56)", kata)
		}
	}
	set := regexp.MustCompile(`UPDATE %S(?: G)? SET ([^\n]*?) WHERE`).FindAllStringSubmatch(s, -1)
	if len(set) != 1 {
		t.Fatalf("harap tepat satu UPDATE (kolom datar json_polis), dapat %d", len(set))
	}
	tanpaKurung := regexp.MustCompile(`\([^)]*\)`).ReplaceAllString(set[0][1], "")
	for _, bagian := range strings.Split(tanpaKurung, ",") {
		kol := strings.TrimSpace(strings.SplitN(bagian, "=", 2)[0])
		switch kol {
		case "IDPEGA", "NOENDORS", "TGL_INPUT", "USERNAME":
		default:
			t.Errorf("UPDATE pemuat menyentuh %s - di luar kolom datar json_polis (ID-21)", kol)
		}
	}
}

// Diagram sheet NB Treaty In Prop F17: "baris generasi lampau TIDAK BOLEH
// disunting - itulah pembekuan OldData (P58)". SETIAP UPDATE atas
// T_GENERAL_POLIS yang SQL-nya dirakit fungsi tersendiri hanya menyentuh
// generasi terbuka (tanpa penerus yang OLD_POLIS_ID-nya menunjuknya);
// `tulisInduk` dan `SetelNomorPolis` merakitnya di badan fungsi dengan
// `syaratTerbuka` yang sama.
func TestUbahGeneralPolisHanyaGenerasiTerbuka(t *testing.T) {
	const tabel = "UJI_SKEMA.T_GENERAL_POLIS"
	penjaga := "NOT EXISTS (SELECT 1 FROM UJI_SKEMA.T_GENERAL_POLIS s WHERE s.OLD_POLIS_ID = g.ID)"
	for nama, q := range map[string]string{
		"pindah posisi":          sqlPindahGenerasi(tabel),
		"kolom datar json_polis": sqlSetelKolomDatarLama(tabel),
	} {
		if !strings.Contains(q, tabel+" g ") || !strings.Contains(q, penjaga) {
			t.Errorf("%s: UPDATE T_GENERAL_POLIS tanpa penjaga generasi terbuka (F17):\n%s", nama, q)
		}
	}
}
