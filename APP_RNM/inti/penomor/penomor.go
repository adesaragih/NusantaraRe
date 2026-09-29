// Package penomor membentuk nomor bisnis bersama - penghitung nomor klaim
// dan nomor PL - beserta periode produksinya.
//
// Refactor bentuk B (30-09-2026): dulu `repository/penomor.go` dan bagian
// `models/polis_nomor.go` / `models/polis_periode.go` yang dipakai Komite.
package penomor

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

	"nusantarare/inti/db"
)

// Penomor membaca dan menaikkan penghitung nomor bisnis.
//
// ⚠️ DIPINDAH DARI `PohonKlaim` 28-09-2026, saat tiket 03 PremiumList
// membutuhkan rantai yang sama. Ketiga methodnya - `AwalanProduksi`,
// `HariClosing`, `UrutNomorBerikut` - tidak pernah menyentuh satu pun tabel
// klaim; ketiganya membaca `KODE_PRODUKSI`, `TANGGAL_CLOSING`, dan
// `GENERATE_SEQUENCE_NUMBER`. Menggantungkannya pada pohon klaim memaksa
// modul kedua memilih antara meminjam tipe yang bukan miliknya atau menyalin
// ketiganya - dan salinan penghitung nomor adalah dua penghitung.
//
// ⛔ Yang TIDAK ikut pindah: bentuk nomornya. `RakitNomorKlaim` dan
// `models.RakitNomorPL` tetap terpisah, sebab keduanya memang berbeda -
// tahun empat angka di klaim, dua angka di premium list.
type Penomor struct {
	db *db.DB
}

// NewPenomor menyusun penomor atas satu basis data.
func NewPenomor(db *db.DB) *Penomor { return &Penomor{db: db} }

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
//
// ⛔ CACAT YANG DIPERBAIKI 28-09-2026, DITEMUKAN SAAT TIKET 03 HENDAK MEMAKAI
// FUNGSI INI. Ronde sebelumnya menggeser bulan dengan `saat.AddDate(0,1,0)`.
// Go MELIMPAHKAN tanggal yang tidak ada: 31 Januari + 1 bulan menjadi 3 MARET,
// sehingga periodenya `03.2026` - FEBRUARI TERLEWAT SAMA SEKALI. Oracle
// `ADD_MONTHS` justru MENJEPIT ke akhir bulan (29 Februari), yaitu `02.2026`.
//
// Cacat itu menyala setiap kali seseorang menerbitkan nomor pada tanggal 29-31
// bulan yang penggantinya lebih pendek, dan hasilnya nomor yang terbukukan ke
// bulan yang salah tanpa satu pun galat. Karena yang dibutuhkan hanya BULAN
// dan TAHUN, penjepitan itu setara dengan menaikkan nomor bulannya - dan
// itulah yang `models.PeriodeProduksi` kerjakan sejak tiket 02.
//
// ⛔ ATURAN PERGESERANNYA SEKARANG SATU, DIPAKAI BERDUA. Menyalinnya untuk
// premium list berarti dua aturan periode yang dapat berselisih; yang
// berselisih diam-diam adalah yang membukukan ke bulan yang salah.
// Perbedaan yang tersisa - cabang cutover - tinggal di sini, sebab hanya
// penomoran yang punya.
func HitungPeriodeNomor(saat time.Time, hariClosing int) (PeriodeNomor, error) {
	// ⛔ Hari dibaca di zona Jakarta, sama dengan pergeserannya. `TRUNC(v_now)`
	// di procedure membaca `SYSDATE`, yaitu jam server - Jakarta. Membaca
	// cutover di satu zona dan pergeseran di zona lain membuat keduanya
	// berselisih sehari di sekitar tengah malam.
	lokal := DiJakarta(saat)
	hari := time.Date(lokal.Year(), lokal.Month(), lokal.Day(), 0, 0, 0, 0, time.UTC)
	if !hari.After(batasCutover) {
		return PeriodeNomor{MMYYYY: periodeCutover, Tahun: tahunCutover}, nil
	}
	periode, err := PeriodeProduksi(hariClosing, saat)
	if err != nil {
		return PeriodeNomor{}, err
	}
	return PeriodeNomor{
		MMYYYY: fmt.Sprintf("%02d.%04d", int(periode.Month()), periode.Year()),
		Tahun:  strconv.Itoa(periode.Year()),
	}, nil
}

// HariClosing membaca hari tutup buku.
//
// `[data DBA]` `SELECT TO_NUMBER(tanggal) FROM POOLDATA.TANGGAL_CLOSING
// WHERE ROWNUM = 1` - kolomnya `VARCHAR2(10)`, dibaca sebagai bilangan.
//
// ⛔ Kosong atau tak terurai GAGAL TERANG. Menebak hari tutup buku menggeser
// periode seluruh nomor yang terbit hari itu.
func (r *Penomor) HariClosing(ctx context.Context, tx *db.Tx) (int, error) {
	tabel, err := r.db.Qualify("TANGGAL_CLOSING")
	if err != nil {
		return 0, err
	}
	q := fmt.Sprintf(`SELECT TANGGAL FROM %s WHERE ROWNUM = 1`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var teks sql.NullString
	if err := tx.QueryRowContext(ctx, q).Scan(&teks); err != nil {
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
func (r *Penomor) UrutNomorBerikut(ctx context.Context, tx *db.Tx,
	class, jenis string, p PeriodeNomor, saat time.Time) (int, error) {

	tabel, err := r.db.Qualify("GENERATE_SEQUENCE_NUMBER")
	if err != nil {
		return 0, err
	}
	qKunci := fmt.Sprintf(
		`SELECT NO_SEQ FROM %s WHERE CLASS = :1 AND JENIS = :2 AND TAHUN = :3 FOR UPDATE`,
		tabel)
	if err := db.PeriksaSQL(qKunci); err != nil {
		return 0, err
	}
	var urut int
	err = tx.QueryRowContext(ctx, qKunci, class, jenis, p.Tahun).Scan(&urut)

	switch {
	case err == sql.ErrNoRows:
		// Belum ada baris kunci ini: urut pertama, lalu INSERT.
		qSisip := fmt.Sprintf(`INSERT INTO %s
			(CLASS, JENIS, TAHUN, NO_SEQ, TANGGAL, MM_YYYY)
			VALUES (:1,:2,:3,:4,:5,:6)`, tabel)
		if err := db.PeriksaSQL(qSisip); err != nil {
			return 0, err
		}
		hasil, err := tx.ExecContext(ctx, qSisip,
			class, jenis, p.Tahun, 1, saat, p.MMYYYY)
		if err != nil {
			return 0, fmt.Errorf("repository: menyisipkan penghitung nomor: %w", err)
		}
		if err := db.PastikanSatuBaris(hasil, "penyisipan penghitung nomor"); err != nil {
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
	if err := db.PeriksaSQL(qUbah); err != nil {
		return 0, err
	}
	hasil, err := tx.ExecContext(ctx, qUbah,
		urut, p.MMYYYY, saat, class, jenis, p.Tahun)
	if err != nil {
		return 0, fmt.Errorf("repository: menaikkan penghitung nomor: %w", err)
	}
	if err := db.PastikanSatuBaris(hasil, "penaikan penghitung nomor"); err != nil {
		return 0, err
	}
	return urut, nil
}

// ErrKodeProduksiKosong - `KODE_PRODUKSI` tidak punya baris untuk tipe itu.
var ErrKodeProduksiKosong = errors.New(
	"repository: KODE_PRODUKSI tidak memberi awalan untuk tipe yang diminta")

// TipeKodeProduksiLife adalah nilai `TYPE` untuk lini Life.
//
// `[terverifikasi]` `Claim Life/RDBList/GetKodeProdLife_SQL.xml` baris 85:
// `SELECT KODE AS "ParamSeq.HASIL3" FROM POOLDATA.KODE_PRODUKSI WHERE
// TYPE ='LIFE'`.
const TipeKodeProduksiLife = "LIFE"

// AwalanProduksi membaca awalan nomor bisnis sebuah lini.
//
// ⛔ DI-LOOKUP, BUKAN KONSTANTA. Ronde sebelumnya menanam `"RNML-"` sebagai
// konstanta Go - dan AC tiket 02 melarangnya dengan kalimat yang tidak dapat
// disalahartikan: *"Prefix diperoleh lewat lookup ke `POOLDATA.KODE_PRODUKSI`
// (`TYPE='LIFE'`), tidak ditanam sebagai konstanta di kode."*
//
// Alasannya bukan kerapian: awalan itu MILIK basis data, dan lingkungan yang
// berbeda dapat memakai awalan berbeda. Konstanta Go membuat seluruh nomor
// yang terbit di lingkungan mana pun memakai awalan lingkungan yang
// kebetulan dipakai saat kode ditulis.
//
// ⚠️ `TYPE` diterima sebagai parameter, bukan ditanam: `KODE_PRODUKSI`
// melayani lebih dari satu lini, dan Claim Prop kelak membaca tabel yang sama.
func (r *Penomor) AwalanProduksi(ctx context.Context, tx *db.Tx,
	tipe string) (string, error) {

	tabel, err := r.db.Qualify("KODE_PRODUKSI")
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT KODE FROM %s WHERE TYPE = :1`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var teks sql.NullString
	if err := tx.QueryRowContext(ctx, q, tipe).Scan(&teks); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%w: %q", ErrKodeProduksiKosong, tipe)
		}
		return "", fmt.Errorf("repository: membaca awalan produksi %q: %w", tipe, err)
	}
	// ⛔ Kosong GAGAL TERANG. Awalan kosong menghasilkan nomor seperti
	// `KL1.08.2026.00936` - yang terlihat sah, tersimpan, dan baru ketahuan
	// salah ketika seseorang mencarinya dan tidak menemukannya.
	if !teks.Valid || strings.TrimSpace(teks.String) == "" {
		return "", fmt.Errorf("%w: %q", ErrKodeProduksiKosong, tipe)
	}
	// ⚠️ TIDAK dipangkas spasinya di dalam: `"RNML-"` berakhir tanda hubung,
	// dan awalan yang berakhir spasi pun milik basis data. Yang dibuang hanya
	// spasi di kedua ujung, sebab kolomnya `VARCHAR2` dan bukan `CHAR`.
	return strings.TrimSpace(teks.String), nil
}
