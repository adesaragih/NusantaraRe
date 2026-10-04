// Package outbox menjalankan efek keluar (berkas, email, Arasapas) tanpa
// menahan alur bisnis, dan menyimpan kegagalannya ke `T_LOG_SERVICE_RNM`
// untuk dicoba ulang oleh pekerja per modul.
//
// Refactor bentuk B (30-09-2026): dulu bagian `services/efekkeluar.go`,
// `services/antrean.go`, dan `repository/efekkeluar.go`.
package outbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/layanan"
)

var (
	// ErrAntreanBelumDiputuskan - tempat antre-ulang belum disahkan.
	ErrAntreanBelumDiputuskan = errors.New(
		"services: tempat antre-ulang efek keluar belum diputuskan work owner")
)

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
		layanan.ErrEndpointTidakDitemukan,
		layanan.ErrResolverBelumDiputuskan,
		ErrAntreanBelumDiputuskan,
		ErrPenyimpananBelumDisetujui,
		ErrEmailBelumDisetujui,
		ErrArasapasBelumDisetujui,
		// Tiket 07 Komite - aditif: keadaan permanen sampai manusia menyetujui.
		ErrKasirBelumDisetujui,
		ErrPengirimStubNonProduksi,
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
	lingkungan inti.Lingkungan
	antrean    Antrean
	efek       []EfekKeluar
}

// NewPenyalur menyusun penyalur atas sejumlah efek keluar.
func NewPenyalur(l inti.Lingkungan, a Antrean, efek ...EfekKeluar) *Penyalur {
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
// Komite memakai ulang efek ini (komite_pengirim.go); sumbernya
// `Komite Claim Life/Activity/SendEmailKlaimLife.xml` langkah 15 b3479
// (`IsPEGAPROD` b3799). ⚠️ Sensus remark 28-09-2026: CC di langkah 1-2
// ter-remark; CC HIDUP ada di langkah 3 (b704-705) dan BCC tetap (b3540) -
// keduanya belum ditiru, dicatat di PARITAS Komite.
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
	Resolver layanan.ResolverEndpoint
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
		return layanan.ErrResolverBelumDiputuskan
	}
	if _, err := layanan.AlamatLayanan(ctx, e.Resolver, layanan.KunciArasapasLife); err != nil {
		return err
	}
	return ErrArasapasBelumDisetujui
}

// ErrPengirimStubNonProduksi - non-produksi TIDAK mengirim, dan TIDAK mengaku terkirim.
//
// ⛔ Temuan /code-review giliran 10: stub yang menuntaskan baris `selesai`
// tidak terbedakan dari kiriman nyata - anti-dobel lalu menganggap Kasir
// SUDAH dibayar bila basis data non-produksi kelak dipromosikan atau flag
// lingkungannya berubah. Galat permanen ini membuat barisnya berhenti di
// `gagal-permanen` dengan sebab yang menyebut dirinya, bukan `selesai`.
var ErrPengirimStubNonProduksi = errors.New(
	"services: lingkungan bukan produksi - efek Komite TIDAK dikirim (pengirim stub); " +
		"baris ini bukan kiriman yang berhasil")

// ErrKasirBelumDisetujui - pemanggilan Kasir nyata menuntut persetujuan.
//
// ⛔ Dan ia perilaku BARU: di korpus REST Kasir ter-remark
// (`HitServiceToKasirKMTLife_Act` b3033), jadi sistem lama tidak pernah
// memanggilnya dari jalur ini (sensus remark 28-09-2026, OQ-K-06).
var ErrKasirBelumDisetujui = errors.New(
	"services: pemanggilan Kasir belum disetujui; menghubungkan layanan pembayaran " +
		"menuntut persetujuan manusia (OQ-002: kontrak Kasir tidak ada di korpus)")
