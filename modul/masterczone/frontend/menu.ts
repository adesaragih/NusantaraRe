// Menu modul CZone - satu tombol di grup MASTER (keputusan work owner 04-10-2026: satu menu per master, "8 modul
// terpisah"). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.

import type { MenuModul } from '../../../inti/frontend/modul'
import { KELOMPOK_MASTERCZONE } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_MASTERCZONE = 'masterczone'

export const HALAMAN_MASTERCZONE = ['masterczone-daftar'] as const
export type HalamanCZone = (typeof HALAMAN_MASTERCZONE)[number]

/** Halaman yang dibuka tombol "CZone" di sidebar. */
export const HALAMAN_AWAL_MASTERCZONE: HalamanCZone = 'masterczone-daftar'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanCZone> = {
  nama: NAMA_MASTERCZONE,
  kelompok: KELOMPOK_MASTERCZONE,
  halaman: HALAMAN_MASTERCZONE,
  halamanAwal: HALAMAN_AWAL_MASTERCZONE,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    masterczone: HalamanCZone
  }
}
