// Menu modul Accumulated Type - satu tombol di grup MASTER (keputusan work owner 04-10-2026: satu menu per master, "8 modul
// terpisah"). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.

import type { MenuModul } from '../../../inti/frontend/modul'
import { KELOMPOK_MASTERACCUMULATEDTYPE } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_MASTERACCUMULATEDTYPE = 'masteraccumulatedtype'

export const HALAMAN_MASTERACCUMULATEDTYPE = ['masteraccumulatedtype-daftar'] as const
export type HalamanAccumulatedType = (typeof HALAMAN_MASTERACCUMULATEDTYPE)[number]

/** Halaman yang dibuka tombol "Accumulated Type" di sidebar. */
export const HALAMAN_AWAL_MASTERACCUMULATEDTYPE: HalamanAccumulatedType = 'masteraccumulatedtype-daftar'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanAccumulatedType> = {
  nama: NAMA_MASTERACCUMULATEDTYPE,
  kelompok: KELOMPOK_MASTERACCUMULATEDTYPE,
  halaman: HALAMAN_MASTERACCUMULATEDTYPE,
  halamanAwal: HALAMAN_AWAL_MASTERACCUMULATEDTYPE,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    masteraccumulatedtype: HalamanAccumulatedType
  }
}
