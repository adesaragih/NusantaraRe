// Butir menu modul Treaty Description - satu modul, satu menu (keputusan work owner 05-10-2026). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 920); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_TD } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_TD = 'treatydescription'
export const HALAMAN_TD = ['treatydescription-daftar'] as const
export type HalamanTD = (typeof HALAMAN_TD)[number]

/** Halaman yang dibuka tombol "Treaty Description" di sidebar. */
export const HALAMAN_AWAL_TD: HalamanTD = 'treatydescription-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanTD> = {
  nama: NAMA_TD,
  kelompok: MENU_TD.kelompok,
  halaman: HALAMAN_TD,
  halamanAwal: HALAMAN_AWAL_TD,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    treatydescription: HalamanTD
  }
}
