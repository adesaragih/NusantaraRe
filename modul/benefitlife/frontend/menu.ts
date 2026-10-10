// Butir menu modul Benefit - satu modul, satu menu (keputusan work owner 08-10-2026 K4). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 945, LABEL "Benefit"); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_BN } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_BN = 'benefitlife'
export const HALAMAN_BN = ['benefitlife-daftar'] as const
export type HalamanBN = (typeof HALAMAN_BN)[number]

/** Halaman yang dibuka tombol "Benefit" di sidebar. */
export const HALAMAN_AWAL_BN: HalamanBN = 'benefitlife-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanBN> = {
  nama: NAMA_BN,
  kelompok: MENU_BN.kelompok,
  halaman: HALAMAN_BN,
  halamanAwal: HALAMAN_AWAL_BN,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    benefitlife: HalamanBN
  }
}
