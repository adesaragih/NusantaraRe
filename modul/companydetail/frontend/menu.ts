import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_CD } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_CD = 'companydetail'
export const HALAMAN_CD = ['companydetail-daftar'] as const
export type HalamanCD = (typeof HALAMAN_CD)[number]

/** Halaman yang dibuka tombol "Company Detail" di sidebar. */
export const HALAMAN_AWAL_CD: HalamanCD = 'companydetail-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanCD> = {
  nama: NAMA_CD,
  kelompok: MENU_CD.kelompok,
  halaman: HALAMAN_CD,
  halamanAwal: HALAMAN_AWAL_CD,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    companydetail: HalamanCD
  }
}
