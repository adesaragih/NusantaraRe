// Menu modul District - satu tombol di grup MASTER (keputusan work owner 04-10-2026: satu menu per master, "8 modul
// terpisah"). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.

import type { MenuModul } from '../../../inti/frontend/modul'
import { KELOMPOK_MASTERDISTRICT } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_MASTERDISTRICT = 'masterdistrict'

export const HALAMAN_MASTERDISTRICT = ['masterdistrict-daftar'] as const
export type HalamanDistrict = (typeof HALAMAN_MASTERDISTRICT)[number]

/** Halaman yang dibuka tombol "District" di sidebar. */
export const HALAMAN_AWAL_MASTERDISTRICT: HalamanDistrict = 'masterdistrict-daftar'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanDistrict> = {
  nama: NAMA_MASTERDISTRICT,
  kelompok: KELOMPOK_MASTERDISTRICT,
  halaman: HALAMAN_MASTERDISTRICT,
  halamanAwal: HALAMAN_AWAL_MASTERDISTRICT,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    masterdistrict: HalamanDistrict
  }
}
