// Butir menu modul Treaty Exchange Yearly - satu modul, satu menu (keputusan work owner 05-10-2026). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 919); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_TEY } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_TEY = 'treatyexchangeyearly'
export const HALAMAN_TEY = ['treatyexchangeyearly-daftar'] as const
export type HalamanTEY = (typeof HALAMAN_TEY)[number]

/** Halaman yang dibuka tombol "Treaty Exchange Yearly" di sidebar. */
export const HALAMAN_AWAL_TEY: HalamanTEY = 'treatyexchangeyearly-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanTEY> = {
  nama: NAMA_TEY,
  kelompok: MENU_TEY.kelompok,
  halaman: HALAMAN_TEY,
  halamanAwal: HALAMAN_AWAL_TEY,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    treatyexchangeyearly: HalamanTEY
  }
}
