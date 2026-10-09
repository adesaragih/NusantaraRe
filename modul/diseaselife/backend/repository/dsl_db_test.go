//go:build db

package repository_test

// Seam repository Disease Life terhadap Oracle NYATA (skema uji; POOLDATA tidak pernah menjadi sasaran -
// `uji/skemauji`). Tanpa ORACLE_DSN MELEWATI. Keadaan DEV sebelum 080 ditiru (fakta WO 08-10-2026): DISEASE_LIFE
// (ID VARCHAR2(100) NULLABLE, ICD_CODE VARCHAR2(100), DISEASE VARCHAR2(1000)) TANPA PK dan TANPA indeks. Migrasi modul
// 080-081 dijalankan dari berkasnya. Fixture UJI-.
//
// ⚠️ Koneksi lewat `skemauji.BukaRepositori()`, BUKAN `skemauji.Buka()` (penjaga Claim Life
// `TestSetiapPemanggilBukaMemeriksaBolehDilewati` mengunci cacah pemanggilnya).
// ⚠️ DITULIS dan DIKOMPILASI (`go vet -tags=db`), BELUM PERNAH DIJALANKAN - executor tidak terhubung ke Oracle.

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul/diseaselife/backend/models"
	"nusantarare/modul/diseaselife/backend/repository"
	"nusantarare/modul/diseaselife/backend/services"
	"nusantarare/uji/skemauji"
)

type ujiDB struct {
	repo  *db.DB
	skema string
	ctx   context.Context
}

// sumberModul - folder `migrations/` modul ini.
func sumberModul(t *testing.T) fstest.MapFS {
	t.Helper()
	entri, err := os.ReadDir(filepath.Join("..", "migrations"))
	if err != nil {
		t.Fatal(err)
	}
	sumber := fstest.MapFS{}
	for _, e := range entri {
		isi, err := os.ReadFile(filepath.Join("..", "migrations", e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		sumber["migrations/"+e.Name()] = &fstest.MapFile{Data: isi}
	}
	return sumber
}

// ddl - pernyataan migrasi modul `kunci` (maju atau mundur), `{skema}` diganti skema uji.
func ddl(t *testing.T, kunci string, mundur bool, skema string) []string {
	t.Helper()
	l, err := migrasi.Daftar(mundur, sumberModul(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range l {
		if migrasi.KunciLangkah(x.Nama) == kunci {
			var out []string
			for _, q := range x.Pernyataan {
				out = append(out, strings.ReplaceAll(q, "{skema}", skema))
			}
			return out
		}
	}
	t.Fatalf("langkah %s tidak ada", kunci)
	return nil
}

// buka - koneksi skema uji (lewati tanpa ORACLE_DSN).
func buka(t *testing.T) *ujiDB {
	t.Helper()
	repo, err := skemauji.BukaRepositori()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	ctx := context.Background()
	if err := repo.Ping(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	return &ujiDB{repo: repo, skema: repo.Skema(), ctx: ctx}
}

// pasang - keadaan DEV sebelum 080 (bentuk DEV, nilai UJI-): tanpa PK, ID angka teks, satu ID bukan angka.
func pasang(t *testing.T) *ujiDB {
	t.Helper()
	u := buka(t)
	u.bongkar(t)
	baris := func(id, icd, nama string) string {
		return `INSERT INTO {s}.DISEASE_LIFE (ID, ICD_CODE, DISEASE) VALUES ('` + id + `', '` + icd + `', '` + nama + `')`
	}
	for _, q := range []string{
		`CREATE TABLE {s}.DISEASE_LIFE (ID VARCHAR2(100), ICD_CODE VARCHAR2(100), DISEASE VARCHAR2(1000))`,
		baris("100001", "UJI01", "UJI KOLERA"),
		baris("100002", "UJI02", "UJI DEMAM TIFOID"),
		baris("197585", "UJI99", "UJI TERTINGGI"),
		baris("UJI-X", "UJI98", "UJI ID BUKAN ANGKA"),
	} {
		u.exec(t, q)
	}
	t.Cleanup(func() {
		u.bongkar(t)
		_ = u.repo.Close()
	})
	return u
}

func (u *ujiDB) bongkar(t *testing.T) {
	t.Helper()
	for _, q := range []string{`DROP TABLE {s}.DISEASE_LIFE PURGE`, `DROP SEQUENCE {s}.SEQ_DISEASE_LIFE`} {
		if _, err := u.repo.ExecContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema)); err != nil &&
			!strings.Contains(err.Error(), "ORA-00942") && !strings.Contains(err.Error(), "ORA-02289") {
			t.Fatalf("membongkar %s: %v", q, err)
		}
	}
}

func (u *ujiDB) exec(t *testing.T, q string) {
	t.Helper()
	if _, err := u.repo.ExecContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema)); err != nil {
		t.Fatalf("%.60s: %v", q, err)
	}
}

// gagalORA - pernyataan `q` HARUS gagal dengan kode Oracle `kode`.
func (u *ujiDB) gagalORA(t *testing.T, q, kode string) {
	t.Helper()
	_, err := u.repo.ExecContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema))
	if err == nil || !strings.Contains(err.Error(), kode) {
		t.Errorf("%.60s: mau %s, dapat %v", q, kode, err)
	}
}

func (u *ujiDB) cacah(t *testing.T, q string) int {
	t.Helper()
	var n int
	if err := u.repo.QueryRowContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema)).Scan(&n); err != nil {
		t.Fatalf("%.60s: %v", q, err)
	}
	return n
}

func (u *ujiDB) periksa(t *testing.T, judul string, mau map[string]int) {
	t.Helper()
	for q, n := range mau {
		if got := u.cacah(t, q); got != n {
			t.Errorf("%s: %d, mau %d: %s", judul, got, n, q)
		}
	}
}

// lari - semua pernyataan langkah `kunci`, dua kali (aman diulang).
func (u *ujiDB) lari(t *testing.T, kunci string, mundur bool) {
	t.Helper()
	for i := 0; i < 2; i++ {
		for _, q := range ddl(t, kunci, mundur, u.skema) {
			if _, err := u.repo.ExecContext(u.ctx, q); err != nil &&
				!(mundur && strings.Contains(err.Error(), "ORA-02289")) {
				t.Fatalf("%s (mundur=%v, putaran %d): %v", kunci, mundur, i, err)
			}
		}
	}
}

// D1.1: SEQ_DISEASE_LIFE mulai ID ANGKA tertinggi + 1 (ID bukan angka diabaikan), NOCACHE; diulang tidak mati di
// ORA-00955; M_DISEASE_LIFE_SEQ tidak dibuat / disentuh.
func TestDBSequenceMulaiSesudahIDTertinggi(t *testing.T) {
	u := pasang(t)
	u.lari(t, "080_seq_disease_life", false)
	u.periksa(t, "sequence", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_SEQUENCES WHERE SEQUENCE_OWNER = UPPER('{s}') AND SEQUENCE_NAME = 'SEQ_DISEASE_LIFE'
			AND LAST_NUMBER = 197586 AND INCREMENT_BY = 1 AND CACHE_SIZE = 0 AND CYCLE_FLAG = 'N'`: 1,
		`SELECT COUNT(*) FROM SYS.ALL_SEQUENCES WHERE SEQUENCE_OWNER = UPPER('{s}') AND SEQUENCE_NAME = 'M_DISEASE_LIFE_SEQ'`: 0,
	})
}

// D1.2: PK BERHENTI dengan galat jelas bila ID kembar (ORA-02437) atau NULL (ORA-01449) - nol baris dihapus, nol PK;
// sesudah data dibereskan (oleh fixture uji ini, di DEV oleh berkas WO) PK berdiri, ID NOT NULL; diulang = nol galat;
// mundur membuang PK dan ID kembali NULLABLE.
func TestDBPKBerhentiBilaKembarAtauNull(t *testing.T) {
	u := pasang(t)
	blok := ddl(t, "081_disease_life_pk", false, u.skema)[0]
	u.exec(t, `INSERT INTO {s}.DISEASE_LIFE (ID, ICD_CODE, DISEASE) VALUES ('100002', 'UJI-TEST123', 'Uji Sakit')`)
	u.gagalORA(t, blok, "ORA-02437")
	u.periksa(t, "kembar: nol perubahan", map[string]int{
		`SELECT COUNT(*) FROM {s}.DISEASE_LIFE`: 5,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'DISEASE_LIFE'`: 0,
	})
	u.exec(t, `DELETE FROM {s}.DISEASE_LIFE WHERE ID = '100002' AND ICD_CODE = 'UJI-TEST123'`)
	u.exec(t, `INSERT INTO {s}.DISEASE_LIFE (ID, ICD_CODE, DISEASE) VALUES (NULL, 'UJI-NUL', 'UJI ID NULL')`)
	u.gagalORA(t, blok, "ORA-01449")
	u.exec(t, `DELETE FROM {s}.DISEASE_LIFE WHERE ID IS NULL`)
	u.lari(t, "081_disease_life_pk", false)
	u.periksa(t, "PK berdiri", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'DISEASE_LIFE'
			AND CONSTRAINT_NAME = 'PK_DISEASE_LIFE' AND CONSTRAINT_TYPE = 'P' AND STATUS = 'ENABLED'`: 1,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'DISEASE_LIFE'
			AND COLUMN_NAME = 'ID' AND NULLABLE = 'N'`: 1,
		`SELECT COUNT(*) FROM {s}.DISEASE_LIFE`: 4,
	})
	u.gagalORA(t, `INSERT INTO {s}.DISEASE_LIFE (ID, ICD_CODE, DISEASE) VALUES ('100001', 'UJI-K', 'UJI')`, "ORA-00001")
	u.lari(t, "081_disease_life_pk", true)
	u.periksa(t, "mundur", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND CONSTRAINT_NAME = 'PK_DISEASE_LIFE'`: 0,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'DISEASE_LIFE'
			AND COLUMN_NAME = 'ID' AND NULLABLE = 'Y'`: 1,
	})
}

func layananOracle(u *ujiDB) *services.Layanan {
	return services.BaruLayanan(repository.Baru(u.repo), func(ctx context.Context, fn func(tx *db.Tx) error) error {
		tx, err := u.repo.Mulai(ctx)
		if err != nil {
			return err
		}
		if err := fn(tx); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

// Seam: Add = TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL) = 197586, huruf besar; Edit; ICD kembar ditolak; grid bersaring dan
// berhalaman di server (bawaan ID menurun); ID dari sequence yang sudah ada = ErrIDTerpakai.
func TestDBLayanan(t *testing.T) {
	u := pasang(t)
	u.lari(t, "080_seq_disease_life", false)
	u.lari(t, "081_disease_life_pk", false)
	l := layananOracle(u)
	penuh := services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	p, err := l.Simpan(u.ctx, penuh, "", models.Isian{ICDCode: " uji03 ", Disease: " uji demam "})
	if err != nil || p != (models.Penyakit{ID: "197586", ICDCode: "UJI03", Disease: "UJI DEMAM"}) {
		t.Fatalf("add %+v %v", p, err)
	}
	if p, err := l.Simpan(u.ctx, penuh, "100002", models.Isian{ICDCode: "UJI02", Disease: "UJI DEMAM TIFOID BERAT"}); err != nil ||
		p.Disease != "UJI DEMAM TIFOID BERAT" {
		t.Errorf("edit %+v %v", p, err)
	}
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{ICDCode: "uji01", Disease: "UJI X"}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("kembar %v", err)
	}
	h, err := l.Daftar(u.ctx, models.Saringan{ICDCode: "uji0", Disease: "demam"})
	if err != nil || h.Total != 2 || h.Daftar[0].ID != "197586" || h.Daftar[1].ID != "100002" {
		t.Errorf("daftar %+v %v", h, err)
	}
	u.exec(t, `INSERT INTO {s}.DISEASE_LIFE (ID, ICD_CODE, DISEASE) VALUES ('197587', 'UJI-PEGA', 'UJI DITULIS PEGA')`)
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{ICDCode: "UJI04", Disease: "UJI Y"}); !errors.Is(err, services.ErrIDTerpakai) {
		t.Errorf("ID sequence sudah ada mau ErrIDTerpakai, dapat %v", err)
	}
}

// salinInti - langkah inti `nama` (maju + `_down`) dari migrasi inti SUNGGUHAN.
func salinInti(t *testing.T, ke fstest.MapFS, nama ...string) {
	t.Helper()
	for _, n := range nama {
		for _, b := range []string{n + ".sql", n + "_down.sql"} {
			isi, err := fs.ReadFile(inti.SumberMigrasi(), "migrations/"+b)
			if err != nil {
				t.Fatal(err)
			}
			ke["migrations/"+b] = &fstest.MapFile{Data: isi}
		}
	}
}

// K0 di skema uji: pelari SUNGGUHAN (`migrasi.Jalankan`) atas migrasi modul ini + 900 / 901 / 909 inti menjalankan 080,
// 081 SEBELUM 900 dan 951 SESUDAHNYA; baris menu lahir MASTER TREATY URUTAN 15 DIMIGRASI '1'; diulang = nol langkah.
func TestDBUrutanPelari080Lalu900Lalu951(t *testing.T) {
	u := pasang(t)
	sumber := sumberModul(t)
	cacahT := func(nama string) int {
		n := 0
		if err := u.repo.QueryRowContext(u.ctx, strings.ReplaceAll(`SELECT COUNT(*) FROM {s}.T_MIGRASI WHERE NAMA = '`+nama+`'`,
			"{s}", u.skema)).Scan(&n); err != nil && !strings.Contains(err.Error(), "ORA-00942") {
			t.Fatalf("membaca T_MIGRASI: %v", err)
		}
		return n
	}
	if cacahT("080_seq_disease_life") > 0 {
		t.Skip("skema uji sudah memuat migrasi modul ini (skemauji.Pasang); uji urutan butuh skema tanpa 080")
	}
	sudah := cacahT("900_m_nav_menu")
	intiSumber := fstest.MapFS{}
	salinInti(t, intiSumber, "900_m_nav_menu", "901_m_nav_menu_datar", "909_m_nav_menu_master_treaty")
	lap, err := migrasi.Jalankan(u.ctx, u.repo, intiSumber, sumber)
	if err != nil {
		t.Fatalf("pelari: %v (laporan %+v)", err, lap)
	}
	t.Cleanup(func() {
		if _, err := migrasi.Bongkar(u.ctx, u.repo, sumber); err != nil {
			t.Errorf("bongkar modul: %v", err)
		}
		if sudah == 0 {
			if _, err := migrasi.Bongkar(u.ctx, u.repo, intiSumber); err != nil {
				t.Errorf("bongkar inti: %v", err)
			}
		}
	})
	i080, i951 := slices.Index(lap.Dijalankan, "080_seq_disease_life"), slices.Index(lap.Dijalankan, "951_menu_diseaselife")
	if i080 < 0 || i951 < i080 {
		t.Errorf("dijalankan %v", lap.Dijalankan)
	}
	if sudah == 0 {
		if i900 := slices.Index(lap.Dijalankan, "900_m_nav_menu"); !(i080 < i900 && i900 < i951) {
			t.Errorf("urutan 080 -> 900 -> 951 dilanggar: %v", lap.Dijalankan)
		}
	}
	for _, x := range lap.ObjekSudahAda {
		if strings.Contains(x, "DISEASE") {
			t.Errorf("dilewati, objeknya sudah ada: %s", x)
		}
	}
	u.periksa(t, "menu", map[string]int{
		`SELECT COUNT(*) FROM {s}.M_NAV_MENU WHERE KODE = 'diseaselife' AND LABEL = 'Disease Life'
			AND GROUPMENU = 'MASTER TREATY' AND MODUL = 'diseaselife' AND URUTAN = 15 AND DIMIGRASI = '1'`: 1,
		`SELECT COUNT(*) FROM {s}.T_MIGRASI WHERE NAMA IN ('080_seq_disease_life', '081_disease_life_pk', '951_menu_diseaselife')`: 3,
		`SELECT COUNT(*) FROM {s}.DISEASE_LIFE`: 4,
	})
	if ulang, err := migrasi.Jalankan(u.ctx, u.repo, intiSumber, sumber); err != nil || len(ulang.Dijalankan) != 0 {
		t.Errorf("diulang: %+v %v", ulang, err)
	}
}
