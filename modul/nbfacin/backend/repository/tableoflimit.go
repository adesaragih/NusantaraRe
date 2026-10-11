package repository

// Popup Choose Class of Construction (tiket 40): RD `BrowseTableOfLimit_RD` atas tabel warisan POOLDATA.TABLEOFLIMIT,
// dan kode bisnis (BIZCODE) case dari POOLDATA.BUSINESS.
//
// `[terverifikasi]` RD (kelas ASM-FW-GISFW-Int-TABLEOFLIMIT): filter `.Tahun = Param.Tahun`, `.Bizcode = Param.Bizcode`,
// `.Category = Param.Category`, `.ID = Param.ID` (A AND B AND C AND D); DISTINCT atas kolom laporan (Bizcode, Category,
// Description, PctLimit, Note); maks 500; urut Category lalu Description. Layar Occupation punya DUA jalan ke RD ini:
// autocomplete Class of Construction (`Section\OccupationItemFacIn_Section.xml`, data page D_BrowseTableOfLimit) mengirim
// `Tahun = pyWorkPage.OfferFacIn.CurrentYear` (= tahun Begin date, SetValidateDate_Act), sedangkan tombol Choose Class of
// Construction (`Section\ChooseClassofContraction.xml`) mengirim Tahun KOSONG (filter dibuang). Endpoint ini mengikuti
// popup tombol - TANPA saringan tahun (A161, menggantikan A153: contoh `DDL\TABLEOFLIMIT.xml` ber-TAHUN 2017, saringan
// tahun Begin date membuang semua baris di DEV; dikonfirmasi work owner 03-10-2026, butir 95).
// DDL `DDL\TABLEOFLIMIT.txt` (03-10-2026): seluruh kolom VARCHAR2(4000 BYTE); nama tabel dari DDL itu.
// BIZCODE = BUSINESS.ID Class of Business case (keputusan work owner butir 89).

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelTableOfLimitWarisan - tabel warisan POOLDATA, baca saja.
	TabelTableOfLimitWarisan = "TABLEOFLIMIT"
	// BatasTableOfLimit - pyMaxRecords RD.
	BatasTableOfLimit = 500
)

// PembacaTableOfLimit - kode bisnis Class of Business dan pilihan Class of Construction.
type PembacaTableOfLimit interface {
	// KodeBisnis - BUSINESS.ID ber-NOTE `nama` di group `grup` (paling banyak dua, cukup untuk mengenali ganda).
	KodeBisnis(ctx context.Context, nama, grup string) ([]string, error)
	// OldIDBisnis - BUSINESS.OLDID baris yang sama dengan KodeBisnis (tiket 62: `BusinessOldId` predikat bentuk B tangga).
	OldIDBisnis(ctx context.Context, nama, grup string) ([]string, error)
	// BisnisLini - ID, OLDID, GROUPPANEL baris BUSINESS ber-GROUPPANEL (tiket 73 P0: lini case; paling banyak dua).
	BisnisLini(ctx context.Context, nama, grup string) ([]models.BarisBisnisLini, error)
	// DaftarTableOfLimit - baris TABLEOFLIMIT ber-BIZCODE `kode`, seluruh tahun (A161); `kategori` kosong = tanpa
	// saringan kategori.
	DaftarTableOfLimit(ctx context.Context, kode, kategori string) ([]models.BarisTableOfLimit, error)
}

// TableOfLimitOracle - PembacaTableOfLimit atas Oracle.
type TableOfLimitOracle struct{ db *db.DB }

// NewTableOfLimitOracle merakit pembaca BUSINESS / TABLEOFLIMIT.
func NewTableOfLimitOracle(d *db.DB) *TableOfLimitOracle { return &TableOfLimitOracle{db: d} }

// sqlKodeBisnis - NOTE dicocokkan persis (nama Class of Business tersimpan = BUSINESS.NOTE pilihan layar, tiket 28/29).
func sqlKodeBisnis(bisnis string) string { return sqlKolomBisnis(bisnis, "ID") }

// sqlOldIDBisnis - `[terverifikasi]` Pega mengisi `pyWorkPage.OfferFacIn.QuotationData.BusinessOldId` dari
// RDBList\GetOldIDBusiness_SQL.xml (`select oldid as CARI1 ... from business where (ID ={OldID.CARI2} OR NOTE =
// {OldID.CARI2}) AND GROUPPANEL IS NOT NULL`, pxResults(1), Activity\SaveJsonOfferFacIn_Act.xml langkah 3-5). Di sini baris
// BUSINESS dikenali seperti KodeBisnis (NOTE + BUSINESSGROUPID, butir 89) - `[keputusan agent]`, tiket 62.
func sqlOldIDBisnis(bisnis string) string { return sqlKolomBisnis(bisnis, "OLDID") }

// sqlBisnisLini - lini case (tiket 73 P0). GROUPPANEL VARCHAR2(4000) `[terverifikasi]` DDL\BUSINESS.txt L5. Saringan
// `GROUPPANEL IS NOT NULL` disamakan dengan Pega `[terverifikasi]` RDBList\GetOldIDBusiness_SQL.xml (`... from business where
// (ID ={OldID.CARI2} OR NOTE ={OldID.CARI2}) AND GROUPPANEL IS NOT NULL`) - keputusan work owner 09-10-2026. ID, OLDID, dan
// GROUPPANEL dibaca dalam SATU kueri supaya ketiganya dari baris yang sama sesudah saringan. Pengenalan baris tetap NOTE +
// BUSINESSGROUPID (butir 89); `ID = ... OR NOTE = ...` Pega TIDAK ditiru (temuan, belum diputuskan).
func sqlBisnisLini(bisnis string) string {
	return "SELECT ID, OLDID, GROUPPANEL FROM " + bisnis +
		" WHERE NOTE = :1 AND BUSINESSGROUPID = :2 AND GROUPPANEL IS NOT NULL ORDER BY ID FETCH FIRST 2 ROWS ONLY"
}

func sqlKolomBisnis(bisnis, kolom string) string {
	return "SELECT " + kolom + " FROM " + bisnis + " WHERE NOTE = :1 AND BUSINESSGROUPID = :2 ORDER BY ID FETCH FIRST 2 ROWS ONLY"
}

// sqlTableOfLimit - DISTINCT atas kolom laporan RD; :1 kode, [:2 kategori], batas terakhir. Tanpa TAHUN (A161) - baris
// yang sama di beberapa tahun menyatu oleh DISTINCT (TAHUN bukan kolom laporan RD).
func sqlTableOfLimit(tol string, denganKategori bool) string {
	syarat, batas := "BIZCODE = :1", ":2"
	if denganKategori {
		syarat, batas = "BIZCODE = :1 AND CATEGORY = :2", ":3"
	}
	return "SELECT DESCRIPTION, PCTLIMIT FROM (SELECT DISTINCT BIZCODE, CATEGORY, DESCRIPTION, PCTLIMIT, NOTE FROM " + tol +
		" WHERE " + syarat + ") ORDER BY CATEGORY, DESCRIPTION, PCTLIMIT, NOTE FETCH FIRST " + batas + " ROWS ONLY"
}

// KodeBisnis - lihat PembacaTableOfLimit.
func (r *TableOfLimitOracle) KodeBisnis(ctx context.Context, nama, grup string) ([]string, error) {
	return r.kolomBisnis(ctx, sqlKodeBisnis, nama, grup)
}

// OldIDBisnis - lihat PembacaTableOfLimit.
func (r *TableOfLimitOracle) OldIDBisnis(ctx context.Context, nama, grup string) ([]string, error) {
	return r.kolomBisnis(ctx, sqlOldIDBisnis, nama, grup)
}

// BisnisLini - lihat PembacaTableOfLimit.
func (r *TableOfLimitOracle) BisnisLini(ctx context.Context, nama, grup string) ([]models.BarisBisnisLini, error) {
	q, err := r.db.Qualify(TabelBisnis)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlBisnisLini(q), nama, grup)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelBisnis, err)
	}
	defer baris.Close()
	hasil := []models.BarisBisnisLini{}
	for baris.Next() {
		var id, old, gp sql.NullString
		if err := baris.Scan(&id, &old, &gp); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelBisnis, err)
		}
		hasil = append(hasil, models.BarisBisnisLini{ID: id.String, OldID: old.String, GroupPanel: gp.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelBisnis, err)
	}
	return hasil, nil
}

func (r *TableOfLimitOracle) kolomBisnis(ctx context.Context, susun func(string) string, nama, grup string) ([]string, error) {
	q, err := r.db.Qualify(TabelBisnis)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, susun(q), nama, grup)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelBisnis, err)
	}
	defer baris.Close()
	hasil := []string{}
	for baris.Next() {
		var id sql.NullString
		if err := baris.Scan(&id); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelBisnis, err)
		}
		hasil = append(hasil, id.String)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelBisnis, err)
	}
	return hasil, nil
}

// DaftarTableOfLimit - lihat PembacaTableOfLimit. PCTLIMIT teks apa adanya (DDL VARCHAR2; koma dan spasi ujung tetap).
func (r *TableOfLimitOracle) DaftarTableOfLimit(ctx context.Context, kode, kategori string) ([]models.BarisTableOfLimit, error) {
	q, err := r.db.Qualify(TabelTableOfLimitWarisan)
	if err != nil {
		return nil, err
	}
	arg := []any{kode, BatasTableOfLimit}
	if kategori != "" {
		arg = []any{kode, kategori, BatasTableOfLimit}
	}
	baris, err := r.db.QueryContext(ctx, sqlTableOfLimit(q, kategori != ""), arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelTableOfLimitWarisan, err)
	}
	defer baris.Close()
	hasil := []models.BarisTableOfLimit{}
	for baris.Next() {
		var ket, batas sql.NullString
		if err := baris.Scan(&ket, &batas); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelTableOfLimitWarisan, err)
		}
		hasil = append(hasil, models.BarisTableOfLimit{Description: ket.String, PctLimit: batas.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelTableOfLimitWarisan, err)
	}
	return hasil, nil
}
