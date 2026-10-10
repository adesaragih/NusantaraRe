// Butir menu modul Claim Prop - satu modul, satu menu (baris M_NAV_MENU `claimprop` yang sudah ada, grup KLAIM; nol
// butir menu baru). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.
//
// ⛔ Tanpa `antreanBeranda`: Claim Prop TIDAK masuk kotak masuk Beranda (bukan bagian Pega; prompt §2).

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_CP } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_CP = 'claimprop'
export const HALAMAN_CP = ['claimprop-daftar'] as const
export type HalamanCP = (typeof HALAMAN_CP)[number]

/** Halaman yang dibuka tombol "Claim Prop" di sidebar. */
export const HALAMAN_AWAL_CP: HalamanCP = 'claimprop-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanCP> = {
  nama: NAMA_CP,
  kelompok: MENU_CP.kelompok,
  halaman: HALAMAN_CP,
  halamanAwal: HALAMAN_AWAL_CP,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    claimprop: HalamanCP
  }
}
