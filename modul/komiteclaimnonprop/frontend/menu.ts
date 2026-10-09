// Butir menu modul Komite Claim Non Prop - satu modul, satu menu (baris M_NAV_MENU `komiteclaimnonprop` yang sudah ada,
// grup KLAIM, migrasi inti 900; nol butir menu baru). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat
// `import.meta.glob`. Menunya disembunyikan lewat data (perintah work owner 09-10-2026, pola Komite Claim Prop): modul
// ini dipasang bagi pemegang menu `claimnonprop` (`MODUL_DIPINJAM` frontend/App.tsx).
//
// Tanpa `antreanBeranda`: kotak masuk Beranda bukan bagian Pega.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_KCNP } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_KCNP = 'komiteclaimnonprop'
export const HALAMAN_KCNP = ['komiteclaimnonprop-daftar'] as const
export type HalamanKCNP = (typeof HALAMAN_KCNP)[number]

/** Halaman awal - daftar kerja penyetuju (worklist KomiteRouter). */
export const HALAMAN_AWAL_KCNP: HalamanKCNP = 'komiteclaimnonprop-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanKCNP> = {
  nama: NAMA_KCNP,
  kelompok: MENU_KCNP.kelompok,
  halaman: HALAMAN_KCNP,
  halamanAwal: HALAMAN_AWAL_KCNP,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    komiteclaimnonprop: HalamanKCNP
  }
}
