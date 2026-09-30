// Menu modul PremiumList Life - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `ENTRI_MENU` di `lib/daftarMenu.ts`.
//
// Struktur tim satu folder per modul (30-09-2026): nama kelompok sidebar modul
// ini tinggal DI SINI, bukan di `inti/frontend/labels.ts`; dan
// `PENDAFTARAN_MENU` di bawah dibaca perakit `frontend/daftar.ts` lewat
// `import.meta.glob` - nol baris per modul di berkas bersama.

import { MENU_MODUL } from '../../../inti/frontend/labels'
import type { ButirMenuModul } from '../../../inti/frontend/lib/daftarMenu'
import type { MenuModul } from '../../../inti/frontend/modul'

/** Nama modul - SAMA dengan `const Nama` di `modul/premiumlistlife/backend/modul.go`. */
export const NAMA_PREMIUMLIST = 'premiumlistlife'

/** Nama kelompok sidebar - nama folder korpus VERBATIM (`D:\XML\RNM_BRD\PremiumList Life`). */
export const KELOMPOK_PREMIUMLIST = 'PremiumList Life'

export const HALAMAN_PREMIUMLIST = ['premiumlist'] as const
export type HalamanPremiumList = (typeof HALAMAN_PREMIUMLIST)[number]

export const MENU_PREMIUMLIST: readonly ButirMenuModul<HalamanPremiumList>[] = [
  { modul: 'premiumlist', label: MENU_MODUL.premiumList, kelompok: KELOMPOK_PREMIUMLIST },
]

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanPremiumList> = {
  nama: NAMA_PREMIUMLIST,
  kelompok: KELOMPOK_PREMIUMLIST,
  halaman: HALAMAN_PREMIUMLIST,
  menu: MENU_PREMIUMLIST,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    premiumlistlife: HalamanPremiumList
  }
}
