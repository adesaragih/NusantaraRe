package repository

// Daftar case NB di portal Opportunity (tiket 32). Case = baris T_WORK_POLIS ber-LINI
// LiniFacIn (case Life tidak tampil) + T_NB_OPPORTUNITY + nama tertanggung T_M_ACCOUNT +
// nama marketing officer dari T_QUOTATIONDATA.MOID (tiket 31). Pencarian "mengandung",
// tidak peka huruf (PolaCari), atas ID case dan BUSINESS_PROSPECT_NAME - placeholder Pega
// "NB-1234 or Name" (`Section\SFAPortal_OpportunitiesList.xml`). Parameter terikat.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// PembacaPortal - satu halaman daftar case NB beserta cacah seluruhnya.
type PembacaPortal interface {
	CariPortal(ctx context.Context, cari string, offset, ukuran int) ([]models.BarisPortal, int, error)
}

// PortalOracle - PembacaPortal atas Oracle.
type PortalOracle struct{ db *db.DB }

// NewPortalOracle merakit pembaca daftar portal.
func NewPortalOracle(d *db.DB) *PortalOracle { return &PortalOracle{db: d} }

// tabelPortal - nama tabel yang sudah dikualifikasi.
type tabelPortal struct{ work, opp, akun, quo, mo string }

// saringPortal - LINI selalu (:1); bila cari, :2/:3 bernilai pola yang sama.
func saringPortal(cari bool) string {
	if !cari {
		return " WHERE w.LINI = :1"
	}
	return ` WHERE w.LINI = :1 AND (UPPER(w.ID) LIKE :2 ESCAPE '\' OR UPPER(o.BUSINESS_PROSPECT_NAME) LIKE :3 ESCAPE '\')`
}

// sqlCariPortal - urut terbaru dulu: TGL_CREATE DESC, ID DESC (pemutus seri). Nama
// tertanggung dan marketing lewat subkueri MAX - T_M_ACCOUNT dan MARKETINGOFFICER tanpa PK
// di DDL-nya (pola A88).
func sqlCariPortal(t tabelPortal, cari bool) string {
	n := 2
	if cari {
		n = 4
	}
	return fmt.Sprintf(`SELECT w.ID, o.BUSINESS_PROSPECT_NAME, o.GROUP_BUSINESS,
  (SELECT MAX(a.INSUREDNAME) FROM %s a WHERE a.ID = o.ACCOUNT_ID),
  (SELECT MAX(m.CLIENTNAME) FROM %s m WHERE m.ID = q.MOID),
  w.STATUS_WORK
FROM %s w
LEFT JOIN %s o ON o.ID = w.ID
LEFT JOIN %s q ON q.PARENT_ID = w.ID%s
ORDER BY w.TGL_CREATE DESC, w.ID DESC OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY`, t.akun, t.mo, t.work, t.opp, t.quo, saringPortal(cari), n, n+1)
}

// sqlCacahPortal - cacah seluruh baris yang cocok (tanpa T_QUOTATIONDATA: satu baris per case).
func sqlCacahPortal(t tabelPortal, cari bool) string {
	return "SELECT COUNT(*) FROM " + t.work + " w LEFT JOIN " + t.opp + " o ON o.ID = w.ID" + saringPortal(cari)
}

// CariPortal - lihat PembacaPortal.
func (r *PortalOracle) CariPortal(ctx context.Context, cari string, offset, ukuran int) ([]models.BarisPortal, int, error) {
	var t tabelPortal
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelOpportunity, &t.opp}, {TabelAkun, &t.akun}, {TabelQuotationData, &t.quo},
		{TabelMarketingOfficer, &t.mo}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return nil, 0, err
		}
		*x.ke = q
	}
	pola := PolaCari(cari)
	arg := []any{LiniFacIn}
	if pola != "" {
		arg = append(arg, pola, pola)
	}
	var total int
	if err := r.db.QueryRowContext(ctx, sqlCacahPortal(t, pola != ""), arg...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: cacah daftar case NB: %w", err)
	}
	baris, err := r.db.QueryContext(ctx, sqlCariPortal(t, pola != ""), append(arg, offset, ukuran)...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: baca daftar case NB: %w", err)
	}
	defer baris.Close()
	hasil := []models.BarisPortal{}
	for baris.Next() {
		var caseID, nama, grup, tertanggung, marketing, status sql.NullString
		if err := baris.Scan(&caseID, &nama, &grup, &tertanggung, &marketing, &status); err != nil {
			return nil, 0, fmt.Errorf("repository: daftar case NB: %w", err)
		}
		hasil = append(hasil, models.BarisPortal{CaseID: caseID.String, Name: nama.String, GroupBusiness: grup.String,
			InsuredName: tertanggung.String, Marketing: marketing.String, Status: status.String})
	}
	if err := baris.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository: daftar case NB: %w", err)
	}
	return hasil, total, nil
}
