// Menu modul Nation - satu tombol di grup MASTER (keputusan work owner 04-10-2026: satu menu per master, "8 modul
// terpisah"). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.

import type { MenuModul } from '../../../inti/frontend/modul'
import { KELOMPOK_MASTERNATION } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_MASTERNATION = 'masternation'

export const HALAMAN_MASTERNATION = ['masternation-daftar'] as const
export type HalamanNation = (typeof HALAMAN_MASTERNATION)[number]

/** Halaman yang dibuka tombol "Nation" di sidebar. */
export const HALAMAN_AWAL_MASTERNATION: HalamanNation = 'masternation-daftar'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanNation> = {
  nama: NAMA_MASTERNATION,
  kelompok: KELOMPOK_MASTERNATION,
  halaman: HALAMAN_MASTERNATION,
  halamanAwal: HALAMAN_AWAL_MASTERNATION,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    masternation: HalamanNation
  }
}
