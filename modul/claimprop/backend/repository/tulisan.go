package repository

// Untuk apa berkas ini: TULISAN KE TABEL WARISAN di titik yang sama dengan XML, ditulis ulang tanpa procedure:
//
//	PEGA_JSON_OS_AKSEP_KLAIM / ..._KLAIMTNP  -> OS_AKSEPTASI_KLAIM (kolom datar + DATA_JSON = DataPega apa adanya -
//	                                            keputusan work owner 08-10-2026, meralat 07-10-2026)
//	PEGA_JSON_KLAIM_PNC (InsertClaimPNC)     -> JSON_KLAIM (UPDATE ... lalu INSERT bila tak ada; DATA_JSON kosong)
//	InsertLogServiceClaim                    -> MONITORING_KLAIM_LOG
//	SaveCatasrtope_Act (Obj-Save)            -> CATASTROPHE
//	efek keluar                              -> T_LOG_SERVICE_RNM (inti/backend/outbox)

import (
	"context"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/claimprop/backend/models"
)

// hariJakarta - `TO_DATE(to_char(sysdate,'dd/MM/yyyy'),'dd/MM/yyyy')`: tanggal hari aksi tanpa jam.
func hariJakarta(t time.Time) time.Time {
	l := t.In(models.Jakarta)
	return time.Date(l.Year(), l.Month(), l.Day(), 0, 0, 0, 0, time.UTC)
}

// SisipOS menyisipkan satu baris OS_AKSEPTASI_KLAIM.
func (g *Gudang) SisipOS(ctx context.Context, tx *db.Tx, b models.BarisOS, saat time.Time) error {
	tabel, err := g.db.Qualify("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return err
	}
	q, args, err := sqlSisipOS(tabel, b, saat)
	if err != nil {
		return err
	}
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("repository: menyisipkan OS_AKSEPTASI_KLAIM: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "baris OS_AKSEPTASI_KLAIM")
}

// sqlSisipOS - INSERT satu baris OS_AKSEPTASI_KLAIM beserta argumennya; kolom angka / tanggal kosong tidak disebut.
func sqlSisipOS(tabel string, b models.BarisOS, saat time.Time) (string, []any, error) {
	var args []any
	ph := func(v ...any) string {
		s := ""
		for i, x := range v {
			args = append(args, x)
			if i > 0 {
				s += ", "
			}
			s += fmt.Sprintf(":%d", len(args))
		}
		return s
	}
	angka := func(j, s string) (string, error) {
		v, err := angkaAtauNil(j, s)
		if err != nil {
			return "", err
		}
		return "(TO_NUMBER(" + ph(v[0]) + ") / POWER(10, " + ph(v[1]) + "))", nil
	}
	kol := []string{"CASEID", "NOCLAIM", "TANGGAL", "NOPOLIS", "STS_REJECT", "MASTERID", "CURRENCYID", "CURRENCY",
		"TYPE", "TYPELOSSID", "TYPELOSSNAME", "CAUSEOFLOSSID", "CAUSEOFLOSS", "ACCEPTEDNO", "INSERT_OP"}
	nilai := []string{ph(b.CaseID), ph(teksAtauNil(b.NoClaim)), ph(hariJakarta(saat)), ph(teksAtauNil(b.NoPolis)),
		ph(teksAtauNil(b.StsReject)), ph(teksAtauNil(b.MasterID)), ph(teksAtauNil(b.CurrencyID)),
		ph(teksAtauNil(b.Currency)), ph(teksAtauNil(b.StsReject)), ph(teksAtauNil(b.TypeLossID)),
		ph(teksAtauNil(b.TypeLoss)), ph(teksAtauNil(b.CauseOfLossID)), ph(teksAtauNil(b.CauseOfLoss)),
		ph(teksAtauNil(b.AcceptedNo)), ph(teksAtauNil(b.InsertOp))}
	for _, x := range []struct{ kol, v string }{{"VALUE", b.Value}, {"GROSSVALUE", b.GrossValue},
		{"KURSVALUE", b.KursValue}, {"PERSENRNM", b.PersenRNM}} {
		if x.v == "" {
			continue
		}
		e, err := angka(x.kol, x.v)
		if err != nil {
			return "", nil, err
		}
		kol = append(kol, x.kol)
		nilai = append(nilai, e)
	}
	if !b.EstimationDate.IsZero() {
		kol = append(kol, "ESTIMATIONDATE")
		nilai = append(nilai, ph(b.EstimationDate.In(models.Jakarta)))
	}
	// DATA_JSON = DataPega (procedure: `INSERT ... DATA_JSON ... VALUES (... DataPega ...)`); CHECK DATA_JSON IS JSON.
	if b.DataJSON != "" {
		kol = append(kol, "DATA_JSON")
		nilai = append(nilai, ph(b.DataJSON))
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tabel, joinKoma(kol), joinKoma(nilai)), args, nil
}

func joinKoma(s []string) string {
	out := ""
	for i, x := range s {
		if i > 0 {
			out += ", "
		}
		out += x
	}
	return out
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
// IDPROD 1). Dua kolom konversi procedure diisi NULL; di sini tidak disebut - katalog DEV 07-10-2026: nullable tanpa
// DATA_DEFAULT, hasilnya sama.
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

// LogLayanan - satu baris MONITORING_KLAIM_LOG (InsertLogServiceClaim).
type LogLayanan struct {
	IDPega, Parameter, JenisService, NoAkseptasi, NoDLA, StsMessage, ResponMessage string
}

// sqlLogLayanan - satu baris MONITORING_KLAIM_LOG.
func sqlLogLayanan(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (TGL_INPUT, IDPEGA, PARAMETER, JN_SERVICE, NO_AKSEPTASI, NO_DLA, STS_MESSAGE,
		RESPON_MESSAGE) VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, tabel)
}

// CatatLogLayanan = InsertLogServiceClaim.
func (g *Gudang) CatatLogLayanan(ctx context.Context, tx *db.Tx, l LogLayanan, saat time.Time) error {
	tabel, err := g.db.Qualify("MONITORING_KLAIM_LOG")
	if err != nil {
		return err
	}
	q := sqlLogLayanan(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, saat, teksAtauNil(l.IDPega), teksAtauNil(potong(l.Parameter, 1000)),
		teksAtauNil(l.JenisService), teksAtauNil(l.NoAkseptasi), teksAtauNil(l.NoDLA),
		teksAtauNil(potong(l.StsMessage, 500)), teksAtauNil(potong(l.ResponMessage, 1000)))
	if err != nil {
		return fmt.Errorf("repository: mencatat log layanan klaim: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "log layanan klaim")
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

// Lini dan modul outbox Claim Prop.
const (
	LiniOutbox  = models.LiniProp
	ModulOutbox = "claimprop"
)

// AntreEfek menulis satu efek keluar ke outbox di transaksi aksi.
func (g *Gudang) AntreEfek(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) (string, error) {
	return outbox.NewPenyimpan(g.db).AntreEfek(ctx, tx, LiniOutbox, ModulOutbox, jenis, rujukan, muatan, saat)
}
