// Menu modul Claim Life - refactor bentuk B (30-09-2026): dipindah apa adanya
// dari `ENTRI_MENU` di `lib/daftarMenu.ts`.
//
// Struktur tim satu folder per modul (30-09-2026): nama modul tinggal DI SINI,
// bukan di `inti/frontend/labels.ts`; dan `PENDAFTARAN_MENU` di bawah dibaca
// perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul
// di berkas bersama.
//
// Menu DATAR (keputusan work owner 30-09-2026, `PROMPT-MENU-DATAR-PER-GROUPMENU.md`):
// satu modul satu menu. Tombol "Claim Life" di sidebar membuka HALAMAN AWAL di
// bawah; butir `inbox` dan `register` dicabut. Register dibuka tombol Register
// di Inbox Claim Life (`onRegister`, `Register_Flow.xml:155`) - XML Claim Life
// tidak memuat jalan masuk Register lain (nol rule Portal/Navigation;
// `Register_Flow.xml` b69 `pyCanStartInteractively` dan b83
// `pyCanCreateWorkObject` false).

import type { MenuModul } from '../../../inti/frontend/modul'

/** Nama modul - SAMA dengan `const Nama` di `modul/claimlife/backend/modul.go`. */
export const NAMA_CLAIMLIFE = 'claimlife'

/**
 * Nama modul - nama folder korpus VERBATIM.
 *
 * `[terverifikasi]` `Flow/Register_Flow.xml:270`
 * `<pyWorkTypeName>ClaimLife</pyWorkTypeName>`.
 */
export const KELOMPOK_CLAIMLIFE = 'Claim Life'

/**
 * Halaman modul ini.
 *
 * ⚠️ `register`, `outstanding`, dan `detail` dibuka DARI DALAM halaman awal:
 * Register dari tombol kepala Inbox, dua lainnya dari baris kasus.
 */
export const HALAMAN_CLAIMLIFE = ['inbox', 'register', 'outstanding', 'detail'] as const
export type HalamanClaimLife = (typeof HALAMAN_CLAIMLIFE)[number]

/** Halaman yang dibuka tombol "Claim Life" - dulu butir pertamanya. */
export const HALAMAN_AWAL_CLAIMLIFE: HalamanClaimLife = 'inbox'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanClaimLife> = {
  nama: NAMA_CLAIMLIFE,
  kelompok: KELOMPOK_CLAIMLIFE,
  halaman: HALAMAN_CLAIMLIFE,
  halamanAwal: HALAMAN_AWAL_CLAIMLIFE,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    claimlife: HalamanClaimLife
  }
}
