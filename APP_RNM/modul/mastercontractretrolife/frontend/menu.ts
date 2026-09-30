// Menu modul Master Contract Retro Life - paket 10.
//
// Menu DATAR (keputusan work owner 30-09-2026, `PROMPT-MENU-DATAR-PER-GROUPMENU.md`): satu modul satu
// menu - tombol "Master Contract Retro Life" di bawah GROUPMENU MASTER (baris `M_NAV_MENU` isi awal 900,
// dinyalakan slot `958_menu_mastercontractretrolife.sql`) membuka halaman awal di bawah. Tiga harness
// lain dan `ViewRate`/`ViewRateTable` dibuka DARI DALAM layar, bukan tombol menu.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul
// di berkas bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_MCRL } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `modul/mastercontractretrolife/backend/modul.go`. */
export const NAMA_MCRL = 'mastercontractretrolife'

/**
 * Halaman modul ini - SATU: `GridRetrocessionLife` (judul `MASTER CONTRACT RETRO LIFE`). Keempat harness
 * korpus adalah popup yang dibuka dari dalamnya (PARITAS §1), jadi bukan halaman sendiri.
 */
export const HALAMAN_MCRL = ['mcrl-tahun'] as const
export type HalamanMCRL = (typeof HALAMAN_MCRL)[number]

/**
 * Halaman awal - bukti: `Section/GridRetrocessionLife.xml` b1017 (judul) menyertakan
 * `InputRetrocessionLife` b2475; tidak dirujuk rule mana pun di korpus modul (dibuka portal), sedangkan
 * keempat harness dibuka `showHarness` popup dari tombol `ReinsType` b11988, `Business List` b12999,
 * `Reinsurer List` b14034, `Security Reinsurer` b12441.
 */
export const HALAMAN_AWAL_MCRL: HalamanMCRL = 'mcrl-tahun'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanMCRL> = {
  nama: NAMA_MCRL,
  kelompok: MENU_MCRL.kelompok,
  halaman: HALAMAN_MCRL,
  halamanAwal: HALAMAN_AWAL_MCRL,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    mastercontractretrolife: HalamanMCRL
  }
}
