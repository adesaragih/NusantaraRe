import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_AG } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_AG = 'aggregate'
export const HALAMAN_AG = ['aggregate-daftar'] as const
export type HalamanAG = (typeof HALAMAN_AG)[number]

/** Halaman yang dibuka tombol "Aggregate" di sidebar. */
export const HALAMAN_AWAL_AG: HalamanAG = 'aggregate-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanAG> = {
  nama: NAMA_AG,
  kelompok: MENU_AG.kelompok,
  halaman: HALAMAN_AG,
  halamanAwal: HALAMAN_AWAL_AG,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    aggregate: HalamanAG
  }
}
