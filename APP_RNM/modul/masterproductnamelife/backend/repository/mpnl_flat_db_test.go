//go:build db

package repository_test

// Tabel flat dan alat pindah terhadap Oracle NYATA (skema uji; POOLDATA tidak pernah menjadi sasaran - `uji/skemauji`).
// Tanpa ORACLE_DSN MELEWATI. Tabel flat dari berkas migrasi 140–147 YANG SAMA dengan produksi (dibaca `migrasi.Daftar`,
// tiket 01 AC 49); kedua tabel JSON warisan ditiru di sini (katalog DEV: `ID VARCHAR2(6)`, `JSONDATA CLOB IS JSON`).
// Fixture UJI-.
//
// ⚠️ Koneksi lewat `skemauji.BukaRepositori()` (seperti uji `db` handlers modul ini), BUKAN `skemauji.Buka()`: penjaga
// Claim Life `TestSetiapPemanggilBukaMemeriksaBolehDilewati` mengunci cacah pemanggilnya di seluruh APP_RNM.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul/masterproductnamelife/backend/repository"
	"nusantarare/uji/skemauji"
)

type ujiFlat struct {
	repo  *db.DB
	skema string
	ctx   context.Context
}

// ddlFlat - pernyataan berkas migrasi 140–147 modul ini (maju atau mundur), `{skema}` diganti skema uji.
func ddlFlat(t *testing.T, mundur bool, skema string) []string {
	t.Helper()
	entri, err := os.ReadDir(filepath.Join("..", "migrations"))
	if err != nil {
		t.Fatal(err)
	}
	sumber := fstest.MapFS{}
	for _, e := range entri {
		if strings.HasPrefix(e.Name(), "14") {
			isi, err := os.ReadFile(filepath.Join("..", "migrations", e.Name()))
			if err != nil {
				t.Fatal(err)
			}
			sumber["migrations/"+e.Name()] = &fstest.MapFile{Data: isi}
		}
	}
	langkah, err := migrasi.Daftar(mundur, sumber)
	if err != nil || len(langkah) != len(repository.DaftarTabelFlat) {
		t.Fatalf("membaca migrasi flat: %d langkah, %v", len(langkah), err)
	}
	var hasil []string
	for _, l := range langkah {
		for _, q := range l.Pernyataan {
			hasil = append(hasil, strings.ReplaceAll(q, "{skema}", skema))
		}
	}
	return hasil
}

func pasangFlat(t *testing.T) *ujiFlat {
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
	u := &ujiFlat{repo: repo, skema: repo.Skema(), ctx: ctx}
	u.bongkar(t)
	for _, q := range ddlFlat(t, false, u.skema) {
		u.exec(t, q)
	}
	for _, q := range []string{
		`CREATE TABLE {s}.M_PRODUCT_LIFE (ID VARCHAR2(6), JSONDATA CLOB CONSTRAINT UJI_FLAT_MPL CHECK (JSONDATA IS JSON),
			RIRISKID VARCHAR2(10), RIRISK VARCHAR2(100))`,
		`CREATE TABLE {s}.M_PRODUCTINWARD_LIFE (ID VARCHAR2(6), JSONDATA CLOB CONSTRAINT UJI_FLAT_MPIL CHECK (JSONDATA IS JSON))`,
	} {
		u.exec(t, q)
	}
	t.Cleanup(func() {
		u.bongkar(t)
		_ = repo.Close()
	})
	return u
}

// bongkar - kedua tiruan JSON dan tabel flat (jalur mundur 147→140); yang belum ada dilewati.
func (u *ujiFlat) bongkar(t *testing.T) {
	t.Helper()
	q := []string{fmt.Sprintf(`DROP TABLE %s.M_PRODUCT_LIFE PURGE`, u.skema),
		fmt.Sprintf(`DROP TABLE %s.M_PRODUCTINWARD_LIFE PURGE`, u.skema)}
	for _, p := range append(q, ddlFlat(t, true, u.skema)...) {
		if _, err := u.repo.ExecContext(u.ctx, p); err != nil && !strings.Contains(err.Error(), "ORA-00942") {
			t.Fatalf("membongkar: %v", err)
		}
	}
}

func (u *ujiFlat) exec(t *testing.T, q string, args ...any) {
	t.Helper()
	if _, err := u.repo.ExecContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema), args...); err != nil {
		t.Fatalf("%.40s: %v", q, err)
	}
}

func (u *ujiFlat) cacah(t *testing.T, tabel string) int {
	t.Helper()
	var n int
	if err := u.repo.QueryRowContext(u.ctx, fmt.Sprintf(`SELECT COUNT(*) FROM %s.%s`, u.skema, tabel)).Scan(&n); err != nil {
		t.Fatalf("mencacah %s: %v", tabel, err)
	}
	return n
}

// isiJSON - produk UJI- di kedua tabel JSON: setiap jenis kolom, K3 (MATURE bukan tanggal), K4 (outward kosong),
// pasangan inward ber-ID lain, desimal di bawah satu (`TM9` = `.5`), stempel komentar.
func (u *ujiFlat) isiJSON(t *testing.T) {
	t.Helper()
	for _, b := range [][2]string{
		{"100044", `{"ID":"100044","PRODUCTNAME":"UJI A","CEDING":"UJI C","RICOMM":"0.5","IsORS":"true",
			"POLICYHODER":"UJI-ORG-1","POLICYHODERNAME":"UJI P","CREATEOP":"UJI-A","UPDATEOP":"UJI-B","pxObjClass":"UJI",
			"LienClause":[{"Usia":"60","Manfaat":"50"}],"DocumentClaim":[{"Document":"UJI DOK"}],
			"PlanList":[{"Plan":"UJI PLAN","PlanID":"P1","RIRATE":"UJI RATE","RIRATEID":"R1"}],
			"FinancialUnderwritingList":[{"MinInsured":"1","MaxInsured":"1000000.25","Employee":"UJI E","Non_Employee":"UJI N"}],
			"UnderwritingLimitList":[{"MinInsured":"0","MaxInsured":"150000000.5","MinAge":"18","MaxAge":"65","Medical":"NM"}],
			"OutwardList":[{"REINSTYPEID":"","OVR_COMM":""},{"REINSTYPEID":"10200","TRANSACTIONYEAR":"2026","OVR_COMM":""}],
			"CommentList":[{"Date":"20241202T065450.847 GMT","OperatorName":"UJI-A","Suggest":"ok","IsApproved":""}]}`},
		{"100045", `{"ID":"100045","PRODUCTNAME":"UJI B","POLICYHODERNAME":"UJI Q"}`},
	} {
		u.exec(t, `INSERT INTO {s}.M_PRODUCT_LIFE (ID, JSONDATA) VALUES (:1, :2)`, b[0], b[1])
	}
	for _, b := range [][2]string{
		{"100044", `{"ID":"100044","PRODUCTID":"100044","BEGIN":"01/03/2026","STNC":"15/03/2026","MATURE":"UJI-BUKAN-TANGGAL",
			"MINAGE":"17","MAXAGE":"65","CEDINGLIMIT":"150000000.5","POLICYHODER":"UJI-ORG-1","POLICYHODERNAME":"UJI P"}`},
		{"100099", `{"ID":"100099","PRODUCTID":"100045","POLICYHODERNAME":"UJI Q","CURRENCY":"IDR"}`},
	} {
		u.exec(t, `INSERT INTO {s}.M_PRODUCTINWARD_LIFE (ID, JSONDATA) VALUES (:1, :2)`, b[0], b[1])
	}
}

func TestDBPindahFlatUjiJalankanDanUlang(t *testing.T) {
	u := pasangFlat(t)
	u.isiJSON(t)
	g := repository.Baru(u.repo)

	// Normalisasi menghentikan -jalankan, nol tulisan.
	u.exec(t, `INSERT INTO {s}.M_PRODUCT_LIFE (ID, JSONDATA) VALUES ('100046', '{"ID":"100046","RICOMM":"12.50"}')`)
	lap, err := g.PindahFlat(u.ctx, true, false)
	if !errors.Is(err, repository.ErrPindahTidakLolos) || lap.Normalisasi["M_PRODUCTNAME_LIFE.RICOMM (nol ekor desimal)"] != 1 ||
		u.cacah(t, repository.TabelFlatInduk) != 0 {
		t.Fatalf("normalisasi harus menolak tanpa tulisan: %v\n%s", err, lap.Teks())
	}
	// -terima-normalisasi: normalisasi tidak menahan; kegagalan tetap menahan (lihat uji murni).
	if lap, err = g.PindahFlat(u.ctx, true, true); err != nil || !lap.Ditulis || u.cacah(t, repository.TabelFlatInduk) != 3 {
		t.Fatalf("-terima-normalisasi: %v\n%s", err, lap.Teks())
	}
	u.exec(t, `DELETE FROM {s}.M_PRODUCTNAME_LIFE`)
	u.exec(t, `DELETE FROM {s}.M_PRODUCT_LIFE WHERE ID = '100046'`)

	// -uji: nol tulisan.
	lap, err = g.PindahFlat(u.ctx, false, false)
	if err != nil || !lap.Lolos() || lap.Ditulis || u.cacah(t, repository.TabelFlatInduk) != 0 {
		t.Fatalf("-uji: %v\n%s", err, lap.Teks())
	}
	if len(lap.K3) != 1 || lap.K4OutwardKosong != 1 || lap.InwardBerIDLain != 1 {
		t.Errorf("K3/K4/inward: %s", lap.Teks())
	}

	// -jalankan: tertulis, cacah sama dengan laporan, baca ulang = bentuk kanonik.
	lap, err = g.PindahFlat(u.ctx, true, false)
	if err != nil || !lap.Ditulis {
		t.Fatalf("-jalankan: %v\n%s", err, lap.Teks())
	}
	for _, tabel := range repository.DaftarTabelFlat {
		if n := u.cacah(t, tabel); n != lap.Baris[tabel] {
			t.Errorf("%s: %d baris, laporan %d", tabel, n, lap.Baris[tabel])
		}
	}
	p, err := g.BacaFlatUji(u.ctx, "100044")
	if err != nil {
		t.Fatal(err)
	}
	if p.Umum.RIComm != "0.5" || !p.Umum.IsORS || p.Inward.Begin != "2026-03-01" || p.Inward.Mature != "" ||
		p.Inward.CedingLimit != "150000000.5" || p.UnderwritingLimit[0].MaxInsured != "150000000.5" ||
		p.FinancialUnderwriting[0].MaxInsured != "1000000.25" || p.CommentList[0].Date != "20241202T065450.847 GMT" ||
		len(p.OutwardList) != 1 || p.OutwardList[0].TransactionYear != "2026" || p.Umum.PolicyHolder != "UJI-ORG-1" {
		t.Errorf("baca ulang Oracle: %+v", p)
	}
	if q, err := g.BacaFlatUji(u.ctx, "100045"); err != nil || q.Inward.Currency != "IDR" || q.Inward.ID != "100045" {
		t.Errorf("inward ber-ID lain pindah ke baris produknya: %+v %v", q.Inward, err)
	}

	// Aman diulang: isi sama → tertulis ulang tanpa beda.
	if lap, err = g.PindahFlat(u.ctx, true, false); err != nil || !lap.Ditulis {
		t.Fatalf("ulang: %v\n%s", err, lap.Teks())
	}
	// Tulisan baru di tabel flat → ditolak, tidak ditimpa.
	u.exec(t, `UPDATE {s}.M_PRODUCTNAME_LIFE SET PRODUCTNAME = 'UJI UBAH' WHERE ID = '100045'`)
	if lap, err = g.PindahFlat(u.ctx, true, true); !errors.Is(err, repository.ErrTulisanFlatBaru) || lap.TulisanFlatBerbeda != 1 {
		t.Errorf("tulisan flat baru harus menolak: %v\n%s", err, lap.Teks())
	}
	var nama string
	if err := u.repo.QueryRowContext(u.ctx, fmt.Sprintf(`SELECT PRODUCTNAME FROM %s.M_PRODUCTNAME_LIFE WHERE ID = '100045'`,
		u.skema)).Scan(&nama); err != nil || nama != "UJI UBAH" {
		t.Errorf("tulisan baru tidak boleh ditimpa: %q %v", nama, err)
	}
	// Tabel JSON tidak pernah disentuh.
	if u.cacah(t, "M_PRODUCT_LIFE") != 2 || u.cacah(t, "M_PRODUCTINWARD_LIFE") != 2 {
		t.Error("tabel JSON warisan berubah")
	}
}

// TestDBNilaiFlatOracleSamaDenganNormalkanFlat - bentuk kanonik yang diramalkan NormalkanFlat = yang dikembalikan
// Oracle, untuk setiap jenis kolom (termasuk NULL, desimal di bawah satu, bulat, tanggal, stempel).
func TestDBNilaiFlatOracleSamaDenganNormalkanFlat(t *testing.T) {
	u := pasangFlat(t)
	u.isiJSON(t)
	g := repository.Baru(u.repo)
	p, err := repository.UraiProduk("100044", jsonProdukUji(t, u), "100044", jsonInwardUji(t, u))
	if err != nil {
		t.Fatal(err)
	}
	p.OutwardList = p.OutwardList[1:]
	p.Inward.Mature = ""
	n, masalah := repository.NormalkanFlat(p)
	if len(masalah) != 0 {
		t.Fatalf("masalah: %+v", masalah)
	}
	tx, err := u.repo.Mulai(u.ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := g.TulisFlatUji(u.ctx, tx, n); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	baca, err := g.BacaFlatUji(u.ctx, "100044")
	if err != nil {
		t.Fatal(err)
	}
	if beda := repository.BedaProduk(n, baca); len(beda) != 0 {
		t.Errorf("Oracle berbeda dari NormalkanFlat: %+v", beda)
	}
}

func jsonProdukUji(t *testing.T, u *ujiFlat) string {
	t.Helper()
	var s string
	if err := u.repo.QueryRowContext(u.ctx, fmt.Sprintf(`SELECT JSONDATA FROM %s.M_PRODUCT_LIFE WHERE ID = '100044'`,
		u.skema)).Scan(&s); err != nil {
		t.Fatal(err)
	}
	return s
}

func jsonInwardUji(t *testing.T, u *ujiFlat) string {
	t.Helper()
	var s string
	if err := u.repo.QueryRowContext(u.ctx, fmt.Sprintf(`SELECT JSONDATA FROM %s.M_PRODUCTINWARD_LIFE WHERE ID = '100044'`,
		u.skema)).Scan(&s); err != nil {
		t.Fatal(err)
	}
	return s
}
