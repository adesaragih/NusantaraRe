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
// jalan autocomplete - saringan tahun Begin date (A153, permintaan sesi 0f).
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
	// DaftarTableOfLimit - baris TABLEOFLIMIT ber-BIZCODE `kode` dan TAHUN `tahun`; `kategori` kosong = tanpa saringan
	// kategori.
	DaftarTableOfLimit(ctx context.Context, kode, tahun, kategori string) ([]models.BarisTableOfLimit, error)
}

// TableOfLimitOracle - PembacaTableOfLimit atas Oracle.
type TableOfLimitOracle struct{ db *db.DB }

// NewTableOfLimitOracle merakit pembaca BUSINESS / TABLEOFLIMIT.
func NewTableOfLimitOracle(d *db.DB) *TableOfLimitOracle { return &TableOfLimitOracle{db: d} }

// sqlKodeBisnis - NOTE dicocokkan persis (nama Class of Business tersimpan = BUSINESS.NOTE pilihan layar, tiket 28/29).
func sqlKodeBisnis(bisnis string) string {
	return "SELECT ID FROM " + bisnis + " WHERE NOTE = :1 AND BUSINESSGROUPID = :2 ORDER BY ID FETCH FIRST 2 ROWS ONLY"
}

// sqlTableOfLimit - DISTINCT atas kolom laporan RD; :1 kode, :2 tahun (teks, DDL VARCHAR2), [:3 kategori], batas terakhir.
func sqlTableOfLimit(tol string, denganKategori bool) string {
	syarat, batas := "BIZCODE = :1 AND TAHUN = :2", ":3"
	if denganKategori {
		syarat, batas = "BIZCODE = :1 AND TAHUN = :2 AND CATEGORY = :3", ":4"
	}
	return "SELECT DESCRIPTION, PCTLIMIT FROM (SELECT DISTINCT BIZCODE, CATEGORY, DESCRIPTION, PCTLIMIT, NOTE FROM " + tol +
		" WHERE " + syarat + ") ORDER BY CATEGORY, DESCRIPTION, PCTLIMIT, NOTE FETCH FIRST " + batas + " ROWS ONLY"
}

// KodeBisnis - lihat PembacaTableOfLimit.
func (r *TableOfLimitOracle) KodeBisnis(ctx context.Context, nama, grup string) ([]string, error) {
	q, err := r.db.Qualify(TabelBisnis)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlKodeBisnis(q), nama, grup)
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
func (r *TableOfLimitOracle) DaftarTableOfLimit(ctx context.Context, kode, tahun, kategori string) ([]models.BarisTableOfLimit, error) {
	q, err := r.db.Qualify(TabelTableOfLimitWarisan)
	if err != nil {
		return nil, err
	}
	arg := []any{kode, tahun, BatasTableOfLimit}
	if kategori != "" {
		arg = []any{kode, tahun, kategori, BatasTableOfLimit}
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
