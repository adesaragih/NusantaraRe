// Menu modul Master Data.
//
// Menu DATAR (keputusan work owner 30-09-2026): satu modul satu menu - tombol "Master Data" di bawah GROUPMENU MASTER
// (baris `M_NAV_MENU` migrasi inti 904, dinyalakan slot `990_menu_masterdata.sql`) membuka halaman awal di bawah.
// Daftar per tabel master dibuka DARI DALAM halaman ini (rencana M-1), bukan tombol menu.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.

import type { MenuModul } from '../../../inti/frontend/modul'

/** Nama modul - SAMA dengan `const Nama` di `modul/masterdata/backend/modul.go`. */
export const NAMA_MASTERDATA = 'masterdata'

/** Nama modul - `Folder korpus` di MODUL.md modul ini (= M_NAV_MENU.LABEL, `MODUL_LUAR_KORPUS.masterData`). */
export const KELOMPOK_MASTERDATA = 'Master Data'

export const HALAMAN_MASTERDATA = ['masterdata-daftar'] as const
export type HalamanMasterData = (typeof HALAMAN_MASTERDATA)[number]

/** Halaman yang dibuka tombol "Master Data" di sidebar. */
export const HALAMAN_AWAL_MASTERDATA: HalamanMasterData = 'masterdata-daftar'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanMasterData> = {
  nama: NAMA_MASTERDATA,
  kelompok: KELOMPOK_MASTERDATA,
  halaman: HALAMAN_MASTERDATA,
  halamanAwal: HALAMAN_AWAL_MASTERDATA,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    masterdata: HalamanMasterData
  }
}
