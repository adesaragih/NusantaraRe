package repository

// Untuk apa berkas ini: TULISAN KE TABEL WARISAN di titik yang sama dengan XML, ditulis ulang tanpa procedure dan tanpa
// COMMIT (isi procedure dibaca dari ALL_SOURCE DEV 10-10-2026; bentuk sama dengan Claim Fac In tahap 1):
//
//	SaveOSClaim_SQL -> PEGA_JSON_OS_AKSEP_KLAIM   OS_AKSEPTASI_KLAIM (CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS,
//	                                              STS_REJECT, STS_KONVERSI, STS_DLA; TGL_PROD trigger)
//	                                              KomitePost_Adjustment S8, KomitePost_Reject S12, KomitePost_CloseClaim S12
//	InsertClaimPNC  -> PEGA_JSON_KLAIM_PNC        JSON_KLAIM (INSERT bila IDPEGA belum ada; DATA_JSON tidak diisi)
//	InsertLogServiceClaim                         MONITORING_KLAIM_LOG (KomitePost_Adjustment S18-S19 "AKSEPATSI")
//	InsertHistoryAkseptasiPega_Sql                HISTORYAKSEPTASIPEGA (S20-S21)
//	UpdateSubProgresKlaim                         SUBPROGRESSCLAIM.POSITION2 (S22-S23)
//	InsertClaimRejected_Sql                       CLAIMREJECTED (KomitePost_Reject S14)
//	efek keluar                                   T_LOG_SERVICE_RNM (inti/backend/outbox)

import (
	"context"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/komiteclaimfacin/backend/models"
)

// hariJakarta - `TO_DATE(to_char(sysdate,'dd/MM/yyyy'),'dd/MM/yyyy')`: tanggal hari aksi tanpa jam.
func hariJakarta(t time.Time) time.Time {
	l := t.In(models.Jakarta)
	return time.Date(l.Year(), l.Month(), l.Day(), 0, 0, 0, 0, time.UTC)
}

// detik - `To_date(..., 'DD/MM/YYYY HH24:MI:SS')`: waktu Jakarta tanpa pecahan detik.
func detik(t time.Time) time.Time {
	l := t.In(models.Jakarta)
	return time.Date(l.Year(), l.Month(), l.Day(), l.Hour(), l.Minute(), l.Second(), 0, time.UTC)
}

// potong - teks dipangkas ke panjang kolom (rune).
func potong(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

// sqlSisipOS - isi PEGA_JSON_OS_AKSEP_KLAIM baris 16 (tanpa COMMIT). MASTERID / CLAIMOLD = NULL (CARI18 / CARI20 tidak
// diisi pemanggil mana pun).
func sqlSisipOS(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, STS_KONVERSI, STS_DLA)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, tabel)
}

// SisipOS menyisipkan satu baris OS_AKSEPTASI_KLAIM.
func (g *Gudang) SisipOS(ctx context.Context, tx *db.Tx, b models.BarisOS, saat time.Time) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
	tabel, err := g.db.Qualify("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return err
	}
	q := sqlSisipOS(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, b.CaseID, db.KosongJadiNil(b.NoClaim), db.KosongJadiNil(b.DataJSON),
		hariJakarta(saat), db.KosongJadiNil(b.NoPolis), db.KosongJadiNil(b.StsReject), db.KosongJadiNil(b.StsKonversi),
		db.KosongJadiNil(b.StsDLA))
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

// SalinJSONKlaim = InsertJsonClaimNonMBU_act S4-S6 (`InsertClaimPNC`): baris IDPEGA yang sudah ada dibiarkan (UPDATE
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

// CatatLogLayanan = InsertLogServiceClaim (KomitePost_Adjustment S18-S19).
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

// CatatRiwayatAkseptasi = KomitePost_Adjustment S20-S21.
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

// sqlUbahSubProgres = UpdateSubProgresKlaim (`UPDATE POOLDATA.SUBPROGRESSCLAIM SET POSITION2 = {CARI11} WHERE IDPEGA =
// {pyWorkPage.pyID}`): baris sub-progres kasus komite "Waiting Committee" (Claim Fac In CreateKMTNo_Act 16).
func sqlUbahSubProgres(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET POSITION2 = :1 WHERE IDPEGA = :2`, tabel)
}

// UbahSubProgres = KomitePost_Adjustment S22-S23 (tingkat akhir atau tolak). Nol baris = sub-progres tidak pernah ditulis
// (bukan galat - XML pun tidak memeriksa).
func (g *Gudang) UbahSubProgres(ctx context.Context, tx *db.Tx, komiteID, posisi string) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
	tabel, err := g.db.Qualify("SUBPROGRESSCLAIM")
	if err != nil {
		return err
	}
	q := sqlUbahSubProgres(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, posisi, models.KunciInstans(komiteID)); err != nil {
		return fmt.Errorf("repository: memperbarui SUBPROGRESSCLAIM: %w", err)
	}
	return nil
}

// sqlSisipKlaimDitolak = InsertClaimRejected_Sql (tanpa COMMIT).
func sqlSisipKlaimDitolak(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (INSKEY, ID, INSNAME, LABEL, STATUSWORK, CREATEOPNAME, CREATEOPERATOR, OBJCLASS,
		UPDATEDATETIME, UPDATEOPNAME, UPDATEOPERATOR, REMARK) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12)`,
		tabel)
}

// SisipKlaimDitolak = KomitePost_Reject S14.
func (g *Gudang) SisipKlaimDitolak(ctx context.Context, tx *db.Tx, k models.KlaimDitolak) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
	tabel, err := g.db.Qualify("CLAIMREJECTED")
	if err != nil {
		return err
	}
	q := sqlSisipKlaimDitolak(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, k.InsKey, db.KosongJadiNil(k.ID), db.KosongJadiNil(k.InsName),
		db.KosongJadiNil(k.Label), db.KosongJadiNil(k.StatusWork), db.KosongJadiNil(potong(k.PembuatNama, 512)),
		db.KosongJadiNil(k.PembuatID), db.KosongJadiNil(k.Kelas), detik(k.Diperbarui),
		db.KosongJadiNil(potong(k.PengubahNama, 512)), db.KosongJadiNil(k.PengubahID),
		db.KosongJadiNil(potong(k.Remark, 1000)))
	if err != nil {
		return fmt.Errorf("repository: menyisipkan CLAIMREJECTED: %w", err)
	}
	return db.PastikanSatuBaris(h, "baris CLAIMREJECTED")
}

// Lini dan modul outbox Komite Claim Fac In.
const (
	LiniOutbox  = models.LiniFacIn
	ModulOutbox = "komiteclaimfacin"
)

// AntreEfek menulis satu efek keluar ke outbox di transaksi Submit.
func (g *Gudang) AntreEfek(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) (string, error) {
	return outbox.NewPenyimpan(g.db).AntreEfek(ctx, tx, LiniOutbox, ModulOutbox, jenis, rujukan, muatan, saat)
}
