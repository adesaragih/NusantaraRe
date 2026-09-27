package services

// Efek keluar asinkron, antre-ulang, dan flag lingkungan - tiket 12.
//
// Untuk apa berkas ini: menjalankan yang terjadi DI LUAR basis data kita -
// berkas, email, Arasapas - tanpa membiarkan kegagalannya menahan alur klaim.
//
// Dibaca sesudah: komite.go.
//
// ⛔ TIGA efek keluar, bukan empat. Konversi ke produksi lewat payload JSON
// DIBUANG `[keputusan work owner 2026-09-16]`: hilir membaca langsung dari
// tabel klaim. Penjaga statik tiket ini memastikan nol payload keluar.
//
// Istilah:
//   - efek keluar : panggilan ke sistem di luar basis data kita.
//   - penyalur    : yang menjalankan semua efek dan mengumpulkan kegagalannya.
//   - antre-ulang : tempat kegagalan disimpan untuk dicoba lagi nanti.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/internal/repository"
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
	// ErrAntreanBelumDiputuskan - tempat antre-ulang belum disahkan.
	ErrAntreanBelumDiputuskan = errors.New(
		"services: tempat antre-ulang efek keluar belum diputuskan work owner")
)

// Lingkungan membedakan produksi dari selainnya.
type Lingkungan int

const (
	// BukanProduksi - efek keluar TIDAK dijalankan.
	BukanProduksi Lingkungan = iota
	// Produksi - efek keluar dijalankan.
	Produksi
)

// AdalahProduksi menyatakan lingkungan ini produksi.
//
// ⚠️ Namanya BUKAN `Produksi`: metode dan konstanta bernama sama membuat
// `l.Produksi()` dan `Produksi` terbaca seolah hal yang sama.
func (l Lingkungan) AdalahProduksi() bool { return l == Produksi }

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

// LayakDicobaUlang membedakan kegagalan konfigurasi dari kegagalan jaringan.
//
// ⛔ Antre-ulang yang memperlakukan keduanya sama akan berputar sia-sia:
// alamat yang tidak ada di `M_LINK_SERVICE` tidak akan muncul karena dicoba
// lagi, dan setiap putaran hanya menambah catatan kegagalan yang sama.
func LayakDicobaUlang(err error) bool {
	// ⛔ Daftar ini memuat SETIAP galat yang keadaannya PERMANEN sampai
	// seorang manusia bertindak. Ronde pertama hanya memuat tiga yang pertama
	// - dan melewatkan justru ketiga galat yang benar-benar diproduksi di
	// produksi hari ini (`…BelumDisetujui`), sehingga antrean akan memutar
	// selamanya keadaan yang tidak mungkin berubah karena dicoba lagi. Persis
	// yang AC tiket ini larang.
	for _, permanen := range []error{
		ErrEndpointTidakDitemukan,
		ErrResolverBelumDiputuskan,
		ErrAntreanBelumDiputuskan,
		ErrPenyimpananBelumDisetujui,
		ErrEmailBelumDisetujui,
		ErrArasapasBelumDisetujui,
	} {
		if errors.Is(err, permanen) {
			return false
		}
	}
	return err != nil
}

// MuatanEfek adalah yang dibawa sebuah efek keluar.
type MuatanEfek struct {
	KlaimID      string
	AdjustmentID string
	AkunID       string
	Waktu        time.Time
}

// EfekKeluar adalah satu panggilan ke sistem di luar basis data kita.
type EfekKeluar interface {
	Nama() string
	Jalankan(ctx context.Context, m MuatanEfek) error
}

// CatatanEfekGagal adalah satu kegagalan efek keluar yang perlu dicoba ulang.
//
// ⚠️ `Sebab` adalah TEKS, bukan galat. Ia disimpan dan dibaca kembali nanti,
// mungkin oleh proses lain; galat Go tidak menyeberang batas proses.
type CatatanEfekGagal struct {
	// MuatanEfek disematkan, bukan disalin medan demi medan. Ronde pertama
	// mengulang keempat medannya di sini dan menulis ulang literalnya di
	// pemanggil - tiga tempat yang harus berubah bersama, dan tidak ada yang
	// memaksanya.
	MuatanEfek
	Nama       string
	Sebab      string
	LayakUlang bool
}

// Antrean menyimpan kegagalan efek keluar untuk dicoba lagi.
//
// ⛔ Antarmuka, bukan tabel. Tempatnya menunggu butir **am** bersama tabel
// jejaknya - kegagalan efek keluar masuk JALUR AUDIT, bukan sekadar log
// layanan (ADR-U-0007).
//
// ⚠️ TANPA `*repository.Tx`. Kegagalan efek keluar justru terjadi SESUDAH
// transaksi klaimnya selesai; mengikatnya ke transaksi itu berarti kegagalan
// efek keluar dapat menggagalkan penyimpanan klaim - persis yang ADR-U-0008
// larang.
type Antrean interface {
	Antre(ctx context.Context, c CatatanEfekGagal) error
}

// AntreanBelumDiputuskan gagal terang selama tempatnya belum disahkan.
type AntreanBelumDiputuskan struct{}

// Antre selalu gagal, dan menyebut apa yang ditunggu.
func (AntreanBelumDiputuskan) Antre(context.Context, CatatanEfekGagal) error {
	return ErrAntreanBelumDiputuskan
}

// HasilSalur adalah ringkasan satu putaran efek keluar.
type HasilSalur struct {
	// Dilewati berarti lingkungannya bukan produksi - BUKAN berarti gagal.
	Dilewati bool
	Gagal    []CatatanEfekGagal
	// GagalDiantre adalah kegagalan yang bahkan tidak dapat diantre. Ia
	// dilaporkan, bukan ditelan: antre-ulang yang gagal diam-diam berarti
	// kegagalan hilang untuk selamanya.
	GagalDiantre []error
}

// Penyalur menjalankan seluruh efek keluar dan mengumpulkan kegagalannya.
type Penyalur struct {
	lingkungan Lingkungan
	antrean    Antrean
	efek       []EfekKeluar
}

// NewPenyalur menyusun penyalur atas sejumlah efek keluar.
func NewPenyalur(l Lingkungan, a Antrean, efek ...EfekKeluar) *Penyalur {
	return &Penyalur{lingkungan: l, antrean: a, efek: efek}
}

// Salurkan menjalankan seluruh efek keluar. Ia TIDAK mengembalikan galat.
//
// ⛔ Nol galat dikembalikan, dan itu disengaja (ADR-U-0008). Bila ia
// mengembalikan galat, setiap pemanggil akan tergoda meneruskannya ke atas -
// dan transisi status klaim akan tertahan oleh layanan luar yang sedang
// gagal. Yang dikembalikan RINGKASAN: pemanggil boleh melaporkannya, tidak
// boleh menggagalkan klaim karenanya.
//
// ⚠️ Konsekuensi yang diterima (ADR-U-0008): sebuah klaim dapat mencapai
// keadaan akhir sementara efek keluarnya masih tertunda.
func (p *Penyalur) Salurkan(ctx context.Context, m MuatanEfek) HasilSalur {
	if !p.lingkungan.AdalahProduksi() {
		// ⛔ Dilewati BUKAN gagal. Tidak ada yang diantre: mengantre yang
		// sengaja dilewati akan membanjiri antrean dengan pekerjaan yang
		// memang tidak diminta.
		return HasilSalur{Dilewati: true}
	}
	var hasil HasilSalur
	for _, e := range p.efek {
		// ⛔ Satu efek gagal BUKAN alasan sisanya tidak dicoba. Email yang
		// gagal tidak boleh menghalangi berkas tersimpan.
		if err := e.Jalankan(ctx, m); err != nil {
			c := CatatanEfekGagal{
				MuatanEfek: m,
				Nama:       e.Nama(),
				Sebab:      err.Error(),
				LayakUlang: LayakDicobaUlang(err),
			}
			hasil.Gagal = append(hasil.Gagal, c)
			if err := p.antrean.Antre(ctx, c); err != nil {
				hasil.GagalDiantre = append(hasil.GagalDiantre,
					fmt.Errorf("mengantre kegagalan %q: %w", e.Nama(), err))
			}
		}
	}
	return hasil
}

// LingkunganDariFlag menerjemahkan flag konfigurasi menjadi lingkungan.
//
// ⚠️ Terpisah dari `DariTingkatProduksi`. Yang satu membaca kolom Pega
// (`pzProductionLevel`), yang lain membaca konfigurasi kita sendiri
// (`IS_PEGA_PROD`, ADR-U-0005). Keduanya sengaja tidak disatukan: menyatukan
// dua sumber kebenaran atas pertanyaan "apakah ini produksi" berarti salah
// satunya diam-diam menang.
func LingkunganDariFlag(produksi bool) Lingkungan {
	if produksi {
		return Produksi
	}
	return BukanProduksi
}

// Nama ketiga efek keluar - `[keputusan work owner 2026-09-16]` TIGA, bukan
// empat. Konversi ke produksi lewat payload JSON dibuang.
const (
	NamaEfekBerkas   = "berkas"
	NamaEfekEmail    = "email"
	NamaEfekArasapas = "arasapas"
)

var (
	// ErrPenyimpananBelumDisetujui - storage nyata menuntut persetujuan manusia.
	ErrPenyimpananBelumDisetujui = errors.New(
		"services: penyimpanan berkas belum disetujui; menghubungkannya dan " +
			"memanggil GET_TOKEN_STORAGE menuntut persetujuan manusia")
	// ErrEmailBelumDisetujui - pengiriman email nyata menuntut persetujuan.
	ErrEmailBelumDisetujui = errors.New(
		"services: pengiriman email belum disetujui; menghubungkan layanan " +
			"email nyata menuntut persetujuan manusia")
	// ErrArasapasBelumDisetujui - pemanggilan Arasapas menuntut persetujuan.
	ErrArasapasBelumDisetujui = errors.New(
		"services: pemanggilan Arasapas belum disetujui; menghubungkannya " +
			"menuntut persetujuan manusia")
)

// EfekBerkas mengunggah dan mengambil berkas pendukung klaim.
//
// `[dugaan]` `Claim Life/Activity/InsertGoogleStorage_Act.xml`,
// `GetUrlGoogleStorage_Act.xml`, `DeleteGoogleStorage_Act.xml`, dan
// `RDBList/GetTokenStorage_SQL.xml` - ketiganya disebut TIKET, dan tiket 12
// belum membacanya baris demi baris. Labelnya `[dugaan]` sampai dibaca;
// `[terverifikasi]` hanya untuk yang dibaca sesi ini dengan path + baris.
//
// ⛔ Gagal terang. Memanggil prosedur itu dan menghubungkan storage nyata
// sama-sama menuntut persetujuan manusia; mengarangnya di sini berarti
// mengirim berkas sungguhan ke tempat yang belum disepakati.
type EfekBerkas struct{}

// Nama menyebut efek ini di catatan kegagalan.
func (EfekBerkas) Nama() string { return NamaEfekBerkas }

// Jalankan selalu gagal, dan menyebut apa yang ditunggu.
func (EfekBerkas) Jalankan(context.Context, MuatanEfek) error {
	return ErrPenyimpananBelumDisetujui
}

// EfekEmail memberi tahu anggota Komite saat kasus diserahkan.
//
// ⛔ TEMUAN A2: ia TIDAK memakai `M_LINK_SERVICE` sama sekali.
// `[terverifikasi]` `Claim Life/Activity/SendEmailKlaimLF.xml` pecahan baris
// 2505 `Call SendEmailWithAttachment`, dengan `smtpPort "587"` (2556) dan
// `smtpHost` berupa LITERAL di korpus (2574) - SMTP langsung, bukan REST.
//
// Brief menduga kuncinya `("SendEmail", …)` dari katalog; membaca
// activity-nya membantah dugaan itu. Inilah sebabnya kunci dibaca dari
// activity, bukan ditebak dari nama.
//
// ⛔ Host SMTP-nya TIDAK DISALIN ke mana pun - bukan ke kode, bukan ke tiket,
// bukan ke komentar ini. `[terbuka — work owner]`: dari mana alamat SMTP
// datang di sistem baru, sebab ADR-U-0013 mengatur `M_LINK_SERVICE` dan
// jalur ini tidak melewatinya.
type EfekEmail struct{}

// Nama menyebut efek ini di catatan kegagalan.
func (EfekEmail) Nama() string { return NamaEfekEmail }

// Jalankan selalu gagal, dan menyebut apa yang ditunggu.
func (EfekEmail) Jalankan(context.Context, MuatanEfek) error {
	return ErrEmailBelumDisetujui
}

// EfekArasapas mengirim klaim ke Arasapas lewat alamat yang di-resolve.
//
// `[terverifikasi]` `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml`
// pecahan baris 497 memanggil `GetLinkService`, baris 540-541 memberinya
// `Kategori_1 = "Klaim"` dan `Kategori_2 = "insertClaimLife"`, lalu baris 610
// `Connect-REST`.
//
// ⚠️ `[terbuka]` OQ-035 - kepemilikan rule itu; satu salinan, dua konteks.
// Tidak memblokir.
type EfekArasapas struct {
	Resolver ResolverEndpoint
}

// Nama menyebut efek ini di catatan kegagalan.
func (EfekArasapas) Nama() string { return NamaEfekArasapas }

// Jalankan me-resolve alamatnya lebih dulu, lalu berhenti terang.
//
// ⛔ Resolusi dijalankan SUNGGUHAN meski pemanggilannya belum disetujui.
// Itu disengaja: kegagalan konfigurasi (kunci tidak ada di `M_LINK_SERVICE`)
// harus terlihat sebagai kegagalan konfigurasi, bukan tersembunyi di balik
// "belum disetujui" - dan keduanya diperlakukan berbeda oleh antre-ulang.
func (e EfekArasapas) Jalankan(ctx context.Context, _ MuatanEfek) error {
	if e.Resolver == nil {
		return ErrResolverBelumDiputuskan
	}
	if _, err := AlamatLayanan(ctx, e.Resolver, KunciArasapasLife); err != nil {
		return err
	}
	return ErrArasapasBelumDisetujui
}

// EfekKeluarClaimLife menyusun KETIGA efek keluar modul ini, berurutan.
//
// ⛔ Tiga, dan penjaga statik memastikan tidak ada yang keempat berupa
// pengiriman payload JSON.
func EfekKeluarClaimLife(r ResolverEndpoint) []EfekKeluar {
	return []EfekKeluar{
		EfekBerkas{},
		EfekEmail{},
		EfekArasapas{Resolver: r},
	}
}

// resolverOracle membaca alamat dari `M_LINK_SERVICE` - A2, ADR-U-0013.
type resolverOracle struct{ pohon *repository.PohonKlaim }

// ResolverLinkServiceOracle menyusun resolver yang membaca saat jalan.
func ResolverLinkServiceOracle(svc *Service) ResolverEndpoint {
	return resolverOracle{pohon: repository.NewPohonKlaim(svc.db)}
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
