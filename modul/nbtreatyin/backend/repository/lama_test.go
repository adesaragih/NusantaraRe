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
	set := regexp.MustCompile(`UPDATE %S SET ([^\n]*?) WHERE`).FindAllStringSubmatch(s, -1)
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
