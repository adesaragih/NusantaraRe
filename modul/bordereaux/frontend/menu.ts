import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_BDX } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_BDX = 'bordereaux'
export const HALAMAN_BDX = ['bordereaux-daftar'] as const
export type HalamanBDX = (typeof HALAMAN_BDX)[number]

/** Halaman yang dibuka tombol "Bordereaux" di sidebar. */
export const HALAMAN_AWAL_BDX: HalamanBDX = 'bordereaux-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanBDX> = {
  nama: NAMA_BDX,
  kelompok: MENU_BDX.kelompok,
  halaman: HALAMAN_BDX,
  halamanAwal: HALAMAN_AWAL_BDX,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    bordereaux: HalamanBDX
  }
}
