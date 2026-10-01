package repository

// Satu baris T_PREMIUM_LIST per status penawaran - migrasi 062.
//
// `[keputusan work owner 01-10-2026]` "Save Offer status Accept -> 1 baris
// Accept; save lagi status Pending -> 1 baris baru Pending; save lagi dengan
// status yang sama -> baris status itu diperbarui, bukan bertambah."
//
// Bentuknya, untuk satu kasus:
//
//	baris utama   ID = ID_PEGA = nomor kasus   STATUS_PENAWARAN = status TERAKHIR
//	baris status  ID = PengenalBarisStatus(kasus, status), ID_PEGA = nomor kasus
//
// ⛔ Baris utama TETAP berkunci nomor kasus dan selalu memuat isian TERAKHIR,
// sebab seluruh pembaca lain (kotak masuk, Premium List Detail, penomoran PL,
// summary, Claim Life, dan FK `PREMIUM_LIST_ID` tabel anak) mencarinya dengan
// `ID = nomor kasus`. Baris status hanya SALINAN isian penawaran; kolom milik
// tahap lain (TYPE, NO_POLIS, PRODUCT_NAME, ...) tidak disalin, jadi baris
// status tidak pernah ditemukan pencarian polis mana pun.
//
// Urutan satu simpanan, di transaksi pemanggil (`PindahkanBarisStatus`, lalu
// `Penawaran.Simpan` menimpa baris utama):
//
//	status lama kosong ATAU sama      -> tidak ada yang dipindah
//	status lama A, status baru B      -> salin baris utama ke baris status A
//	                                     (ganti bila sudah ada), lalu buang
//	                                     baris status B (isinya pindah ke baris
//	                                     utama) - tepat satu baris per status
//
// ⛔ Nol `COMMIT` (ADR-U-0029); `Qualify` di setiap query (ADR-U-0033).

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// kolomSalinanBarisStatus - kolom yang disalin dari baris utama ke baris status.
//
// ⛔ Diturunkan dari `kolomBacaPenawaran` (isian layar), ditambah identitas
// kasus - satu daftar, bukan salinan kedua yang akan berselisih.
func kolomSalinanBarisStatus() []string {
	kolom := []string{"ID_PEGA", "TGL_INPUT", "CREATE_OP_NAME"}
	for _, k := range kolomBacaPenawaran {
		kolom = append(kolom, k.nama)
	}
	return kolom
}

// PengenalBarisStatus - `ID` baris status: 32 heksa, deterministik.
//
// ⚠️ Pola `PengenalPesertaUnggah`. Tidak pernah sama dengan nomor kasus
// (`NBLF-<n>` bukan 32 heksa), jadi baris status tidak menimpa baris utama.
func PengenalBarisStatus(kasusID, status string) string {
	sum := md5.Sum([]byte(kasusID + "\x00status-penawaran\x00" + status))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

func sqlKunciStatusPenawaran(polis string) string {
	return fmt.Sprintf(`SELECT STATUS_PENAWARAN FROM %s WHERE ID = :1 FOR UPDATE`, polis)
}

func sqlBuangBarisStatus(polis string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND ID_PEGA = :2`, polis)
}

func sqlSalinKeBarisStatus(polis string) string {
	kolom := strings.Join(kolomSalinanBarisStatus(), ", ")
	return fmt.Sprintf(`INSERT INTO %s (ID, %s) SELECT :1, %s FROM %s WHERE ID = :2`,
		polis, kolom, kolom, polis)
}

// StatusTerkunci membaca status terakhir baris utama DAN menguncinya sampai
// commit - dua simpanan serentak atas kasus yang sama berurutan di sini.
func (r *Penawaran) StatusTerkunci(ctx context.Context, tx *db.Tx, kasusID string) (string, error) {
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return "", err
	}
	q := sqlKunciStatusPenawaran(polis)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var status sql.NullString
	err = tx.QueryRowContext(ctx, q, kasusID).Scan(&status)
	if err == sql.ErrNoRows {
		return "", ErrHeaderPolisTidakAda
	}
	if err != nil {
		return "", fmt.Errorf("repository: mengunci header polis: %w", err)
	}
	return strings.TrimSpace(status.String), nil
}

// PindahkanBarisStatus menyiapkan baris status sebelum baris utama ditimpa.
//
// Mengembalikan true bila isian lama dipindah ke baris status tersendiri.
func (r *Penawaran) PindahkanBarisStatus(ctx context.Context, tx *db.Tx,
	kasusID, statusLama, statusBaru string) (bool, error) {

	if statusLama == "" || statusLama == statusBaru {
		return false, nil
	}
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return false, err
	}
	buang, salin := sqlBuangBarisStatus(polis), sqlSalinKeBarisStatus(polis)
	for _, q := range []string{buang, salin} {
		if err := db.PeriksaSQL(q); err != nil {
			return false, err
		}
	}
	idLama := PengenalBarisStatus(kasusID, statusLama)
	// 1. Baris status lama diganti salinan baris utama (isian status lama).
	if _, err := tx.ExecContext(ctx, buang, idLama, kasusID); err != nil {
		return false, fmt.Errorf("repository: membuang baris status %q: %w", statusLama, err)
	}
	hasil, err := tx.ExecContext(ctx, salin, idLama, kasusID)
	if err != nil {
		return false, fmt.Errorf("repository: menyalin baris status %q: %w", statusLama, err)
	}
	if err := db.PastikanSatuBaris(hasil, "salinan baris status penawaran"); err != nil {
		return false, err
	}
	// 2. Baris status baru (bila pernah ada) dibuang: isinya kini baris utama.
	if _, err := tx.ExecContext(ctx, buang, PengenalBarisStatus(kasusID, statusBaru), kasusID); err != nil {
		return false, fmt.Errorf("repository: membuang baris status %q: %w", statusBaru, err)
	}
	return true, nil
}
