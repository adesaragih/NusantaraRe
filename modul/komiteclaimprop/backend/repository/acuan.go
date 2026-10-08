package repository

// Untuk apa berkas ini: ACUAN - bacaan baca-saja RDB komite, ditulis ulang dari `RDBList/` korpus Komite Claim Prop:
//
//	TreatyYearTreatyin_SQL  TREATYYEAR (TREATYGROUPID, :ymd BETWEEN STARTDATE AND ENDDATE)  SaveAcceptationTreaty_TKMT S2
//	GetLimitPLATreatyin     PROPORTIONALARRG kolom RP (alias CARI7), TREATYDESCID '10001'      S4.2
//	GetListRetro_Sql        TREATYREINSURER (REINSTYPEID, TREATYYEAR, TREATYGROUPID)           S4.3
//	Obj-Browse BANKACCOUNT  IDOFBANK                                                           HitServiceToKasirKMT_Act S12
//	GetEmailCeding_SQL      GL.F_GET_EMAIL (HANYA saat muatan Kasir disusun di produksi)       HitServiceToKasirKMT_Act S14.1.2
//	OperatorID.pyUserName   M_LOGIN_GO.NAME                                                    S32 InsertHistory.CARI4
//
//	getStatusKonversi_SQL   REINSURANCE.TRLOSS_DETAIL_T (HANYA produksi, gerbang IsPEGAPROD)  HitServiceToKasirKMT_Act S3
//
// ⚠️ Dua bacaan lintas skema tanpa hak di akun DEV (TRLOSS_DETAIL_T ORA-00942, F_GET_EMAIL ORA-00904) - keduanya hanya
// berjalan di produksi (keputusan work owner 08-10-2026); hak DBA diminta.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/komiteclaimprop/backend/models"
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

func dec(kol string) string { return fmt.Sprintf(db.FmtDesimal, kol) }

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

// TahunTreaty = TreatyYearTreatyin_SQL (`TreatyYearOut.pxResults(1).CARI1`; STARTDATE/ENDDATE teks "yyyyMMdd").
func (a *Acuan) TahunTreaty(ctx context.Context, grup, ymd string) (string, error) {
	t, err := a.q("TREATYYEAR")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT TREATYYEAR FROM %s WHERE TREATYGROUPID = :1
		AND :2 BETWEEN STARTDATE AND ENDDATE FETCH FIRST 1 ROWS ONLY`, t), grup, ymd)
	return v, err
}

// LimitPLA = GetLimitPLATreatyin: kolom RP baris pertama (`LimitDla.pxResults(1).CARI7`).
func (a *Acuan) LimitPLA(ctx context.Context, tahun, grup, reins string) (string, error) {
	t, err := a.q("PROPORTIONALARRG")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT %s FROM %s WHERE TREATYDESCID = '10001' AND TREATYYEAR = :1
		AND TREATYGROUPID = :2 AND REINSTYPEID = :3 FETCH FIRST 1 ROWS ONLY`, dec("RP"), t), tahun, grup, reins)
	return v, err
}

// DaftarRetro = GetListRetro_Sql (`ListInsurerDla`): reinsurerid / name / ricomm / pctshare.
func (a *Acuan) DaftarRetro(ctx context.Context, tahun, grup, reins string) ([]models.BarisRetro, error) {
	t, err := a.q("TREATYREINSURER")
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT REINSURERID, NAME, %s, %s FROM %s
		WHERE REINSTYPEID = :1 AND TREATYYEAR = :2 AND TREATYGROUPID = :3 ORDER BY ID`, dec("RICOMM"), dec("PCTSHARE"), t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := a.g.db.QueryContext(ctx, q, reins, tahun, grup)
	if err != nil {
		return nil, fmt.Errorf("repository: bacaan retro komite: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []models.BarisRetro
	for rows.Next() {
		var n [4]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3]); err != nil {
			return nil, fmt.Errorf("repository: memindai retro komite: %w", err)
		}
		out = append(out, models.BarisRetro{ReinsurerID: n[0].String, Name: n[1].String,
			RiComm: rapikanDesimal(strings.TrimSpace(n[2].String)), PctShare: rapikanDesimal(strings.TrimSpace(n[3].String))})
	}
	return out, rows.Err()
}

// IDBankRekening = Obj-Browse BANKACCOUNT (S12-S13, `BankAccount.pxResults(1).IDOFBANK`; saringan Obj-Browse tidak
// diekspor - dibaca seperti Claim Prop: nama bank, cabang, nomor rekening).
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
