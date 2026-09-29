package services

// Efek keluar PremiumList Life - tiket 06.
//
// Untuk apa berkas ini: yang terjadi SESUDAH premium list tersimpan -
// kiriman ke Arasapas, dan alarm bila kiriman itu gagal. Keduanya berjalan
// di luar transaksi simpan, dan kegagalannya tidak pernah membatalkan simpan
// (ADR-U-0008).
//
// `[terverifikasi]` `Activity/InsertJsonPolisLife_Act.xml` langkah 14
// `SendEmailNotification` ("Kalau blm, email errornya") dan 15
// `serviceInsertArasapasLife_act`, keduanya dijaga `IsPEGAPROD` (b5128,
// b5236). Langkah 14 dijaga juga `PL_NUMBER==""` (b5105): email adalah ALARM,
// bukan notifikasi bisnis `[keputusan work owner]`.
//
// ⛔ PEMICU ALARMNYA BERGESER, dengan sadar (AC relasional tiket 06, spec §6).
// Di Pega alarm menyala bila pembacaan balik menemukan simpan yang separuh.
// Di sini simpan atomik - keadaan separuh mustahil - jadi alarm menyala pada
// satu-satunya kegagalan yang tersisa: kegagalan efek keluar itu sendiri.
//
// ⛔ Kedua efek masih STUB yang gagal terang - memanggil Arasapas dan
// mengirim email nyata menuntut persetujuan manusia (CLAUDE.md §10). Yang
// dibangun tiket ini: urutan, gerbang lingkungan, outbox, dan kunci alamat.
//
// Dibaca sesudah: efekkeluar.go, antrean.go (milik Claim Life, dipakai ulang).

import (
	"context"

	"nusantarare/inti"
	"nusantarare/inti/layanan"
	"nusantarare/inti/outbox"
)

// ModulPremiumListLife mengisi kolom `MODUL` outbox `T_LOG_SERVICE_RNM`.
//
// ⚠️ Sejajar `ModulClaimLife` (antrean.go). Outbox berbagi tabel dengan Claim
// Life; kolom `MODUL` yang memisahkan baris kedua modul.
const ModulPremiumListLife = "PREMIUMLISTLIFE"

// KunciArasapasPremiumList adalah kunci endpoint Arasapas modul ini.
//
// `[terverifikasi]` `Activity/serviceInsertArasapasLife_act.xml` langkah 4
// `GetLinkService` (b799), parameter `Kategori_1 = "Production"` dan
// `Kategori_2 = "convertJsonNusareToProduction"` (b844-845 pecahan baris).
//
// ⛔ BERBEDA dari `KunciArasapasLife` Claim Life (`"Klaim"`,
// `"insertClaimLife"`) - dua activity, dua kunci, dibaca masing-masing.
//
// ⛔ OQ-PL-11 `[terbuka — work owner / pemilik Arasapas]`. Nama layanannya
// `convertJsonNusareToProduction`, dan `ConnectREST/ConvertJsonNusareToProduction.xml`
// hanya mengirim TIGA pengenal (`.pzInsKey`, `.OfferFacIn.PolicyData.PolicyNo`,
// `.pxCreateDateTime`) - bukan isi polis. Layanan hilir tampaknya membaca
// `JSON_POLIS` sendiri, dan tabel itu TIDAK LAGI ditulis (pl1). Sebelum stub
// ini diganti panggilan nyata, pemilik layanan harus menyatakan dari mana ia
// membaca polis sesudah JSON dibuang.
var KunciArasapasPremiumList = layanan.KunciLayanan{
	Kategori1: "Production",
	Kategori2: "convertJsonNusareToProduction",
}

// Nama kedua efek - masuk kolom `JENIS_EFEK` outbox.
const (
	NamaEfekArasapasPolis = "arasapas-polis"
	NamaEfekAlarmPolis    = "email-alarm-polis"
)

// EfekArasapasPolis mengirim polis ke Arasapas lewat alamat yang di-resolve.
//
// ⛔ Resolusi dijalankan SUNGGUHAN (ADR-U-0013) meski panggilannya belum
// disetujui - sebab yang sama dengan `EfekArasapas` Claim Life: kunci yang
// tidak ada di `M_LINK_SERVICE` harus terlihat sebagai kegagalan konfigurasi.
type EfekArasapasPolis struct {
	Resolver layanan.ResolverEndpoint
}

// Nama menyebut efek ini di catatan kegagalan.
func (EfekArasapasPolis) Nama() string { return NamaEfekArasapasPolis }

// Jalankan me-resolve alamatnya lebih dulu, lalu berhenti terang.
func (e EfekArasapasPolis) Jalankan(ctx context.Context, _ outbox.MuatanEfek) error {
	if e.Resolver == nil {
		return layanan.ErrResolverBelumDiputuskan
	}
	if _, err := layanan.AlamatLayanan(ctx, e.Resolver, KunciArasapasPremiumList); err != nil {
		return err
	}
	return outbox.ErrArasapasBelumDisetujui
}

// EfekAlarmEmailPolis memberi tahu tim operasi bahwa efek keluar gagal.
//
// ⚠️ `[belum terverifikasi]` penerima dan isi email: `SendEmailNotification`
// langkah 14 belum dibaca baris demi baris untuk alamat tujuannya, dan alamat
// orang tidak disalin ke kode mana pun (CLAUDE.md §4 butir 10).
type EfekAlarmEmailPolis struct{}

// Nama menyebut efek ini di catatan kegagalan.
func (EfekAlarmEmailPolis) Nama() string { return NamaEfekAlarmPolis }

// Jalankan selalu gagal, dan menyebut apa yang ditunggu.
func (EfekAlarmEmailPolis) Jalankan(context.Context, outbox.MuatanEfek) error {
	return outbox.ErrEmailBelumDisetujui
}

// PenyalurPolis menjalankan Arasapas, lalu alarm HANYA bila Arasapas gagal.
//
// ⛔ Dua `Penyalur` bersarang, bukan satu berisi dua efek. `Penyalur`
// menjalankan SETIAP efeknya; alarm yang ikut di dalamnya akan terkirim pada
// setiap simpan yang berhasil - persis yang AC 27 spec larang.
type PenyalurPolis struct {
	kirim *outbox.Penyalur
	alarm *outbox.Penyalur
}

// NewPenyalurPolis menyusunnya. Gerbang lingkungan milik `Penyalur` - SATU
// tempat (AC tiket 06), dipakai kedua modul.
func NewPenyalurPolis(l inti.Lingkungan, a outbox.Antrean, kirim, alarm outbox.EfekKeluar) *PenyalurPolis {
	return &PenyalurPolis{
		kirim: outbox.NewPenyalur(l, a, kirim),
		alarm: outbox.NewPenyalur(l, a, alarm),
	}
}

// Salurkan menjalankan kiriman, lalu alarm bila kiriman gagal. Nol galat -
// alasannya sama dengan `Penyalur.Salurkan`.
func (p *PenyalurPolis) Salurkan(ctx context.Context, m outbox.MuatanEfek) outbox.HasilSalur {
	h := p.kirim.Salurkan(ctx, m)
	if h.Dilewati || len(h.Gagal) == 0 {
		return h
	}
	a := p.alarm.Salurkan(ctx, m)
	h.Gagal = append(h.Gagal, a.Gagal...)
	h.GagalDiantre = append(h.GagalDiantre, a.GagalDiantre...)
	return h
}

// penyalurPolisBawaan gagal terang di produksi, dan dilewati di selainnya.
func penyalurPolisBawaan(s *Service) *PenyalurPolis {
	return NewPenyalurPolis(s.Lingkungan(), outbox.AntreanBelumDiputuskan{},
		EfekArasapasPolis{Resolver: layanan.ResolverBelumDiputuskan{}}, EfekAlarmEmailPolis{})
}

// PenyalurPremiumListOracle menyusun penyalur dengan outbox dan resolver Oracle.
//
// ⚠️ Kedua EFEKnya tetap stub (`…BelumDisetujui`); yang tersambung ke Oracle
// adalah TEMPAT kegagalannya mendarat (`T_LOG_SERVICE_RNM`, `MODUL =
// PREMIUMLISTLIFE`) dan pembacaan `M_LINK_SERVICE`.
func PenyalurPremiumListOracle(svc *Service) *PenyalurPolis {
	return NewPenyalurPolis(svc.Lingkungan(), outbox.AntreanEfekOracleModul(svc, ModulPremiumListLife),
		EfekArasapasPolis{Resolver: layanan.ResolverLinkServiceOracle(svc)}, EfekAlarmEmailPolis{})
}
