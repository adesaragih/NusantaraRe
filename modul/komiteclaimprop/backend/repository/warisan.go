package repository

// Untuk apa berkas ini: TULISAN KE TABEL WARISAN di titik yang sama dengan XML, ditulis ulang tanpa procedure dan tanpa
// COMMIT (keputusan work owner 08-10-2026; isi procedure dibaca dari ALL_SOURCE DEV). DATA_JSON hanya diisi di
// OS_AKSEPTASI_KLAIM (keputusan work owner 08-10-2026 "isi json nya khusus os_akseptasi_klaim", meralat prompt §3);
// JSON_KLAIM tetap tanpa DATA_JSON:
//
//	SaveOSClaim_SQL -> PEGA_JSON_OS_AKSEP_KLAIM   OS_AKSEPTASI_KLAIM (kolom datar + DATA_JSON halaman) S17
//	InsertClaimPNC  -> PEGA_JSON_KLAIM_PNC        JSON_KLAIM (INSERT bila IDPEGA belum ada)            S28
//	InsertLogServiceClaim                         MONITORING_KLAIM_LOG                                  S30-S31
//	InsertHistoryAkseptasiPega_Sql                HISTORYAKSEPTASIPEGA                                  S32-S33
//	efek keluar                                   T_LOG_SERVICE_RNM (inti/backend/outbox)               S29/S34/S35

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/komiteclaimprop/backend/models"
)

// hariJakarta - `TRUNC(SYSDATE)`: tanggal hari aksi tanpa jam.
func hariJakarta(t time.Time) time.Time {
	l := t.In(models.Jakarta)
	return time.Date(l.Year(), l.Month(), l.Day(), 0, 0, 0, 0, time.UTC)
}

// potong - teks dipangkas ke panjang kolom (rune).
func potong(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

// pecahAngka - angka eksak sebagai (koefisien, skala): `TO_NUMBER(:k) / POWER(10, :s)`, nol float.
func pecahAngka(d *apd.Decimal) (string, int64) {
	s := d.Text('f')
	tanda := ""
	if strings.HasPrefix(s, "-") {
		tanda, s = "-", s[1:]
	}
	bulat, pecahan, _ := strings.Cut(s, ".")
	koef := strings.TrimLeft(bulat+pecahan, "0")
	if koef == "" {
		return "0", 0
	}
	return tanda + koef, int64(len(pecahan))
}

// sqlSisipOS - INSERT satu baris OS_AKSEPTASI_KLAIM beserta argumennya; kolom angka / tanggal kosong tidak disebut
// (NULL, sama dengan procedure yang tidak mengisinya).
func sqlSisipOS(tabel string, b models.BarisOSAkseptasi, saat time.Time) (string, []any, error) {
	var args []any
	ph := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf(":%d", len(args))
	}
	kol := []string{"CASEID", "NOCLAIM", "TANGGAL", "NOPOLIS", "STS_REJECT", "MASTERID", "CAUSEOFLOSS", "CAUSEOFLOSSID",
		"CURRENCYID", "CURRENCY", "TYPE", "ACCEPTEDNO", "PAYMENTTYPE"}
	nilai := []string{ph(b.CaseID), ph(db.KosongJadiNil(b.NoClaim)), ph(hariJakarta(saat)), ph(db.KosongJadiNil(b.NoPolis)),
		ph(db.KosongJadiNil(b.StsReject)), ph(db.KosongJadiNil(b.MasterID)), ph(db.KosongJadiNil(potong(b.CauseOfLoss, 1000))),
		ph(db.KosongJadiNil(b.CauseOfLossID)), ph(db.KosongJadiNil(b.CurrencyID)), ph(db.KosongJadiNil(b.Currency)),
		ph(db.KosongJadiNil(b.Type)), ph(db.KosongJadiNil(b.AcceptedNo)), ph(db.KosongJadiNil(b.PaymentType))}
	for _, x := range []struct{ kol, v string }{{"GROSSVALUE", b.GrossValue}, {"PERSENRNM", b.PersenRNM},
		{"VALUE", b.Value}, {"KURSVALUE", b.KursValue}} {
		if strings.TrimSpace(x.v) == "" {
			continue
		}
		d, err := utils.ParseDecimal(strings.TrimSpace(x.v))
		if err != nil {
			return "", nil, fmt.Errorf("repository: %s bukan angka (%q): %w", x.kol, x.v, err)
		}
		k, s := pecahAngka(d)
		kol = append(kol, x.kol)
		nilai = append(nilai, "(TO_NUMBER("+ph(k)+") / POWER(10, "+ph(s)+"))")
	}
	if !b.EstimationDate.IsZero() {
		kol = append(kol, "ESTIMATIONDATE")
		nilai = append(nilai, ph(b.EstimationDate.In(models.Jakarta)))
	}
	// DATA_JSON = halaman TempOSAkseptasi (procedure: `INSERT ... DATA_JSON ... VALUES (... DataPega ...)`); CHECK
	// `DATA_JSON IS JSON`; teks biasa.
	if b.DataJSON != "" {
		kol = append(kol, "DATA_JSON")
		nilai = append(nilai, ph(b.DataJSON))
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tabel, strings.Join(kol, ", "), strings.Join(nilai, ", ")),
		args, nil
}

// SisipOS = SaveAcceptation_Act S5 (`SaveOSClaim_SQL`).
func (g *Gudang) SisipOS(ctx context.Context, tx *db.Tx, b models.BarisOSAkseptasi, saat time.Time) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
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
	h, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("repository: menyisipkan OS_AKSEPTASI_KLAIM: %w", err)
	}
	return db.PastikanSatuBaris(h, "baris OS_AKSEPTASI_KLAIM")
}

// sqlAdaJSONKlaim / sqlSisipJSONKlaim - isi PEGA_JSON_KLAIM_PNC tanpa DATA_JSON.
func sqlAdaJSONKlaim(tabel string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE IDPEGA = :1`, tabel)
}

func sqlSisipJSONKlaim(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (MNK_NO_KLAIM, IDPEGA, TGL_INPUT, NOPOLIS, IDPROD) VALUES (:1, :2, :3, :4, 1)`,
		tabel)
}

// SalinJSONKlaim = InsertJsonClaimTreaty_act S4 (`InsertClaimPNC`): baris IDPEGA yang sudah ada dibiarkan (UPDATE
// procedure hanya menyetel DATA_JSON, yang tidak diisi), selainnya INSERT (MNK_NO_KLAIM, IDPEGA, TGL_INPUT, NOPOLIS,
// IDPROD 1).
func (g *Gudang) SalinJSONKlaim(ctx context.Context, tx *db.Tx, idPega, noKlaim, noPolis string, saat time.Time) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
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
	h, err := tx.ExecContext(ctx, q, db.KosongJadiNil(noKlaim), idPega, saat, db.KosongJadiNil(noPolis))
	if err != nil {
		return fmt.Errorf("repository: menyisipkan JSON_KLAIM: %w", err)
	}
	return db.PastikanSatuBaris(h, "baris JSON_KLAIM")
}

// sqlLogLayanan - satu baris MONITORING_KLAIM_LOG.
func sqlLogLayanan(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (TGL_INPUT, IDPEGA, PARAMETER, JN_SERVICE, NO_AKSEPTASI, NO_DLA, STS_MESSAGE,
		RESPON_MESSAGE) VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, tabel)
}

// CatatLogLayanan = InsertLogServiceClaim (S31).
func (g *Gudang) CatatLogLayanan(ctx context.Context, tx *db.Tx, l models.LogLayanan, saat time.Time) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
	tabel, err := g.db.Qualify("MONITORING_KLAIM_LOG")
	if err != nil {
		return err
	}
	q := sqlLogLayanan(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, saat, db.KosongJadiNil(l.IDPega), db.KosongJadiNil(potong(l.Parameter, 1000)),
		db.KosongJadiNil(l.JenisService), db.KosongJadiNil(l.NoAkseptasi), db.KosongJadiNil(l.NoDLA),
		db.KosongJadiNil(potong(l.StsMessage, 500)), db.KosongJadiNil(potong(l.ResponMessage, 1000)))
	if err != nil {
		return fmt.Errorf("repository: mencatat log layanan klaim: %w", err)
	}
	return db.PastikanSatuBaris(h, "log layanan klaim")
}

// sqlRiwayatAkseptasi - InsertHistoryAkseptasiPega_Sql (`Tgl_Transfer = sysdate`; OPERATORID tidak disebut).
func sqlRiwayatAkseptasi(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID_PEGA, TGL_TRANSFER, STATUS, USERNAME, WORKBASKET, ID_KOMITE)
		VALUES (:1, :2, :3, :4, :5, :6)`, tabel)
}

// CatatRiwayatAkseptasi = KomitePostAdjustment S32-S33.
func (g *Gudang) CatatRiwayatAkseptasi(ctx context.Context, tx *db.Tx, r models.RiwayatAkseptasi, saat time.Time) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
	tabel, err := g.db.Qualify("HISTORYAKSEPTASIPEGA")
	if err != nil {
		return err
	}
	q := sqlRiwayatAkseptasi(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, db.KosongJadiNil(r.IDPega), saat, db.KosongJadiNil(potong(r.Status, 150)),
		db.KosongJadiNil(potong(r.Username, 150)), db.KosongJadiNil(r.Workbasket), db.KosongJadiNil(r.IDKomite))
	if err != nil {
		return fmt.Errorf("repository: mencatat HISTORYAKSEPTASIPEGA: %w", err)
	}
	return db.PastikanSatuBaris(h, "baris HISTORYAKSEPTASIPEGA")
}

// Lini dan modul outbox Komite Claim Prop.
const (
	LiniOutbox  = models.LiniProp
	ModulOutbox = "komiteclaimprop"
)

// AntreEfek menulis satu efek keluar ke outbox di transaksi Submit.
func (g *Gudang) AntreEfek(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) (string, error) {
	return outbox.NewPenyimpan(g.db).AntreEfek(ctx, tx, LiniOutbox, ModulOutbox, jenis, rujukan, muatan, saat)
}
