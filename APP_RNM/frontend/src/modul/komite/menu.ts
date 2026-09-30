// Menu modul Komite Claim Life - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `ENTRI_MENU` di `lib/daftarMenu.ts`.

import { MENU_MODUL, MODUL } from '../../inti/labels'
import type { ButirMenuModul } from '../../inti/lib/daftarMenu'

/** Nama modul - SAMA dengan `const Nama` di `modul/komiteclaimlife/modul.go`. */
export const NAMA_KOMITE = 'komiteclaimlife'

export const HALAMAN_KOMITE = ['komite'] as const
export type HalamanKomite = (typeof HALAMAN_KOMITE)[number]

export const MENU_KOMITE: readonly ButirMenuModul<HalamanKomite>[] = [
  { modul: 'komite', label: MENU_MODUL.inboxKomite, kelompok: MODUL.komiteClaimLife },
]
