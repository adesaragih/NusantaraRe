package repository

// Tanggal tutup buku - `POOLDATA.TANGGAL_CLOSING`, tiket 02 PremiumList.
//
// ⛔ TABEL WARISAN, DIBACA SAJA. Nol tulisan ke sana di berkas ini, dan tidak
// boleh ada: aturan periode bukan milik modul ini, dan `[data DBA]`
// `PROC_GENERATE_SEQUENCE_NUMBER` juga menggulir periode lewat tabel yang
// sama - dua penulis atas satu aturan berarti dua aturan.
//
// ⛔ DIBACA SETIAP KALI DIBUTUHKAN, nol cache. Tanggal tutup buku berubah,
// dan cache yang tidak pernah kedaluwarsa membukukan transaksi ke periode
// kemarin tanpa satu pun galat (AC 7 tiket 02).
//
// Sumbernya `RDBList/GETTanggalClosing_SQL.xml` b85:
//
//	SELECT * FROM POOLDATA.TANGGAL_CLOSING
//
// ⚠️ `SELECT *` ditiru sebagai `SELECT TANGGAL` - kolom yang DIPAKAI
// `SubmitPremiumList_Act` b918 (`TglProd.pxResults(1).TANGGAL`). (Sensus remark
// 28-09-2026: di Pega nilai b918 hanya mengalir ke langkah 15 yang ter-remark;
// pembaca hidup `TANGGAL_CLOSING` adalah `PROC_GENERATE_SEQUENCE_NUMBER`.
// Pembaca ini ada karena `[keputusan work owner]` "ikuti yang dari DB".) Membawa
// seluruh kolom berarti perubahan tabel di hulu mengubah bentuk baris kami
// tanpa ada yang memintanya.
//
// Dibaca sesudah: polis_inbox.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/db"
)

// TutupBuku membaca tanggal tutup buku yang berlaku.
type TutupBuku struct{ db *db.DB }

// NewTutupBuku menyusun pembacanya.
func NewTutupBuku(db *db.DB) *TutupBuku { return &TutupBuku{db: db} }

// sqlTanggalTutupBuku merakit query-nya.
//
// ⚠️ `FETCH FIRST 1 ROWS ONLY` ditambahkan; rule aslinya membaca
// `pxResults(1)` saja, yaitu baris pertama apa pun cacahnya. Menuliskannya
// membuat niat itu terbaca - dan menahan tabel yang ternyata berisi lebih
// dari satu baris tanpa ada yang menyadarinya.
func sqlTanggalTutupBuku(tabel string) string {
	return fmt.Sprintf(
		`SELECT TANGGAL FROM %s FETCH FIRST 1 ROWS ONLY`, tabel)
}

// Tanggal mengembalikan tanggal tutup buku, atau 0 bila tabelnya kosong.
//
// ⛔ Nol BUKAN nilai pengganti - ia penanda "tidak ada", dan
// `models.PeriodeProduksi` menolaknya. Yang mengembalikan 25 di sini akan
// membuat seluruh disiplin tiket 02 sia-sia dalam satu baris.
func (r *TutupBuku) Tanggal(ctx context.Context) (int, error) {
	tabel, err := r.db.Qualify("TANGGAL_CLOSING")
	if err != nil {
		return 0, err
	}
	q := sqlTanggalTutupBuku(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var tgl sql.NullInt64
	err = r.db.QueryRowContext(ctx, q).Scan(&tgl)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("repository: membaca POOLDATA.TANGGAL_CLOSING: %w", err)
	}
	if !tgl.Valid {
		return 0, nil
	}
	return int(tgl.Int64), nil
}
