package services

// Antre-ulang efek keluar dan pekerjanya - butir aq, A2.
//
// Untuk apa berkas ini: menutup stub `AntreanBelumDiputuskan`. Kegagalan efek
// keluar disimpan ke `T_LOG_SERVICE_RNM` dengan jadwal percobaan berikutnya, dan
// seorang pekerja memungutnya kembali satu per satu.
//
// ⛔ Kegagalan PERMANEN tidak pernah dijadwalkan ulang. `LayakDicobaUlang`
// sudah memisahkannya; yang berubah di sini hanya akibatnya - `gagal-permanen`
// tanpa jadwal. Antrean yang menjadwalkan ulang kegagalan konfigurasi akan
// memutari keadaan yang tidak mungkin berubah karena dicoba lagi.
//
// Dibaca sesudah: efekkeluar.go.

import (
	"nusantarare/inti/layanan"
	"nusantarare/inti/outbox"
)

// ModulClaimLife mengisi kolom `MODUL` baris outbox milik Claim Life.
const ModulClaimLife = "CLAIMLIFE"

// AntreanEfekOracle menyusun antre-ulang yang menulis ke `T_LOG_SERVICE_RNM`.
func AntreanEfekOracle(svc *Service) outbox.Antrean {
	return outbox.AntreanEfekOracleModul(svc, ModulClaimLife)
}

// NewPekerjaEfek menyusun pekerjanya.
//
// ⚠️ Pelaksana satu-satunya hari ini (`PelaksanaBerkasLokal`) milik Claim
// Life, jadi pekerjanya pun Claim Life. Baris PREMIUMLISTLIFE menunggu pekerja
// modulnya sendiri - yang belum ada (tiket 06).
//
// ⚠️ Penjejaknya kini ikut membawa `ModulClaimLife`; dulu kosong. Medan itu
// tidak dibaca jalur penjejak (`rekamMenyerah` hanya menulis jejak), jadi
// nol perilaku berubah.
func NewPekerjaEfek(svc *Service, p outbox.PelaksanaEfek) *outbox.PekerjaEfek {
	return outbox.NewPekerjaEfekModul(svc, p, ModulClaimLife)
}

// PenyalurClaimLifeOracle menyusun penyalur dengan KETIGA ketergantungan yang
// sudah punya jalan ke Oracle - antre-ulang dan resolver endpoint.
//
// ⛔ Lingkungannya tetap datang dari `Service`, bukan dari sini. Penyalur yang
// memutuskan sendiri "ini produksi" adalah penyalur yang suatu hari mengirim
// email kepada orang sungguhan dari mesin pengembang.
//
// ⚠️ Ketiga EFEKnya masih gagal terang (`…BelumDisetujui`) - menghubungkan
// storage, email, dan Arasapas nyata menuntut persetujuan manusia. Yang
// ditutup butir **aq** adalah TEMPAT kegagalannya mendarat, bukan izin
// memanggil layanannya.
func PenyalurClaimLifeOracle(svc *Service) *outbox.Penyalur {
	return outbox.NewPenyalur(svc.Lingkungan(), AntreanEfekOracle(svc),
		EfekKeluarClaimLife(layanan.ResolverLinkServiceOracle(svc))...)
}
