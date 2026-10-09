// Butir menu modul Claim Non Prop - satu modul, satu menu (baris M_NAV_MENU `claimnonprop` yang sudah ada, grup KLAIM;
// nol butir menu baru). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.
//
// ⛔ Tanpa `antreanBeranda`: Claim Non Prop TIDAK masuk kotak masuk Beranda (pola Claim Prop; bukan bagian Pega).

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_CNP } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_CNP = 'claimnonprop'
export const HALAMAN_CNP = ['claimnonprop-daftar'] as const
export type HalamanCNP = (typeof HALAMAN_CNP)[number]

/** Halaman yang dibuka tombol "Claim Non Prop" di sidebar. */
export const HALAMAN_AWAL_CNP: HalamanCNP = 'claimnonprop-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanCNP> = {
  nama: NAMA_CNP,
  kelompok: MENU_CNP.kelompok,
  halaman: HALAMAN_CNP,
  halamanAwal: HALAMAN_AWAL_CNP,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    claimnonprop: HalamanCNP
  }
}
