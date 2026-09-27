package repository

// Pembaca alamat layanan keluar - A2, ADR-U-0013.
//
// Untuk apa berkas ini: menerjemahkan kunci `(KATEGORI_1, KATEGORI_2)` menjadi
// alamat, SAAT JALAN, dari `M_LINK_SERVICE`.
//
// `[terverifikasi]` `Claim Life/Activity/GetLinkService.xml` pecahan baris 371
// `Obj-Browse`, disaring `.KATEGORI_1` (491) dan `.KATEGORI_2` (517),
// mengambil `.URL` (393), hasilnya `linkService.pxResults(1).URL` (705-706).
//
// ⛔ PENYIMPANGAN SADAR pada hasil KOSONG. Langkah yang membaca
// `pxResults(1).URL` ber-`pyStepsPreCondition` KOSONG (baris 701), sehingga
// `Obj-Browse` yang tidak menemukan apa pun menghasilkan URL kosong tanpa satu
// pun galat - dan `Connect-REST` sesudahnya menembak alamat kosong. Di sini
// nol baris GAGAL TERANG.
//
// ⛔ `URL` dan `USERNAME` TIDAK PERNAH masuk log maupun artefak. Galat di
// berkas ini menyebut KUNCInya, tidak pernah nilainya.
//
// Dibaca sesudah: roster.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrAlamatLayananTidakAda - kunci kategori tidak ada di `M_LINK_SERVICE`.
var ErrAlamatLayananTidakAda = errors.New(
	"repository: kunci kategori tidak ada di M_LINK_SERVICE")

// AlamatLayanan adalah satu baris hasil resolusi.
//
// ⚠️ Kedua medannya RAHASIA OPERASIONAL. Ia dipakai pemanggil lalu dibuang;
// nol pencetakan, nol penyimpanan ke artefak.
type AlamatLayanan struct {
	URL      string
	Username string
}

// sqlAlamatLayanan merakit pernyataannya.
func sqlAlamatLayanan(tabel string) string {
	return fmt.Sprintf(
		`SELECT URL, USERNAME FROM %s WHERE KATEGORI_1 = :1 AND KATEGORI_2 = :2`,
		tabel)
}

// AmbilAlamatLayanan mencari alamat sebuah kunci kategori.
//
// ⚠️ Nol baris GAGAL TERANG - lihat penyimpangan sadar di kepala berkas.
func (r *PohonKlaim) AmbilAlamatLayanan(ctx context.Context,
	kategori1, kategori2 string) (AlamatLayanan, error) {

	tabel, err := r.db.Qualify("M_LINK_SERVICE")
	if err != nil {
		return AlamatLayanan{}, err
	}
	q := sqlAlamatLayanan(tabel)
	if err := PeriksaSQL(q); err != nil {
		return AlamatLayanan{}, err
	}
	var url, user sql.NullString
	err = r.db.sql.QueryRowContext(ctx, q, kategori1, kategori2).Scan(&url, &user)
	if err == sql.ErrNoRows {
		// ⛔ Kuncinya disebut, nilainya tidak - tidak ada nilai untuk disebut.
		return AlamatLayanan{}, fmt.Errorf("%w: (%s, %s)",
			ErrAlamatLayananTidakAda, kategori1, kategori2)
	}
	if err != nil {
		// ⛔ Galat dibungkus TANPA menyertakan barisnya: pesan galat driver
		// dapat memuat nilai kolom yang gagal dibaca.
		return AlamatLayanan{}, fmt.Errorf(
			"repository: membaca alamat layanan (%s, %s): %w",
			kategori1, kategori2, err)
	}
	return AlamatLayanan{URL: url.String, Username: user.String}, nil
}
