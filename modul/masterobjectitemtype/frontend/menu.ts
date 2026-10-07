// Menu modul Object Item Type - satu tombol di grup MASTER (keputusan work owner 04-10-2026: satu menu per master, "8 modul
// terpisah"). `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.

import type { MenuModul } from '../../../inti/frontend/modul'
import { KELOMPOK_MASTEROBJECTITEMTYPE } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_MASTEROBJECTITEMTYPE = 'masterobjectitemtype'

export const HALAMAN_MASTEROBJECTITEMTYPE = ['masterobjectitemtype-daftar'] as const
export type HalamanObjectItemType = (typeof HALAMAN_MASTEROBJECTITEMTYPE)[number]

/** Halaman yang dibuka tombol "Object Item Type" di sidebar. */
export const HALAMAN_AWAL_MASTEROBJECTITEMTYPE: HalamanObjectItemType = 'masterobjectitemtype-daftar'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanObjectItemType> = {
  nama: NAMA_MASTEROBJECTITEMTYPE,
  kelompok: KELOMPOK_MASTEROBJECTITEMTYPE,
  halaman: HALAMAN_MASTEROBJECTITEMTYPE,
  halamanAwal: HALAMAN_AWAL_MASTEROBJECTITEMTYPE,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    masterobjectitemtype: HalamanObjectItemType
  }
}
