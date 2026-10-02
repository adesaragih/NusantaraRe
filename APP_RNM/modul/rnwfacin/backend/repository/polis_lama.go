// Package repository memuat antarmuka baca modul RNW Fac In.
package repository

import (
	"context"
	"errors"

	"nusantarare/modul/rnwfacin/backend/models"
)

// PembacaPolisLama membaca data kasus polis yang akan diperpanjang, berkunci
// nomor polis (`OldPolicyNo`, tiket R05).
//
// `[pertanyaan terbuka]` sumbernya: gerbang masuk renewal ada di konfigurasi
// work type/portal yang tidak terekspor (tiket R01); `GetData_ACT` hanya membuka
// kasus yang sudah ada. Implementasi menunggu peta tabel polis lama (DBA).
type PembacaPolisLama interface {
	Baca(ctx context.Context, nomorPolis string) (models.Kasus, error)
}

// ErrPolisLamaTakDitemukan - tidak ada polis dengan nomor itu.
var ErrPolisLamaTakDitemukan = errors.New("repository: polis lama tidak ditemukan")

// SumberDaftarRenewal - seluruh kasus renewal sebelum disaring (tiket R05).
// Penyaringan, urutan, dan batas baris RenewalList_RD ada di layanan, bukan di sini.
//
// `[pertanyaan terbuka]` sumbernya: RenewalList_RD berkelas
// `ASM-FW-GISFW-Work-Renewal`, irisan pelaporan atas kelas kerja NB; tabel
// Oracle-nya tidak ada di korpus (DBA).
type SumberDaftarRenewal interface {
	KasusRenewal(ctx context.Context) ([]models.BarisDaftarRenewal, error)
}
