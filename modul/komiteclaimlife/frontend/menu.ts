// Menu modul Komite Claim Life - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `ENTRI_MENU` di `lib/daftarMenu.ts`.
//
// Struktur tim satu folder per modul (30-09-2026): nama modul tinggal DI SINI,
// bukan di `inti/frontend/labels.ts`; dan `PENDAFTARAN_MENU` di bawah dibaca
// perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul
// di berkas bersama.
//
// Menu DATAR (keputusan work owner 30-09-2026): satu modul satu menu - tombol
// "Komite Claim Life" membuka halaman awal di bawah; butir `komite` dicabut.

import type { MenuModul } from '../../../inti/frontend/modul'

/** Nama modul - SAMA dengan `const Nama` di `modul/komiteclaimlife/backend/modul.go`. */
export const NAMA_KOMITE = 'komiteclaimlife'

/** Nama modul - nama folder korpus VERBATIM (`D:\XML\RNM_BRD\Komite Claim Life`). */
export const KELOMPOK_KOMITE = 'Komite Claim Life'

export const HALAMAN_KOMITE = ['komite'] as const
export type HalamanKomite = (typeof HALAMAN_KOMITE)[number]

/** Halaman yang dibuka tombol "Komite Claim Life" - dulu butir pertamanya. */
export const HALAMAN_AWAL_KOMITE: HalamanKomite = 'komite'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanKomite> = {
  nama: NAMA_KOMITE,
  kelompok: KELOMPOK_KOMITE,
  halaman: HALAMAN_KOMITE,
  halamanAwal: HALAMAN_AWAL_KOMITE,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    komiteclaimlife: HalamanKomite
  }
}
