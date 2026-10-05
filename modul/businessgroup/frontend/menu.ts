// Butir menu modul Business Group - satu modul, satu menu (keputusan work owner 05-10-2026). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 918); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_BG } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_BG = 'businessgroup'
export const HALAMAN_BG = ['businessgroup-daftar'] as const
export type HalamanBG = (typeof HALAMAN_BG)[number]

/** Halaman yang dibuka tombol "Business Group" di sidebar. */
export const HALAMAN_AWAL_BG: HalamanBG = 'businessgroup-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanBG> = {
  nama: NAMA_BG,
  kelompok: MENU_BG.kelompok,
  halaman: HALAMAN_BG,
  halamanAwal: HALAMAN_AWAL_BG,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    businessgroup: HalamanBG
  }
}
