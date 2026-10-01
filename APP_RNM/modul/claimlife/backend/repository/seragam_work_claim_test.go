package repository_test

// Migrasi 023 - `T_WORK_CLAIM` diseragamkan (keputusan work owner 01-10-2026,
// lanjutan brief seragam kolom T_WORK_POLIS). TANPA Oracle.
//
// Untuk apa berkas ini: menagih bentuk AKHIR `T_WORK_CLAIM` dan
// `T_GENERAL_CLAIM` sesudah SELURUH migrasi maju setiap modul, serta urutan
// langkah 023 dan jalur mundurnya.
//
// Keputusan work owner 01-10-2026:
//   - `PY_POSITION` -> `POSITION` (isinya tetap `pyPosition`, nama peran);
//   - `TYPE` pindah ke header klaim `T_GENERAL_CLAIM` - `TypeKlaim` membacanya
//     dari sana;
//   - `CASE_ID` dibuang, `ID` dipakai: kasus baru memang `CASE_ID = ID`, dan
//     migrasi klaim lama memakai CASEID warisan sebagai `ID`;
//   - `TAHAP` TETAP.
//
// Dibaca sesudah: migrations/023_seragam_kolom_t_work_claim.sql.

import (
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
	"nusantarare/uji/skemauji"
)

const berkas023 = "023_seragam_kolom_t_work_claim.sql"

var (
	polaBuatTabelKlaim   = regexp.MustCompile(`(?s)^CREATE TABLE \{skema\}\.(\w+) \((.*)\n\)$`)
	polaTambahKolomKlaim = regexp.MustCompile(`(?s)^ALTER TABLE \{skema\}\.(\w+) ADD \((.*)\n\)$`)
	polaKolomTipeKlaim   = regexp.MustCompile(`^\s*([A-Z0-9_]+)\s+([A-Z0-9_]+(?:\(\d+(?:,\d+)?\))?)`)
)

// bentukAkhirKlaim membaca kolom dan tipe tabel yang diminta sesudah seluruh
// migrasi maju SETIAP modul (`skemauji.SumberMigrasi`), urut nomor: CREATE,
// ADD, lalu DROP COLUMN - juga yang di dalam blok berpelindung katalog.
func bentukAkhirKlaim(t *testing.T, tabel ...string) map[string]map[string]string {
	t.Helper()
	dicari := map[string]bool{}
	for _, n := range tabel {
		dicari[n] = true
	}
	langkah, err := migrasi.Daftar(false, skemauji.SumberMigrasi()...)
	if err != nil {
		t.Fatal(err)
	}
	hasil := map[string]map[string]string{}
	baca := func(tab, isi string) {
		if hasil[tab] == nil {
			hasil[tab] = map[string]string{}
		}
		for _, b := range strings.Split(isi, "\n") {
			if m := polaKolomTipeKlaim.FindStringSubmatch(b); m != nil && m[1] != "CONSTRAINT" && m[1] != "REFERENCES" {
				hasil[tab][m[1]] = m[2]
			}
		}
	}
	for _, m := range langkah {
		for _, p := range m.Pernyataan {
			p = strings.TrimSpace(p)
			if c := polaBuatTabelKlaim.FindStringSubmatch(p); c != nil && dicari[c[1]] {
				baca(c[1], c[2])
				continue
			}
			if c := polaTambahKolomKlaim.FindStringSubmatch(p); c != nil && dicari[c[1]] {
				baca(c[1], c[2])
				continue
			}
			if tab, kolom := migrasi.KolomAlterBuang(p); dicari[tab] {
				for _, k := range kolom {
					delete(hasil[tab], k)
				}
			}
		}
	}
	return hasil
}

func TestMigrasi023BentukAkhirWorkClaim(t *testing.T) {
	b := bentukAkhirKlaim(t, "T_WORK_CLAIM", "T_GENERAL_CLAIM")
	mau := map[string]string{
		"ID": "VARCHAR2(32)", "COVER_KEY": "VARCHAR2(32)", "LINI": "VARCHAR2(16)",
		"POSITION": "VARCHAR2(64)", "SENDTO_ADMIN": "VARCHAR2(8)", "SENDTO_MEDICAL": "VARCHAR2(8)",
		"CREATE_OP": "VARCHAR2(64)", "CREATE_OP_NAME": "VARCHAR2(128)", "TGL_UPDATE": "DATE",
		"TAHAP": "VARCHAR2(32)", "TGL_CREATE": "DATE", "STATUS_WORK": "VARCHAR2(32)",
	}
	if !reflect.DeepEqual(b["T_WORK_CLAIM"], mau) {
		t.Errorf("bentuk akhir T_WORK_CLAIM:\n dapat %v\n mau   %v", b["T_WORK_CLAIM"], mau)
	}
	// TYPE pindah ke header klaim, tipe yang sama dengan asalnya (001).
	if got := b["T_GENERAL_CLAIM"]["TYPE"]; got != "VARCHAR2(32)" {
		t.Errorf("T_GENERAL_CLAIM.TYPE = %q, mau VARCHAR2(32)", got)
	}
}

// bacaPerintah023 membaca pernyataan satu arah 023 dari folder modul ini.
func bacaPerintah023(t *testing.T, nama string) []string {
	t.Helper()
	p, err := migrasi.PernyataanLangkah(os.DirFS(".."), nama)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// Urutan maju. Pengaman CASE_ID PALING DEPAN: bila satu baris saja punya
// CASE_ID yang berbeda dari ID, Oracle menolak CHECK-nya (ORA-02293) dan
// migrasi berhenti SEBELUM apa pun berubah - data tidak hilang diam-diam.
func TestMigrasi023LangkahMaju(t *testing.T) {
	p := bacaPerintah023(t, berkas023)
	if len(p) != 8 {
		t.Fatalf("023 terbaca %d pernyataan, mau 8", len(p))
	}
	blok := func(i int) migrasi.PerintahKatalog {
		t.Helper()
		pk, ok := migrasi.BacaPerintahKatalog(p[i])
		if !ok {
			t.Fatalf("pernyataan %d bukan blok berpelindung katalog:\n%s", i+1, p[i])
		}
		return pk
	}
	if pk := blok(0); pk.Tabel != "T_WORK_CLAIM" || pk.Objek != "CASE_ID" || !pk.BilaAda ||
		pk.Perintah != "ALTER TABLE {skema}.T_WORK_CLAIM ADD CHECK (CASE_ID IS NULL OR CASE_ID = ID)" {
		t.Errorf("pengaman CASE_ID: %+v", pk)
	}
	for i, m := range []struct {
		tabel string
		kolom []string
	}{{"T_WORK_CLAIM", []string{"POSITION"}}, {"T_GENERAL_CLAIM", []string{"TYPE"}}} {
		tab, kol := migrasi.KolomAlterTambah(p[1+i])
		if tab != m.tabel || !reflect.DeepEqual(kol, m.kolom) {
			t.Errorf("pernyataan %d menambah %s %v, mau %s %v", i+2, tab, kol, m.tabel, m.kolom)
		}
	}
	mau := []struct{ tabel, objek, perintah string }{
		{"T_WORK_CLAIM", "PY_POSITION", "UPDATE {skema}.T_WORK_CLAIM SET POSITION = PY_POSITION"},
		{"T_WORK_CLAIM", "PY_POSITION", "ALTER TABLE {skema}.T_WORK_CLAIM DROP COLUMN PY_POSITION"},
		{"T_WORK_CLAIM", "TYPE", "UPDATE {skema}.T_GENERAL_CLAIM g SET g.TYPE = (SELECT w.TYPE FROM {skema}.T_WORK_CLAIM w WHERE w.ID = g.ID) WHERE g.TYPE IS NULL"},
		{"T_WORK_CLAIM", "TYPE", "ALTER TABLE {skema}.T_WORK_CLAIM DROP COLUMN TYPE"},
		{"T_WORK_CLAIM", "CASE_ID", "ALTER TABLE {skema}.T_WORK_CLAIM DROP COLUMN CASE_ID CASCADE CONSTRAINTS"},
	}
	for i, m := range mau {
		pk := blok(3 + i)
		if pk.Katalog != "ALL_TAB_COLUMNS" || pk.Tabel != m.tabel || pk.Objek != m.objek || !pk.BilaAda || pk.Perintah != m.perintah {
			t.Errorf("pernyataan %d: %+v, mau %+v", 4+i, pk, m)
		}
	}
	semua := strings.Join(p, "\n")
	for _, w := range []string{"RENAME", "COMMIT", "TAHAP"} {
		if strings.Contains(semua, w) {
			t.Errorf("023 maju memuat %q", w)
		}
	}
}

// Jalur mundur: seluruhnya berpelindung, kebalikan urutan maju. CASE_ID dan
// TYPE kembali beserta isinya (CASE_ID = ID untuk baris klaim; TYPE baris
// komite dari klaim induknya - begitulah `pxAddChildWork` menyalinnya).
func TestMigrasi023LangkahMundur(t *testing.T) {
	p := bacaPerintah023(t, strings.TrimSuffix(berkas023, ".sql")+"_down.sql")
	mau := []struct {
		tabel, objek string
		bilaAda      bool
		perintah     string
	}{
		{"T_WORK_CLAIM", "CASE_ID", false, "ALTER TABLE {skema}.T_WORK_CLAIM ADD (CASE_ID VARCHAR2(64))"},
		{"T_WORK_CLAIM", "CASE_ID", true, "UPDATE {skema}.T_WORK_CLAIM w SET w.CASE_ID = w.ID WHERE w.CASE_ID IS NULL AND EXISTS (SELECT 1 FROM {skema}.T_GENERAL_CLAIM g WHERE g.ID = w.ID)"},
		{"T_WORK_CLAIM", "TYPE", false, "ALTER TABLE {skema}.T_WORK_CLAIM ADD (TYPE VARCHAR2(32))"},
		{"T_WORK_CLAIM", "TYPE", true, "UPDATE {skema}.T_WORK_CLAIM w SET w.TYPE = COALESCE((SELECT g.TYPE FROM {skema}.T_GENERAL_CLAIM g WHERE g.ID = w.ID), (SELECT g.TYPE FROM {skema}.T_GENERAL_CLAIM g WHERE g.ID = w.COVER_KEY)) WHERE w.TYPE IS NULL"},
		{"T_GENERAL_CLAIM", "TYPE", true, "ALTER TABLE {skema}.T_GENERAL_CLAIM DROP COLUMN TYPE"},
		{"T_WORK_CLAIM", "POSITION", true, "ALTER TABLE {skema}.T_WORK_CLAIM RENAME COLUMN POSITION TO PY_POSITION"},
	}
	if len(p) != len(mau) {
		t.Fatalf("023_down terbaca %d pernyataan, mau %d", len(p), len(mau))
	}
	for i, m := range mau {
		pk, ok := migrasi.BacaPerintahKatalog(p[i])
		if !ok || pk.Katalog != "ALL_TAB_COLUMNS" || pk.Tabel != m.tabel || pk.Objek != m.objek ||
			pk.BilaAda != m.bilaAda || pk.Perintah != m.perintah {
			t.Errorf("pernyataan mundur %d: %+v (ok=%v), mau %+v", i+1, pk, ok, m)
		}
	}
}
