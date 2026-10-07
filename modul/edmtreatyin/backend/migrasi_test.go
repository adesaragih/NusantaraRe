package backend

// Penjaga migrasi modul EDM Treaty In: empat tabel proyeksi selisih (360-363) SEPAKAT dengan katalog
// `models/katalog_selisih.go` dan dokumen `docs/STRUKTUR-TABEL-EDM-TREATY-IN.md`; nol tabel dasar dibuat ulang
// (tabel generasi milik nbtreatyin); nol prosedur / COMMIT; slot menu 970 hanya satu UPDATE DIMIGRASI.

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"nusantarare/modul/edmtreatyin/backend/models"
)

func bacaMigrasi(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := fs.WalkDir(berkasMigrasi, "migrations", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := berkasMigrasi.ReadFile(p)
		if err != nil {
			return err
		}
		out[filepath.Base(p)] = string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

var reCreate = regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.(\w+) \((.*?)\n\)\n/`)

// kolomDDL - tabel -> kolom (urut) dari seluruh migrasi maju.
func kolomDDL(t *testing.T) map[string][]string {
	t.Helper()
	out := map[string][]string{}
	for nama, isi := range bacaMigrasi(t) {
		if strings.HasSuffix(nama, "_down.sql") {
			continue
		}
		for _, m := range reCreate.FindAllStringSubmatch(isi, -1) {
			for _, baris := range strings.Split(m[2], "\n") {
				b := strings.TrimSpace(baris)
				if b == "" || strings.HasPrefix(b, "CONSTRAINT") || strings.HasPrefix(b, "--") {
					continue
				}
				out[m[1]] = append(out[m[1]], strings.Fields(b)[0])
			}
		}
	}
	return out
}

func TestMigrasiHanyaEmpatTabelProyeksiSelisih(t *testing.T) {
	ddl := kolomDDL(t)
	var nama []string
	for n := range ddl {
		nama = append(nama, n)
	}
	sort.Strings(nama)
	harap := []string{"T_POLIS_DIFFERENCE", "T_POLIS_DIFFERENCE_INSTALMENT", "T_POLIS_DIFFERENCE_SPREADING", "T_POLIS_XOL_LAYER_DIFFERENCE"}
	if strings.Join(nama, ",") != strings.Join(harap, ",") {
		t.Fatalf("tabel dibuat = %v, harap %v (AC 9, 29, 30: nol tabel OldData, nol induk XOL, nol rincian angsuran)", nama, harap)
	}
}

// kolomKunci - kolom yang ditulis repository di luar katalog.
var kolomKunci = map[string][]string{
	"T_POLIS_DIFFERENCE":            {"ID", "POLIS_ID", "NOPOLIS", "PRODKE", "EDM_NO", "IDPEGA", "SUMBER"},
	"T_POLIS_DIFFERENCE_SPREADING":  {"ID", "DIFFERENCE_ID", "NOURUT", "PASANGAN_BERGESER", "RUMUS_BERLAPIS"},
	"T_POLIS_DIFFERENCE_INSTALMENT": {"ID", "DIFFERENCE_ID", "NOURUT", "PASANGAN_BERGESER", "RUMUS_BERLAPIS"},
	"T_POLIS_XOL_LAYER_DIFFERENCE":  {"ID", "DIFFERENCE_ID", "NOURUT", "PASANGAN_BERGESER"},
}

func TestKatalogSelisihSepakatDenganDDL(t *testing.T) {
	ddl := kolomDDL(t)
	for _, tb := range models.SemuaTabelSelisih {
		harap := map[string]bool{}
		for _, k := range kolomKunci[tb.Nama] {
			harap[k] = true
		}
		for _, k := range tb.Kolom {
			harap[k.Kolom] = true
		}
		ada := map[string]bool{}
		for _, k := range ddl[tb.Nama] {
			ada[k] = true
			if !harap[k] {
				t.Errorf("%s: kolom DDL %s tidak ada di katalog/kunci", tb.Nama, k)
			}
		}
		for k := range harap {
			if !ada[k] {
				t.Errorf("%s: kolom katalog %s tidak dibuat DDL", tb.Nama, k)
			}
		}
	}
}

func TestStrukturMemuatSetiapKolomDDL(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("..", "docs", "STRUKTUR-TABEL-EDM-TREATY-IN.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	for tabel, kolom := range kolomDDL(t) {
		bab := strings.Index(doc, "## "+tabel+"\n")
		if bab < 0 {
			t.Fatalf("STRUKTUR tanpa bab %s", tabel)
		}
		isi := doc[bab+len(tabel)+4:]
		if j := strings.Index(isi, "\n## "); j >= 0 {
			isi = isi[:j]
		}
		for _, k := range kolom {
			if !strings.Contains(isi, "| `"+k+"` |") {
				t.Errorf("STRUKTUR %s tanpa kolom %s", tabel, k)
			}
		}
	}
}

func TestUangSelisihBertipeNumber3810TanpaFloat(t *testing.T) {
	for nama, isi := range bacaMigrasi(t) {
		for _, larang := range []string{"FLOAT", "BINARY_DOUBLE", "BINARY_FLOAT", "COMMIT", "PROCEDURE", "TRIGGER"} {
			if strings.Contains(strings.ToUpper(isi), larang) {
				t.Errorf("%s memuat %s (AC 48; migrasi tanpa prosedur / COMMIT)", nama, larang)
			}
		}
	}
	for _, tb := range models.SemuaTabelSelisih {
		for _, k := range tb.Kolom {
			if !k.Golongan.Desimal() {
				continue
			}
			for nama, isi := range bacaMigrasi(t) {
				re := regexp.MustCompile(`\n\s+` + k.Kolom + `\s+(\S+),?\n`)
				if m := re.FindStringSubmatch(isi); m != nil && strings.Contains(isi, "TABLE {skema}."+tb.Nama+" (") &&
					strings.TrimSuffix(m[1], ",") != "NUMBER(38,10)" {
					t.Errorf("%s %s.%s = %s, harap NUMBER(38,10)", nama, tb.Nama, k.Kolom, m[1])
				}
			}
		}
	}
}

func TestTabelGenerasiNBTidakDibuatUlang(t *testing.T) {
	for nama, isi := range bacaMigrasi(t) {
		for _, t0 := range []string{"T_GENERAL_POLIS_TREATY (", "T_WORK_POLIS (", "T_POLIS_SPREADING (", "T_POLIS_INSTALMENT ("} {
			if strings.Contains(isi, "CREATE TABLE {skema}."+t0) {
				t.Errorf("%s membuat tabel milik modul lain %s (AC 54)", nama, t0)
			}
		}
	}
}

func TestSlotMenu970HanyaUpdateDimigrasi(t *testing.T) {
	m := bacaMigrasi(t)
	naik, turun := m["970_menu_edmtreatyin.sql"], m["970_menu_edmtreatyin_down.sql"]
	if !strings.Contains(naik, "UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1'") || strings.Contains(naik, "INSERT INTO") {
		t.Fatalf("970 harus satu UPDATE DIMIGRASI '1', nol INSERT:\n%s", naik)
	}
	if !strings.Contains(naik, "WHERE KODE = 'edmtreatyin'") || !strings.Contains(turun, "DIMIGRASI = '0'") {
		t.Fatalf("970 / 970_down tidak menunjuk baris edmtreatyin")
	}
}
