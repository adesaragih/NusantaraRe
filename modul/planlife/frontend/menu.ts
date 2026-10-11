// Butir menu modul Plan - satu modul, satu menu (keputusan work owner 08-10-2026 K6). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 949, LABEL "Plan"); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_PL } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_PL = 'planlife'
export const HALAMAN_PL = ['planlife-daftar'] as const
export type HalamanPL = (typeof HALAMAN_PL)[number]

/** Halaman yang dibuka tombol "Plan" di sidebar. */
export const HALAMAN_AWAL_PL: HalamanPL = 'planlife-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanPL> = {
  nama: NAMA_PL,
  kelompok: MENU_PL.kelompok,
  halaman: HALAMAN_PL,
  halamanAwal: HALAMAN_AWAL_PL,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    planlife: HalamanPL
  }
}
