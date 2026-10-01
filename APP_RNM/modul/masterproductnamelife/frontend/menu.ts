// Menu modul Master Product Name Life - paket 10.
//
// Menu DATAR (keputusan work owner 30-09-2026): satu modul satu menu - tombol "Master Product Name Life" di
// bawah GROUPMENU MASTER (baris `M_NAV_MENU` isi awal 900, dinyalakan slot `960_menu_masterproductnamelife.sql`)
// membuka halaman awal di bawah. Form produk, ketujuh pemilih, dialog Save/Edit, View Rate, dan lampiran dibuka
// DARI DALAM halaman itu, bukan tombol menu.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di
// berkas bersama.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_MPNL } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `modul/masterproductnamelife/backend/modul.go`. */
export const NAMA_MPNL = 'masterproductnamelife'

/**
 * Halaman modul ini - SATU: `InboxProductName` (grid produk; form, pemilih, dan dialog dibuka dari dalamnya).
 * Harness `InwardProductName` tidak dibangun (tombol `Inward` b75368 `OTHER 1=2`, RALAT R11).
 */
export const HALAMAN_MPNL = ['mpnl-produk'] as const
export type HalamanMPNL = (typeof HALAMAN_MPNL)[number]

/**
 * Halaman awal - bukti: `Section/InboxProductName.xml` tidak dibuka rule mana pun di korpus (dibuka portal);
 * ia hanya dirujuk sebagai sasaran `refresh otherSection` oleh enam section pemilih (`Ceding_Section.xml` b2285,
 * `SOB_Section.xml` b2300, `PolicyHolder_Section.xml` b2318, `Currency_Section.xml` b2323, `RIRISK_Section.xml`
 * b2307, `CauseOfLoss_Section.xml` b2242). Satu-satunya harness, `InwardProductName`, dibuka hanya oleh tombol
 * `Inward` b75368 yang bervisibilitas `OTHER 1=2` (PARITAS §1).
 */
export const HALAMAN_AWAL_MPNL: HalamanMPNL = 'mpnl-produk'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanMPNL> = {
  nama: NAMA_MPNL,
  kelompok: MENU_MPNL.kelompok,
  halaman: HALAMAN_MPNL,
  halamanAwal: HALAMAN_AWAL_MPNL,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    masterproductnamelife: HalamanMPNL
  }
}
