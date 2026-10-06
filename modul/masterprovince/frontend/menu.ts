// Menu modul Province - satu tombol di grup MASTER (keputusan work owner 04-10-2026: satu menu per master, "8 modul
// terpisah"). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.

import type { MenuModul } from '../../../inti/frontend/modul'
import { KELOMPOK_MASTERPROVINCE } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_MASTERPROVINCE = 'masterprovince'

export const HALAMAN_MASTERPROVINCE = ['masterprovince-daftar'] as const
export type HalamanProvince = (typeof HALAMAN_MASTERPROVINCE)[number]

/** Halaman yang dibuka tombol "Province" di sidebar. */
export const HALAMAN_AWAL_MASTERPROVINCE: HalamanProvince = 'masterprovince-daftar'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanProvince> = {
  nama: NAMA_MASTERPROVINCE,
  kelompok: KELOMPOK_MASTERPROVINCE,
  halaman: HALAMAN_MASTERPROVINCE,
  halamanAwal: HALAMAN_AWAL_MASTERPROVINCE,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    masterprovince: HalamanProvince
  }
}
