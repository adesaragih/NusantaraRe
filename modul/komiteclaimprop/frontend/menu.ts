// Butir menu modul Komite Claim Prop - satu modul, satu menu (baris M_NAV_MENU `komiteclaimprop` yang sudah ada, grup
// KLAIM, migrasi inti 900; nol butir menu baru). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat
// `import.meta.glob`.
//
// ⛔ Tanpa `antreanBeranda`: kotak masuk Beranda bukan bagian Pega (prompt §2).

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_KCP } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_KCP = 'komiteclaimprop'
export const HALAMAN_KCP = ['komiteclaimprop-daftar'] as const
export type HalamanKCP = (typeof HALAMAN_KCP)[number]

/** Halaman yang dibuka tombol "Komite Claim Prop" di sidebar - daftar kerja penyetuju (worklist KomiteRouter). */
export const HALAMAN_AWAL_KCP: HalamanKCP = 'komiteclaimprop-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanKCP> = {
  nama: NAMA_KCP,
  kelompok: MENU_KCP.kelompok,
  halaman: HALAMAN_KCP,
  halamanAwal: HALAMAN_AWAL_KCP,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    komiteclaimprop: HalamanKCP
  }
}
