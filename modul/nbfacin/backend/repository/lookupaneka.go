package repository

// Lookup form objek Aneka (tiket 73 P1, keputusan work owner 09-10-2026: hanya kapal dan tipe kendaraan; Object HE / Merk / Draft
// Wording = teks bebas sampai DDL sumbernya ada, K73-10 / K73-11). Baca saja, seluruh nilai parameter terikat.
//
// `[terverifikasi]` `NB FacIn\ReportDefinition\BrowseShip_RD.xml`: kelas Int-SHIP, filter B `.Name Contains Param.Name`
// (pyCaseInsensitive, pyUseNullIfEmpty=false), G `.Remark != "Not Active"` (peka huruf) OR H `.Remark IS NULL`, filter A/C-F
// dikirim kosong oleh form (`InputDtlObjectAneka_FacIn.xml` L10575-10610); urut `.ID` DESC, pyMaxRecords 500. Tabel
// `[dugaan kuat]` POOLDATA.SHIP (`D:\migrasi\RNM\DDL\SHIP.txt`, kolom cocok 1:1; semua VARCHAR2 - ID diurut sebagai teks).
// RD ber-pyUseAlternateDB: DB alternatifnya `[pertanyaan terbuka]`, dibaca lewat skema biasa (db.Qualify).
//
// `[terverifikasi]` `NB FacIn\ReportDefinition\BrowseType_RD.xml`: kelas Int-BRANDDETAIL, filter B `.Type = "ANEKA"` (peka huruf);
// A `.ID` dan C `.TypeName` dikirim kosong oleh form (L27338-27352) - TIDAK disaring merk / objek; urut `.ID` DESC, maks 500.
// Sumber view POOLDATA.BRANDDETAIL (`DDL\BRANDDETAIL.txt`, atas JSON m_branddetail; tipe kolom belum terverifikasi).

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelShip / ViewBrandDetail - ejaan DDL.
	TabelShip       = "SHIP"
	ViewBrandDetail = "BRANDDETAIL"
	// BatasLookupAneka - pyMaxRecords kedua RD.
	BatasLookupAneka = 500
	// remarkTidakAktif - literal filter G BrowseShip_RD.
	remarkTidakAktif = "Not Active"
	// tipeAneka - literal filter B BrowseType_RD.
	tipeAneka = "ANEKA"
)

// PembacaLookupAneka - pilihan lookup form objek Aneka.
type PembacaLookupAneka interface {
	// CariKapal - kata kosong = tanpa saringan nama (Pega mengabaikan filter kosong `[dugaan]`).
	CariKapal(ctx context.Context, kata string) ([]models.BarisKapal, error)
	DaftarTipeKendaraan(ctx context.Context) ([]models.BarisTipeKendaraan, error)
	// CariOkupasiAneka - popup Choose Occupation Aneka: businessName / kata kosong = saringan itu tidak dipakai.
	CariOkupasiAneka(ctx context.Context, businessName, kata string) ([]models.PilihanOkupasiAneka, error)
	// CariCoverageAneka - Coverage Note tab Coverage Aneka (tiket 73 P2): kodeBisnis / kata kosong = saringan itu tidak dipakai.
	CariCoverageAneka(ctx context.Context, kodeBisnis, kata string) ([]models.BarisCoverage, error)
	// IndemnityPct - INDEMNITY_PCT M_BI_INDEMNITY berkunci Indemnity + Unit (indemnity.go); nol baris = larik kosong.
	IndemnityPct(ctx context.Context, indemnity, unit string) ([]*apd.Decimal, error)
}

// LookupAnekaOracle - PembacaLookupAneka atas Oracle.
type LookupAnekaOracle struct{ db *db.DB }

// NewLookupAnekaOracle merakit pembaca lookup Aneka.
func NewLookupAnekaOracle(d *db.DB) *LookupAnekaOracle { return &LookupAnekaOracle{db: d} }

// sqlCariKapal - :1 literal Remark tidak aktif, [:2 pola nama], batas terakhir.
func sqlCariKapal(tabel string, denganKata bool) string {
	syarat, batas := "", ":2"
	if denganKata {
		syarat, batas = ` AND UPPER(NAME) LIKE :2 ESCAPE '\'`, ":3"
	}
	return "SELECT ID, NAME, FLAGNAME, YEARMAKE1, NRT, GRT, DWT, CLASSIFICATION, CONSTRUCTION, SHIPTYPEID, SHIPNAME, REGISTER, IMO," +
		" YEARMAKE2, REMARK, FORMERSHIPNAME, REBUILT FROM " + tabel +
		" WHERE (REMARK != :1 OR REMARK IS NULL)" + syarat + " ORDER BY ID DESC FETCH FIRST " + batas + " ROWS ONLY"
}

// sqlTipeKendaraan - :1 literal Type, :2 batas.
func sqlTipeKendaraan(view string) string {
	return "SELECT ID, TYPENAME FROM " + view + " WHERE TYPE = :1 ORDER BY ID DESC FETCH FIRST :2 ROWS ONLY"
}

// CariKapal - lihat PembacaLookupAneka.
func (r *LookupAnekaOracle) CariKapal(ctx context.Context, kata string) ([]models.BarisKapal, error) {
	q, err := r.db.Qualify(TabelShip)
	if err != nil {
		return nil, err
	}
	sqlq, arg := sqlCariKapal(q, false), []any{remarkTidakAktif, BatasLookupAneka}
	if pola := PolaCari(kata); pola != "" {
		sqlq, arg = sqlCariKapal(q, true), []any{remarkTidakAktif, pola, BatasLookupAneka}
	}
	baris, err := r.db.QueryContext(ctx, sqlq, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelShip, err)
	}
	defer baris.Close()
	hasil := []models.BarisKapal{}
	for baris.Next() {
		var v [17]sql.NullString
		tuju := make([]any, len(v))
		for i := range v {
			tuju[i] = &v[i]
		}
		if err := baris.Scan(tuju...); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelShip, err)
		}
		hasil = append(hasil, models.BarisKapal{ID: v[0].String, Nama: v[1].String, Flag: v[2].String, Tahun: v[3].String,
			NRT: v[4].String, GRT: v[5].String, DWT: v[6].String, Klasifikasi: v[7].String, Konstruksi: v[8].String,
			ShipTypeID: v[9].String, ShipTypeName: v[10].String, Register: v[11].String, Imo: v[12].String, YMake2: v[13].String,
			Remark: v[14].String, FormerShipName: v[15].String, Rebuilt: v[16].String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelShip, err)
	}
	return hasil, nil
}

// DaftarTipeKendaraan - lihat PembacaLookupAneka.
func (r *LookupAnekaOracle) DaftarTipeKendaraan(ctx context.Context) ([]models.BarisTipeKendaraan, error) {
	q, err := r.db.Qualify(ViewBrandDetail)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlTipeKendaraan(q), tipeAneka, BatasLookupAneka)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", ViewBrandDetail, err)
	}
	defer baris.Close()
	hasil := []models.BarisTipeKendaraan{}
	for baris.Next() {
		var id, nama sql.NullString
		if err := baris.Scan(&id, &nama); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", ViewBrandDetail, err)
		}
		hasil = append(hasil, models.BarisTipeKendaraan{ID: id.String, TypeName: nama.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", ViewBrandDetail, err)
	}
	return hasil, nil
}

// Choose Occupation Aneka `[terverifikasi]`: tombol `NB FacIn\Section\InputOkupasiAneka_FacIn.xml` L1224 -> HARNESS
// `Harness\ChooseOccupation.xml` (bukan section FIRE tiket 40): grid `D_Occupation.pxResults` (L2307), parameter DP (L3803-3806 /
// L3814-3817) BUSINESSNAME = pyWorkPage.Quotation.BusinessName, TYPE = ANEKA (TANPA kutip - FIRE `"FIRE"` berkutip; ditafsir
// literal `[dugaan]`), NAME = InputParam.CARI5 (kotak Search), ID kosong. `DataPage\D_Occupation.xml` -> RD
// `ReportDefinition\BrowseOccupationANEKA_RD.xml`: kelas Int-OCCUPATION, filter A AND B AND C AND D - A `.ID =`, B `.Type =`,
// C `.BusinessName =` (peka huruf), D `.Name Contains` (pyCaseInsensitive, pyUseNullIfEmpty=false); kolom ID, Name, BusinessCode,
// Type, BusinessName, OldID; urut `.ID` DESC; maks 500. Tabel POOLDATA.OCCUPATION (`DDL\OCCUPATION.txt`, semua VARCHAR2 -> ID
// diurut teks). Filter kosong diabaikan `[dugaan]` (pola lookup lain).
// tipeOkupasiAneka - nilai parameter TYPE harness.
const tipeOkupasiAneka = "ANEKA"

// sqlOkupasiAneka - :1 TYPE, [BUSINESSNAME], [pola nama], batas; bind berurutan.
func sqlOkupasiAneka(tabel string, denganBisnis, denganKata bool) string {
	syarat, n := "TYPE = :1", 2
	if denganBisnis {
		syarat += " AND BUSINESSNAME = :" + strconv.Itoa(n)
		n++
	}
	if denganKata {
		syarat += " AND UPPER(NAME) LIKE :" + strconv.Itoa(n) + ` ESCAPE '\'`
		n++
	}
	return "SELECT ID, NAME FROM " + tabel + " WHERE " + syarat + " ORDER BY ID DESC FETCH FIRST :" + strconv.Itoa(n) + " ROWS ONLY"
}

// CariOkupasiAneka - lihat PembacaLookupAneka.
func (r *LookupAnekaOracle) CariOkupasiAneka(ctx context.Context, businessName, kata string) ([]models.PilihanOkupasiAneka, error) {
	q, err := r.db.Qualify(TabelOccupation)
	if err != nil {
		return nil, err
	}
	arg := []any{tipeOkupasiAneka}
	if businessName != "" {
		arg = append(arg, businessName)
	}
	pola := PolaCari(kata)
	if pola != "" {
		arg = append(arg, pola)
	}
	arg = append(arg, BatasLookupAneka)
	baris, err := r.db.QueryContext(ctx, sqlOkupasiAneka(q, businessName != "", pola != ""), arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelOccupation, err)
	}
	defer baris.Close()
	hasil := []models.PilihanOkupasiAneka{}
	for baris.Next() {
		var id, nama sql.NullString
		if err := baris.Scan(&id, &nama); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelOccupation, err)
		}
		hasil = append(hasil, models.PilihanOkupasiAneka{ID: id.String, Name: nama.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelOccupation, err)
	}
	return hasil, nil
}
