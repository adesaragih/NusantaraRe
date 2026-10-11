// Package repository adalah satu-satunya lapisan modul Endorsment Fac In yang menulis SQL. Setiap nilai diikat (`:n`,
// go-ora mengikat menurut URUTAN kemunculan, penampung unik - `TestNolPenampungBerulangDiSQLBerpembatasBaris`); nama
// objek lewat `Qualify` (ADR-U-0033). Nol COMMIT (ADR-U-0029). Nol CALL procedure (ADR-0043).
//
// Irisan I2 (periksa sebelum Create) hanya MEMBACA. Tabel warisan POOLDATA yang dibaca - tidak pernah dibuat modul
// ini: JSON_POLIS, OS_AKSEPTASI_KLAIM, FACINPRODUCTION, FACOUTPRODUCTION; lintas skema: ARASAPAS.DETAIL_INVOICE
// (hanya `arasapas.go`). Tabel milik modul ini: M_EDM_PENGECUALIAN_POLIS (migrasi 261).
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// Nama objek. Tabel warisan = nama di `pyBrowseSQL` rule Pega (tanpa skema di rule; `pooldata.` di GetStartDate).
const (
	TabelJSONPolis        = "JSON_POLIS"
	TabelKlaim            = "OS_AKSEPTASI_KLAIM"
	TabelFacinProduction  = "FACINPRODUCTION"
	TabelFacoutProduction = "FACOUTPRODUCTION"
	// TabelPengecualian - milik modul ini (261), diisi DBA.
	TabelPengecualian = "M_EDM_PENGECUALIAN_POLIS"
)

// ErrBelumAda - tabel/kolom/sequence tidak ada di skema ini (ORA-00942 / ORA-00904 / ORA-02289).
var ErrBelumAda = errors.New("repository: tabel, kolom, atau sequence tidak ada di skema ini")

// Gudang - akses Oracle modul ini.
type Gudang struct{ db *db.DB }

// Baru membuat gudang.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

func bungkus(err error, apa string) error {
	if err == nil {
		return nil
	}
	s := err.Error()
	if strings.Contains(s, "ORA-00942") || strings.Contains(s, "ORA-00904") || strings.Contains(s, "ORA-02289") {
		return fmt.Errorf("%w: %s: %v", ErrBelumAda, apa, err)
	}
	return fmt.Errorf("repository: %s: %w", apa, err)
}

// satuTeks - kolom pertama baris pertama; tanpa baris = ("", false).
func (g *Gudang) satuTeks(ctx context.Context, apa, q string, args ...any) (string, bool, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return "", false, err
	}
	var v sql.NullString
	err := g.db.QueryRowContext(ctx, q, args...).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, bungkus(err, apa)
	}
	return v.String, true, nil
}

// SqlNilaiPolisTerakhir - satu medan `DATA_JSON` baris JSON_POLIS bernomor polis itu yang `PRODKE`-nya = jumlah
// barisnya - 1. Bentuk WHERE = `RDBList/GetEDMStatus_SQL.xml` dan `RDBList/GetBusinessType_Sql.xml`
// `[terverifikasi]`:
//
//	select b.DATA_JSON.QuotationData.EdmType AS CARI20 FROM JSON_POLIS b WHERE NOPOLIS={InputData.CARI17}
//	  and PRODKE=(SELECT COUNT(NOPOLIS)-1 FROM JSON_POLIS b WHERE NOPOLIS={InputData.CARI17})
//
// ⚠️ `[dugaan]` Notasi titik Oracle (`b.DATA_JSON.QuotationData.EdmType`) ditulis sebagai `JSON_VALUE(… NULL ON
// ERROR)`: keduanya memberi teks skalar; notasi titik menuntut kolom dikenal JSON, sedangkan constraint `VALID_JSON`
// JSON_POLIS berstatus DISABLE (`D:\migrasi\RNM\DDL\JSON_POLIS.txt`). Pola modul endorsementlife. `PRODKE` VARCHAR2(5)
// dibandingkan dengan angka, apa adanya (konversi implisit Oracle, sama dengan rule). Pega membaca `pxResults(1)` -
// baris pertama tanpa ORDER BY; di sini FETCH FIRST 1.
func SqlNilaiPolisTerakhir(t, jalurJSON string) string {
	return fmt.Sprintf(`SELECT JSON_VALUE(b.DATA_JSON, '%s' NULL ON ERROR) FROM %s b
	  WHERE b.NOPOLIS = :1 AND b.PRODKE = (SELECT COUNT(c.NOPOLIS) - 1 FROM %s c WHERE c.NOPOLIS = :2)
	  FETCH FIRST 1 ROWS ONLY`, jalurJSON, t, t)
}

// Jalur JSON yang dibaca (teks tetap, bukan masukan pengguna).
const (
	JalurEdmType      = "$.QuotationData.EdmType"
	JalurBusinessType = "$.QuotationData.BusinessType"
)

func (g *Gudang) nilaiPolisTerakhir(ctx context.Context, apa, jalur, polis string) (string, error) {
	t, err := g.db.Qualify(TabelJSONPolis)
	if err != nil {
		return "", err
	}
	v, _, err := g.satuTeks(ctx, apa, SqlNilaiPolisTerakhir(t, jalur), polis, polis)
	return v, err
}

// StatusEDM - gerbang 1, `GetEDMStatus_SQL` → `OutputData1.pxResults(1).CARI20`; tanpa baris = "".
func (g *Gudang) StatusEDM(ctx context.Context, polis string) (string, error) {
	return g.nilaiPolisTerakhir(ctx, "membaca EdmType polis (GetEDMStatus_SQL)", JalurEdmType, polis)
}

// JenisBisnis - langkah 7, `GetBusinessType_Sql` → `OutData.pxResults(1).CARI2`; tanpa baris = "".
func (g *Gudang) JenisBisnis(ctx context.Context, polis string) (string, error) {
	return g.nilaiPolisTerakhir(ctx, "membaca BusinessType polis (GetBusinessType_Sql)", JalurBusinessType, polis)
}

// SqlKlaimTerakhir - `RDBList/GetDataClaim_SQL.xml` `[terverifikasi]`: `Select CaseID, sts_reject as CARI1 from
// os_akseptasi_klaim where nopolis = {InputData.CARI17} order by tanggal desc fetch first row only`. CaseID tidak
// dibaca langkah mana pun. `TO_CHAR`: tipe STS_REJECT tidak ada di DDL kiriman - teks dan angka sama-sama terbaca teks.
func SqlKlaimTerakhir(t string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(k.STS_REJECT) FROM %s k WHERE k.NOPOLIS = :1
	  ORDER BY k.TANGGAL DESC FETCH FIRST 1 ROWS ONLY`, t)
}

// Klaim - gerbang 3, `ListClaim.pxResults(*).CARI1`: nol atau satu nilai.
func (g *Gudang) Klaim(ctx context.Context, polis string) ([]string, error) {
	t, err := g.db.Qualify(TabelKlaim)
	if err != nil {
		return nil, err
	}
	v, ada, err := g.satuTeks(ctx, "membaca klaim polis (GetDataClaim_SQL)", SqlKlaimTerakhir(t), polis)
	if err != nil || !ada {
		return nil, err
	}
	return []string{v}, nil
}

// SqlAdaRenewal - `RDBList/GetListRNWbyNopolis_SQL.xml` `[terverifikasi]`: `Select IDPEGA as HASIL1 from json_polis
// where oldnopolis = {InputData.CARI17}`. Gerbang 5 hanya menguji `@LengthOfPageList(ListRNW.pxResults)>0` - cukup
// ada/tidak.
func SqlAdaRenewal(t string) string {
	return fmt.Sprintf(`SELECT 1 FROM %s j WHERE j.OLDNOPOLIS = :1 FETCH FIRST 1 ROWS ONLY`, t)
}

// AdaRenewal - gerbang 5.
func (g *Gudang) AdaRenewal(ctx context.Context, polis string) (bool, error) {
	t, err := g.db.Qualify(TabelJSONPolis)
	if err != nil {
		return false, err
	}
	_, ada, err := g.satuTeks(ctx, "membaca renewal polis (GetListRNWbyNopolis_SQL)", SqlAdaRenewal(t), polis)
	return ada, err
}

// SqlAdaFacout - `RDBList/GetFacoutList_SQL.xml` `[terverifikasi]`: `select IDPEGA as HASIL1,policyno HASIL2 from
// facoutproduction where policyno = {InputData.CARI17}`. Gerbang 6 menguji `@LengthOfPageList(ListFacout.pxResults)=0`.
func SqlAdaFacout(t string) string {
	return fmt.Sprintf(`SELECT 1 FROM %s f WHERE f.POLICYNO = :1 FETCH FIRST 1 ROWS ONLY`, t)
}

// AdaFacout - gerbang 6 (hanya EDM retro).
func (g *Gudang) AdaFacout(ctx context.Context, polis string) (bool, error) {
	t, err := g.db.Qualify(TabelFacoutProduction)
	if err != nil {
		return false, err
	}
	_, ada, err := g.satuTeks(ctx, "membaca facout polis (GetFacoutList_SQL)", SqlAdaFacout(t), polis)
	return ada, err
}

// PeriodeTeks - satu baris `GetStartDate`: BEGINDATE / ENDDATE sebagai teks `YYYYMMDDHH24MISS` (jam dinding DATE
// Oracle, tanpa zona). Kosong = NULL. Tafsir zonanya milik services.
type PeriodeTeks struct {
	Mulai, Akhir string
}

// SqlPeriodeProduksi - `RDBList/GetStartDate.xml` `[terverifikasi]`, bentuk sama:
//
//	select distinct BEGINDATE AS CARI1, ENDDATE AS CARI2 from pooldata.facinproduction a where nopolis = {InputSearch.CARI2}
//	and idpega in (select distinct idpega from pooldata.facinproduction where nopolis = a.nopolis
//	and spread_date = (select max(SPREAD_DATE) from pooldata.facinproduction where nopolis = a.nopolis))FETCH FIRST 1 ROW ONLY
//
// BEGINDATE/ENDDATE bertipe DATE (`D:\migrasi\RNM\DDL\FACINPRODUCTION.txt`); dibaca `TO_CHAR` supaya jam dindingnya
// tidak ditafsir zona sesi driver.
func SqlPeriodeProduksi(t string) string {
	return fmt.Sprintf(`SELECT DISTINCT TO_CHAR(a.BEGINDATE, 'YYYYMMDDHH24MISS'), TO_CHAR(a.ENDDATE, 'YYYYMMDDHH24MISS')
	  FROM %s a WHERE a.NOPOLIS = :1
	  AND a.IDPEGA IN (SELECT DISTINCT b.IDPEGA FROM %s b WHERE b.NOPOLIS = a.NOPOLIS
	    AND b.SPREAD_DATE = (SELECT MAX(c.SPREAD_DATE) FROM %s c WHERE c.NOPOLIS = a.NOPOLIS))
	  FETCH FIRST 1 ROWS ONLY`, t, t, t)
}

// PeriodeProduksi - `CheckEDMPolisDate` langkah 3; nil = tanpa baris.
func (g *Gudang) PeriodeProduksi(ctx context.Context, polis string) (*PeriodeTeks, error) {
	t, err := g.db.Qualify(TabelFacinProduction)
	if err != nil {
		return nil, err
	}
	q := SqlPeriodeProduksi(t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	var mulai, akhir sql.NullString
	err = g.db.QueryRowContext(ctx, q, polis).Scan(&mulai, &akhir)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, bungkus(err, "membaca periode polis (GetStartDate)")
	}
	return &PeriodeTeks{Mulai: mulai.String, Akhir: akhir.String}, nil
}

// SqlPengecualian - nomor polis satu JENIS (`TANGGAL` / `BATAL`), urut.
func SqlPengecualian(t string) string {
	return fmt.Sprintf(`SELECT p.NOPOLIS FROM %s p WHERE p.JENIS = :1 ORDER BY p.NOPOLIS`, t)
}

// Pengecualian - daftar nomor polis satu JENIS. Tabel belum dibuat DBA → ErrBelumAda (pemanggil memutuskan).
func (g *Gudang) Pengecualian(ctx context.Context, jenis string) ([]string, error) {
	t, err := g.db.Qualify(TabelPengecualian)
	if err != nil {
		return nil, err
	}
	q := SqlPengecualian(t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, jenis)
	if err != nil {
		return nil, bungkus(err, "membaca "+TabelPengecualian)
	}
	defer func() { _ = rows.Close() }()
	out := []string{}
	for rows.Next() {
		var v sql.NullString
		if err := rows.Scan(&v); err != nil {
			return nil, bungkus(err, "memindai "+TabelPengecualian)
		}
		out = append(out, v.String)
	}
	return out, bungkus(rows.Err(), "membaca "+TabelPengecualian)
}
