// Menu modul Claim Life - refactor bentuk B (30-09-2026): dipindah apa adanya
// dari `ENTRI_MENU` di `lib/daftarMenu.ts`.

import { MENU, MODUL } from '../../../inti/frontend/labels'
import type { ButirMenuModul } from '../../../inti/frontend/lib/daftarMenu'

/** Nama modul - SAMA dengan `const Nama` di `modul/claimlife/modul.go`. */
export const NAMA_CLAIMLIFE = 'claimlife'

/**
 * Halaman modul ini.
 *
 * ⚠️ `outstanding` dan `detail` BUKAN butir menu: di Pega keduanya dibuka
 * DARI DALAM kasus (flow action), bukan dari navigasi.
 */
export const HALAMAN_CLAIMLIFE = ['inbox', 'register', 'outstanding', 'detail'] as const
export type HalamanClaimLife = (typeof HALAMAN_CLAIMLIFE)[number]

export const MENU_CLAIMLIFE: readonly ButirMenuModul<HalamanClaimLife>[] = [
  { modul: 'inbox', label: MENU.inbox, kelompok: MODUL.claimLife },
  { modul: 'register', label: MENU.register, kelompok: MODUL.claimLife },
]
