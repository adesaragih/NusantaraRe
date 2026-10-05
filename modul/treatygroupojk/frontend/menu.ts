// Butir menu modul Treaty Group OJK - satu modul, satu menu (keputusan work owner 05-10-2026). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 916); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_TGO } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_TGO = 'treatygroupojk'
export const HALAMAN_TGO = ['treatygroupojk-daftar'] as const
export type HalamanTGO = (typeof HALAMAN_TGO)[number]

/** Halaman yang dibuka tombol "Treaty Group OJK" di sidebar. */
export const HALAMAN_AWAL_TGO: HalamanTGO = 'treatygroupojk-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanTGO> = {
  nama: NAMA_TGO,
  kelompok: MENU_TGO.kelompok,
  halaman: HALAMAN_TGO,
  halamanAwal: HALAMAN_AWAL_TGO,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    treatygroupojk: HalamanTGO
  }
}
