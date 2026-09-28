package repository

// Penulisan keputusan satu tingkat Komite - tiket 02 Komite Claim Life.
//
// `[terverifikasi]` `Komite Claim Life/Activity/KomitePostAdjustment.xml`:
// langkah 3 menulis `KomiteList(Local.Komite)` - `KomiteAproval =
// .AcceptStatus`, `KomiteComment = .Comment`, `DateApprove =
// @CurrentDateTime()`; langkah 13 `KomiteCount = KomiteCount + 1`.
//
// ⛔ DUA baris, SATU transaksi milik pemanggil, dan KEDUANYA bersyarat:
//
//	anak tangga  WHERE DATA_KOMITE_ID, KOMITE_URUT, KOMITE_OPERATORID, APPROVAL = 0
//	kepala kasus WHERE ID, KOMITE_COUNT = count yang dibaca
//
// Keputusan dibaca di luar transaksi; dua anggota (atau dua klik) yang
// memutuskan bersamaan tidak boleh sama-sama menang. Yang kalah menerima
// `ErrKeputusanKomiteBersamaan` dan transaksinya batal utuh.
//
// ⛔ Baris dipilih lewat `KOMITE_URUT` + `DATA_KOMITE_ID` (AC 34 spec) -
// bukan indeks posisi `KomiteList(n)` Pega.

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrKeputusanKomiteBersamaan - keadaan tangga berubah sejak dibaca.
var ErrKeputusanKomiteBersamaan = errors.New(
	"repository: tangga komite berubah sejak dibaca; muat ulang kasusnya")

// sqlCatatAnakTangga menulis keputusan satu anggota.
func sqlCatatAnakTangga(list string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET KOMITE_APPROVAL = :1, KOMITE_COMMENT = :2, DATE_APPROVE = :3
	 WHERE DATA_KOMITE_ID = :4 AND KOMITE_URUT = :5
	   AND KOMITE_OPERATORID = :6 AND KOMITE_APPROVAL = :7`, list)
}

// sqlMajukanTangga menaikkan tingkat dan menyimpan `AcceptStatus`.
//
// ⚠️ `ACCEPT_STATUS` = `pyWorkPage.AcceptStatus` terakhir - bahan
// `IsKomiteLoop` (lihat `models.KasusDiTangga`).
func sqlMajukanTangga(gen string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET ACCEPT_STATUS = :1, KOMITE_COUNT = :2
	 WHERE ID = :3 AND KOMITE_COUNT = :4`, gen)
}

// CatatKeputusan menulis keputusan tingkat `urut` lalu menaikkan tingkat.
func (r *InboxKomite) CatatKeputusan(ctx context.Context, tx *Tx, kasusID string,
	urut int, operatorID, keputusan, komentar string, saat time.Time) error {

	if tx == nil {
		return errors.New("repository: keputusan komite menuntut transaksi")
	}
	gen, err := r.db.Qualify("T_GENERAL_KOMITE")
	if err != nil {
		return err
	}
	list, err := r.db.Qualify("T_KOMITE_KOMITELIST")
	if err != nil {
		return err
	}
	q1 := sqlCatatAnakTangga(list)
	if err := PeriksaSQL(q1); err != nil {
		return err
	}
	h, err := tx.tx.ExecContext(ctx, q1, keputusan, kosongJadiNil(komentar), saat,
		kasusID, urut, operatorID, ApprovalKomiteMenunggu)
	if err != nil {
		return fmt.Errorf("repository: menulis keputusan anak tangga: %w", err)
	}
	if n, err := h.RowsAffected(); err != nil {
		return fmt.Errorf("repository: membaca cacah anak tangga: %w", err)
	} else if n != 1 {
		return fmt.Errorf("%w: anak tangga %d kasus %q", ErrKeputusanKomiteBersamaan, urut, kasusID)
	}

	q2 := sqlMajukanTangga(gen)
	if err := PeriksaSQL(q2); err != nil {
		return err
	}
	h, err = tx.tx.ExecContext(ctx, q2, keputusan, urut+1, kasusID, urut)
	if err != nil {
		return fmt.Errorf("repository: memajukan tangga komite: %w", err)
	}
	if n, err := h.RowsAffected(); err != nil {
		return fmt.Errorf("repository: membaca cacah kepala komite: %w", err)
	} else if n != 1 {
		return fmt.Errorf("%w: kepala kasus %q", ErrKeputusanKomiteBersamaan, kasusID)
	}
	return nil
}

// sqlLewatiAnakTangga mengosongkan anak tangga yang dilewati eskalasi.
//
// ⛔ Approval menjadi NULL - "tidak pernah memberi keputusan", persis kalimat
// keputusan work owner - BUKAN `1`/`2`. Nilai keputusan untuk tingkat yang
// tidak memutuskan adalah keputusan karangan.
func sqlLewatiAnakTangga(list string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET KOMITE_APPROVAL = NULL, DATE_APPROVE = NULL
	 WHERE DATA_KOMITE_ID = :1 AND KOMITE_URUT = :2 AND KOMITE_APPROVAL = :3`, list)
}

// sqlNaikkanTingkat menaikkan `KOMITE_COUNT` TANPA menyentuh `ACCEPT_STATUS`.
func sqlNaikkanTingkat(gen string) string {
	return fmt.Sprintf(`UPDATE %s SET KOMITE_COUNT = :1 WHERE ID = :2 AND KOMITE_COUNT = :3`, gen)
}

// Eskalasi melewati anak tangga `urut` dan menaikkan tingkat ke `urut+1`.
func (r *InboxKomite) Eskalasi(ctx context.Context, tx *Tx, kasusID string, urut int) error {
	if tx == nil {
		return errors.New("repository: eskalasi komite menuntut transaksi")
	}
	gen, err := r.db.Qualify("T_GENERAL_KOMITE")
	if err != nil {
		return err
	}
	list, err := r.db.Qualify("T_KOMITE_KOMITELIST")
	if err != nil {
		return err
	}
	q1 := sqlLewatiAnakTangga(list)
	if err := PeriksaSQL(q1); err != nil {
		return err
	}
	h, err := tx.tx.ExecContext(ctx, q1, kasusID, urut, ApprovalKomiteMenunggu)
	if err != nil {
		return fmt.Errorf("repository: melewati anak tangga: %w", err)
	}
	if n, err := h.RowsAffected(); err != nil {
		return fmt.Errorf("repository: membaca cacah anak tangga: %w", err)
	} else if n != 1 {
		return fmt.Errorf("%w: anak tangga %d kasus %q", ErrKeputusanKomiteBersamaan, urut, kasusID)
	}
	q2 := sqlNaikkanTingkat(gen)
	if err := PeriksaSQL(q2); err != nil {
		return err
	}
	h, err = tx.tx.ExecContext(ctx, q2, urut+1, kasusID, urut)
	if err != nil {
		return fmt.Errorf("repository: menaikkan tingkat komite: %w", err)
	}
	if n, err := h.RowsAffected(); err != nil {
		return fmt.Errorf("repository: membaca cacah kepala komite: %w", err)
	} else if n != 1 {
		return fmt.Errorf("%w: kepala kasus %q", ErrKeputusanKomiteBersamaan, kasusID)
	}
	return nil
}
