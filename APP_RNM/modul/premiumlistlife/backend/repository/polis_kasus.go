package repository

// Kasus polis baru - GILIRAN-13 paket 1 (pl3, bn).
//
// Untuk apa berkas ini: melahirkan satu work object polis - baris
// `T_WORK_POLIS` beserta baris `T_PREMIUM_LIST` kosong yang layar berikutnya
// butuhkan - di dalam transaksi MILIK PEMANGGIL.
//
// ⛔ PENGENAL `NBLF-<n>`, tanpa nol di depan - dari DATA, bukan dikarang.
// Korpus tidak memuat awalan kelas `ASM-FW-GISFW-Work-LIFE` (nol
// `pyWorkIDPrefix` untuk kelas itu).
//
// `[terverifikasi — sampel baca-saja DEV 29-09-2026]` 400/400 nilai berbentuk
// `ASM-FW-GISFW-WORK NBLF-<1..5 digit>`, nol berawalan nol. Perintah audit:
// `SELECT IDPEGA FROM POOLDATA.<tabel> WHERE ROWNUM <= 200` untuk
// `JSON_OFFER_LIFE` dan `M_LIFE_PREMIUM_SUMMARY`; nilainya DIUBAH menjadi
// bentuk di mesin (angka -> `<n>`) sebelum dicetak, nol nilai disimpan.
//
// ⚠️ SAMPEL, bukan agregat seperti yang pl3 minta: kueri `GROUP BY` atas
// seluruh tabel dihentikan sesudah 300 detik (terlalu berat untuk DEV), dan
// sampel 200 baris per tabel menggantikannya. Bentuk di luar sampel itu
// `[belum terverifikasi]`.
//
// Pengenal disimpan berbentuk `pyID` - tanpa awalan kelas - sama dengan
// `T_WORK_CLAIM` (`CLM-…`).
//
// Dibaca sesudah: polis_work.go, pengenalwork.go.

import (
	"context"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
)

// AwalanWorkPolis adalah awalan pengenal work object polis Life.
const AwalanWorkPolis = "NBLF-"

// RakitPengenalWorkPolis menyusun pengenal dari angka urutannya.
//
// ⛔ TANPA padding - berbeda dengan `RakitPengenalWork` klaim (enam digit,
// butir aa). Pengenal warisan polis tidak berawalan nol, dan pengenal baru
// yang dipad akan terbaca sebagai jenis nomor yang lain.
func RakitPengenalWorkPolis(urut string) string {
	return AwalanWorkPolis + strings.TrimSpace(urut)
}

// sqlSisipKasusPolis - baris kerja polis baru, kelima kolomnya.
func sqlSisipKasusPolis(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, LINI, POSITION, STATUS_WORK, FLAG_ONGOING_POLICY)
		  VALUES (:1, :2, :3, :4, :5)`, tabel)
}

func argSisipKasusPolis(id string, k models.KasusPolisBaru) []any {
	return []any{id, k.Lini, k.Posisi, k.Status, k.Flag}
}

// sqlSisipPremiumListKosong - header polis kosong, berbagi pengenal.
//
// ⛔ `ID` = `ID_PEGA` = pengenal work. `ID` karena shared PK (050/051);
// `ID_PEGA` karena kotak masuk menggabung `p.ID_PEGA = w.ID`. Penampung
// unik, pengenal dikirim dua kali.
//
// `CREATE_OP_NAME` - padanan `pxCreateOpName` (kolom `Create Operator Name`
// kotak masuk, InboxPremiumList b750). ⛔ Yang ditulis PENGENAL AKUN pembuat,
// bukan nama orangnya: `inti.Pelaku` sengaja tidak membawa nama (ADR-U-0030).
func sqlSisipPremiumListKosong(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, ID_PEGA, TGL_INPUT, CREATE_OP_NAME) VALUES (:1, :2, :3, :4)`, tabel)
}

// PengenalBerikut menerbitkan satu pengenal work polis dari SEQ_WORK_POLIS.
func (r *WorkPolis) PengenalBerikut(ctx context.Context, tx *db.Tx) (string, error) {
	// Pembaca sequence yang SAMA dengan pengenal klaim - satu tempat untuk
	// `NEXTVAL`, bukan salinan kedua.
	urut, err := r.db.NomorBerikut(ctx, tx, "SEQ_WORK_POLIS")
	if err != nil {
		return "", err
	}
	return RakitPengenalWorkPolis(urut), nil
}

// SisipKasusBaru menulis baris kerja dan header polis kosong.
//
// ⛔ Keduanya di transaksi PEMANGGIL: kasus tanpa header tidak tampil lengkap
// di kotak masuk, dan header tanpa kasus adalah baris yatim.
func (r *WorkPolis) SisipKasusBaru(ctx context.Context, tx *db.Tx, id string,
	k models.KasusPolisBaru, pembuat string, saat time.Time) error {

	kerja, err := r.db.Qualify("T_WORK_POLIS")
	if err != nil {
		return err
	}
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return err
	}
	for _, l := range []struct {
		nama string
		q    string
		args []any
	}{
		{"kasus polis", sqlSisipKasusPolis(kerja), argSisipKasusPolis(id, k)},
		{"header polis", sqlSisipPremiumListKosong(polis), []any{id, id, saat, db.KosongJadiNil(pembuat)}},
	} {
		if err := db.PeriksaSQL(l.q); err != nil {
			return err
		}
		hasil, err := tx.ExecContext(ctx, l.q, l.args...)
		if err != nil {
			return fmt.Errorf("repository: menyisipkan %s: %w", l.nama, err)
		}
		if err := db.PastikanSatuBaris(hasil, l.nama); err != nil {
			return err
		}
	}
	return nil
}
