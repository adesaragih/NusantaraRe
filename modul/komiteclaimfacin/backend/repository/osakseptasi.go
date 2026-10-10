package repository

// Untuk apa berkas ini: OS_AKSEPTASI_KLAIM - SaveOSClaim_SQL -> PEGA_JSON_OS_AKSEP_KLAIM ditulis ulang tanpa procedure
// dan tanpa COMMIT (isi procedure dibaca dari ALL_SOURCE DEV 10-10-2026): CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS,
// STS_REJECT, status konversi, STS_DLA (TGL_PROD = trigger; MASTERID / CLAIMOLD tidak diisi pemanggil mana pun).
// KomitePost_Adjustment S8, KomitePost_Reject S12 (SaveReject_ACT_KMT), KomitePost_CloseClaim S12.
//
// ⛔ Satu-satunya berkas modul ini yang menyebut kolom status konversi (penjaga Claim Life
// `repository/migrasibatas_test.go` `berkasBolehMenyebutKolomTakDibawa`, izin work owner 10-10-2026 OQ-KCFI-08):
// kolom itu hanya ikut INSERT bila terisi (SaveReject_ACT_KMT S17, Value kosong -> "1"); kosong = tidak disebut
// (preseden Claim Prop).

import (
	"context"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/komiteclaimfacin/backend/models"
)

// sqlSisipOS - INSERT OS_AKSEPTASI_KLAIM; `konversi` = kolom status konversi ikut disebut (nilai terisi).
func sqlSisipOS(tabel string, konversi bool) string {
	if konversi {
		return fmt.Sprintf(`INSERT INTO %s (CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, STS_KONVERSI, STS_DLA)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, tabel)
	}
	return fmt.Sprintf(`INSERT INTO %s (CASEID, NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, STS_DLA)
		VALUES (:1, :2, :3, :4, :5, :6, :7)`, tabel)
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
	konversi := b.StsKonversi != ""
	q := sqlSisipOS(tabel, konversi)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	args := []any{b.CaseID, db.KosongJadiNil(b.NoClaim), db.KosongJadiNil(b.DataJSON), hariJakarta(saat),
		db.KosongJadiNil(b.NoPolis), db.KosongJadiNil(b.StsReject)}
	if konversi {
		args = append(args, b.StsKonversi)
	}
	h, err := tx.ExecContext(ctx, q, append(args, db.KosongJadiNil(b.StsDLA))...)
	if err != nil {
		return fmt.Errorf("repository: menyisipkan OS_AKSEPTASI_KLAIM: %w", err)
	}
	return db.PastikanSatuBaris(h, "baris OS_AKSEPTASI_KLAIM")
}
