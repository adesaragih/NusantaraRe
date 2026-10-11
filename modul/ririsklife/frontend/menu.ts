// Butir menu modul R/I Risk - satu modul, satu menu (keputusan work owner 08-10-2026 K4). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 941, LABEL "R/I Risk"); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_RK } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_RK = 'ririsklife'
export const HALAMAN_RK = ['ririsklife-daftar'] as const
export type HalamanRK = (typeof HALAMAN_RK)[number]

/** Halaman yang dibuka tombol "R/I Risk" di sidebar. */
export const HALAMAN_AWAL_RK: HalamanRK = 'ririsklife-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanRK> = {
  nama: NAMA_RK,
  kelompok: MENU_RK.kelompok,
  halaman: HALAMAN_RK,
  halamanAwal: HALAMAN_AWAL_RK,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    ririsklife: HalamanRK
  }
}
