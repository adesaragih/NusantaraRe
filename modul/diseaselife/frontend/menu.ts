// Butir menu modul Disease Life - satu modul, satu menu (keputusan work owner 08-10-2026 D4). Labelnya datang dari
// `M_NAV_MENU.LABEL` (slot menu modul 951, LABEL "Disease Life"); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_DSL } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_DSL = 'diseaselife'
export const HALAMAN_DSL = ['diseaselife-daftar'] as const
export type HalamanDSL = (typeof HALAMAN_DSL)[number]

/** Halaman yang dibuka tombol "Disease Life" di sidebar. */
export const HALAMAN_AWAL_DSL: HalamanDSL = 'diseaselife-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanDSL> = {
  nama: NAMA_DSL,
  kelompok: MENU_DSL.kelompok,
  halaman: HALAMAN_DSL,
  halamanAwal: HALAMAN_AWAL_DSL,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    diseaselife: HalamanDSL
  }
}
