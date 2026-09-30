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
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
)

// EfekKeluarClaimLife menyusun KETIGA efek keluar modul ini, berurutan.
//
// ⛔ Tiga, dan penjaga statik memastikan tidak ada yang keempat berupa
// pengiriman payload JSON.
func EfekKeluarClaimLife(r layanan.ResolverEndpoint) []outbox.EfekKeluar {
	return []outbox.EfekKeluar{
		outbox.EfekBerkas{},
		outbox.EfekEmail{},
		outbox.EfekArasapas{Resolver: r},
	}
}
