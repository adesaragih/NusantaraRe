// Butir menu modul Claim Fac In - satu modul, satu menu (baris M_NAV_MENU `claimfacin` yang sudah ada sejak 900, grup
// KLAIM, DIMIGRASI lewat slot 978; nol butir menu baru). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat
// `import.meta.glob`. Pola `modul/claimnonprop/frontend/menu.ts` (disalin, bukan impor).
//
// ⛔ Tanpa `antreanBeranda`: Claim Fac In TIDAK masuk kotak masuk Beranda (pola Claim Prop / Claim Non Prop).

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_CFI } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_CFI = 'claimfacin'
export const HALAMAN_CFI = ['claimfacin-daftar'] as const
export type HalamanCFI = (typeof HALAMAN_CFI)[number]

/** Halaman yang dibuka tombol "Claim Fac In" di sidebar. */
export const HALAMAN_AWAL_CFI: HalamanCFI = 'claimfacin-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanCFI> = {
  nama: NAMA_CFI,
  kelompok: MENU_CFI.kelompok,
  halaman: HALAMAN_CFI,
  halamanAwal: HALAMAN_AWAL_CFI,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    claimfacin: HalamanCFI
  }
}
