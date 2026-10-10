package repository

// Untuk apa berkas ini: ACUAN - bacaan baca-saja RDB komite, ditulis ulang dari `RDBList/` / Obj-Browse korpus Komite
// Claim FacIn (pola Komite Claim Prop / Non Prop):
//
//	Obj-Browse EMAILKOMITE  roster FACIN aktif (GetKomite)                                   ApprovalKomite_Act S4 / S5.2
//	CurrencyStandard        POOLDATA.GETCURRENCYSTANDARD(:cur, SYSDATE) - kurs `.CurrencyDol` SetValueKomite S11.1.1.4
//	                        (medan tanpa kolom di Claim Fac In; sumber sama dengan SetNilaiResikoSendiri)
//	BrowseReinsuranceType_RD REINSURANCETYPE (ID, NOTE) - kolom Treaty Type SpreadingDetail / DetailAdjustmentFac
//	Obj-Browse BANKACCOUNT  IDOFBANK                                                         HitServiceToKasirKMT_Act S12
//	GetEmailCeding_SQL      GL.F_GET_EMAIL (HANYA produksi)                                  HitServiceToKasirKMT_Act S14.2
//	OperatorID.pyUserName   M_LOGIN_GO.NAME                                                  S20 InsertHistory.CARI4
//	pyEmailAddress          M_LOGIN_GO.EMAIL (inti 904) / anggota workbasket                 SendEmailKlaim_KMT S12-S15
//
// GL.F_GET_EMAIL tanpa hak di akun DEV - hanya berjalan di produksi (keputusan work owner 08-10-2026, pola Komite Claim
// Prop).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/komiteclaimfacin/backend/models"
)

// Acuan - bacaan acuan.
type Acuan struct {
	g *Gudang
	// Produksi - padanan `When/IsPEGAPROD`.
	Produksi bool
}

// AcuanDari menyusun acuan atas gudang.
func AcuanDari(g *Gudang, produksi bool) *Acuan { return &Acuan{g: g, Produksi: produksi} }

func (a *Acuan) q(objek string) (string, error) { return a.g.db.Qualify(objek) }

func rapikanDesimal(s string) string {
	switch {
	case strings.HasPrefix(s, "."):
		return "0" + s
	case strings.HasPrefix(s, "-."):
		return "-0" + s[1:]
	}
	return s
}

func dec(kol string) string { return fmt.Sprintf(db.FmtDesimal, kol) }

// satu - satu kolom satu baris; tidak ada baris = "" tanpa galat.
func (a *Acuan) satu(ctx context.Context, q string, args ...any) (string, bool, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return "", false, err
	}
	var v sql.NullString
	err := a.g.db.QueryRowContext(ctx, q, args...).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("repository: bacaan acuan komite: %w", err)
	}
	return rapikanDesimal(strings.TrimSpace(v.String)), true, nil
}

// banyak - kueri n kolom; setiap baris = []string.
func (a *Acuan) banyak(ctx context.Context, q string, n int, args ...any) ([][]string, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := a.g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: bacaan acuan komite: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out [][]string
	for rows.Next() {
		v := make([]sql.NullString, n)
		t := make([]any, n)
		for i := range v {
			t[i] = &v[i]
		}
		if err := rows.Scan(t...); err != nil {
			return nil, fmt.Errorf("repository: memindai acuan komite: %w", err)
		}
		s := make([]string, n)
		for i := range v {
			s[i] = rapikanDesimal(strings.TrimSpace(v[i].String))
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// sqlRosterKomite - Obj-Browse GetKomite (ApprovalKomite_Act S4 / S5.2): EMAILKOMITE `STS_AKTIF = '1'`,
// `STS_KLAIM = 'FACIN'`; saringan DEGREE > 1 / LIMIT_BOTTOM < total diterapkan models (`PerluasTangga`).
func sqlRosterKomite(t string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(ID), OPERATOR_ID, EMAIL, JABATAN, DEGREE, %s
		  FROM %s WHERE STS_KLAIM = :1 AND STS_AKTIF = '1' ORDER BY ID`, dec("LIMIT_BOTTOM"), t)
}

// RosterKomite membaca roster komite FACIN aktif.
func (a *Acuan) RosterKomite(ctx context.Context) ([]models.AnggotaRoster, error) {
	t, err := a.q("EMAILKOMITE")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, sqlRosterKomite(t), 6, models.LiniFacIn)
	if err != nil {
		return nil, err
	}
	out := []models.AnggotaRoster{}
	for _, r := range rows {
		out = append(out, models.AnggotaRoster{ID: r[0], OperatorID: r[1], Email: r[2], Jabatan: r[3], Degree: r[4],
			LimitBottom: r[5]})
	}
	return out, nil
}

// KursStandar = CurrencyStandard (`POOLDATA.GETCURRENCYSTANDARD(:cur, SYSDATE)`).
func (a *Acuan) KursStandar(ctx context.Context, cur string) (string, error) {
	f, err := a.q("GETCURRENCYSTANDARD")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT %s FROM DUAL`, dec(f+"(:1, SYSDATE)")), cur)
	return v, err
}

// NamaJenisReas = BrowseReinsuranceType_RD (`REINSURANCETYPE` ID -> NOTE).
func (a *Acuan) NamaJenisReas(ctx context.Context) (map[string]string, error) {
	t, err := a.q("REINSURANCETYPE")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT TO_CHAR(ID), NOTE FROM %s`, t), 2)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, r := range rows {
		out[r[0]] = r[1]
	}
	return out, nil
}

// IDBankRekening = Obj-Browse BANKACCOUNT (S12-S13, saringan NAMEOFBANK, BRANCHOFBANK, ACCOUNTNO).
func (a *Acuan) IDBankRekening(ctx context.Context, bank, cabang, akun string) (string, error) {
	t, err := a.q("BANKACCOUNT")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT IDOFBANK FROM %s WHERE NAMEOFBANK = :1 AND BRANCHOFBANK = :2
		AND ACCOUNTNO = :3 FETCH FIRST 1 ROWS ONLY`, t), bank, cabang, akun)
	return v, err
}

// StatusKonversi = getStatusKonversi_Act (HitServiceToKasirKMT_Act S3): S2 nomor akseptasi tanpa titik, S3 RDB
// `getStatusKonversi_SQL` (`COUNT(1) FROM reinsurance.trloss_detail_t WHERE NO_AKSEP`) HANYA bila IsPEGAPROD dan nomor
// terisi, S4 `> 0` -> "1". Di luar produksi kosong (S3 tidak berjalan) - kasir keluar di transisi S3.
func (a *Acuan) StatusKonversi(ctx context.Context, noAksep string) (string, error) {
	if !a.Produksi || noAksep == "" {
		return "", nil
	}
	v, _, err := a.satu(ctx, `SELECT TO_CHAR(COUNT(1)) FROM REINSURANCE.TRLOSS_DETAIL_T WHERE NO_AKSEP = :1`, noAksep)
	if err != nil {
		return "", err
	}
	if v != "" && v != "0" {
		return "1", nil
	}
	return v, nil
}

// EmailCeding = GetEmailCeding_SQL (`gl.f_get_email(:ceding) FROM DUAL`) - hanya di produksi.
func (a *Acuan) EmailCeding(ctx context.Context, ceding string) (string, error) {
	if !a.Produksi {
		return "", nil
	}
	v, _, err := a.satu(ctx, `SELECT GL.F_GET_EMAIL(:1) FROM DUAL`, ceding)
	return v, err
}

// NamaPelaku - `OperatorID.pyUserName` (M_LOGIN_GO.NAME); akun tanpa nama = ID akun.
func (a *Acuan) NamaPelaku(ctx context.Context, akun string) (string, error) {
	t, err := a.q("M_LOGIN_GO")
	if err != nil {
		return "", err
	}
	v, ada, err := a.satu(ctx, fmt.Sprintf(`SELECT NAME FROM %s WHERE LOGIN_ID = :1`, t), akun)
	if err != nil {
		return "", err
	}
	if !ada || v == "" {
		return akun, nil
	}
	return v, nil
}

// EmailPelaku - `Data-Admin-Operator-ID.pyEmailAddress` (S13 Obj-Browse pembuat) = M_LOGIN_GO.EMAIL.
func (a *Acuan) EmailPelaku(ctx context.Context, akun string) (string, error) {
	t, err := a.q("M_LOGIN_GO")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT EMAIL FROM %s WHERE LOGIN_ID = :1`, t), akun)
	return v, err
}

// sqlEmailAnggotaWorkbasket - email akun aktif pemegang workbasket aktif `:1`, unik, urut alfabet.
func sqlEmailAnggotaWorkbasket(lwb, wb, login string) string {
	return fmt.Sprintf(`SELECT DISTINCT m.EMAIL FROM %s l
		  JOIN %s w ON w.WORKBASKET_ID = l.WORKBASKET_ID
		  JOIN %s m ON m.LOGIN_ID = l.LOGIN_ID
		 WHERE l.WORKBASKET_ID = :1 AND w.IS_ACTIVE = 1 AND m.IS_ACTIVE = '1' AND m.EMAIL IS NOT NULL
		 ORDER BY m.EMAIL`, lwb, wb, login)
}

// EmailAnggotaWorkbasket - SendEmailKlaim_KMT S12 penyetuju berikut ber-KomiteID workbasket: email semua anggotanya
// (pola Komite Claim Prop 09-10-2026; EMAIL roster kosong sejak migrasi 640).
func (a *Acuan) EmailAnggotaWorkbasket(ctx context.Context, workbasket string) ([]string, error) {
	var nama [3]string
	for i, t := range []string{"M_LOGIN_GO_WORKBASKET", "M_WORKBASKET", "M_LOGIN_GO"} {
		q, err := a.q(t)
		if err != nil {
			return nil, err
		}
		nama[i] = q
	}
	rows, err := a.banyak(ctx, sqlEmailAnggotaWorkbasket(nama[0], nama[1], nama[2]), 1, workbasket)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, r := range rows {
		if r[0] != "" {
			out = append(out, r[0])
		}
	}
	return out, nil
}
