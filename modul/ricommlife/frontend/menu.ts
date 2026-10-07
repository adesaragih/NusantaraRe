// Butir menu modul R/I Comm Life - satu modul, satu menu (perintah work owner 06-10-2026). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 925); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_RC } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_RC = 'ricommlife'
export const HALAMAN_RC = ['ricommlife-daftar'] as const
export type HalamanRC = (typeof HALAMAN_RC)[number]

/** Halaman yang dibuka tombol "R/I Comm Life" di sidebar. */
export const HALAMAN_AWAL_RC: HalamanRC = 'ricommlife-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanRC> = {
  nama: NAMA_RC,
  kelompok: MENU_RC.kelompok,
  halaman: HALAMAN_RC,
  halamanAwal: HALAMAN_AWAL_RC,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    ricommlife: HalamanRC
  }
}
