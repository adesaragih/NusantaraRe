// Butir menu modul Treaty Group - satu modul, satu menu (keputusan work owner 05-10-2026). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 917); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_TG } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_TG = 'treatygroup'
export const HALAMAN_TG = ['treatygroup-daftar'] as const
export type HalamanTG = (typeof HALAMAN_TG)[number]

/** Halaman yang dibuka tombol "Treaty Group" di sidebar. */
export const HALAMAN_AWAL_TG: HalamanTG = 'treatygroup-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanTG> = {
  nama: NAMA_TG,
  kelompok: MENU_TG.kelompok,
  halaman: HALAMAN_TG,
  halamanAwal: HALAMAN_AWAL_TG,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    treatygroup: HalamanTG
  }
}
