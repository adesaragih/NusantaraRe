// Menu modul City - satu tombol di grup MASTER (keputusan work owner 04-10-2026: satu menu per master, "8 modul
// terpisah"). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.

import type { MenuModul } from '../../../inti/frontend/modul'
import { KELOMPOK_MASTERCITY } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_MASTERCITY = 'mastercity'

export const HALAMAN_MASTERCITY = ['mastercity-daftar'] as const
export type HalamanCity = (typeof HALAMAN_MASTERCITY)[number]

/** Halaman yang dibuka tombol "City" di sidebar. */
export const HALAMAN_AWAL_MASTERCITY: HalamanCity = 'mastercity-daftar'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanCity> = {
  nama: NAMA_MASTERCITY,
  kelompok: KELOMPOK_MASTERCITY,
  halaman: HALAMAN_MASTERCITY,
  halamanAwal: HALAMAN_AWAL_MASTERCITY,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    mastercity: HalamanCity
  }
}
