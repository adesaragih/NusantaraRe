import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_MO } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_MO = 'marketingofficer'
export const HALAMAN_MO = ['marketingofficer-daftar'] as const
export type HalamanMO = (typeof HALAMAN_MO)[number]

/** Halaman yang dibuka tombol "Marketing Officer" di sidebar. */
export const HALAMAN_AWAL_MO: HalamanMO = 'marketingofficer-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanMO> = {
  nama: NAMA_MO,
  kelompok: MENU_MO.kelompok,
  halaman: HALAMAN_MO,
  halamanAwal: HALAMAN_AWAL_MO,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    marketingofficer: HalamanMO
  }
}
