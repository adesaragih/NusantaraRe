// Butir menu modul Reinsurance Type - satu modul, satu menu (keputusan work owner 05-10-2026). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 921); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_RT } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_RT = 'reinsurancetype'
export const HALAMAN_RT = ['reinsurancetype-daftar'] as const
export type HalamanRT = (typeof HALAMAN_RT)[number]

/** Halaman yang dibuka tombol "Reinsurance Type" di sidebar. */
export const HALAMAN_AWAL_RT: HalamanRT = 'reinsurancetype-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanRT> = {
  nama: NAMA_RT,
  kelompok: MENU_RT.kelompok,
  halaman: HALAMAN_RT,
  halamanAwal: HALAMAN_AWAL_RT,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    reinsurancetype: HalamanRT
  }
}
