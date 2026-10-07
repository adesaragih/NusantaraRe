package repository

// Submit tahap marketing (tiket 57). Bukti `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\`:
//   - riwayat: `RDBList\InsertViewSuggest_SQL.xml` (POOLDATA.HISTORYAKSEPTASIPRODUCTION, mapping `Activity\SaveViewSuggest.xml`),
//     `RDBList\InsertHistoryAkseptasiPega_Sql.xml` (HISTORYAKSEPTASIPEGA tanpa awalan skema = POOLDATA, data DBA
//     2026-09-18 komiteclaimprop K6), `RDBList\InsertDataFacinPerformance_SQL.xml` (POOLDATA.FACINPERFORMANCE);
//   - pilihan: `RDBList\GetAksepBanding_SQL.xml` (verbatim), `RDBList\GetCauseOfDecline.xml` / `GetCauseOfDecline2.xml`;
//   - proteksi: `RDBList\CheckZipCode_SQL.xml` + Obj-Browse kelas Int-RW (ProtectFIREMBUPA_Act 9.2.11-16),
//     `RDBList\GetCategoryOccupation.xml`.
// Tabel lama TIDAK dibuat ulang (keputusan work owner 07-10-2026 butir 4, ADR-0042 butir 2); kolom = DDL `D:\migrasi\RNM\DDL\`.
// Nol COMMIT (Pega memakai BEGIN ... COMMIT; END; - di sini satu transaksi milik pemanggil, ADR-U-0029). Nol CALL (ADR-0043).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelRiwayatProduksi, TabelRiwayatPega, TabelPerformanceFacIn - tabel lama POOLDATA, ditulis Submit (tiket 57).
	TabelRiwayatProduksi  = "HISTORYAKSEPTASIPRODUCTION"
	TabelRiwayatPega      = "HISTORYAKSEPTASIPEGA"
	TabelPerformanceFacIn = "FACINPERFORMANCE"
	// TabelLimitProperty - POOLDATA.M_LIMIT_PROPERTYY (pilihan Banding to), TabelSebabDecline - POOLDATA.CAUSEOFDECLINE.
	TabelLimitProperty = "M_LIMIT_PROPERTYY"
	TabelSebabDecline  = "CAUSEOFDECLINE"
	// WorkbasketMarketing - PositionNote tahap marketing (Assignment12), dipakai saringan GetAksepBanding_SQL.
	WorkbasketMarketing = "ReasFacInMarketing"
)

// PenyimpanSubmit - baca / tulis Submit tahap marketing.
type PenyimpanSubmit interface {
	// KonteksSubmit - keadaan kasus; ErrKasusTidakAda bila case FAC tidak ada.
	KonteksSubmit(ctx context.Context, id string) (models.KonteksSubmit, error)
	// AkunSubmit - NAME / DIVISION_CODE M_LOGIN_GO; akun tak ada -> kosong.
	AkunSubmit(ctx context.Context, akunID string) (models.AkunSubmit, error)
	// RiwayatAkseptasi - baris HISTORYAKSEPTASIPRODUCTION ber-IDPEGA itu, terbaru dulu.
	RiwayatAkseptasi(ctx context.Context, idPega string) ([]models.RiwayatAkseptasi, error)
	// OpsiBanding - GetAksepBanding_SQL: JABATAN M_LIMIT_PROPERTYY untuk workbasket non-marketing TERAKHIR kasus dan
	// TEAM_GROUP-nya, urut hasil kueri.
	OpsiBanding(ctx context.Context, idPega, teamGroup string) ([]string, error)
	// SebabDecline - seluruh baris CAUSEOFDECLINE; saringan GetCauseOfDecline / GetCauseOfDecline2 / opsi per T/C diterapkan
	// layanan atas baris yang sama.
	SebabDecline(ctx context.Context) ([]models.BarisCauseOfDecline, error)
	// CekRW - kode pos di RW, dan pasangannya dengan CITYNAME / DISTRICTNAME / PROVINCENAME.
	CekRW(ctx context.Context, zip, kota, kecamatan, provinsi string) (models.HasilCekRW, error)
	// KategoriOkupasi - OCCUPATION.CATEGORY per OCCUPATIONCODE (baris pertama urut CATEGORY); kode tak ada -> tidak di peta.
	KategoriOkupasi(ctx context.Context, kode []string) (map[string]string, error)
	// TulisSubmit - riwayat, performance, posisi, NB_STATUS, kunci korespondensi; ErrPosisiBerubah bila posisi kasus sudah
	// bukan PosisiLama.
	TulisSubmit(ctx context.Context, tx *db.Tx, t models.TulisanSubmit) error
}

// ErrPosisiBerubah - posisi kasus berubah di antara baca dan tulis (Submit serentak).
var ErrPosisiBerubah = errors.New("repository: posisi kasus sudah berubah")

// SubmitOracle - PenyimpanSubmit atas Oracle.
type SubmitOracle struct{ db *db.DB }

// NewSubmitOracle merakit penyimpan Submit.
func NewSubmitOracle(d *db.DB) *SubmitOracle { return &SubmitOracle{db: d} }

// tabelSubmit - nama berkualifikasi skema (ADR-U-0033).
type tabelSubmit struct {
	work, opp, general, quo, akun, mo, login, produksi, pega, perf, limit, sebab, rw, okupasi, kor, pen string
}

func (r *SubmitOracle) tabel() (tabelSubmit, error) {
	var t tabelSubmit
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelOpportunity, &t.opp}, {TabelGeneralPolis, &t.general}, {TabelQuotationData, &t.quo},
		{TabelAkun, &t.akun}, {TabelMarketingOfficer, &t.mo}, {TabelLoginGo, &t.login}, {TabelRiwayatProduksi, &t.produksi},
		{TabelRiwayatPega, &t.pega}, {TabelPerformanceFacIn, &t.perf}, {TabelLimitProperty, &t.limit},
		{TabelSebabDecline, &t.sebab}, {TabelRW, &t.rw}, {TabelOccupation, &t.okupasi}, {TabelKorespondensi, &t.kor},
		{TabelPenerimaKorespondensi, &t.pen}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

// ---- SQL ---------------------------------------------------------------------------------------------------------

// sqlKonteksSubmit - kasus + marketing (MAX: MARKETINGOFFICER.ID tidak ber-PK di DDL-nya). MOLEADER memakai saringan
// GetDataMarketing_SQL (MOSTATUS '1', BRANCHSTATUS kosong, CLIENTID2 bukan LEADER) atas ID = MOID, bukan CLIENTNAME LIKE.
func sqlKonteksSubmit(t tabelSubmit) string {
	return `SELECT w.POSITION, w.FLAG_ONGOING_POLICY, q.MOID,
 (SELECT MAX(m.CLIENTNAME) FROM ` + t.mo + ` m WHERE m.ID = q.MOID),
 (SELECT MAX(m.TEAMGROUP) FROM ` + t.mo + ` m WHERE m.ID = q.MOID),
 (SELECT MAX(m.MOLEADER) FROM ` + t.mo + ` m WHERE m.ID = q.MOID AND m.MOSTATUS = '1' AND m.BRANCHSTATUS IS NULL AND m.CLIENTID2 <> 'LEADER'),
 (SELECT MAX(a.INSUREDNAME) FROM ` + t.akun + ` a WHERE a.ID = o.ACCOUNT_ID),
 q.SOB_NAME, o.CLASS_OF_BUSINESS, ` + angkaKeluar("g.PERCENT_SHARE") + `, ` + angkaKeluar("g.PAY_RI_COMMISION") + `
FROM ` + t.work + ` w
LEFT JOIN ` + t.opp + ` o ON o.ID = w.ID
LEFT JOIN ` + t.general + ` g ON g.ID = w.ID
LEFT JOIN ` + t.quo + ` q ON q.PARENT_ID = w.ID
WHERE w.ID = :1 AND w.LINI = :2`
}

func sqlAkunSubmit(login string) string {
	return "SELECT NAME, DIVISION_CODE FROM " + login + " WHERE LOGIN_ID = :1"
}

// sqlRiwayatAkseptasi - grid ViewSuggest; TGL_INP ditulis sebagai jam dinding WIB (To_date teks Asia/Jakarta).
func sqlRiwayatAkseptasi(produksi string) string {
	return "SELECT TO_CHAR(TGL_INP, 'DD-MM-YYYY HH24:MI'), PIC, APPROVAL, POSISI, KETERANGAN FROM " + produksi +
		" WHERE IDPEGA = :1 ORDER BY TGL_INP DESC, NOURUT DESC"
}

// sqlOpsiBanding - GetAksepBanding_SQL verbatim (bind menggantikan {pyWorkPage.pzInsKey} / {...TeamGroup}).
func sqlOpsiBanding(limit, pega string) string {
	return `SELECT JABATAN FROM ` + limit + ` WHERE WORKBASKET IN (SELECT WORKBASKET FROM (
  SELECT WORKBASKET FROM ` + pega + ` WHERE ID_PEGA = :1 AND WORKBASKET != :2 ORDER BY TGL_TRANSFER DESC
) WHERE ROWNUM = 1) AND TEAM_GROUP = :3`
}

// sqlSebabDecline - seluruh baris; GetCauseOfDecline (`DESCRIPTION IS NULL`) dan GetCauseOfDecline2 (`TERMCONDITION2 IS NOT
// NULL`) diterapkan layanan atas baris yang sama dengan opsi per T/C.
func sqlSebabDecline(sebab string) string {
	return "SELECT ID, TERMCONDITION, DESCRIPTION, TERMCONDITION2 FROM " + sebab
}

// sqlCekRW - CheckZipCode_SQL (`select count(*) from rw where zipcode = ...`) + tiga Obj-Browse Int-RW (ZipCode dan satu
// kolom lain) dalam satu kueri. Bind berurutan menurut kemunculan: kota, kecamatan, provinsi, kode pos.
func sqlCekRW(rw string) string {
	return `SELECT COUNT(*),
 COUNT(CASE WHEN CITYNAME = :1 THEN 1 END),
 COUNT(CASE WHEN DISTRICTNAME = :2 THEN 1 END),
 COUNT(CASE WHEN PROVINCENAME = :3 THEN 1 END)
FROM ` + rw + ` WHERE ZIPCODE = :4`
}

// sqlKategoriOkupasi - GetCategoryOccupation (`select distinct CATEGORY ... where OCCUPATIONCODE = ...`) untuk :1..:n.
func sqlKategoriOkupasi(okupasi string, n int) string {
	b := make([]string, n)
	for i := range b {
		b[i] = fmt.Sprintf(":%d", i+1)
	}
	return "SELECT DISTINCT OCCUPATIONCODE, CATEGORY FROM " + okupasi + " WHERE OCCUPATIONCODE IN (" + strings.Join(b, ", ") +
		") ORDER BY OCCUPATIONCODE, CATEGORY"
}

// sqlSisipRiwayatProduksi - InsertViewSuggest_SQL; KETERANGAN dipotong 3990 seperti Pega (substr).
func sqlSisipRiwayatProduksi(produksi string) string {
	return `INSERT INTO ` + produksi + ` (IDPEGA, TYPE_POLIS, NOURUT, POSISI, PIC, TGL_INP, DIV, TYPE, PUTARAN, APPROVAL, KETERANGAN,
 AKSES_LOGIN, B2B, BUSINESS_CODE, PERCENT_RNM)
SELECT :1, :2, NVL(MAX(NOURUT), 0) + 1, :3, :4, TO_DATE(:5, 'DD/MM/YYYY HH24:MI:SS'), :6, :7, :8, :9, SUBSTR(:10, 1, 3990),
 :11, NULL, NULL, ` + angkaMasuk(":12") + `
FROM ` + produksi + ` WHERE IDPEGA = :13`
}

// sqlSisipRiwayatPega - InsertHistoryAkseptasiPega_Sql (ID_KOMITE tidak diisi jalur ini: InsertHistory.CARI6 tidak pernah
// di-set di Activity\InsertHistoryAkseptasiPega.xml).
func sqlSisipRiwayatPega(pega string) string {
	return "INSERT INTO " + pega + " (ID_PEGA, TGL_TRANSFER, STATUS, USERNAME, WORKBASKET) VALUES (:1, SYSDATE, :2, :3, :4)"
}

// sqlSisipPerformance - InsertDataFacinPerformance_SQL (RISLIPRECEIVED NULL).
func sqlSisipPerformance(perf string) string {
	return `INSERT INTO ` + perf + ` (IDPEGA, POSITIONDOC, STATUSDOC, DATETIME, USERNAME, INSUREDNAME, SOB, MOLEADER, COB, RISLIPRECEIVED)
VALUES (:1, :2, :3, SYSDATE, :4, :5, :6, :7, :8, NULL)`
}

// sqlPosisiSubmit - POSITION baru; STATUS_WORK baru bila diisi (NULL = tetap).
func sqlPosisiSubmit(work string) string {
	return "UPDATE " + work + " SET POSITION = :1, STATUS_WORK = NVL(:2, STATUS_WORK), TGL_UPDATE = SYSDATE" +
		" WHERE ID = :3 AND LINI = :4 AND POSITION = :5"
}

// sqlNBStatus / sqlSisipNBStatus - T_GENERAL_POLIS.NB_STATUS (migrasi 208); baris General belum ada (UPDATE 0 baris) ->
// dibuat ber-ID saja.
func sqlNBStatus(general string) string {
	return "UPDATE " + general + " SET NB_STATUS = :1 WHERE ID = :2"
}

func sqlSisipNBStatus(general string) string {
	return "INSERT INTO " + general + " (ID, NB_STATUS) VALUES (:1, :2)"
}

// sqlKunciKorespondensi / sqlKunciPenerima - SetCorrespondenceReadOnly: `.OfferFacIn.CorrespondenceList(*).IsEditable = "0"`
// (penerima ikut induk, pola TandaiTerkirim).
func sqlKunciKorespondensi(kor string) string {
	return "UPDATE " + kor + " SET IS_EDITABLE = '0' WHERE PARENT_ID = :1 AND IS_EDITABLE <> '0'"
}

func sqlKunciPenerima(pen, kor string) string {
	return "UPDATE " + pen + " SET IS_EDITABLE = '0' WHERE PARENT_ID IN (SELECT ID FROM " + kor + " WHERE PARENT_ID = :1) AND IS_EDITABLE <> '0'"
}

// ---- baca --------------------------------------------------------------------------------------------------------

// KonteksSubmit - lihat PenyimpanSubmit.
func (r *SubmitOracle) KonteksSubmit(ctx context.Context, id string) (models.KonteksSubmit, error) {
	t, err := r.tabel()
	if err != nil {
		return models.KonteksSubmit{}, err
	}
	var v [11]sql.NullString
	tujuan := make([]any, len(v))
	for i := range v {
		tujuan[i] = &v[i]
	}
	err = r.db.QueryRowContext(ctx, sqlKonteksSubmit(t), id, LiniFacIn).Scan(tujuan...)
	if errors.Is(err, sql.ErrNoRows) {
		return models.KonteksSubmit{}, ErrKasusTidakAda
	}
	if err != nil {
		return models.KonteksSubmit{}, fmt.Errorf("repository: baca kasus Submit: %w", err)
	}
	k := models.KonteksSubmit{Position: v[0].String, FlagOnGoingPolicy: v[1].String, MOID: v[2].String,
		MarketingName: v[3].String, TeamGroup: v[4].String, MOLeader: v[5].String, InsuredName: v[6].String,
		SobName: v[7].String, ClassOfBusiness: v[8].String}
	if k.PercentShare, err = bacaDesimal(id, "T_GENERAL_POLIS.PERCENT_SHARE", &v[9]); err != nil {
		return models.KonteksSubmit{}, err
	}
	if k.RICommision, err = bacaDesimal(id, "T_GENERAL_POLIS.PAY_RI_COMMISION", &v[10]); err != nil {
		return models.KonteksSubmit{}, err
	}
	return k, nil
}

// AkunSubmit - lihat PenyimpanSubmit.
func (r *SubmitOracle) AkunSubmit(ctx context.Context, akunID string) (models.AkunSubmit, error) {
	login, err := r.db.Qualify(TabelLoginGo)
	if err != nil {
		return models.AkunSubmit{}, err
	}
	var nama, divisi sql.NullString
	err = r.db.QueryRowContext(ctx, sqlAkunSubmit(login), akunID).Scan(&nama, &divisi)
	if errors.Is(err, sql.ErrNoRows) {
		return models.AkunSubmit{}, nil
	}
	if err != nil {
		return models.AkunSubmit{}, fmt.Errorf("repository: baca %s: %w", TabelLoginGo, err)
	}
	return models.AkunSubmit{Nama: nama.String, Divisi: divisi.String}, nil
}

// teksKolom - satu kolom teks per baris (NULL dilewati).
func (r *SubmitOracle) teksKolom(ctx context.Context, tabel, q string, arg ...any) ([]string, error) {
	baris, err := r.db.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", tabel, err)
	}
	defer baris.Close()
	hasil := []string{}
	for baris.Next() {
		var s sql.NullString
		if err := baris.Scan(&s); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", tabel, err)
		}
		if s.Valid {
			hasil = append(hasil, s.String)
		}
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", tabel, err)
	}
	return hasil, nil
}

// RiwayatAkseptasi - lihat PenyimpanSubmit.
func (r *SubmitOracle) RiwayatAkseptasi(ctx context.Context, idPega string) ([]models.RiwayatAkseptasi, error) {
	produksi, err := r.db.Qualify(TabelRiwayatProduksi)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlRiwayatAkseptasi(produksi), idPega)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelRiwayatProduksi, err)
	}
	defer baris.Close()
	hasil := []models.RiwayatAkseptasi{}
	for baris.Next() {
		var tgl, pic, approval, posisi, ket sql.NullString
		if err := baris.Scan(&tgl, &pic, &approval, &posisi, &ket); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelRiwayatProduksi, err)
		}
		hasil = append(hasil, models.RiwayatAkseptasi{Tanggal: tgl.String, PIC: pic.String, Approval: approval.String,
			Status: posisi.String, Comment: ket.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelRiwayatProduksi, err)
	}
	return hasil, nil
}

// OpsiBanding - lihat PenyimpanSubmit.
func (r *SubmitOracle) OpsiBanding(ctx context.Context, idPega, teamGroup string) ([]string, error) {
	t, err := r.tabel()
	if err != nil {
		return nil, err
	}
	return r.teksKolom(ctx, TabelLimitProperty, sqlOpsiBanding(t.limit, t.pega), idPega, WorkbasketMarketing, teamGroup)
}

// SebabDecline - lihat PenyimpanSubmit.
func (r *SubmitOracle) SebabDecline(ctx context.Context) ([]models.BarisCauseOfDecline, error) {
	sebab, err := r.db.Qualify(TabelSebabDecline)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlSebabDecline(sebab))
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelSebabDecline, err)
	}
	defer baris.Close()
	hasil := []models.BarisCauseOfDecline{}
	for baris.Next() {
		var id, tc, desc, tc2 sql.NullString
		if err := baris.Scan(&id, &tc, &desc, &tc2); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelSebabDecline, err)
		}
		hasil = append(hasil, models.BarisCauseOfDecline{ID: id.String, TermCondition: tc.String, Description: desc.String,
			TermCondition2: tc2.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelSebabDecline, err)
	}
	return hasil, nil
}

// CekRW - lihat PenyimpanSubmit.
func (r *SubmitOracle) CekRW(ctx context.Context, zip, kota, kecamatan, provinsi string) (models.HasilCekRW, error) {
	rw, err := r.db.Qualify(TabelRW)
	if err != nil {
		return models.HasilCekRW{}, err
	}
	var n [4]int
	if err := r.db.QueryRowContext(ctx, sqlCekRW(rw), kota, kecamatan, provinsi, zip).Scan(&n[0], &n[1], &n[2], &n[3]); err != nil {
		return models.HasilCekRW{}, fmt.Errorf("repository: baca %s: %w", TabelRW, err)
	}
	return models.HasilCekRW{Zip: n[0] > 0, City: n[1] > 0, District: n[2] > 0, Province: n[3] > 0}, nil
}

// KategoriOkupasi - lihat PenyimpanSubmit.
func (r *SubmitOracle) KategoriOkupasi(ctx context.Context, kode []string) (map[string]string, error) {
	hasil := map[string]string{}
	if len(kode) == 0 {
		return hasil, nil
	}
	okupasi, err := r.db.Qualify(TabelOccupation)
	if err != nil {
		return nil, err
	}
	arg := make([]any, len(kode))
	for i, k := range kode {
		arg[i] = k
	}
	baris, err := r.db.QueryContext(ctx, sqlKategoriOkupasi(okupasi, len(kode)), arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelOccupation, err)
	}
	defer baris.Close()
	for baris.Next() {
		var k, c sql.NullString
		if err := baris.Scan(&k, &c); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelOccupation, err)
		}
		if _, ada := hasil[k.String]; !ada {
			hasil[k.String] = c.String
		}
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelOccupation, err)
	}
	return hasil, nil
}

// ---- tulis -------------------------------------------------------------------------------------------------------

// TulisSubmit - lihat PenyimpanSubmit. Urutan: posisi (penjaga Submit serentak) -> riwayat produksi -> riwayat Pega ->
// performance -> NB_STATUS (UPDATE, atau INSERT bila baris General belum ada) -> kunci korespondensi.
func (r *SubmitOracle) TulisSubmit(ctx context.Context, tx *db.Tx, w models.TulisanSubmit) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	k := db.KosongJadiNil
	hasil, err := jalankan(ctx, tx, sqlPosisiSubmit(t.work), "posisi Submit", w.PosisiBaru, k(w.StatusWorkBaru), w.CaseID, LiniFacIn,
		w.PosisiLama)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return fmt.Errorf("repository: posisi Submit: %w", err)
	} else if n != 1 {
		return ErrPosisiBerubah
	}
	p := w.Riwayat
	if _, err := jalankan(ctx, tx, sqlSisipRiwayatProduksi(t.produksi), "riwayat "+TabelRiwayatProduksi, w.IDPega, k(p.TypePolis),
		k(p.Posisi), k(p.PIC), p.TglInp, k(p.Div), k(p.Type), k(p.Putaran), k(p.Approval), k(p.Keterangan), k(p.AksesLogin),
		ikatDesimal(p.PercentRNM), w.IDPega); err != nil {
		return err
	}
	if _, err := jalankan(ctx, tx, sqlSisipRiwayatPega(t.pega), "riwayat "+TabelRiwayatPega, w.IDPega, k(w.Pega.Status),
		k(w.Pega.Username), k(w.Pega.Workbasket)); err != nil {
		return err
	}
	if f := w.Performance; f != nil {
		if _, err := jalankan(ctx, tx, sqlSisipPerformance(t.perf), TabelPerformanceFacIn, w.IDPega, k(f.PositionDoc), k(f.StatusDoc),
			k(f.Username), k(f.InsuredName), k(f.SOB), k(f.MOLeader), k(f.COB)); err != nil {
			return err
		}
	}
	hasil, err = jalankan(ctx, tx, sqlNBStatus(t.general), "NB_STATUS", k(w.NBStatus), w.CaseID)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return fmt.Errorf("repository: NB_STATUS: %w", err)
	} else if n == 0 {
		if _, err := jalankan(ctx, tx, sqlSisipNBStatus(t.general), "NB_STATUS", w.CaseID, k(w.NBStatus)); err != nil {
			return err
		}
	}
	if _, err := jalankan(ctx, tx, sqlKunciPenerima(t.pen, t.kor), "kunci penerima korespondensi", w.CaseID); err != nil {
		return err
	}
	_, err = jalankan(ctx, tx, sqlKunciKorespondensi(t.kor), "kunci korespondensi", w.CaseID)
	return err
}
