// Package layanan menyelesaikan alamat layanan keluar dari `M_LINK_SERVICE`
// (ADR-U-0013) dan menerbitkan token penyimpanan dari `GCP_IMAGE`.
//
// Refactor bentuk B (30-09-2026): dulu bagian `services/efekkeluar.go`,
// `services/tokenstorage.go`, `repository/linkservice.go`, dan
// `repository/tokenstorage.go`.
package layanan

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
)

// Kunci kategori - SATU-SATUNYA hal tentang alamat yang boleh jadi konstanta.
//
// ⛔ ADR-U-0013: alamat endpoint di-resolve saat jalan lewat `M_LINK_SERVICE`,
// dan TIDAK boleh ada URL sebagai literal, konstanta, maupun env var di kode.
// Pemisahan dev-prod terjadi lewat ISI TABEL per-database, bukan lewat
// percabangan di kode. Penjaga statik berkas ini menegakkannya.
//
// `[terverifikasi]` `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml`
// pecahan baris 540-541: `Kategori_1 = "Klaim"`, `Kategori_2 =
// "insertClaimLife"`.
const (
	KategoriKlaim           = "Klaim"
	KategoriInsertClaimLife = "insertClaimLife"
	KategoriGoogle          = "Google"
	KategoriUpload          = "upload"
	KategoriGetURL          = "geturl"
	KategoriHapus           = "delete"
)

// Kunci per efek, DIBACA dari activity-nya masing-masing - bukan ditebak dari
// nama katalog.
//
// `[terverifikasi]` pecahan baris:
//
//	InsertGoogleStorage_Act  1775 · 1776  ("Google", "upload")
//	GetUrlGoogleStorage_Act  1711 · 1712  ("Google", "geturl")
//	DeleteGoogleStorage_Act  1349 · 1351  ("Google", "delete")
//	serviceInsertArasapas…    540 ·  541  ("Klaim",  "insertClaimLife")
var (
	KunciUnggahBerkas = KunciLayanan{Kategori1: KategoriGoogle, Kategori2: KategoriUpload}
	KunciURLBerkas    = KunciLayanan{Kategori1: KategoriGoogle, Kategori2: KategoriGetURL}
	KunciHapusBerkas  = KunciLayanan{Kategori1: KategoriGoogle, Kategori2: KategoriHapus}
)

// KunciLayanan adalah pasangan `(KATEGORI_1, KATEGORI_2)`.
//
// `[terverifikasi]` `Claim Life/Activity/GetLinkService.xml` pecahan baris 371
// `Obj-Browse` atas `M_LINK_SERVICE`, disaring `.KATEGORI_1` (baris 491) dan
// `.KATEGORI_2` (baris 517), mengambil `.URL` (baris 393).
type KunciLayanan struct {
	Kategori1 string
	Kategori2 string
}

// KunciArasapasLife adalah kunci endpoint Arasapas untuk Claim Life.
var KunciArasapasLife = KunciLayanan{
	Kategori1: KategoriKlaim,
	Kategori2: KategoriInsertClaimLife,
}

var (
	// ErrEndpointTidakDitemukan - kunci kategori tidak ada di M_LINK_SERVICE.
	ErrEndpointTidakDitemukan = errors.New(
		"services: kunci kategori tidak ditemukan di M_LINK_SERVICE")
	// ErrResolverBelumDiputuskan - pembacaan M_LINK_SERVICE belum disahkan.
	ErrResolverBelumDiputuskan = errors.New(
		"services: pembacaan M_LINK_SERVICE belum disahkan work owner")
)

// ⛔ `DariTingkatProduksi` DIBUANG. Ia menerjemahkan `pzProductionLevel`
// - `[terverifikasi]` `Claim Life/When/IsPEGAPROD.xml` pecahan baris 164 dan
// 318, `pxProcess.pzProductionLevel = "5"` - tetapi aplikasi ini TIDAK PERNAH
// membaca kolom itu; lingkungannya datang dari `config.IsPegaProd`. Kode mati
// yang menyimpan pengetahuan tetap kode mati: pengetahuannya hidup di bab
// pembacaan XML tiket 12, lengkap dengan nomor barisnya.

// ResolverEndpoint mencari alamat sebuah layanan dari kunci kategorinya.
//
// ⛔ Antarmuka, bukan query. `M_LINK_SERVICE` tabel PRODUKSI; membacanya
// menuntut persetujuan manusia.
type ResolverEndpoint interface {
	Resolve(ctx context.Context, kunci KunciLayanan) (string, error)
}

// ResolverBelumDiputuskan gagal terang selama tabelnya belum boleh dibaca.
type ResolverBelumDiputuskan struct{}

// Resolve selalu gagal, dan menyebut apa yang ditunggu.
func (ResolverBelumDiputuskan) Resolve(context.Context, KunciLayanan) (string, error) {
	return "", ErrResolverBelumDiputuskan
}

// AlamatLayanan mencari alamat, dan menolak alamat kosong.
//
// ⛔ PENYIMPANGAN SADAR. `[terverifikasi]` `GetLinkService.xml` pecahan baris
// 705-706 menyetel `ResponLink.URL = linkService.pxResults(1).URL` pada
// langkah ber-`pyStepsPreCondition` KOSONG (baris 701): `Obj-Browse` yang
// tidak menemukan apa pun menghasilkan URL KOSONG tanpa satu pun galat, dan
// `Connect-REST` sesudahnya menembak alamat kosong. Di sini itu gagal terang.
func AlamatLayanan(ctx context.Context, r ResolverEndpoint,
	kunci KunciLayanan) (string, error) {

	alamat, err := r.Resolve(ctx, kunci)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(alamat) == "" {
		return "", fmt.Errorf("%w: (%s, %s)",
			ErrEndpointTidakDitemukan, kunci.Kategori1, kunci.Kategori2)
	}
	return alamat, nil
}

// resolverOracle membaca alamat dari `M_LINK_SERVICE` - A2, ADR-U-0013.
type resolverOracle struct {
	pohon *PembacaLinkService
}

// ResolverLinkServiceOracle menyusun resolver yang membaca saat jalan.
func ResolverLinkServiceOracle(svc inti.Akar) ResolverEndpoint {
	return resolverOracle{pohon: NewPembacaLinkService(svc.DB())}
}

// Resolve menerjemahkan kunci kategori menjadi alamat.
//
// ⛔ Yang dikembalikan hanya URL-nya. `USERNAME` sengaja TIDAK diteruskan:
// belum ada pemanggil yang memerlukannya, dan nilai rahasia yang beredar
// tanpa pemakai adalah nilai rahasia yang akhirnya tercetak di suatu tempat.
func (r resolverOracle) Resolve(ctx context.Context, kunci KunciLayanan) (string, error) {
	alamat, err := r.pohon.AmbilAlamatLayanan(ctx, kunci.Kategori1, kunci.Kategori2)
	if err != nil {
		return "", err
	}
	return alamat.URL, nil
}
