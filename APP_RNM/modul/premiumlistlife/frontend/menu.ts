// Menu modul PremiumList Life - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `ENTRI_MENU` di `lib/daftarMenu.ts`.

import { MENU_MODUL, MODUL } from '../../../inti/frontend/labels'
import type { ButirMenuModul } from '../../../inti/frontend/lib/daftarMenu'

/** Nama modul - SAMA dengan `const Nama` di `modul/premiumlistlife/modul.go`. */
export const NAMA_PREMIUMLIST = 'premiumlistlife'

export const HALAMAN_PREMIUMLIST = ['premiumlist'] as const
export type HalamanPremiumList = (typeof HALAMAN_PREMIUMLIST)[number]

export const MENU_PREMIUMLIST: readonly ButirMenuModul<HalamanPremiumList>[] = [
  { modul: 'premiumlist', label: MENU_MODUL.premiumList, kelompok: MODUL.premiumListLife },
]
