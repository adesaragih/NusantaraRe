package repository

// Pengenal work object - identitas baris T_WORK_CLAIM.
//
// Pemilik: tiket 02 (register klaim).
//
// Untuk apa berkas ini: seluruh tabel T_CLAIMLF_* memakai identitas angka dari
// sequence (ADR-U-0006), tetapi T_WORK_CLAIM tidak - identitasnya TEKS
// BERFORMAT, `CLM-xxxxxx` untuk baris klaim dan `KMTLF-xxxxxx` untuk baris
// komite. Itu penyimpangan sadar yang dicatat tiket 14 AC 34, sebab nomor itu
// dibaca manusia dan muncul di layar.
//
// `[keputusan work owner 26-09-2026, butir aa]` sesudah tiket 14 AC 35 terbuka
// tiga sesi: angka urutannya datang dari SEQ_WORK_CLAIM (langkah migrasi 009),
// dan awalan serta LPAD dirakit DI SINI - bukan di DDL - sebab awalan
// bergantung jenis baris, dan itu aturan dagang.
//
// Dibaca sesudah: pohonklaim.go.

import (
	"context"
	"fmt"
	"strings"

	"nusantarare/inti/db"
)

// Awalan pengenal work object, satu per jenis baris.
const (
	AwalanKlaim = "CLM-"
	// ⛔ RALAT A2, 27-09-2026: `KMT-` menjadi `KMTLF-`.
	//
	// Sumbernya `.scratch/komite-claim-life/STRUKTUR-TABEL-KOMITE-CLAIM-LIFE.md`,
	// yang brief lanjutan 4 §1 sebut sebagai bentuk yang dipakai:
	// `T_GENERAL_KOMITE.ID` = `T_WORK_CLAIM.ID` baris komite, **berawalan
	// `KMTLF-`**.
	//
	// Alasannya terbaca dari dokumen yang sama: `T_GENERAL_KOMITE` LINTAS-LINI
	// (LIFE dan PROP), sehingga awalan yang tidak menyebut lini membuat dua
	// modul berbagi ruang nomor yang sama.
	AwalanKomite = "KMTLF-"
)

// lebarUrutanWork adalah jumlah digit sesudah awalan.
//
// Enam digit menampung 999.999 work object per jenis. ⛔ TANPA reset tahunan:
// korpus tidak memuat satu pun bukti bahwa Pega me-reset urutan ini per tahun,
// dan mengarang reset membuat dua baris bernomor sama pada tahun berbeda.
const lebarUrutanWork = 6

// PengenalWorkBerikut membentuk satu identitas work object baru.
//
// Awalannya diperiksa: nilai di luar kedua awalan yang dikenal ditolak, sebab
// pengenal berawalan karangan akan lolos ke basis data dan baru terlihat salah
// berbulan-bulan kemudian, saat seseorang mencarinya dan tidak menemukannya.
func (r *PohonKlaim) PengenalWorkBerikut(ctx context.Context, tx *db.Tx, awalan string) (string, error) {
	switch awalan {
	case AwalanKlaim, AwalanKomite:
	default:
		return "", fmt.Errorf("repository: awalan pengenal work %q tidak dikenal; "+
			"yang sah hanya %q dan %q", awalan, AwalanKlaim, AwalanKomite)
	}
	urut, err := r.db.NomorBerikut(ctx, tx, "SEQ_WORK_CLAIM")
	if err != nil {
		return "", err
	}
	return RakitPengenalWork(awalan, urut), nil
}

// RakitPengenalWork menyusun pengenal dari awalan dan angka urutannya.
//
// Dipisah dari pembacaan sequence supaya bentuknya dapat diuji tanpa Oracle.
// Urutan yang lebih panjang dari enam digit TIDAK dipotong: memotongnya akan
// membuat pengenal berulang, dan pengenal berulang jauh lebih buruk daripada
// pengenal yang kepanjangan.
func RakitPengenalWork(awalan, urut string) string {
	urut = strings.TrimSpace(urut)
	if len(urut) < lebarUrutanWork {
		urut = strings.Repeat("0", lebarUrutanWork-len(urut)) + urut
	}
	return awalan + urut
}
