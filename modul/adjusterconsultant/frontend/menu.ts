// Butir menu modul Adjuster Consultant - satu modul, satu menu (keputusan work owner 30-09-2026). Labelnya datang dari
// `M_NAV_MENU.LABEL` (migrasi inti 915); `kelompok` di sini harus sama dengannya.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_ADJ } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_ADJ = 'adjusterconsultant'
export const HALAMAN_ADJ = ['adjusterconsultant-daftar'] as const
export type HalamanADJ = (typeof HALAMAN_ADJ)[number]

/** Halaman yang dibuka tombol "Adjuster Consultant" di sidebar. */
export const HALAMAN_AWAL_ADJ: HalamanADJ = 'adjusterconsultant-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanADJ> = {
  nama: NAMA_ADJ,
  kelompok: MENU_ADJ.kelompok,
  halaman: HALAMAN_ADJ,
  halamanAwal: HALAMAN_AWAL_ADJ,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    adjusterconsultant: HalamanADJ
  }
}
