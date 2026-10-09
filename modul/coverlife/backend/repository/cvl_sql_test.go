package repository

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func satuBaris(q string) string { return strings.Join(strings.Fields(q), " ") }

func memuat(t *testing.T, nama, q string, bagian ...string) {
	t.Helper()
	s := satuBaris(q)
	for _, b := range bagian {
		if !strings.Contains(s, b) {
			t.Errorf("%s tanpa %q:\n%s", nama, b, s)
		}
	}
}

// Grid BrowseCoverLife_RD: kolom ID, COVER (+ NOTE untuk Edit), tanpa saring (b4080), ID ANGKA menaik tetap
// (b3967 / b3972), 50 lewat bind; Cover kembar tanpa beda huruf; tulis = INSERT / UPDATE kolom bernama.
func TestSql(t *testing.T) {
	const t1 = "S.M_COVER_LIFE"
	memuat(t, "daftar", SqlDaftar(t1), "SELECT ID, COVER, NOTE FROM S.M_COVER_LIFE ORDER BY",
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY")
	if strings.Contains(SqlDaftar(t1), "WHERE") || strings.Contains(SqlJumlah(t1), "WHERE") {
		t.Error("grid tanpa saring (pyGridFiltering false b4080)")
	}
	memuat(t, "jumlah", SqlJumlah(t1), "SELECT COUNT(*) FROM S.M_COVER_LIFE")
	memuat(t, "ambil", SqlAmbil(t1), "SELECT ID, COVER, NOTE FROM S.M_COVER_LIFE WHERE ID = :1")
	memuat(t, "pemakai", SqlPemakaiNama(t1), "WHERE UPPER(TRIM(COVER)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID")
	memuat(t, "nomor", SqlNomorBaru("S.M_COVER_LIFE_SEQ"), "SELECT TO_CHAR(S.M_COVER_LIFE_SEQ.NEXTVAL) FROM DUAL")
	memuat(t, "ada", SqlAdaID(t1), "SELECT COUNT(*) FROM S.M_COVER_LIFE WHERE ID = :1")
	memuat(t, "sisip", SqlSisip(t1), "INSERT INTO S.M_COVER_LIFE (ID, COVER, NOTE) VALUES (:1, :2, :3)")
	memuat(t, "ubah", SqlUbah(t1), "UPDATE S.M_COVER_LIFE SET COVER = :1, NOTE = :2 WHERE ID = :3")
}

// Lapis penjaga: tulis hanya ke M_COVER_LIFE; objek milik aplikasi lain (COVERAGE*, COVERNOTE*) ditolak.
func TestPeriksaTulis(t *testing.T) {
	if !slices.Equal(DaftarTabelDitulis, []string{"M_COVER_LIFE"}) {
		t.Errorf("tabel ditulis %v", DaftarTabelDitulis)
	}
	for _, objek := range []string{"COVERAGE", "COVERAGE_FACIN", "COVERAGETRAVEL", "COVERNOTE", "M_COVER_LIFE_SEQ"} {
		if err := PeriksaTulis(objek, "UPDATE X SET A = 1"); !errors.Is(err, ErrBacaSaja) {
			t.Errorf("%s ditulis: %v", objek, err)
		}
	}
	if err := PeriksaTulis("DUAL", "SELECT 1 FROM DUAL"); err != nil {
		t.Errorf("SELECT: %v", err)
	}
}

// ORA-00001 (PK SYS_C009203) = ErrKembar; tabel / sequence tidak ada = ErrBelumAda.
func TestBungkus(t *testing.T) {
	if err := Bungkus(errors.New("ORA-00001: unique constraint violated"), "x"); !errors.Is(err, ErrKembar) {
		t.Errorf("ORA-00001 %v", err)
	}
	for _, kode := range []string{"ORA-00942", "ORA-00904", "ORA-02289"} {
		if err := Bungkus(errors.New(kode+": x"), "x"); !errors.Is(err, ErrBelumAda) {
			t.Errorf("%s %v", kode, err)
		}
	}
}

// polaViewLama - nama view lama sebagai KATA UTUH (bukan bagian `M_COVER_LIFE` / `M_COVER_LIFE_SEQ`).
var polaViewLama = regexp.MustCompile(`(^|[^A-Za-z0-9_])COVER_LIFE([^A-Za-z0-9_]|$)`)

// C1: view COVER_LIFE dibuang 086 - TIDAK ADA kode modul ini (Go backend selain migrasi, TS / TSX frontend; bukan uji,
// bukan komentar) yang menyebut COVER_LIFE sebagai tabel / view. Hanya M_COVER_LIFE (dan sequence-nya). Juga nol JSON
// dan nol DELETE di baris kode produksi repository; kolom = kolom view lama.
func TestKodeTidakMenyebutViewLama(t *testing.T) {
	if !slices.Equal(KolomTabel, []string{"ID", "COVER", "NOTE"}) {
		t.Errorf("kolom %v", KolomTabel)
	}
	if !polaViewLama.MatchString("SELECT ID FROM S.COVER_LIFE") || polaViewLama.MatchString("FROM S.M_COVER_LIFE") ||
		polaViewLama.MatchString("S.M_COVER_LIFE_SEQ.NEXTVAL") {
		t.Fatal("pola view lama tidak menggigit / terlalu lebar")
	}
	akar := filepath.Join("..", "..")
	diperiksa := 0
	err := filepath.WalkDir(akar, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "migrations" || d.Name() == "docs" {
				return filepath.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(p)
		if (ext != ".go" && ext != ".ts" && ext != ".tsx") || strings.HasSuffix(p, "_test.go") || strings.HasSuffix(p, ".test.ts") {
			return nil
		}
		isi, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		diperiksa++
		for i, baris := range strings.Split(string(isi), "\n") {
			kode := strings.TrimSpace(baris)
			if strings.HasPrefix(kode, "//") || strings.HasPrefix(kode, "*") || strings.HasPrefix(kode, "/*") {
				continue
			}
			if polaViewLama.MatchString(kode) {
				t.Errorf("%s:%d menyebut view lama COVER_LIFE: %s", p, i+1, kode)
			}
			if filepath.Base(filepath.Dir(p)) == "repository" && (strings.Contains(kode, "JSON") || strings.Contains(kode, "DELETE")) {
				t.Errorf("%s:%d memakai JSON / DELETE: %s", p, i+1, kode)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if diperiksa < 8 {
		t.Fatalf("hanya %d berkas diperiksa; pembacanya yang rusak", diperiksa)
	}
}
