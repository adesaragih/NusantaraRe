package repository

// Uji murni bagian repository pemuat dokumen lama endorsemen (tiket EDM 10 / 09) - tanpa Oracle: bentuk SQL dan
// penjaga "satu antarmuka penyimpanan" (spec-penyimpanan ID-3, AC 44): berkas pemuat TIDAK punya penulisan sendiri
// ke tabel katalog maupun proyeksi selisih; penanda migrasi hanya menyentuh baris SUMBER 'PEGA' (AC 43).

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

const ujiSkema = "UJI_SKEMA."

func TestSQLPemuatEDMBerskemaTanpaCommit(t *testing.T) {
	gen := ujiSkema + "T_GENERAL_POLIS_TREATY"
	for nama, q := range map[string]string{
		"kunci":          sqlPmKunciJSONPolisEDM(ujiSkema + "JSON_POLIS"),
		"baca":           sqlPmBacaJSONPolis(ujiSkema + "JSON_POLIS"),
		"baca NB":        sqlPmBacaJSONPolisNB(ujiSkema + "JSON_POLIS"),
		"kunci generasi": sqlPmKunciGenerasi(gen),
		"datar":          sqlPmSetelKolomDatarLama(gen),
		"penanda 361":    sqlPmSetelPenanda(ujiSkema+"T_POLIS_DIFFERENCE_SPREADING", ujiSkema+"T_POLIS_DIFFERENCE", true),
		"penanda 363":    sqlPmSetelPenanda(ujiSkema+"T_POLIS_XOL_LAYER_DIFFERENCE", ujiSkema+"T_POLIS_DIFFERENCE", false),
		"riwayat":        sqlPmAdaRiwayatIDPega(ujiSkema+"HISTORYAKSEPTASIPRODUCTION", ujiSkema+"HISTORYAKSEPTASIPEGA"),
	} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
		// Setiap rujukan tabel berskema eksplisit (AC 46).
		for _, m := range regexp.MustCompile(`(?i)\b(FROM|UPDATE|INTO|JOIN)\s+([A-Z_.]+)`).FindAllStringSubmatch(q, -1) {
			if !strings.HasPrefix(m[2], ujiSkema) {
				t.Errorf("%s: %s %s tanpa skema eksplisit:\n%s", nama, m[1], m[2], q)
			}
		}
	}
	// Generasi NB (PRODKE '0' / kosong) bukan lingkup pemuat EDM: kunci menyaringnya keluar.
	if q := sqlPmKunciJSONPolisEDM(ujiSkema + "JSON_POLIS"); !strings.Contains(q, "PRODKE IS NOT NULL") || !strings.Contains(q, "<> '0'") {
		t.Errorf("kunci JSON_POLIS endorsemen tanpa saringan PRODKE: %s", q)
	}
}

// TestPenandaHanyaBarisSumberPega - AC 43 / ID-37: UPDATE penanda menyaring induk SUMBER di SQL; penampungnya urut
// dan tidak berulang (penjaga inti `TestNolPenampungBerulangDiSQLBerpembatasBaris`).
func TestPenandaHanyaBarisSumberPega(t *testing.T) {
	for berlapis, harap := range map[bool][]string{
		true:  {"PASANGAN_BERGESER = :1, RUMUS_BERLAPIS = :2", "a.NOURUT = :3", "d.POLIS_ID = :4", "d.SUMBER = :5"},
		false: {"PASANGAN_BERGESER = :1", "a.NOURUT = :2", "d.POLIS_ID = :3", "d.SUMBER = :4"},
	} {
		q := sqlPmSetelPenanda("UJI.ANAK", "UJI.INDUK", berlapis)
		for _, s := range harap {
			if !strings.Contains(q, s) {
				t.Errorf("berlapis=%v: tanpa %q:\n%s", berlapis, s, q)
			}
		}
		if !berlapis && strings.Contains(q, "RUMUS_BERLAPIS") {
			t.Error("363 tidak punya kolom RUMUS_BERLAPIS")
		}
	}
	isi, err := os.ReadFile("lama_edm.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), "args = append(args, i+1, polisID, models.SumberPega)") {
		t.Error("SetelPenandaMigrasi tidak lagi mengikat SUMBER 'PEGA' pada penyaring induk (AC 43)")
	}
}

// TestPemuatEDMTanpaJalurTulisTerpisah - AC 44: generasi dan proyeksi selisih ditulis pemuat HANYA lewat
// SisipKasus / SimpanHalaman / SetelNomorPolisSelesai / SimpanSelisih / CatatUsulan / TutupKasus. Berkas pemuat
// hanya boleh MEMBACA JSON_POLIS dan memperbarui kolom datar json_polis (TGL_INPUT, USERNAME) serta kedua penanda.
func TestPemuatEDMTanpaJalurTulisTerpisah(t *testing.T) {
	isi, err := os.ReadFile("lama_edm.go")
	if err != nil {
		t.Fatal(err)
	}
	s := strings.ToUpper(string(isi))
	for _, kata := range []string{"INSERT ", "DELETE ", "MERGE ", "JSON_POLIS SET", "CREATE ", "DROP ", "TRUNCATE "} {
		if strings.Contains(s, kata) {
			t.Errorf("lama_edm.go memuat %q - jalur tulis terpisah dilarang (AC 44)", kata)
		}
	}
	set := regexp.MustCompile(`(?m)UPDATE %S (?:G|A) SET ([^\n]*?)(?: WHERE|$)`).FindAllStringSubmatch(s, -1)
	if len(set) != 2 {
		t.Fatalf("harap tepat dua UPDATE (kolom datar json_polis, penanda migrasi), dapat %d", len(set))
	}
	boleh := map[string]bool{"TGL_INPUT": true, "USERNAME": true, "PASANGAN_BERGESER": true, "%S": true}
	for _, x := range set {
		tanpaKurung := regexp.MustCompile(`\([^)]*\)`).ReplaceAllString(x[1], "")
		for _, bagian := range strings.Split(tanpaKurung, ",") {
			kol := strings.TrimSpace(strings.SplitN(bagian, "=", 2)[0])
			if !boleh[kol] {
				t.Errorf("UPDATE pemuat menyentuh %q - di luar kolom datar json_polis dan penanda migrasi", kol)
			}
		}
	}
	if !strings.Contains(s, `", RUMUS_BERLAPIS = :2"`) {
		t.Error("penanda RUMUS_BERLAPIS tidak lagi ditulis di 361/362")
	}
}

// Generasi lampau TIDAK BOLEH disunting (ID-13): UPDATE kolom datar hanya menyentuh generasi terbuka ber-PRODKE >= 1.
func TestUbahGeneralPolisPemuatHanyaGenerasiTerbuka(t *testing.T) {
	const tabel = ujiSkema + "T_GENERAL_POLIS_TREATY"
	q := sqlPmSetelKolomDatarLama(tabel)
	if !strings.Contains(q, tabel+" g ") || !strings.Contains(q, syaratTerbuka(tabel)) || !strings.Contains(q, "g.PRODKE >= 1") {
		t.Errorf("UPDATE kolom datar tanpa penjaga generasi terbuka endorsemen:\n%s", q)
	}
}
