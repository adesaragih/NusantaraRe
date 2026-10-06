// Butir menu modul R/I Rate Life - satu modul, satu menu (perintah work owner 05-10-2026). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 922); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_RR } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_RR = 'riratelife'
export const HALAMAN_RR = ['riratelife-daftar'] as const
export type HalamanRR = (typeof HALAMAN_RR)[number]

/** Halaman yang dibuka tombol "R/I Rate Life" di sidebar. */
export const HALAMAN_AWAL_RR: HalamanRR = 'riratelife-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanRR> = {
  nama: NAMA_RR,
  kelompok: MENU_RR.kelompok,
  halaman: HALAMAN_RR,
  halamanAwal: HALAMAN_AWAL_RR,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    riratelife: HalamanRR
  }
}
