// Menu modul Claim Life - refactor bentuk B (30-09-2026): dipindah apa adanya
// dari `ENTRI_MENU` di `lib/daftarMenu.ts`.
//
// Struktur tim satu folder per modul (30-09-2026): nama kelompok sidebar modul
// ini tinggal DI SINI, bukan di `inti/frontend/labels.ts`; dan
// `PENDAFTARAN_MENU` di bawah dibaca perakit `frontend/daftar.ts` lewat
// `import.meta.glob` - nol baris per modul di berkas bersama.

import { MENU } from '../../../inti/frontend/labels'
import type { ButirMenuModul } from '../../../inti/frontend/lib/daftarMenu'
import type { MenuModul } from '../../../inti/frontend/modul'

/** Nama modul - SAMA dengan `const Nama` di `modul/claimlife/backend/modul.go`. */
export const NAMA_CLAIMLIFE = 'claimlife'

/**
 * Nama kelompok sidebar - nama folder korpus VERBATIM.
 *
 * `[terverifikasi]` `Flow/Register_Flow.xml:270`
 * `<pyWorkTypeName>ClaimLife</pyWorkTypeName>`.
 */
export const KELOMPOK_CLAIMLIFE = 'Claim Life'

/**
 * Halaman modul ini.
 *
 * ⚠️ `outstanding` dan `detail` BUKAN butir menu: di Pega keduanya dibuka
 * DARI DALAM kasus (flow action), bukan dari navigasi.
 */
export const HALAMAN_CLAIMLIFE = ['inbox', 'register', 'outstanding', 'detail'] as const
export type HalamanClaimLife = (typeof HALAMAN_CLAIMLIFE)[number]

export const MENU_CLAIMLIFE: readonly ButirMenuModul<HalamanClaimLife>[] = [
  { modul: 'inbox', label: MENU.inbox, kelompok: KELOMPOK_CLAIMLIFE },
  { modul: 'register', label: MENU.register, kelompok: KELOMPOK_CLAIMLIFE },
]

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanClaimLife> = {
  nama: NAMA_CLAIMLIFE,
  kelompok: KELOMPOK_CLAIMLIFE,
  halaman: HALAMAN_CLAIMLIFE,
  menu: MENU_CLAIMLIFE,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    claimlife: HalamanClaimLife
  }
}
