// Menu modul Accumulation - satu tombol di grup MASTER (keputusan work owner 04-10-2026: satu menu per master, "8 modul
// terpisah"). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.

import type { MenuModul } from '../../../inti/frontend/modul'
import { KELOMPOK_MASTERACCUMULATION } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_MASTERACCUMULATION = 'masteraccumulation'

export const HALAMAN_MASTERACCUMULATION = ['masteraccumulation-daftar'] as const
export type HalamanAccumulation = (typeof HALAMAN_MASTERACCUMULATION)[number]

/** Halaman yang dibuka tombol "Accumulation" di sidebar. */
export const HALAMAN_AWAL_MASTERACCUMULATION: HalamanAccumulation = 'masteraccumulation-daftar'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanAccumulation> = {
  nama: NAMA_MASTERACCUMULATION,
  kelompok: KELOMPOK_MASTERACCUMULATION,
  halaman: HALAMAN_MASTERACCUMULATION,
  halamanAwal: HALAMAN_AWAL_MASTERACCUMULATION,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    masteraccumulation: HalamanAccumulation
  }
}
