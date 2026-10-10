// Butir menu modul Cover Life - satu modul, satu menu (keputusan work owner 08-10-2026 C4). Labelnya datang dari
// `M_NAV_MENU.LABEL` (slot menu modul 957, LABEL "Cover Life"); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_CVL } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_CVL = 'coverlife'
export const HALAMAN_CVL = ['coverlife-daftar'] as const
export type HalamanCVL = (typeof HALAMAN_CVL)[number]

/** Halaman yang dibuka tombol "Cover Life" di sidebar. */
export const HALAMAN_AWAL_CVL: HalamanCVL = 'coverlife-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanCVL> = {
  nama: NAMA_CVL,
  kelompok: MENU_CVL.kelompok,
  halaman: HALAMAN_CVL,
  halamanAwal: HALAMAN_AWAL_CVL,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    coverlife: HalamanCVL
  }
}
