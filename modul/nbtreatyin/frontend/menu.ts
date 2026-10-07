// Menu modul NB Treaty In.
//
// ⛔ SATU halaman: portal realisasi. Pega punya SATU titik masuk untuk alur ini -
// Harness `SFAPortalOpportunities` (daftar + tombol Create); layar admin dan
// layar atasan dibuka DARI daftar itu (INVENTARIS-XML.md bab 3: "sistem baru
// tidak boleh menambah menu di luar itu").

import type { MenuModul } from '../../../inti/frontend/modul'
import { kotakMasuk } from './api'
import { daftarBeranda } from './beranda'

/** Nama modul - SAMA dengan `const Nama` di `modul/nbtreatyin/backend/modul.go`. */
export const NAMA_NBTREATYIN = 'nbtreatyin'

/** Nama modul - nama folder korpus VERBATIM (= M_NAV_MENU.LABEL barisnya). */
export const KELOMPOK_NBTREATYIN = 'NB Treaty In'

export const HALAMAN_NBTREATYIN = ['nbtreatyin-portal'] as const
export type HalamanNBTreatyIn = (typeof HALAMAN_NBTREATYIN)[number]

/** Halaman yang dibuka tombol "NB Treaty In" di sidebar. */
export const HALAMAN_AWAL_NBTREATYIN: HalamanNBTreatyIn = 'nbtreatyin-portal'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanNBTreatyIn> = {
  nama: NAMA_NBTREATYIN,
  kelompok: KELOMPOK_NBTREATYIN,
  halaman: HALAMAN_NBTREATYIN,
  halamanAwal: HALAMAN_AWAL_NBTREATYIN,
  // kotak masuk Beranda per workbasket yang dipegang akun (keputusan work owner 06-10-2026)
  antreanBeranda: kotakMasuk,
  // daftar berkas kotak masuk di Beranda, tanpa masuk menu ini (keputusan work owner 06-10-2026)
  daftarBeranda,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    nbtreatyin: HalamanNBTreatyIn
  }
}
