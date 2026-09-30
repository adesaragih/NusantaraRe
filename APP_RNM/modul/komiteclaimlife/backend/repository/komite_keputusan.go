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
	"database/sql"
	"errors"
	"fmt"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/komiteclaimlife/backend/models"
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
func (r *InboxKomite) CatatKeputusan(ctx context.Context, tx *db.Tx, kasusID string,
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
	if err := db.PeriksaSQL(q1); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q1, keputusan, db.KosongJadiNil(komentar), saat,
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
	if err := db.PeriksaSQL(q2); err != nil {
		return err
	}
	h, err = tx.ExecContext(ctx, q2, keputusan, urut+1, kasusID, urut)
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

// sqlTanggaSebelumDitimpa - keadaan tangga SEBELUM 5.1 menimpanya, dibaca di
// dalam transaksi keputusan: keputusan tingkat akhir yang baru saja ditulis
// `CatatKeputusan` (langkah 3) ikut terbaca, beserta komentarnya.
func sqlTanggaSebelumDitimpa(list string) string {
	return fmt.Sprintf(`SELECT KOMITE_URUT, KOMITE_OPERATORID, KOMITE_APPROVAL, KOMITE_COMMENT
	  FROM %s WHERE DATA_KOMITE_ID = :1 ORDER BY KOMITE_URUT`, list)
}

// TanggaSebelumDitimpa membaca tangga satu kasus di dalam `tx` (OQ-K-05).
func (r *InboxKomite) TanggaSebelumDitimpa(ctx context.Context, tx *db.Tx, kasusID string) ([]AnggotaKasus, error) {
	if tx == nil {
		return nil, errors.New("repository: membaca tangga sebelum penimpaan menuntut transaksi")
	}
	list, err := r.db.Qualify("T_KOMITE_KOMITELIST")
	if err != nil {
		return nil, err
	}
	q := sqlTanggaSebelumDitimpa(list)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, q, kasusID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca tangga komite: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []AnggotaKasus
	for rows.Next() {
		var a AnggotaKasus
		var op, ap, kom sql.NullString
		if err := rows.Scan(&a.Urut, &op, &ap, &kom); err != nil {
			return nil, fmt.Errorf("repository: memindai tangga komite: %w", err)
		}
		a.OperatorID, a.Approval, a.Komentar = op.String, ap.String, kom.String
		out = append(out, a)
	}
	return out, rows.Err()
}

// sqlTimpaTanggaTolakAkhir - `KomitePostAdjustment` langkah 5.1 "Set Reject
// Komite berjenjang" (b5784; `pyStepsBlockName` kosong b5795, hidup):
// perulangan SELURUH `KomiteList` baris adjustment itu (repeat b6127) -
// `KomiteAproval = 2` b5899, `KomiteComment = ""` b5945, `DateApprove =
// @CurrentDateTime()` b5965. OQ-K-05 DITUTUP 29-09-2026 (GILIRAN-17).
//
// ⛔ `KOMITE_APPROVAL IS NOT NULL`: tingkat yang DILEWATI eskalasi (NULL,
// tiket 03) tidak pernah memberi keputusan - menimpanya menjadi Tolak adalah
// keputusan karangan. Pega tidak mengenal eskalasi, jadi di Pega seluruh
// tingkat memang berkeputusan saat 5.1 berjalan.
func sqlTimpaTanggaTolakAkhir(list string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET KOMITE_APPROVAL = :1, KOMITE_COMMENT = NULL, DATE_APPROVE = :2
	 WHERE DATA_KOMITE_ID = :3 AND KOMITE_APPROVAL IS NOT NULL`, list)
}

// TimpaTanggaTolakAkhir menjalankan langkah 5.1 di dalam `tx`. Nol baris =
// galat: tingkat akhir yang menolak pasti sudah menulis keputusannya.
func (r *InboxKomite) TimpaTanggaTolakAkhir(ctx context.Context, tx *db.Tx, kasusID string, saat time.Time) error {
	if tx == nil {
		return errors.New("repository: penimpaan tangga komite menuntut transaksi")
	}
	list, err := r.db.Qualify("T_KOMITE_KOMITELIST")
	if err != nil {
		return err
	}
	q := sqlTimpaTanggaTolakAkhir(list)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, models.KeputusanKomiteTolak, saat, kasusID)
	if err != nil {
		return fmt.Errorf("repository: menimpa tangga komite: %w", err)
	}
	if n, err := h.RowsAffected(); err != nil {
		return fmt.Errorf("repository: membaca cacah penimpaan tangga: %w", err)
	} else if n < 1 {
		return fmt.Errorf("%w: kasus %q tanpa tingkat berkeputusan", ErrKeputusanKomiteBersamaan, kasusID)
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
func (r *InboxKomite) Eskalasi(ctx context.Context, tx *db.Tx, kasusID string, urut int) error {
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
	if err := db.PeriksaSQL(q1); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q1, kasusID, urut, ApprovalKomiteMenunggu)
	if err != nil {
		return fmt.Errorf("repository: melewati anak tangga: %w", err)
	}
	if n, err := h.RowsAffected(); err != nil {
		return fmt.Errorf("repository: membaca cacah anak tangga: %w", err)
	} else if n != 1 {
		return fmt.Errorf("%w: anak tangga %d kasus %q", ErrKeputusanKomiteBersamaan, urut, kasusID)
	}
	q2 := sqlNaikkanTingkat(gen)
	if err := db.PeriksaSQL(q2); err != nil {
		return err
	}
	h, err = tx.ExecContext(ctx, q2, urut+1, kasusID, urut)
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

// sqlNomorAksepDiAdjustment mencacah baris adjustment yang memakai nomor itu.
//
// ⛔ Tiket 04a: nomor akseptasi Komite dan Claim Life berbagi BENTUK tetapi
// tidak berbagi PENGHITUNG. `NomorAkseptasiDipakai` (Claim Life) memeriksa
// tabel datar warisan; jalur kedua yang kini menulis ke `T_CLAIMLF_ADJUSTMENT`
// juga harus diperiksa di sana.
func sqlNomorAksepDiAdjustment(adj string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ACCEPTED_NO = :1`, adj)
}

// NomorAkseptasiDipakaiDiAdjustment menjawab apakah nomor sudah dipakai baris lain.
func (r *InboxKomite) NomorAkseptasiDipakaiDiAdjustment(ctx context.Context, tx *db.Tx,
	nomor string) (bool, error) {

	if tx == nil {
		return false, errors.New("repository: pemeriksaan nomor akseptasi menuntut transaksi")
	}
	adj, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return false, err
	}
	q := sqlNomorAksepDiAdjustment(adj)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := tx.QueryRowContext(ctx, q, nomor).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: memeriksa nomor akseptasi: %w", err)
	}
	return n > 0, nil
}

// sqlRekamAkhirWarisan memperbarui baris datar warisan baris yang diputuskan.
//
// ⚠️ PENYIMPANGAN SADAR dari `UpdateOsAkseptasiClaimLife_sql` (Komite Claim
// Life `RDBList/`, rule yang SAMA di Claim Life): rule itu meng-`INSERT` baris
// BARU (55 kolom, `COMMIT;` b117) di tingkat akhir. Di sistem ini baris datar
// untuk adjustment itu SUDAH ada sejak pendaftaran (`PohonKlaim.Simpan`,
// `ID = adjustment ID`) - `INSERT` kedua menggandakannya. Yang berubah di
// tingkat akhir karena itu DIPERBARUI: status, nomor, tanggal.
// namaTabelLamaKM adalah tabel datar warisan yang ikut ditulis keputusan akhir.
//
// Refactor bentuk B (30-09-2026): dulu dipinjam dari konstanta Claim Life
// (`pohonklaim.go`). Tabelnya satu, penulisnya dua; tiap modul menyebutnya
// sendiri.
const namaTabelLama = "OS_AKSEPTASI_KLAIM_LIFE"

func sqlRekamAkhirWarisan(datar string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET STS_REJECT = :1, NO_ACCEPTATION = :2, ACCEPTATION_DATE = :3
	 WHERE ID = :4`, datar)
}

// ErrStatusAkhirKomiteTidakSah - status rekam akhir di luar {1, 2}.
var ErrStatusAkhirKomiteTidakSah = errors.New(
	"repository: status rekam akhir komite hanya 1 (aksep) atau 2 (tolak)")

// RekamAkhirWarisan - SATU jalur berparameter status (tiket 04b, AC 16 spec):
// aksep (`1`, dengan nomor) dan tolak (`2`, tanpa nomor) lewat fungsi ini.
//
// ⛔ Status dipagari SEBELUM bind: `STS_REJECT` di tabel warisan `NUMBER(38)`
// (`[data DBA]`), jadi teks selain `1`/`2` akan menjadi ORA-01722 - pola
// `PeriksaNilaiWarisan` (butir s1).
func (r *InboxKomite) RekamAkhirWarisan(ctx context.Context, tx *db.Tx, adjID, status,
	nomor string, saat time.Time) error {

	if tx == nil {
		return errors.New("repository: rekam akhir komite menuntut transaksi")
	}
	if status != "1" && status != "2" {
		return fmt.Errorf("%w: %q", ErrStatusAkhirKomiteTidakSah, status)
	}
	datar, err := r.db.Qualify(namaTabelLama)
	if err != nil {
		return err
	}
	q := sqlRekamAkhirWarisan(datar)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, status, db.KosongJadiNil(nomor), saat, adjID)
	if err != nil {
		return fmt.Errorf("repository: merekam akhir komite ke baris datar: %w", err)
	}
	n, err := h.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: membaca cacah baris datar: %w", err)
	}
	// ⛔ Nol baris = baris datar adjustment itu tidak ada; hilir (Arasapas)
	// membaca tabel ini, jadi keputusan yang tidak tercermin di sana GAGAL.
	if n != 1 {
		return fmt.Errorf("repository: rekam akhir komite menyentuh %d baris datar untuk %q, mau 1",
			n, adjID)
	}
	return nil
}
