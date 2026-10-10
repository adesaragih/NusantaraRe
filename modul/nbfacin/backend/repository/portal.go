package repository

// Daftar case NB di portal Opportunity (tiket 32). Case = baris T_WORK_POLIS ber-LINI
// LiniFacIn (case Life tidak tampil) + T_NB_OPPORTUNITY + nama tertanggung T_M_ACCOUNT +
// nama marketing officer dari T_QUOTATIONDATA.MOID (tiket 31). Pencarian "mengandung",
// tidak peka huruf (PolaCari), atas ID case dan BUSINESS_PROSPECT_NAME - placeholder Pega
// "NB-1234 or Name" (`Section\SFAPortal_OpportunitiesList.xml`). Parameter terikat. Hanya case buatan pelaku sendiri
// (T_WORK_POLIS.CREATE_OP = pengenal akun; keputusan work owner 07-10-2026 "list yang muncul hanya list yg dia buat sndri").

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// PembacaPortal - satu halaman daftar case NB buatan akun `pembuat` beserta cacah seluruhnya (sesudah saringan).
type PembacaPortal interface {
	CariPortal(ctx context.Context, pembuat, cari string, offset, ukuran int) ([]models.BarisPortal, int, error)
}

// PortalOracle - PembacaPortal atas Oracle.
type PortalOracle struct{ db *db.DB }

// NewPortalOracle merakit pembaca daftar portal.
func NewPortalOracle(d *db.DB) *PortalOracle { return &PortalOracle{db: d} }

// tabelPortal - nama tabel yang sudah dikualifikasi.
type tabelPortal struct{ work, opp, akun, quo, mo, general string }

// saringPortal - LINI (:1) dan pembuat (:2) selalu; bila cari, :3/:4 bernilai pola yang sama.
func saringPortal(cari bool) string {
	if !cari {
		return " WHERE w.LINI = :1 AND w.CREATE_OP = :2"
	}
	return ` WHERE w.LINI = :1 AND w.CREATE_OP = :2 AND (UPPER(w.ID) LIKE :3 ESCAPE '\' OR UPPER(o.BUSINESS_PROSPECT_NAME) LIKE :4 ESCAPE '\')`
}

// sqlCariPortal - urut terbaru dulu: TGL_CREATE DESC, ID DESC (pemutus seri). Nama
// tertanggung dan marketing lewat subkueri MAX - T_M_ACCOUNT dan MARKETINGOFFICER tanpa PK
// di DDL-nya (pola A88). NB status = T_GENERAL_POLIS.NB_STATUS (PK = ID kasus; WAJIB migrasi 208 - tanpa kolomnya kueri
// ini ORA-00904).
func sqlCariPortal(t tabelPortal, cari bool) string {
	n := 3
	if cari {
		n = 5
	}
	return fmt.Sprintf(`%s%s
ORDER BY w.TGL_CREATE DESC, w.ID DESC OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY`, kepalaBarisCase(t), saringPortal(cari), n, n+1)
}

// kepalaBarisCase - SELECT + JOIN satu baris case NB (bentuk BarisPortal) untuk portal dan kotak masuk Beranda; dibaca
// bacaBarisCase.
func kepalaBarisCase(t tabelPortal) string {
	return fmt.Sprintf(`SELECT w.ID, o.BUSINESS_PROSPECT_NAME, o.GROUP_BUSINESS,
  (SELECT MAX(a.INSUREDNAME) FROM %s a WHERE a.ID = o.ACCOUNT_ID),
  (SELECT MAX(m.CLIENTNAME) FROM %s m WHERE m.ID = q.MOID),
  w.STATUS_WORK, g.NB_STATUS
FROM %s w
LEFT JOIN %s o ON o.ID = w.ID
LEFT JOIN %s q ON q.PARENT_ID = w.ID
LEFT JOIN %s g ON g.ID = w.ID`, t.akun, t.mo, t.work, t.opp, t.quo, t.general)
}

// tabelBarisCase - nama tabel kepalaBarisCase berkualifikasi skema.
func tabelBarisCase(d *db.DB) (tabelPortal, error) {
	var t tabelPortal
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelOpportunity, &t.opp}, {TabelAkun, &t.akun}, {TabelQuotationData, &t.quo},
		{TabelMarketingOfficer, &t.mo}, {TabelGeneralPolis, &t.general}} {
		q, err := d.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

// bacaBarisCase - baris kepalaBarisCase -> BarisPortal (NULL = "").
func bacaBarisCase(baris *sql.Rows) ([]models.BarisPortal, error) {
	defer baris.Close()
	hasil := []models.BarisPortal{}
	for baris.Next() {
		var caseID, nama, grup, tertanggung, marketing, status, nbStatus sql.NullString
		if err := baris.Scan(&caseID, &nama, &grup, &tertanggung, &marketing, &status, &nbStatus); err != nil {
			return nil, fmt.Errorf("repository: daftar case NB: %w", err)
		}
		hasil = append(hasil, models.BarisPortal{CaseID: caseID.String, Name: nama.String, GroupBusiness: grup.String,
			InsuredName: tertanggung.String, Marketing: marketing.String, Status: status.String, NBStatus: nbStatus.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: daftar case NB: %w", err)
	}
	return hasil, nil
}

// sqlCacahPortal - cacah seluruh baris yang cocok (tanpa T_QUOTATIONDATA: satu baris per case).
func sqlCacahPortal(t tabelPortal, cari bool) string {
	return "SELECT COUNT(*) FROM " + t.work + " w LEFT JOIN " + t.opp + " o ON o.ID = w.ID" + saringPortal(cari)
}

// CariPortal - lihat PembacaPortal.
func (r *PortalOracle) CariPortal(ctx context.Context, pembuat, cari string, offset, ukuran int) ([]models.BarisPortal, int, error) {
	t, err := tabelBarisCase(r.db)
	if err != nil {
		return nil, 0, err
	}
	pola := PolaCari(cari)
	arg := []any{LiniFacIn, pembuat}
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
	hasil, err := bacaBarisCase(baris)
	if err != nil {
		return nil, 0, err
	}
	return hasil, total, nil
}
