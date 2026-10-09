package repository

// Untuk apa berkas ini: TULISAN KE TABEL WARISAN di titik yang sama dengan XML, ditulis ulang tanpa procedure (isi
// procedure dibaca dari ALL_SOURCE DEV 09-10-2026):
//
//	PEGA_JSON_OS_AKSEP_KLAIMTNP (SaveDataToOsAkseptasiNP) -> OS_AKSEPTASI_KLAIM: SELALU INSERT CASEID, NOCLAIM, MASTERID,
//	                                                          DATA_JSON, TANGGAL (hari ini), NOPOLIS, STS_REJECT,
//	                                                          STS_KONVERSI (cabang UPDATE procedure dikomentari)
//	PEGA_JSON_KLAIM_PNC (InsertClaimPNC)                  -> JSON_KLAIM (baris IDPEGA yang ada dibiarkan, selainnya
//	                                                          INSERT; DATA_JSON tidak diisi - keputusan work owner)
//	SaveCatasrtope_Act (Obj-Save)                         -> CATASTROPHE
//	efek keluar                                           -> T_LOG_SERVICE_RNM (inti/backend/outbox)

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/claimnonprop/backend/models"
)

// hariJakarta - `TO_DATE(to_char(sysdate,'dd/MM/yyyy'),'dd/MM/yyyy')`: tanggal hari aksi tanpa jam.
func hariJakarta(t time.Time) time.Time {
	l := t.In(models.Jakarta)
	return time.Date(l.Year(), l.Month(), l.Day(), 0, 0, 0, 0, time.UTC)
}

// sqlSisipOS - isi PEGA_JSON_OS_AKSEP_KLAIMTNP baris 16 (tanpa COMMIT).
func sqlSisipOS(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (CASEID, NOCLAIM, MASTERID, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, STS_KONVERSI)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, tabel)
}

// SisipOS menyisipkan satu baris OS_AKSEPTASI_KLAIM (TGL_PROD diisi trigger TRG_TLG_PROD_OS_AKSEPTASI).
func (g *Gudang) SisipOS(ctx context.Context, tx *db.Tx, b models.BarisOS, saat time.Time) error {
	tabel, err := g.db.Qualify("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return err
	}
	q := sqlSisipOS(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, b.CaseID, teksAtauNil(b.NoClaim), teksAtauNil(b.MasterID),
		teksAtauNil(b.DataJSON), hariJakarta(saat), teksAtauNil(b.NoPolis), teksAtauNil(b.StsReject),
		teksAtauNil(b.StsKonversi))
	if err != nil {
		return fmt.Errorf("repository: menyisipkan OS_AKSEPTASI_KLAIM: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "baris OS_AKSEPTASI_KLAIM")
}

// sqlJumlahOS - GetDataOS / GetDataCNPOS: Σ nilai DATA_JSON baris STS 0 satu kasus per layer x mata uang. Kasus sistem
// baru dan kasus Pega lama dicari keduanya (`kunci`, `kunciLama`).
func sqlJumlahOS(tabel string) string {
	jv := func(p string) string {
		return fmt.Sprintf(`NVL(JSON_VALUE(a.DATA_JSON, '$.%s' RETURNING NUMBER), 0)`, p)
	}
	return fmt.Sprintf(`SELECT %s, %s, %s, %s, %s FROM %s a
		WHERE a.CASEID IN (:1, :2) AND JSON_VALUE(a.DATA_JSON, '$.TypeLoss' RETURNING VARCHAR2(400)) = :3
		  AND JSON_VALUE(a.DATA_JSON, '$.Currency' RETURNING VARCHAR2(100)) = :4 AND a.STS_REJECT = 0`,
		dec("SUM("+jv("Adjusterfee")+")"), dec("SUM("+jv("GrossValue")+")"), dec("SUM("+jv("CNPOthersFee")+")"),
		dec("SUM("+jv("Salvage")+")"), dec("SUM("+jv("Value")+")"), tabel)
}

// JumlahOS = GetDataOS / GetDataCNPOS (SaveDataToOSAksep_Act 15.8.2, SaveToOS 6.3.3, GetSelisihActual_Act 6.2).
func (g *Gudang) JumlahOS(ctx context.Context, tx *db.Tx, kasusID, layer, mataUang string) (models.NilaiOS, error) {
	tabel, err := g.db.Qualify("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return models.NilaiOS{}, err
	}
	q := sqlJumlahOS(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return models.NilaiOS{}, err
	}
	var n [5]sql.NullString
	args := []any{models.KunciInstans(kasusID), models.KunciPegaLama(kasusID), layer, mataUang}
	var row interface{ Scan(...any) error }
	if tx != nil {
		row = tx.QueryRowContext(ctx, q, args...)
	} else {
		row = g.db.QueryRowContext(ctx, q, args...)
	}
	if err := row.Scan(&n[0], &n[1], &n[2], &n[3], &n[4]); err != nil {
		return models.NilaiOS{}, fmt.Errorf("repository: menjumlah OS_AKSEPTASI_KLAIM: %w", err)
	}
	t := func(v sql.NullString) string { return rapikanDesimal(v.String) }
	return models.NilaiOS{Adjusterfee: t(n[0]), GrossValue: t(n[1]), CNPOthersFee: t(n[2]), Salvage: t(n[3]),
		Value: t(n[4])}, nil
}

// sqlNomorKlaimOS - CekOSClaimNonProp_SQL (`noclaim FROM OS_AKSEPTASI_KLAIM WHERE CASEID ORDER BY TANGGAL DESC`).
func sqlNomorKlaimOS(tabel string) string {
	return fmt.Sprintf(`SELECT NOCLAIM FROM %s WHERE CASEID IN (:1, :2) AND NOCLAIM IS NOT NULL ORDER BY TANGGAL DESC
		FETCH FIRST 1 ROWS ONLY`, tabel)
}

// NomorKlaimOS = CekOSClaimNonProp_SQL (SaveDataToOSAksep_Act langkah 10-12): nomor klaim terakhir di OS kasus ini.
func (g *Gudang) NomorKlaimOS(ctx context.Context, tx *db.Tx, kasusID string) (string, error) {
	tabel, err := g.db.Qualify("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return "", err
	}
	q := sqlNomorKlaimOS(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var v sql.NullString
	err = tx.QueryRowContext(ctx, q, models.KunciInstans(kasusID), models.KunciPegaLama(kasusID)).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca nomor klaim OS: %w", err)
	}
	return v.String, nil
}

// sqlAdaJSONKlaim / sqlSisipJSONKlaim - isi PEGA_JSON_KLAIM_PNC tanpa DATA_JSON.
func sqlAdaJSONKlaim(tabel string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE IDPEGA = :1`, tabel)
}

func sqlSisipJSONKlaim(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (MNK_NO_KLAIM, IDPEGA, TGL_INPUT, NOPOLIS, IDPROD) VALUES (:1, :2, :3, :4, 1)`, tabel)
}

// SalinJSONKlaim = PEGA_JSON_KLAIM_PNC (InsertJsonClaimTreaty_act): baris IDPEGA yang sudah ada dibiarkan (UPDATE
// procedure hanya menyetel DATA_JSON, yang tidak diisi), selainnya INSERT (MNK_NO_KLAIM, IDPEGA, TGL_INPUT, NOPOLIS,
// IDPROD 1). TGL_KONVERSI / STS_KONVERSI procedure = NULL, tidak disebut.
func (g *Gudang) SalinJSONKlaim(ctx context.Context, tx *db.Tx, idPega, noKlaim, noPolis string, saat time.Time) error {
	tabel, err := g.db.Qualify("JSON_KLAIM")
	if err != nil {
		return err
	}
	qAda, q := sqlAdaJSONKlaim(tabel), sqlSisipJSONKlaim(tabel)
	for _, x := range []string{qAda, q} {
		if err := db.PeriksaSQL(x); err != nil {
			return err
		}
	}
	var n int
	if err := tx.QueryRowContext(ctx, qAda, idPega).Scan(&n); err != nil {
		return fmt.Errorf("repository: memeriksa JSON_KLAIM: %w", err)
	}
	if n > 0 {
		return nil
	}
	hasil, err := tx.ExecContext(ctx, q, teksAtauNil(noKlaim), idPega, saat, teksAtauNil(noPolis))
	if err != nil {
		return fmt.Errorf("repository: menyisipkan JSON_KLAIM: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "baris JSON_KLAIM")
}

// sqlSisipKatastrofe - satu baris CATASTROPHE.
func sqlSisipKatastrofe(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, NOTE, KLAIMTYPE, TGL_INPUT, USER_INPUT, STSKATASTROFE, NONKATASTROFETYPE)
		VALUES (:1, :2, :3, :4, :5, :6, :7)`, tabel)
}

// SisipKatastrofe = SaveCatasrtope_Act langkah 3 (Obj-Save CATASTROPHE).
func (g *Gudang) SisipKatastrofe(ctx context.Context, tx *db.Tx, k models.KatastrofeBaru, saat time.Time) error {
	tabel, err := g.db.Qualify("CATASTROPHE")
	if err != nil {
		return err
	}
	q := sqlSisipKatastrofe(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, k.ID, teksAtauNil(potong(k.Note, 2000)), k.KlaimType, saat,
		teksAtauNil(k.UserInput), teksAtauNil(k.StsKatastrofe), teksAtauNil(k.NonKatastrofeType))
	if err != nil {
		return fmt.Errorf("repository: menyisipkan katastrofe: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "baris katastrofe")
}

// Lini dan modul outbox Claim Non Prop.
const (
	LiniOutbox  = models.LiniNonProp
	ModulOutbox = "claimnonprop"
)

// AntreEfek menulis satu efek keluar ke outbox di transaksi aksi.
func (g *Gudang) AntreEfek(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) (string, error) {
	return outbox.NewPenyimpan(g.db).AntreEfek(ctx, tx, LiniOutbox, ModulOutbox, jenis, rujukan, muatan, saat)
}
