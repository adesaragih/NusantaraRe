package repository

// Penghitung nomor klaim - butir o1, A2.
//
// Untuk apa berkas ini: logika `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`
// DITULIS ULANG DI GO.
//
// ⛔ `[keputusan work owner]` *"jangan ada lagi pemanggilan procedure, segala
// procedure hardcode dalam skrip"* — prinsip **o**. Procedure-nya tidak
// dipanggil; yang tetap di Oracle hanya `SELECT … FOR UPDATE`, sebab kunci
// baris memang milik basis data dan menirunya di aplikasi berarti menulis
// ulang penguncian.
//
// `[data DBA — belum dikonfirmasi DBA]` Sumbernya
// `.scratch/claim-life/SUMBER-PENOMORAN-DBA.md`, yang melabeli dirinya sendiri
// *"dibaca sendiri dari katalog instance pengembangan, belum dikonfirmasi
// DBA"*. Label di sini mengikuti label sumbernya — tidak dinaikkan.
//
// ⚠️ TULISAN KE TABEL WARISAN YANG DISENGAJA. `GENERATE_SEQUENCE_NUMBER`
// di-`UPDATE`/`INSERT` di sini, dan itu inheren pada butir o1: penghitung yang
// tidak disimpan bukan penghitung. Dicatat, bukan disembunyikan.
//
// Dibaca sesudah: pengenalwork.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrTanggalClosingKosong - `TANGGAL_CLOSING` tidak memberi hari tutup buku.
	ErrTanggalClosingKosong = errors.New(
		"repository: TANGGAL_CLOSING tidak memberi hari tutup buku")
	// ErrTanggalClosingTakTerurai - isinya bukan bilangan.
	ErrTanggalClosingTakTerurai = errors.New(
		"repository: TANGGAL_CLOSING bukan bilangan")
)

// batasCutover adalah tanggal khusus di procedure.
//
// `[data DBA]` `IF TRUNC(v_now) <= TO_DATE('02/01/2026','DD/MM/YYYY') THEN
// v_mm_yyyy := '12.2025'; v_tahun := 2025;`
//
// ⚠️ Cabang ini sudah LEWAT hari ini dan tidak akan menyala lagi. Ia
// dipertahankan apa adanya sebab migrasi data (tiket 13) menguraikan nomor
// LAMA, dan nomor yang lahir sebelum tanggal itu memakai periode `12.2025`
// meski dibuat Januari 2026. Menghapusnya membuat pengurai salah baca.
var batasCutover = time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)

const (
	periodeCutover = "12.2025"
	tahunCutover   = "2025"
)

// PeriodeNomor adalah keluaran penentuan periode.
type PeriodeNomor struct {
	// MMYYYY berbentuk `MM.YYYY`, dipakai di badan nomor.
	MMYYYY string
	// Tahun adalah kunci ketiga `(CLASS, JENIS, TAHUN)`. TEKS: kolomnya
	// `VARCHAR2(5)` (ADR-U-0022).
	Tahun string
}

// HitungPeriodeNomor menentukan periode dan tahun kunci.
//
// `[data DBA]` `SUMBER-PENOMORAN-DBA.md` baris 48-66: bila hari `saat`
// melewati hari tutup buku, periode digeser satu bulan `ADD_MONTHS(+1)`.
//
// ⛔ `ADD_MONTHS` menggeser BULAN beserta tahunnya - Desember menjadi Januari
// tahun berikutnya. Ini BERBEDA dari `SaveAdjustment_Act`, yang tahunnya tidak
// ikut bergeser karena kedua cabang `@if`-nya identik (lihat akseptasi.go).
// Dua penomoran, dua perilaku, dan keduanya ditiru apa adanya masing-masing.
func HitungPeriodeNomor(saat time.Time, hariClosing int) PeriodeNomor {
	hari := time.Date(saat.Year(), saat.Month(), saat.Day(), 0, 0, 0, 0, time.UTC)
	if !hari.After(batasCutover) {
		return PeriodeNomor{MMYYYY: periodeCutover, Tahun: tahunCutover}
	}
	periode := saat
	if saat.Day() > hariClosing {
		periode = saat.AddDate(0, 1, 0)
	}
	return PeriodeNomor{
		MMYYYY: fmt.Sprintf("%02d.%04d", int(periode.Month()), periode.Year()),
		Tahun:  strconv.Itoa(periode.Year()),
	}
}

// HariClosing membaca hari tutup buku.
//
// `[data DBA]` `SELECT TO_NUMBER(tanggal) FROM POOLDATA.TANGGAL_CLOSING
// WHERE ROWNUM = 1` - kolomnya `VARCHAR2(10)`, dibaca sebagai bilangan.
//
// ⛔ Kosong atau tak terurai GAGAL TERANG. Menebak hari tutup buku menggeser
// periode seluruh nomor yang terbit hari itu.
func (r *PohonKlaim) HariClosing(ctx context.Context, tx *Tx) (int, error) {
	tabel, err := r.db.Qualify("TANGGAL_CLOSING")
	if err != nil {
		return 0, err
	}
	q := fmt.Sprintf(`SELECT TANGGAL FROM %s WHERE ROWNUM = 1`, tabel)
	if err := PeriksaSQL(q); err != nil {
		return 0, err
	}
	var teks sql.NullString
	if err := tx.tx.QueryRowContext(ctx, q).Scan(&teks); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrTanggalClosingKosong
		}
		return 0, fmt.Errorf("repository: membaca hari tutup buku: %w", err)
	}
	if !teks.Valid || strings.TrimSpace(teks.String) == "" {
		return 0, ErrTanggalClosingKosong
	}
	n, err := strconv.Atoi(strings.TrimSpace(teks.String))
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrTanggalClosingTakTerurai, teks.String)
	}
	return n, nil
}

// UrutNomorBerikut mengunci baris penghitung, menaikkannya, dan menyimpannya.
//
// `[data DBA]` `SUMBER-PENOMORAN-DBA.md` baris 70-99, diringkas tanpa
// menyalin nama tabelnya telanjang (ADR-U-0033 dijaga penjaga statik, dan ia
// membaca komentar juga):
//
//	kunci baris (CLASS, JENIS, TAHUN) dengan FOR UPDATE
//	  ada      -> urut := urut + 1, lalu perbarui barisnya
//	  tak ada  -> urut := 1, lalu sisipkan baris pertama
//
// ⛔ `FOR UPDATE` dipertahankan: dua pemanggil serentak diserialkan oleh kunci
// baris itu, dan tanpa kunci keduanya membaca urut yang sama lalu menerbitkan
// nomor kembar.
//
// ⛔ Nol `COMMIT`. `[data DBA]` procedure-nya pun tidak punya - yang terlihat
// di rule Pega berada di blok pemanggil (OQ-013). Transaksinya milik
// pendaftaran.
func (r *PohonKlaim) UrutNomorBerikut(ctx context.Context, tx *Tx,
	class, jenis string, p PeriodeNomor, saat time.Time) (int, error) {

	tabel, err := r.db.Qualify("GENERATE_SEQUENCE_NUMBER")
	if err != nil {
		return 0, err
	}
	qKunci := fmt.Sprintf(
		`SELECT NO_SEQ FROM %s WHERE CLASS = :1 AND JENIS = :2 AND TAHUN = :3 FOR UPDATE`,
		tabel)
	if err := PeriksaSQL(qKunci); err != nil {
		return 0, err
	}
	var urut int
	err = tx.tx.QueryRowContext(ctx, qKunci, class, jenis, p.Tahun).Scan(&urut)

	switch {
	case err == sql.ErrNoRows:
		// Belum ada baris kunci ini: urut pertama, lalu INSERT.
		qSisip := fmt.Sprintf(`INSERT INTO %s
			(CLASS, JENIS, TAHUN, NO_SEQ, TANGGAL, MM_YYYY)
			VALUES (:1,:2,:3,:4,:5,:6)`, tabel)
		if err := PeriksaSQL(qSisip); err != nil {
			return 0, err
		}
		hasil, err := tx.tx.ExecContext(ctx, qSisip,
			class, jenis, p.Tahun, 1, saat, p.MMYYYY)
		if err != nil {
			return 0, fmt.Errorf("repository: menyisipkan penghitung nomor: %w", err)
		}
		if err := pastikanSatuBaris(hasil, "penyisipan penghitung nomor"); err != nil {
			return 0, err
		}
		return 1, nil

	case err != nil:
		return 0, fmt.Errorf("repository: mengunci penghitung nomor: %w", err)
	}

	urut++
	qUbah := fmt.Sprintf(
		`UPDATE %s SET NO_SEQ = :1, MM_YYYY = :2, TANGGAL = :3
		  WHERE CLASS = :4 AND JENIS = :5 AND TAHUN = :6`, tabel)
	if err := PeriksaSQL(qUbah); err != nil {
		return 0, err
	}
	hasil, err := tx.tx.ExecContext(ctx, qUbah,
		urut, p.MMYYYY, saat, class, jenis, p.Tahun)
	if err != nil {
		return 0, fmt.Errorf("repository: menaikkan penghitung nomor: %w", err)
	}
	if err := pastikanSatuBaris(hasil, "penaikan penghitung nomor"); err != nil {
		return 0, err
	}
	return urut, nil
}
