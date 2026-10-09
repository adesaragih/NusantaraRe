package repository

// Untuk apa berkas ini: ACUAN - bacaan baca-saja RDB komite, ditulis ulang dari `RDBList/` korpus Komite Claim Non Prop
// (pola Komite Claim Prop):
//
//	Obj-Browse BANKACCOUNT  IDOFBANK                                                           HitServiceToKasirKMT_Act S12
//	GetEmailCeding_SQL      GL.F_GET_EMAIL (HANYA saat muatan Kasir disusun di produksi)       HitServiceToKasirKMT_Act S10.2
//	OperatorID.pyUserName   M_LOGIN_GO.NAME                                                    S8-S10, S22 InsertHistory.CARI4
//	pyEmailAddress          M_LOGIN_GO.EMAIL (inti 904)                                        SendEmailKlaim_KMT
//
//	getStatusKonversi_SQL   REINSURANCE.TRLOSS_DETAIL_T (HANYA produksi, gerbang IsPEGAPROD)  HitServiceToKasirKMT_Act S3
//
// Dua bacaan lintas skema tanpa hak di akun DEV (TRLOSS_DETAIL_T ORA-00942, F_GET_EMAIL ORA-00904) - keduanya hanya
// berjalan di produksi (keputusan work owner 08-10-2026, pola Komite Claim Prop); hak DBA diminta.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
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

// IDBankRekening = Obj-Browse BANKACCOUNT (S12-S13, `BankAccount.pxResults(1).IDOFBANK`; saringan NAMEOFBANK,
// BRANCHOFBANK, ACCOUNTNO).
func (a *Acuan) IDBankRekening(ctx context.Context, bank, cabang, akun string) (string, error) {
	t, err := a.q("BANKACCOUNT")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT IDOFBANK FROM %s WHERE NAMEOFBANK = :1 AND BRANCHOFBANK = :2
		AND ACCOUNTNO = :3 FETCH FIRST 1 ROWS ONLY`, t), bank, cabang, akun)
	return v, err
}

// StatusKonversi = getStatusKonversi_Act (HitServiceToKasirKMT_Act S3): S3 RDB `getStatusKonversi_SQL`
// (REINSURANCE.TRLOSS_DETAIL_T) HANYA bila IsPEGAPROD dan nomor tidak kosong - di DEV tidak dijalankan (keputusan work
// owner 08-10-2026 "tambahkan when ispegaprod untuk dev jangan jalankan"); S4 `> 0` -> "1".
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

// EmailPelaku - `Data-Admin-Operator-ID.pyEmailAddress` (SendEmailKlaim_KMT Obj-Browse pembuat) = M_LOGIN_GO.EMAIL;
// akun tanpa email = kosong.
func (a *Acuan) EmailPelaku(ctx context.Context, akun string) (string, error) {
	t, err := a.q("M_LOGIN_GO")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, sqlEmailPelaku(t), akun)
	return v, err
}

func sqlEmailPelaku(t string) string {
	return fmt.Sprintf(`SELECT EMAIL FROM %s WHERE LOGIN_ID = :1`, t)
}

// sqlEmailAnggotaWorkbasket - email akun aktif pemegang workbasket aktif `:1`, unik, urut alfabet.
func sqlEmailAnggotaWorkbasket(lwb, wb, login string) string {
	return fmt.Sprintf(`SELECT DISTINCT m.EMAIL FROM %s l
		  JOIN %s w ON w.WORKBASKET_ID = l.WORKBASKET_ID
		  JOIN %s m ON m.LOGIN_ID = l.LOGIN_ID
		 WHERE l.WORKBASKET_ID = :1 AND w.IS_ACTIVE = 1 AND m.IS_ACTIVE = '1' AND m.EMAIL IS NOT NULL
		 ORDER BY m.EMAIL`, lwb, wb, login)
}

// EmailAnggotaWorkbasket - SendEmailKlaim_KMT penyetuju berikut ber-KomiteID workbasket: email semua anggotanya
// (keputusan work owner 09-10-2026; `.KomiteEmail` roster kosong sejak migrasi claimnonprop 611).
func (a *Acuan) EmailAnggotaWorkbasket(ctx context.Context, workbasket string) ([]string, error) {
	var nama [3]string
	for i, t := range []string{"M_LOGIN_GO_WORKBASKET", "M_WORKBASKET", "M_LOGIN_GO"} {
		q, err := a.q(t)
		if err != nil {
			return nil, err
		}
		nama[i] = q
	}
	q := sqlEmailAnggotaWorkbasket(nama[0], nama[1], nama[2])
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := a.g.db.QueryContext(ctx, q, workbasket)
	if err != nil {
		return nil, fmt.Errorf("repository: email anggota workbasket: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []string
	for rows.Next() {
		var v sql.NullString
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("repository: memindai email anggota workbasket: %w", err)
		}
		if s := strings.TrimSpace(v.String); s != "" {
			out = append(out, s)
		}
	}
	return out, rows.Err()
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
