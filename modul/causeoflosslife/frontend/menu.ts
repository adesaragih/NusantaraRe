// Butir menu modul Cause Of Loss Life - satu modul, satu menu (keputusan work owner 08-10-2026 K5). Labelnya datang dari
// `M_NAV_MENU.LABEL` (slot menu modul 955, LABEL "Cause Of Loss Life"); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_COL } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_COL = 'causeoflosslife'
export const HALAMAN_COL = ['causeoflosslife-daftar'] as const
export type HalamanCOL = (typeof HALAMAN_COL)[number]

/** Halaman yang dibuka tombol "Cause Of Loss Life" di sidebar. */
export const HALAMAN_AWAL_COL: HalamanCOL = 'causeoflosslife-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanCOL> = {
  nama: NAMA_COL,
  kelompok: MENU_COL.kelompok,
  halaman: HALAMAN_COL,
  halamanAwal: HALAMAN_AWAL_COL,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    causeoflosslife: HalamanCOL
  }
}
