// Menu modul Endorsement Life.
//
// Menu DATAR (keputusan work owner 30-09-2026, `PROMPT-MENU-DATAR-PER-GROUPMENU.md`): satu modul satu
// menu - tombol "Endorsement Life" di bawah GROUPMENU TREATY (baris `M_NAV_MENU` isi awal 900,
// dinyalakan slot `976_menu_endorsementlife.sql`) membuka halaman awal di bawah. Layar buat
// endorsement, rincian kasus, dan popup dibuka DARI DALAM halaman ini, bukan tombol menu.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob`.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_EDM } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `modul/endorsementlife/backend/modul.go`. */
export const NAMA_EDM = 'endorsementlife'

/**
 * Halaman modul ini - SATU: kotak masuk `InboxEndorsementLife`. Harness `EndorsmentLife_harnes`
 * (`Create Addendum`), layar kasus `InputEDMLife`, dan popup dibuka dari dalamnya (PARITAS §1).
 */
export const HALAMAN_EDM = ['edm-inbox'] as const
export type HalamanEDM = (typeof HALAMAN_EDM)[number]

/**
 * Halaman awal - bukti: `Harness/InboxEndorsementLife.xml` b27 `DATA-PORTAL!INBOXENDORSEMENTLIFE`, b91 kelas
 * `Data-Portal` (satu-satunya harness portal modul ini), b142 label `Endorsement Life`, b843 `pyInclude`
 * `InboxEndorsementLife`. Tujuh harness lain ber-kelas work/int dan dibuka aksi (`showHarness`).
 */
export const HALAMAN_AWAL_EDM: HalamanEDM = 'edm-inbox'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanEDM> = {
  nama: NAMA_EDM,
  kelompok: MENU_EDM.kelompok,
  halaman: HALAMAN_EDM,
  halamanAwal: HALAMAN_AWAL_EDM,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    endorsementlife: HalamanEDM
  }
}
