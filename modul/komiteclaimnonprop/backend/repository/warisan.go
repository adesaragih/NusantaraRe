package repository

// Untuk apa berkas ini: TULISAN KE TABEL WARISAN di titik yang sama dengan XML, ditulis ulang tanpa procedure dan tanpa
// COMMIT (pola Komite Claim Prop; isi procedure dibaca dari ALL_SOURCE DEV 09-10-2026). DATA_JSON diisi di
// OS_AKSEPTASI_KLAIM dan OS_AKSEPTASI_SUBJECTIVITY (halaman TempOSAkseptasi / TempOSSubjectivity); JSON_KLAIM tetap
// tanpa DATA_JSON (keputusan work owner 08-10-2026, pola Claim Prop):
//
//	SaveOSClaim_SQL         -> PEGA_JSON_OS_AKSEP_KLAIM         OS_AKSEPTASI_KLAIM         InsertOSKlaimCNP S6
//	SaveXOLClaim_SQL        -> XOL2_AKSEP_KLAIM                 CLAIMXOL2                  InsertXOLKlaimCNP S5.2
//	SaveOSSubjectivity_SQL  -> PEGA_JSON_OS_AKSEP_SUBJECTIVITY  OS_AKSEPTASI_SUBJECTIVITY  InsertOSSubjectivityCNP S4
//	InsertClaimPNC          -> PEGA_JSON_KLAIM_PNC              JSON_KLAIM                 KomitePostAdjustment S24
//	InsertHistoryAkseptasiPega_Sql                              HISTORYAKSEPTASIPEGA       KomitePostAdjustment S22-S23
//	efek keluar                                                 T_LOG_SERVICE_RNM (inti/backend/outbox) S14.22 / S19

import (
	"context"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/komiteclaimnonprop/backend/models"
)

// hariJakarta - `TO_DATE(TO_CHAR(SYSDATE, 'dd/MM/yyyy'))`: tanggal hari aksi tanpa jam.
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

// sqlSisipOS - PEGA_JSON_OS_AKSEP_KLAIM baris 16 tanpa STS_KONVERSI / STS_DLA / CLAIMOLD (tidak diisi InsertOSKlaimCNP:
// NULL). Procedure menghitung CASEID lebih dulu tetapi tetap INSERT - ditiru (tanpa hitungan).
func sqlSisipOS(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, MASTERID)
		VALUES (:1, :2, :3, :4, :5, :6, :7)`, tabel)
}

// SisipOS = InsertOSKlaimCNP S6 (`SaveOSClaim_SQL`).
func (g *Gudang) SisipOS(ctx context.Context, tx *db.Tx, b models.BarisOSAkseptasi, saat time.Time) error {
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
	h, err := tx.ExecContext(ctx, q, b.CaseID, db.KosongJadiNil(b.NoClaim), db.KosongJadiNil(b.DataJSON), hariJakarta(saat),
		db.KosongJadiNil(b.NoPolis), db.KosongJadiNil(b.StsReject), db.KosongJadiNil(b.MasterID))
	if err != nil {
		return fmt.Errorf("repository: menyisipkan OS_AKSEPTASI_KLAIM: %w", err)
	}
	return db.PastikanSatuBaris(h, "baris OS_AKSEPTASI_KLAIM")
}

// sqlSisipXOL2 - XOL2_AKSEP_KLAIM (kolom ber-kutip VERBATIM, semuanya VARCHAR2 kecuali TANGGAL).
func sqlSisipXOL2(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (CASEID, "GrossAdjustment", "CNPReinstatement", "Currency", "KursIDR", XOL, TANGGAL)
		VALUES (:1, :2, :3, :4, :5, :6, :7)`, tabel)
}

// SisipXOL2 = InsertXOLKlaimCNP S5.1-S5.2 (`SaveXOLClaim_SQL`), satu baris per layer. TANGGAL = RDB
// `To_date(sysdate, 'DD/MM/RRRR')` - tanggal hari aksi (CARI55 diisi activity tetapi tidak diikat RDB).
func (g *Gudang) SisipXOL2(ctx context.Context, tx *db.Tx, b models.BarisXOL2, saat time.Time) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
	tabel, err := g.db.Qualify("CLAIMXOL2")
	if err != nil {
		return err
	}
	q := sqlSisipXOL2(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, b.CaseID, db.KosongJadiNil(b.GrossAdjustment), db.KosongJadiNil(b.CNPReinstatement),
		db.KosongJadiNil(b.Currency), db.KosongJadiNil(b.KursIDR), db.KosongJadiNil(b.XOL), hariJakarta(saat))
	if err != nil {
		return fmt.Errorf("repository: menyisipkan CLAIMXOL2: %w", err)
	}
	return db.PastikanSatuBaris(h, "baris CLAIMXOL2")
}

// sqlAdaOSSubjectivity / sqlUbahOSSubjectivity / sqlSisipOSSubjectivity - PEGA_JSON_OS_AKSEP_SUBJECTIVITY baris 7, 19,
// 31 (STS_KONVERSI / STS_DLA tidak diisi InsertOSSubjectivityCNP: NULL).
func sqlAdaOSSubjectivity(tabel string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE CASEID = :1`, tabel)
}

func sqlUbahOSSubjectivity(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET STS_SUBJECTIVITY = :1, DATA_JSON = :2, TANGGAL = :3 WHERE CASEID = :4`, tabel)
}

func sqlSisipOSSubjectivity(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, STS_SUBJECTIVITY)
		VALUES (:1, :2, :3, :4, :5, :6, :7)`, tabel)
}

// SimpanOSSubjectivity = InsertOSSubjectivityCNP S4 (`SaveOSSubjectivity_SQL`): CASEID yang sudah ada diubah
// (STS_SUBJECTIVITY, DATA_JSON, TANGGAL), selainnya disisipkan.
func (g *Gudang) SimpanOSSubjectivity(ctx context.Context, tx *db.Tx, b models.BarisOSSubjectivity, saat time.Time) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
	tabel, err := g.db.Qualify("OS_AKSEPTASI_SUBJECTIVITY")
	if err != nil {
		return err
	}
	qAda, qUbah, qSisip := sqlAdaOSSubjectivity(tabel), sqlUbahOSSubjectivity(tabel), sqlSisipOSSubjectivity(tabel)
	for _, x := range []string{qAda, qUbah, qSisip} {
		if err := db.PeriksaSQL(x); err != nil {
			return err
		}
	}
	var n int
	if err := tx.QueryRowContext(ctx, qAda, b.CaseID).Scan(&n); err != nil {
		return fmt.Errorf("repository: memeriksa OS_AKSEPTASI_SUBJECTIVITY: %w", err)
	}
	hari, dj := hariJakarta(saat), db.KosongJadiNil(b.DataJSON)
	if n > 0 {
		if _, err := tx.ExecContext(ctx, qUbah, db.KosongJadiNil(b.StsSubjectivity), dj, hari, b.CaseID); err != nil {
			return fmt.Errorf("repository: mengubah OS_AKSEPTASI_SUBJECTIVITY: %w", err)
		}
		return nil
	}
	h, err := tx.ExecContext(ctx, qSisip, b.CaseID, db.KosongJadiNil(b.NoClaim), dj, hari, db.KosongJadiNil(b.NoPolis),
		db.KosongJadiNil(b.StsReject), db.KosongJadiNil(b.StsSubjectivity))
	if err != nil {
		return fmt.Errorf("repository: menyisipkan OS_AKSEPTASI_SUBJECTIVITY: %w", err)
	}
	return db.PastikanSatuBaris(h, "baris OS_AKSEPTASI_SUBJECTIVITY")
}

// sqlAdaJSONKlaim / sqlSisipJSONKlaim - isi PEGA_JSON_KLAIM_PNC tanpa DATA_JSON.
func sqlAdaJSONKlaim(tabel string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE IDPEGA = :1`, tabel)
}

func sqlSisipJSONKlaim(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (MNK_NO_KLAIM, IDPEGA, TGL_INPUT, NOPOLIS, IDPROD) VALUES (:1, :2, :3, :4, 1)`,
		tabel)
}

// SalinJSONKlaim = InsertJsonClaimTreatyNonProp_act S3 (`InsertClaimPNC`): baris IDPEGA yang sudah ada dibiarkan (UPDATE
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

// sqlRiwayatAkseptasi - InsertHistoryAkseptasiPega_Sql (`Tgl_Transfer = sysdate`; OPERATORID tidak disebut).
func sqlRiwayatAkseptasi(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID_PEGA, TGL_TRANSFER, STATUS, USERNAME, WORKBASKET, ID_KOMITE)
		VALUES (:1, :2, :3, :4, :5, :6)`, tabel)
}

// CatatRiwayatAkseptasi = KomitePostAdjustment S22-S23.
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

// Lini dan modul outbox Komite Claim Non Prop.
const (
	LiniOutbox  = models.LiniNonProp
	ModulOutbox = "komiteclaimnonprop"
)

// AntreEfek menulis satu efek keluar ke outbox di transaksi Submit.
func (g *Gudang) AntreEfek(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) (string, error) {
	return outbox.NewPenyimpan(g.db).AntreEfek(ctx, tx, LiniOutbox, ModulOutbox, jenis, rujukan, muatan, saat)
}
