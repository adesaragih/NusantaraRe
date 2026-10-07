// Menu modul EDM Treaty In. Asal pola: `modul/nbtreatyin/frontend/menu.ts` (06-10-2026).
//
// ⛔ SATU halaman: portal endorsemen. Pega punya SATU titik masuk - Harness `SFAPortal_Endorsement_Treaty` (daftar +
// tombol "Create New Addendum Treaty"); layar Create (`TreatyCreateEdm`, popup) dan layar kasus dibuka DARI portal.

import type { MenuModul } from '../../../inti/frontend/modul'
import { kotakMasuk } from './api'
import { daftarBeranda } from './beranda'

/** Nama modul - SAMA dengan `const Nama` di `modul/edmtreatyin/backend/modul.go`. */
export const NAMA_EDMTREATYIN = 'edmtreatyin'

/** Nama modul - nama folder korpus VERBATIM (= M_NAV_MENU.LABEL barisnya). */
export const KELOMPOK_EDMTREATYIN = 'EDM Treaty In'

export const HALAMAN_EDMTREATYIN = ['edmtreatyin-portal'] as const
export type HalamanEDMTreatyIn = (typeof HALAMAN_EDMTREATYIN)[number]

/** Halaman yang dibuka tombol "EDM Treaty In" di sidebar. */
export const HALAMAN_AWAL_EDMTREATYIN: HalamanEDMTreatyIn = 'edmtreatyin-portal'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanEDMTreatyIn> = {
  nama: NAMA_EDMTREATYIN,
  kelompok: KELOMPOK_EDMTREATYIN,
  halaman: HALAMAN_EDMTREATYIN,
  halamanAwal: HALAMAN_AWAL_EDMTREATYIN,
  // Kotak masuk Beranda per workbasket (keputusan work owner 07-10-2026 "YA"; pola NB): Sec Head / Dept Head
  // membuka berkas EDM dari sini - portal XML hanya menautkan kasus di posisi admin.
  antreanBeranda: kotakMasuk, // GET /kotak-masuk
  daftarBeranda, // GET /kotak-masuk/kasus
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    edmtreatyin: HalamanEDMTreatyIn
  }
}
